package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"middleware-platform/internal/model"
	"middleware-platform/internal/repository"
	"os"
	"sync"
	"time"
	"log"
	"golang.org/x/crypto/ssh"
)

type HostService struct {
	repo *repository.HostRepository
	// 添加同步任务管理
	syncTasks map[uint]chan struct{} // key: FileSyncID, value: control channel
	mu        sync.RWMutex
}

func NewHostService(repo *repository.HostRepository) *HostService {
	return &HostService{
		repo:      repo,
		syncTasks: make(map[uint]chan struct{}),
	}
}

func (s *HostService) GetAll() ([]model.Host, error) {
	return s.repo.FindAll()
}

func (s *HostService) GetByID(id uint) (*model.Host, error) {
	return s.repo.FindByID(id)
}

func (s *HostService) Create(host *model.Host) error {
	// 测试SSH连接
	if err := s.testConnection(host); err != nil {
		return fmt.Errorf("failed to connect to host: %v", err)
	}
	host.Status = "connected"
	return s.repo.Create(host)
}

func (s *HostService) Update(host *model.Host) error {
	return s.repo.Update(host)
}

func (s *HostService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *HostService) testConnection(host *model.Host) error {
	config := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	if host.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(host.Password))
	}

	if host.SSHKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(host.SSHKey))
		if err != nil {
			return fmt.Errorf("failed to parse SSH key: %v", err)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), config)
	if err != nil {
		log.Printf("Failed to ssh host: %v", err)
		return err
	}
	defer client.Close()
	return nil
}

// 计算文件MD5
func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// 获取文件信息
func getFileInfo(filePath string) (int64, int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Failed to get file %s info: %v", filePath, err)
		return 0, 0, err
	}
	return info.ModTime().Unix(), info.Size(), nil
}

// 添加进度跟踪的 io.Reader 包装器
type progressReader struct {
	io.Reader
	total    int64
	current  int64
	progress chan<- float64
	speed    chan<- float64
	start    time.Time
	control  <-chan struct{}
	fileSync *model.FileSync
	repo     *repository.HostRepository
}

func (pr *progressReader) Read(p []byte) (int, error) {
	// 检查是否取消
	select {
	case <-pr.control:
		return 0, fmt.Errorf("sync cancelled")
	default:
	}

	// 检查是否暂停
	if pr.fileSync.IsPaused {
		time.Sleep(time.Second)
		return 0, nil
	}

	n, err := pr.Reader.Read(p)
	if n > 0 {
		pr.current += int64(n)
		progress := float64(pr.current) * 100 / float64(pr.total)
		pr.progress <- progress

		elapsed := time.Since(pr.start).Seconds()
		if elapsed > 0 {
			speed := float64(pr.current) / elapsed
			pr.speed <- speed
		}

		// 更新已同步大小
		pr.fileSync.SyncedSize = pr.current
		pr.repo.UpdateFileSync(pr.fileSync)
	}
	return n, err
}

// PauseSync 暂停同步
func (s *HostService) PauseSync(fileSyncID uint) error {
	fileSync, err := s.repo.FindFileSyncByID(fileSyncID)
	if err != nil {
		return err
	}

	if fileSync.Status != "syncing" {
		return fmt.Errorf("file sync is not in progress")
	}

	fileSync.IsPaused = true
	fileSync.Status = "paused"
	return s.repo.UpdateFileSync(fileSync)
}

// ResumeSync 恢复同步
func (s *HostService) ResumeSync(fileSyncID uint) error {
	fileSync, err := s.repo.FindFileSyncByID(fileSyncID)
	if err != nil {
		return err
	}

	if fileSync.Status != "paused" {
		return fmt.Errorf("file sync is not paused")
	}

	fileSync.IsPaused = false
	fileSync.Status = "syncing"
	go s.SyncFile(fileSync) // 重新开始同步
	return s.repo.UpdateFileSync(fileSync)
}

// CancelSync 取消同步
func (s *HostService) CancelSync(fileSyncID uint) error {
	s.mu.Lock()
	if ch, exists := s.syncTasks[fileSyncID]; exists {
		close(ch)
		delete(s.syncTasks, fileSyncID)
	}
	s.mu.Unlock()

	fileSync, err := s.repo.FindFileSyncByID(fileSyncID)
	if err != nil {
		return err
	}

	fileSync.Status = "cancelled"
	return s.repo.UpdateFileSync(fileSync)
}

