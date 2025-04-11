import api from './api';
import { LogLevel, LogCategory, LogType } from '../types/log';

// 日志查询参数接口
export interface LogQueryParams {
  level?: LogLevel;
  category?: LogCategory;
  start_time?: string;
  end_time?: string;
  user_id?: string;
  request_id?: string;
  method?: string;
  path?: string;
  status_code?: number;
  min_duration?: number;
  max_duration?: number;
  error_code?: string;
  model_name?: string;
  project_id?: string;
  app_id?: string;
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
  type?: LogType;
}

// 获取日志列表
export const getLogs = (params: LogQueryParams) => {
  return api.get('/logs', { params });
};

// 获取日志详情
export const getLogById = (id: string) => {
  return api.get(`/logs/${id}`);
};

// 标记错误为已解决
export const markErrorAsResolved = (id: string) => {
  return api.patch(`/logs/${id}/resolve`);
};

// 获取日志统计信息
export const getLogStats = (startTime?: string, endTime?: string) => {
  const params = { start_time: startTime, end_time: endTime };
  return api.get('/logs/stats', { params });
};

// 获取系统监控指标
export const getMetrics = (timeRange: string = '1h') => {
  return api.get('/logs/metrics', { params: { range: timeRange } });
}; 