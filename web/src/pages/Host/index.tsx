import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, InputNumber, message, Space, Popconfirm, Upload, Checkbox, Progress, Tabs, Tag } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PlusOutlined, EditOutlined, DeleteOutlined, SyncOutlined, UploadOutlined } from '@ant-design/icons';
import { getHostList, createHost, updateHost, deleteHost, syncFile, Host, FileSync, FileSyncHistory, getSyncHistory, getFileSyncs, pauseSync, resumeSync, cancelSync } from '../../api/host';

const HostPage: React.FC = () => {
  const [list, setList] = useState<Host[]>([]);
  const [loading, setLoading] = useState(false);
  const [hostModalVisible, setHostModalVisible] = useState(false);
  const [syncModalVisible, setSyncModalVisible] = useState(false);
  const [currentHost, setCurrentHost] = useState<Host | null>(null);
  const [hostForm] = Form.useForm();
  const [syncForm] = Form.useForm();
  const [syncingFiles, setSyncingFiles] = useState<{[key: number]: FileSync}>({});

  const fetchList = async () => {
    try {
      setLoading(true);
      const data = await getHostList();
      setList(data);
    } catch (err) {
      message.error('Failed to fetch host list');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchList();
  }, []);

  // 定期轮询同步状态
  useEffect(() => {
    const pollSyncStatus = async () => {
      const syncingIds = Object.keys(syncingFiles).map(Number);
      if (syncingIds.length === 0) return;

      try {
        for (const hostId of syncingIds) {
          const syncs = await getFileSyncs(hostId);
          const updatedSyncs = syncs.reduce((acc, sync) => {
            acc[sync.id] = sync;
            return acc;
          }, {} as {[key: number]: FileSync});
          setSyncingFiles(prev => ({...prev, ...updatedSyncs}));

          // 如果所有文件都同步完成，停止轮询
          const allCompleted = syncs.every(sync => 
            sync.status === 'completed' || sync.status === 'failed'
          );
          if (allCompleted) {
            delete syncingFiles[hostId];
            setSyncingFiles({...syncingFiles});
          }
        }
      } catch (err) {
        console.error('Failed to poll sync status:', err);
      }
    };

    const timer = setInterval(pollSyncStatus, 1000);
    return () => clearInterval(timer);
  }, [syncingFiles]);

  const handleHostSubmit = async (values: Partial<Host>) => {
    try {
      if (currentHost) {
        await updateHost(currentHost.id, values);
        message.success('Updated successfully');
      } else {
        await createHost(values);
        message.success('Created successfully');
      }
      setHostModalVisible(false);
      hostForm.resetFields();
      setCurrentHost(null);
      fetchList();
    } catch (err) {
      message.error(currentHost ? 'Failed to update host' : 'Failed to create host');
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteHost(id);
      message.success('Deleted successfully');
      fetchList();
    } catch (err) {
      message.error('Failed to delete host');
    }
  };

  const handleSync = async (values: Partial<FileSync>) => {
    try {
      await syncFile(values);
      message.success('File sync started');
      setSyncModalVisible(false);
      syncForm.resetFields();

      // 添加到同步文件列表
      if (currentHost) {
        setSyncingFiles(prev => ({
          ...prev,
          [currentHost.id]: {
            ...values,
            progress: 0,
            status: 'syncing',
          } as FileSync,
        }));
      }
    } catch (err) {
      message.error('Failed to sync file');
    }
  };

  const columns: ColumnsType<Host> = [
    { 
      title: 'Name',
      dataIndex: 'name',
      key: 'name'
    },
    { 
      title: 'IP',
      dataIndex: 'ip',
      key: 'ip'
    },
    { 
      title: 'Port',
      dataIndex: 'port',
      key: 'port'
    },
    { 
      title: 'Username',
      dataIndex: 'username',
      key: 'username'
    },
    { 
      title: 'Status',
      dataIndex: 'status',
      key: 'status'
    },
    {
      title: 'Action',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => {
              setCurrentHost(record);
              hostForm.setFieldsValue(record);
              setHostModalVisible(true);
            }}
          >
            Edit
          </Button>
          <Button
            type="link"
            icon={<SyncOutlined />}
            onClick={() => {
              setCurrentHost(record);
              setSyncModalVisible(true);
            }}
          >
            Sync Files
          </Button>
          <Popconfirm
            title="Are you sure to delete this host?"
            onConfirm={() => handleDelete(record.id)}
            okText="Yes"
            cancelText="No"
          >
            <Button type="link" danger icon={<DeleteOutlined />}>
              Delete
            </Button>
          </Popconfirm>
        </Space>
      ),
    }
  ];

  return (
    <div>
      <Tabs defaultActiveKey="hosts">
        <Tabs.TabPane tab="Host Management" key="hosts">
          <Space style={{ marginBottom: 16 }}>
            <Button 
              type="primary" 
              icon={<PlusOutlined />}
              onClick={() => setHostModalVisible(true)}
            >
              Add Host
            </Button>
          </Space>

          <Table<Host> 
            loading={loading}
            columns={columns} 
            dataSource={list}
            rowKey="id"
            expandable={{
              expandedRowRender: (record) => (
                <div>
                  <Space style={{ marginBottom: 16 }}>
                    <Button
                      type="primary"
                      icon={<SyncOutlined />}
                      onClick={() => {
                        setCurrentHost(record);
                        setSyncModalVisible(true);
                      }}
                    >
                      Sync New File
                    </Button>
                  </Space>
                  <FileSyncList hostId={record.id} />
                </div>
              ),
            }}
          />
        </Tabs.TabPane>

        <Tabs.TabPane tab="File Sync Status" key="syncs">
          {/* 添加同步进度显示 */}
          <div>
            <h3>Active File Syncs</h3>
            {Object.values(syncingFiles).length > 0 ? (
              Object.values(syncingFiles).map(sync => (
                <div key={sync.id} style={{ marginBottom: 16, padding: 16, border: '1px solid #f0f0f0', borderRadius: 8 }}>
                  <h4>Host: {list.find(h => h.id === sync.hostId)?.name}</h4>
                  <p>
                    <strong>Source:</strong> {sync.sourcePath}<br/>
                    <strong>Target:</strong> {sync.targetPath}
                  </p>
                  <SyncProgress sync={sync} />
                </div>
              ))
            ) : (
              <div>No active file syncs</div>
            )}
          </div>
        </Tabs.TabPane>
      </Tabs>

      <Modal
        title={currentHost ? "Edit Host" : "Add Host"}
        open={hostModalVisible}
        onCancel={() => {
          setHostModalVisible(false);
          hostForm.resetFields();
          setCurrentHost(null);
        }}
        onOk={() => hostForm.submit()}
      >
        <Form form={hostForm} onFinish={handleHostSubmit}>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="ip" label="IP" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="port" label="Port" rules={[{ required: true }]}>
            <InputNumber min={1} max={65535} />
          </Form.Item>
          <Form.Item name="username" label="Username" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="password" label="Password">
            <Input.Password />
          </Form.Item>
          <Form.Item name="sshKey" label="SSH Key">
            <Input.TextArea />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="Sync Files"
        open={syncModalVisible}
        onCancel={() => {
          setSyncModalVisible(false);
          syncForm.resetFields();
        }}
        onOk={() => syncForm.submit()}
      >
        <Form 
          form={syncForm} 
          onFinish={(values) => handleSync({ ...values, hostId: currentHost?.id })}
        >
          <Form.Item name="sourcePath" label="Source Path" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="targetPath" label="Target Path" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="description" label="Description">
            <Input.TextArea />
          </Form.Item>
          <Form.Item name="isIncremental" valuePropName="checked">
            <Checkbox>Incremental Sync</Checkbox>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

// 修改 FileSyncList 组件
const FileSyncList: React.FC<{ hostId: number }> = ({ hostId }) => {
  const [syncs, setSyncs] = useState<FileSync[]>([]);
  const [loading, setLoading] = useState(false);
  const [historyModalVisible, setHistoryModalVisible] = useState(false);
  const [syncHistory, setSyncHistory] = useState<FileSyncHistory[]>([]);

  const fetchSyncs = async () => {
    try {
      setLoading(true);
      const data = await getFileSyncs(hostId);
      setSyncs(data);
    } catch (err) {
      message.error('Failed to fetch sync list');
    } finally {
      setLoading(false);
    }
  };

  const showHistory = async (fileSyncId: number) => {
    try {
      const data = await getSyncHistory(fileSyncId);
      setSyncHistory(data);
      setHistoryModalVisible(true);
    } catch (err) {
      message.error('Failed to fetch sync history');
    }
  };

  useEffect(() => {
    fetchSyncs();
  }, [hostId]);

  const columns: ColumnsType<FileSync> = [
    { 
      title: 'Source Path',
      dataIndex: 'sourcePath',
      key: 'sourcePath',
    },
    { 
      title: 'Target Path',
      dataIndex: 'targetPath',
      key: 'targetPath',
    },
    { 
      title: 'Status',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={
          status === 'completed' ? 'success' :
          status === 'failed' ? 'error' :
          status === 'syncing' ? 'processing' :
          status === 'paused' ? 'warning' :
          'default'
        }>
          {status}
        </Tag>
      ),
    },
    { 
      title: 'Last Sync',
      dataIndex: 'lastSyncAt',
      key: 'lastSyncAt',
    },
    {
      title: 'Action',
      key: 'action',
      render: (_, record) => (
        <Space>
          {record.status === 'syncing' && (
            <Button size="small" onClick={() => pauseSync(record.id)}>Pause</Button>
          )}
          {record.status === 'paused' && (
            <Button size="small" type="primary" onClick={() => resumeSync(record.id)}>Resume</Button>
          )}
          <Button size="small" onClick={() => showHistory(record.id)}>History</Button>
          <Button 
            size="small" 
            type="primary"
            onClick={() => {
              syncFile({
                hostId,
                sourcePath: record.sourcePath,
                targetPath: record.targetPath,
                isIncremental: true,
              });
            }}
          >
            Resync
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <>
      <Table<FileSync>
        loading={loading}
        columns={columns}
        dataSource={syncs}
        rowKey="id"
        size="small"
        pagination={false}
        title={() => <strong>Sync History</strong>}
      />

      <Modal
        title="Sync History"
        open={historyModalVisible}
        onCancel={() => setHistoryModalVisible(false)}
        footer={null}
      >
        <Table
          dataSource={syncHistory}
          columns={[
            { title: 'Time', dataIndex: 'createdAt', key: 'createdAt' },
            { title: 'Status', dataIndex: 'status', key: 'status' },
            { title: 'Type', dataIndex: 'syncType', key: 'syncType' },
            { title: 'Message', dataIndex: 'message', key: 'message' },
            { title: 'MD5', dataIndex: 'md5', key: 'md5' },
            { 
              title: 'Size', 
              dataIndex: 'fileSize', 
              key: 'fileSize',
              render: (size: number) => `${(size / 1024 / 1024).toFixed(2)} MB`
            },
          ]}
          pagination={false}
        />
      </Modal>
    </>
  );
};

// 添加同步进度显示组件
const SyncProgress: React.FC<{sync: FileSync}> = ({sync}) => {
  const handlePause = async () => {
    try {
      await pauseSync(sync.id);
      message.success('Sync paused');
    } catch (err) {
      message.error('Failed to pause sync');
    }
  };

  const handleResume = async () => {
    try {
      await resumeSync(sync.id);
      message.success('Sync resumed');
    } catch (err) {
      message.error('Failed to resume sync');
    }
  };

  const handleCancel = async () => {
    try {
      await cancelSync(sync.id);
      message.success('Sync cancelled');
    } catch (err) {
      message.error('Failed to cancel sync');
    }
  };

  const formatSpeed = (speed?: number) => {
    if (!speed) return '0 B/s';
    if (speed < 1024) return `${speed.toFixed(2)} B/s`;
    if (speed < 1024 * 1024) return `${(speed / 1024).toFixed(2)} KB/s`;
    return `${(speed / 1024 / 1024).toFixed(2)} MB/s`;
  };

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
        <Progress 
          percent={Math.round(sync.progress || 0)} 
          status={
            sync.status === 'failed' ? 'exception' :
            sync.status === 'cancelled' ? 'exception' :
            sync.status === 'paused' ? 'normal' :
            'active'
          }
          style={{ flex: 1 }}
        />
        <Space>
          {sync.status === 'syncing' && (
            <Button size="small" onClick={handlePause}>
              Pause
            </Button>
          )}
          {sync.status === 'paused' && (
            <Button size="small" type="primary" onClick={handleResume}>
              Resume
            </Button>
          )}
          {['syncing', 'paused'].includes(sync.status) && (
            <Button size="small" danger onClick={handleCancel}>
              Cancel
            </Button>
          )}
        </Space>
      </div>
      <div>Speed: {formatSpeed(sync.speed)}</div>
      <div>Status: {sync.status}</div>
      {sync.syncedSize && sync.fileSize && (
        <div>
          Progress: {(sync.syncedSize / 1024 / 1024).toFixed(2)}MB / 
          {(sync.fileSize / 1024 / 1024).toFixed(2)}MB
        </div>
      )}
    </div>
  );
};

export default HostPage; 