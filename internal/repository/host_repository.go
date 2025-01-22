package repository

import (
	"middleware-platform/internal/model"

	"gorm.io/gorm"
)

type HostRepository struct {
	db *gorm.DB
}

func NewHostRepository(db *gorm.DB) *HostRepository {
	db.AutoMigrate(&model.Host{}, &model.FileSync{})
	return &HostRepository{db: db}
}

func (r *HostRepository) FindAll() ([]model.Host, error) {
	var hosts []model.Host
	result := r.db.Find(&hosts)
	return hosts, result.Error
}

func (r *HostRepository) FindByID(id uint) (*model.Host, error) {
	var host model.Host
	if err := r.db.First(&host, id).Error; err != nil {
		return nil, err
	}
	return &host, nil
}

func (r *HostRepository) Create(host *model.Host) error {
	return r.db.Create(host).Error
}

func (r *HostRepository) Update(host *model.Host) error {
	return r.db.Save(host).Error
}

func (r *HostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Host{}, id).Error
}

func (r *HostRepository) CreateFileSync(fileSync *model.FileSync) error {
	return r.db.Create(fileSync).Error
}

func (r *HostRepository) UpdateFileSync(fileSync *model.FileSync) error {
	return r.db.Save(fileSync).Error
}

func (r *HostRepository) FindFileSyncsByHostID(hostID uint) ([]model.FileSync, error) {
	var fileSyncs []model.FileSync
	result := r.db.Where("host_id = ?", hostID).Find(&fileSyncs)
	return fileSyncs, result.Error
}

func (r *HostRepository) CreateSyncHistory(history *model.FileSyncHistory) error {
	return r.db.Create(history).Error
}

func (r *HostRepository) FindSyncHistoryByFileSyncID(fileSyncID uint) ([]model.FileSyncHistory, error) {
	var histories []model.FileSyncHistory
	result := r.db.Where("file_sync_id = ?", fileSyncID).Order("created_at DESC").Find(&histories)
	return histories, result.Error
}

// FindFileSyncByID 根据ID查找文件同步任务
func (r *HostRepository) FindFileSyncByID(id uint) (*model.FileSync, error) {
	var fileSync model.FileSync
	if err := r.db.First(&fileSync, id).Error; err != nil {
		return nil, err
	}
	return &fileSync, nil
} 