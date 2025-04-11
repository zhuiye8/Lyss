// 日志级别枚举
export enum LogLevel {
  DEBUG = 'debug',
  INFO = 'info',
  WARN = 'warn',
  ERROR = 'error',
  FATAL = 'fatal',
}

// 日志类别枚举
export enum LogCategory {
  SYSTEM = 'system',
  USER = 'user',
  AUTH = 'auth',
  API = 'api',
  MODEL = 'model',
  DATABASE = 'database',
  PERFORMANCE = 'performance',
}

// 日志类型枚举
export enum LogType {
  ALL = 'all',
  API = 'api',
  ERROR = 'error',
  MODEL_CALL = 'model_call',
}

// 用户信息接口
export interface User {
  id: string;
  username: string;
  email: string;
}

// 基础日志接口
export interface Log {
  id: string;
  level: LogLevel;
  category: LogCategory;
  message: string;
  user_id?: string;
  user?: User;
  metadata?: any;
  created_at: string;
}

// API日志接口
export interface APILog extends Log {
  method: string;
  path: string;
  status_code: number;
  ip?: string;
  user_agent?: string;
  duration: number;
  request_id?: string;
}

// 错误日志接口
export interface ErrorLog extends Log {
  stack_trace?: string;
  error_code?: string;
  source?: string;
  resolved_at?: string;
  resolved_by?: string;
}

// 模型调用日志接口
export interface ModelCallLog extends Log {
  model_name: string;
  prompt_tokens?: number;
  comp_tokens?: number;
  total_tokens?: number;
  duration: number;
  application_id?: string;
  project_id?: string;
  success: boolean;
}

// 系统指标接口
export interface SystemMetric {
  id: string;
  metric_name: string;
  metric_value: number;
  unit?: string;
  tags?: any;
  created_at: string;
}

// 日志响应接口，用于统一处理后端返回的日志数据
export interface LogResponse {
  data: Log | APILog | ErrorLog | ModelCallLog;
}

// 日志列表响应接口
export interface LogListResponse {
  data: (Log | APILog | ErrorLog | ModelCallLog)[];
  meta: {
    total: number;
    page: number;
    size: number;
  };
}

// 日志统计响应接口
export interface LogStatsResponse {
  data: {
    total_logs: number;
    error_count: number;
    api_count: number;
    model_call_count: number;
    avg_api_duration: number;
    avg_model_call_duration: number;
    error_by_category: Record<string, number>;
    logs_by_level: Record<string, number>;
  };
  meta: {
    start_time: string;
    end_time: string;
  };
}

// 系统指标响应接口
export interface MetricsResponse {
  data: {
    cpu_usage: SystemMetric[];
    memory_usage: SystemMetric[];
    disk_usage: SystemMetric[];
    network_traffic: SystemMetric[];
    api_latency: SystemMetric[];
    model_latency: SystemMetric[];
  };
} 