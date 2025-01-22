import request from '../utils/request';

export interface Middleware {
  id: number;
  name: string;
  type: string;
  version: string;
  host: string;
  port: string;
  status: string;
}

interface ApiResponse<T> {
  code: number;
  data: T;
  message: string;
}

export async function getMiddlewareList(): Promise<Middleware[]> {
  const response = await request.get<ApiResponse<Middleware[]>>('/api/v1/middleware/list');
  return response.data.data || [];
}

export async function createMiddleware(data: Partial<Middleware>): Promise<void> {
  await request.post<ApiResponse<void>>('/api/v1/middleware/create', data);
}

export async function updateMiddleware(id: number, data: Partial<Middleware>): Promise<void> {
  await request.put<ApiResponse<void>>(`/api/v1/middleware/${id}`, data);
}

export async function deleteMiddleware(id: number): Promise<void> {
  await request.delete<ApiResponse<void>>(`/api/v1/middleware/${id}`);
}

export async function exportMiddlewareList(type?: string): Promise<Blob> {
  const response = await request.get('/api/v1/middleware/export', {
    params: type ? { type } : {},
    responseType: 'blob'
  });
  return response.data;
} 