// SyncFile 同步文件到远程主机
func (s *HostService) SyncFile(fileSync *model.FileSync) error {
	// 创建控制通道
	s.mu.Lock()
	controlCh := make(chan struct{})
	s.syncTasks[fileSync.ID] = controlCh
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		delete(s.syncTasks, fileSync.ID)
		s.mu.Unlock()
	}()

	// 获取源文件信息
	modTime, size, err := getFileInfo(fileSync.SourcePath)
	if err != nil {
		log.Printf("Failed to get file info: %v", err)
		return err
	}

	// 计算MD5
	md5sum, err := calculateFileMD5(fileSync.SourcePath)
	if err != nil {
		log.Printf("Failed to calculate file MD5: %v", err)
		return err
	}

	// 检查是否需要同步
	if fileSync.IsIncremental && fileSync.MD5 == md5sum {
		// 文件未变化，不需要同步
		return nil
	}

	// 执行同步
	if err := s.doSync(fileSync, size, controlCh); err != nil {
		// 记录同步失败历史
		log.Printf("Failed to sync file: %v", err)
		s.repo.CreateSyncHistory(&model.FileSyncHistory{
			FileSyncID: fileSync.ID,
			Status:     "failed",
			Message:    err.Error(),
			MD5:        md5sum,
			FileSize:   size,
			SyncType:   getSyncType(fileSync.IsIncremental),
		})
		return err
	}

	// 更新同步状态
	fileSync.Status = "completed"
	fileSync.LastSyncAt = time.Now().Format("2006-01-02 15:04:05")
	fileSync.MD5 = md5sum
	fileSync.ModifiedTime = modTime
	fileSync.FileSize = size

	if err := s.repo.UpdateFileSync(fileSync); err != nil {
		log.Printf("Failed to update file sync: %v", err)
			return err
	}

	// 记录同步成功历史
	return s.repo.CreateSyncHistory(&model.FileSyncHistory{
		FileSyncID: fileSync.ID,
		Status:     "success",
		Message:    "sync completed",
		MD5:        md5sum,
		FileSize:   size,
		SyncType:   getSyncType(fileSync.IsIncremental),
	})
}

// doSync 执行实际的文件同步
func (s *HostService) doSync(fileSync *model.FileSync, size int64, controlCh chan struct{}) error {
	host, err := s.repo.FindByID(1)
	if err != nil {
		log.Printf("FindByID Failed to find host: %v", err)
		return err
	}

	// 创建SSH客户端配置
	config := &ssh.ClientConfig{
		User: host.Username,
		Auth: []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	if host.Password != "" {
		config.Auth = append(config.Auth, ssh.Password(host.Password))
	}

	if host.SSHKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(host.SSHKey))
		if err != nil {
			log.Printf("Failed to parse SSH key: %v", err)
			return fmt.Errorf("failed to parse SSH key: %v", err)
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	}

	// 连接到远程主机
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.IP, host.Port), config)
	if err != nil {
		log.Printf("Failed to ssh host: %v", err)
		return err
	}
	defer client.Close()

	// 创建SFTP会话
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	// 读取源文件
	src, err := os.Open(fileSync.SourcePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// 创建目标文件
	dst, err := session.StdinPipe()
	if err != nil {
		return err
	}

	// 启动scp命令
	if err := session.Start(fmt.Sprintf("scp -t %s", fileSync.TargetPath)); err != nil {
		return err
	}

	// 创建进度和速度通道
	progressCh := make(chan float64)
	speedCh := make(chan float64)
	defer close(progressCh)
	defer close(speedCh)

	// 启动goroutine更新进度
	go func() {
		for {
			select {
			case progress := <-progressCh:
				fileSync.Progress = progress
				s.repo.UpdateFileSync(fileSync)
			case speed := <-speedCh:
				fileSync.Speed = speed
				s.repo.UpdateFileSync(fileSync)
			case <-time.After(5 * time.Second):
				return
			}
		}
	}()

	// 创建带进度的Reader
	reader := &progressReader{
		Reader:   src,
		total:    size,
		progress: progressCh,
		speed:    speedCh,
		start:    time.Now(),
		control:  controlCh,
		fileSync: fileSync,
		repo:     s.repo,
	}

	// 复制文件内容
	_, err = io.Copy(dst, reader)
	return err
}

// getSyncType 根据是否增量同步返回同步类型
func getSyncType(isIncremental bool) string {
	if isIncremental {
		return "incremental"
	}
	return "full"
}

// GetSyncHistory 获取同步历史记录
func (s *HostService) GetSyncHistory(fileSyncID uint) ([]model.FileSyncHistory, error) {
	return s.repo.FindSyncHistoryByFileSyncID(fileSyncID)
}

// GetFileSyncsByHost 获取主机的文件同步任务
func (s *HostService) GetFileSyncsByHost(hostID uint) ([]model.FileSync, error) {
	return s.repo.FindFileSyncsByHostID(hostID)
} 