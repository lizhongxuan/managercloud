import React, { useState, useEffect } from 'react';
import { Card, Row, Col, Select } from 'antd';
import { Line } from '@ant-design/plots';
import { getPerformanceMetrics } from '../../api/metrics';

const MonitoringPage: React.FC = () => {
  const [metrics, setMetrics] = useState<any>({});
  const [selectedMiddleware, setSelectedMiddleware] = useState<number>(1);

  const fetchMetrics = async () => {
    try {
      const data = await getPerformanceMetrics(selectedMiddleware);
      setMetrics(data);
    } catch (err) {
      console.error('Failed to fetch metrics:', err);
    }
  };

  useEffect(() => {
    fetchMetrics();
    const timer = setInterval(fetchMetrics, 30000); // 每30秒刷新一次
    return () => clearInterval(timer);
  }, [selectedMiddleware]);

  return (
    <div>
      <Select
        style={{ width: 200, marginBottom: 16 }}
        value={selectedMiddleware}
        onChange={setSelectedMiddleware}
      >
        {/* 中间件列表选项 */}
      </Select>

      <Row gutter={16}>
        <Col span={12}>
          <Card title="Memory Usage">
            <Line
              data={metrics.memory_usage || []}
              xField="timestamp"
              yField="value"
              seriesField="type"
            />
          </Card>
        </Col>
        <Col span={12}>
          <Card title="Connections">
            <Line
              data={metrics.connections || []}
              xField="timestamp"
              yField="value"
              seriesField="type"
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default MonitoringPage; 