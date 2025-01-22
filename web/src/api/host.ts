import request from '../utils/request';
import { ApiResponse } from '../types/api';

export interface Host {
  id: number;
  name: string;
  ip: string;
  port: number;
  username: string;
  password?: string;
  sshKey?: string;
  status: string;
  description?: string;
}

export interface FileSync {
  id: number;
  hostId: number;
  sourcePath: string;
  targetPath: string;
  status: string;
  lastSyncAt?: string;
  description?: string;
  md5?: string;
  modifiedTime?: number;
  fileSize?: number;
  isIncremental?: boolean;
  progress?: number;    // 同步进度 0-100
  speed?: number;      // 传输速度 bytes/s
  syncedSize?: number; // 已同步大小
  isPaused?: boolean;  // 是否暂停
}

export interface FileSyncHistory {
  id: number;
  fileSyncId: number;
  status: string;
  message: string;
  md5: string;
  fileSize: number;
  syncType: string;
  createdAt: string;
}

export async function getHostList() {
  const response = await request.get<ApiResponse<Host[]>>('/api/v1/hosts/list');
  return response.data.data || [];
}

export async function createHost(data: Partial<Host>) {
  await request.post<ApiResponse<void>>('/api/v1/hosts/create', data);
}

export async function updateHost(id: number, data: Partial<Host>) {
  await request.put<ApiResponse<void>>(`/api/v1/hosts/${id}`, data);
}

export async function deleteHost(id: number) {
  await request.delete<ApiResponse<void>>(`/api/v1/hosts/${id}`);
}

export async function syncFile(data: Partial<FileSync>) {
  await request.post<ApiResponse<void>>('/api/v1/hosts/sync', data);
}

export async function getFileSyncs(hostId: number) {
  const response = await request.get<ApiResponse<FileSync[]>>(`/api/v1/hosts/${hostId}/syncs`);
  return response.data.data || [];
}

export async function getSyncHistory(fileSyncId: number) {
  const response = await request.get<ApiResponse<FileSyncHistory[]>>(`/api/v1/hosts/syncs/${fileSyncId}/history`);
  return response.data.data || [];
}

export async function pauseSync(fileSyncId: number) {
  await request.post<ApiResponse<void>>(`/api/v1/hosts/syncs/${fileSyncId}/pause`);
}

export async function resumeSync(fileSyncId: number) {
  await request.post<ApiResponse<void>>(`/api/v1/hosts/syncs/${fileSyncId}/resume`);
}

export async function cancelSync(fileSyncId: number) {
  await request.post<ApiResponse<void>>(`/api/v1/hosts/syncs/${fileSyncId}/cancel`);
} 