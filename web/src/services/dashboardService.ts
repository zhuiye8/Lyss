import api from './api';
import { 
  IStatisticData, 
  IUsageData, 
  ITopAgent, 
  IRecentActivity,
  IAgentData,
  IModelData,
  ISystemSettings
} from '../types/dashboard';

// 仪表盘数据
export const getStatistics = async (): Promise<IStatisticData> => {
  const response = await api.get('/api/v1/dashboard/statistics');
  return response.data.data;
};

export const getUsageTrend = async (days: number): Promise<IUsageData[]> => {
  const response = await api.get(`/api/v1/dashboard/trends?days=${days}`);
  return response.data.data;
};

export const getTopAgents = async (limit: number): Promise<ITopAgent[]> => {
  const response = await api.get(`/api/v1/dashboard/top-agents?limit=${limit}`);
  return response.data.data;
};

export const getRecentActivities = async (limit: number): Promise<IRecentActivity[]> => {
  const response = await api.get(`/api/v1/dashboard/activities?limit=${limit}`);
  return response.data.data;
};

// 应用管理
export const getAgents = async (page: number = 1, pageSize: number = 10, searchQuery?: string): Promise<{
  data: IAgentData[];
  total: number;
}> => {
  const params = { page, pageSize, searchQuery };
  const response = await api.get('/agents', { params });
  return response.data;
};

export const deleteAgent = async (id: string): Promise<void> => {
  await api.delete(`/agents/${id}`);
};

export const updateAgentStatus = async (id: string, status: 'active' | 'inactive'): Promise<void> => {
  await api.patch(`/agents/${id}/status`, { status });
};

// 模型管理
export const getModels = async (): Promise<IModelData[]> => {
  const response = await api.get('/api/v1/models');
  return response.data.data;
};

export const addModel = async (model: Partial<IModelData>): Promise<IModelData> => {
  const response = await api.post('/api/v1/models', model);
  return response.data.data;
};

export const updateModel = async (id: string, model: Partial<IModelData>): Promise<IModelData> => {
  const response = await api.put(`/api/v1/models/${id}`, model);
  return response.data.data;
};

export const deleteModel = async (id: string): Promise<void> => {
  await api.delete(`/api/v1/models/${id}`);
};

export const testModelConnection = async (model: Partial<IModelData>): Promise<boolean> => {
  try {
    const response = await api.post('/api/v1/models/test-connection', model);
    return response.data.success;
  } catch (error) {
    return false;
  }
};

// 系统设置
export const getSystemSettings = async (): Promise<ISystemSettings> => {
  const response = await api.get('/api/v1/settings/system');
  return response.data.data;
};

export const updateSystemSettings = async (settings: Partial<ISystemSettings>): Promise<ISystemSettings> => {
  const response = await api.put('/api/v1/settings/system', settings);
  return response.data.data;
}; 