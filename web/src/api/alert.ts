import request from '../utils/request';

export interface AlertRule {
  id: number;
  type: string;
  target: string;
  threshold: string;
  operator: string;
  status: string;
}

export async function getAlertsList() {
  const response = await request.get('/api/v1/alerts/list');
  return response.data;
}

export async function createAlertRule(data: Partial<AlertRule>) {
  return request.post('/api/v1/alerts/rules', data);
}

export async function updateAlertRule(id: number, data: Partial<AlertRule>) {
  return request.put(`/api/v1/alerts/rules/${id}`, data);
} 