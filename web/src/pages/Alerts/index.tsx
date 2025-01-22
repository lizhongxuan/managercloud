import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Select, message } from 'antd';
import { getAlertsList, createAlertRule } from '../../api/alert';

const AlertsPage: React.FC = () => {
  const [alerts, setAlerts] = useState([]);
  const [visible, setVisible] = useState(false);
  const [form] = Form.useForm();

  const fetchAlerts = async () => {
    try {
      const data = await getAlertsList();
      setAlerts(data);
    } catch (err) {
      message.error('Failed to fetch alerts');
    }
  };

  useEffect(() => {
    fetchAlerts();
  }, []);

  const columns = [
    { title: 'Type', dataIndex: 'type' },
    { title: 'Target', dataIndex: 'target' },
    { title: 'Threshold', dataIndex: 'threshold' },
    { title: 'Status', dataIndex: 'status' },
  ];

  const handleCreate = async (values: any) => {
    try {
      await createAlertRule(values);
      message.success('Created successfully');
      setVisible(false);
      fetchAlerts();
    } catch (err) {
      message.error('Failed to create alert rule');
    }
  };

  return (
    <div>
      <Button 
        type="primary" 
        onClick={() => setVisible(true)}
        style={{ marginBottom: 16 }}
      >
        Add Alert Rule
      </Button>

      <Table 
        columns={columns} 
        dataSource={alerts}
        rowKey="id"
      />

      <Modal
        title="Add Alert Rule"
        visible={visible}
        onCancel={() => setVisible(false)}
        onOk={() => form.submit()}
      >
        <Form form={form} onFinish={handleCreate}>
          <Form.Item name="type" label="Type" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="cpu_usage">CPU Usage</Select.Option>
              <Select.Option value="memory_usage">Memory Usage</Select.Option>
              <Select.Option value="connections">Connections</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="target" label="Target" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="threshold" label="Threshold" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
          <Form.Item name="operator" label="Operator" rules={[{ required: true }]}>
            <Select>
              <Select.Option value=">">Greater Than</Select.Option>
              <Select.Option value="<">Less Than</Select.Option>
              <Select.Option value="=">Equal To</Select.Option>
            </Select>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AlertsPage; 