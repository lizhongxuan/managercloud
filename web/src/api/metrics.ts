import request from '../utils/request';

export interface Metrics {
  type: string;
  value: number;
  unit: string;
  timestamp: string;
}

export async function getMetricsStatus(middlewareId: number) {
  const response = await request.get(`/api/v1/metrics/status?id=${middlewareId}`);
  return response.data;
}

export async function getPerformanceMetrics(middlewareId: number) {
  const response = await request.get(`/api/v1/metrics/performance?id=${middlewareId}`);
  return response.data;
} 