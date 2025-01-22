import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message, Space, Popconfirm } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import { PlusOutlined, EditOutlined, DeleteOutlined, DownloadOutlined } from '@ant-design/icons';
import { getMiddlewareList, createMiddleware, updateMiddleware, deleteMiddleware, exportMiddlewareList, Middleware } from '../../api/middleware';

const MiddlewarePage: React.FC = () => {
  const [list, setList] = useState<Middleware[]>([]);
  const [visible, setVisible] = useState(false);
  const [loading, setLoading] = useState(false);
  const [currentItem, setCurrentItem] = useState<Middleware | null>(null);
  const [form] = Form.useForm();
  const [exportModalVisible, setExportModalVisible] = useState(false);
  const [exportForm] = Form.useForm();

  const fetchList = async () => {
    try {
      setLoading(true);
      const data = await getMiddlewareList();
      const formattedData = data.map(item => ({
        ...item,
        key: item.id.toString()
      }));
      setList(formattedData);
    } catch (err) {
      message.error('Failed to fetch middleware list');
      setList([]);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchList();
  }, []);

  const handleEdit = (record: Middleware) => {
    setCurrentItem(record);
    form.setFieldsValue(record);
    setVisible(true);
  };

  const handleDelete = async (id: number) => {
    try {
      await deleteMiddleware(id);
      message.success('Deleted successfully');
      fetchList();
    } catch (err) {
      message.error('Failed to delete middleware');
    }
  };

  const handleExport = async (values: { type?: string }) => {
    try {
      const blob = await exportMiddlewareList(values.type);
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = `middleware_list${values.type ? `_${values.type}` : ''}.csv`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
      setExportModalVisible(false);
      exportForm.resetFields();
    } catch (err) {
      message.error('Failed to export data');
    }
  };

  const columns: ColumnsType<Middleware> = [
    { 
      title: 'Name',
      dataIndex: 'name',
      key: 'name'
    },
    { 
      title: 'Type',
      dataIndex: 'type',
      key: 'type'
    },
    { 
      title: 'Version',
      dataIndex: 'version',
      key: 'version'
    },
    { 
      title: 'Host',
      dataIndex: 'host',
      key: 'host'
    },
    { 
      title: 'Port',
      dataIndex: 'port',
      key: 'port'
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
            onClick={() => handleEdit(record)}
          >
            Edit
          </Button>
          <Popconfirm
            title="Are you sure to delete this middleware?"
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

  const handleSubmit = async (values: Partial<Middleware>) => {
    try {
      if (currentItem) {
        await updateMiddleware(currentItem.id, values);
        message.success('Updated successfully');
      } else {
        await createMiddleware(values);
        message.success('Created successfully');
      }
      setVisible(false);
      form.resetFields();
      setCurrentItem(null);
      fetchList();
    } catch (err) {
      message.error(currentItem ? 'Failed to update middleware' : 'Failed to create middleware');
    }
  };

  const handleCancel = () => {
    setVisible(false);
    form.resetFields();
    setCurrentItem(null);
  };

  return (
    <div>
      <h1>Middleware Management</h1>
      <Space style={{ marginBottom: 16 }}>
        <Button 
          type="primary" 
          icon={<PlusOutlined />}
          onClick={() => setVisible(true)}
        >
          Add Middleware
        </Button>
        <Button
          icon={<DownloadOutlined />}
          onClick={() => setExportModalVisible(true)}
        >
          Export to CSV
        </Button>
      </Space>

      <Table<Middleware> 
        loading={loading}
        columns={columns} 
        dataSource={list}
        pagination={{ defaultPageSize: 10 }}
      />

      <Modal
        title={currentItem ? "Edit Middleware" : "Add Middleware"}
        open={visible}
        onCancel={handleCancel}
        onOk={() => form.submit()}
      >
        <Form form={form} onFinish={handleSubmit}>
          <Form.Item name="name" label="Name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="type" label="Type" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="redis">Redis</Select.Option>
              <Select.Option value="mysql">MySQL</Select.Option>
              <Select.Option value="postgresql">PostgreSQL</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="version" label="Version">
            <Input />
          </Form.Item>
          <Form.Item name="host" label="Host" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="port" label="Port" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>

      <Modal
        title="Export Middleware List"
        open={exportModalVisible}
        onCancel={() => {
          setExportModalVisible(false);
          exportForm.resetFields();
        }}
        onOk={() => exportForm.submit()}
      >
        <Form form={exportForm} onFinish={handleExport}>
          <Form.Item name="type" label="Type">
            <Select allowClear>
              <Select.Option value="redis">Redis</Select.Option>
              <Select.Option value="mysql">MySQL</Select.Option>
              <Select.Option value="postgresql">PostgreSQL</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default MiddlewarePage; 