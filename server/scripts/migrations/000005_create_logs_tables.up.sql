-- 创建基础日志表
CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY,
    level VARCHAR(10) NOT NULL,
    category VARCHAR(20) NOT NULL,
    message TEXT NOT NULL,
    user_id UUID,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 创建API日志表，继承自logs表
CREATE TABLE IF NOT EXISTS api_logs (
    method VARCHAR(10) NOT NULL,
    path VARCHAR(255) NOT NULL,
    status_code INTEGER NOT NULL,
    ip VARCHAR(50),
    user_agent VARCHAR(255),
    duration BIGINT NOT NULL,
    request_id VARCHAR(36)
) INHERITS (logs);

-- 创建错误日志表，继承自logs表
CREATE TABLE IF NOT EXISTS error_logs (
    stack_trace TEXT,
    error_code VARCHAR(50),
    source VARCHAR(100),
    resolved_at TIMESTAMP WITH TIME ZONE,
    resolved_by UUID
) INHERITS (logs);

-- 创建模型调用日志表，继承自logs表
CREATE TABLE IF NOT EXISTS model_call_logs (
    model_name VARCHAR(100) NOT NULL,
    prompt_tokens INTEGER,
    comp_tokens INTEGER,
    total_tokens INTEGER,
    duration BIGINT NOT NULL,
    application_id UUID,
    project_id UUID,
    success BOOLEAN NOT NULL
) INHERITS (logs);

-- 创建系统指标表
CREATE TABLE IF NOT EXISTS system_metrics (
    id UUID PRIMARY KEY,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DOUBLE PRECISION NOT NULL,
    unit VARCHAR(20),
    tags JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 创建索引
CREATE INDEX idx_logs_level ON logs(level);
CREATE INDEX idx_logs_category ON logs(category);
CREATE INDEX idx_logs_user_id ON logs(user_id);
CREATE INDEX idx_logs_created_at ON logs(created_at);

CREATE INDEX idx_api_logs_method ON api_logs(method);
CREATE INDEX idx_api_logs_path ON api_logs(path);
CREATE INDEX idx_api_logs_status_code ON api_logs(status_code);
CREATE INDEX idx_api_logs_request_id ON api_logs(request_id);

CREATE INDEX idx_error_logs_error_code ON error_logs(error_code);
CREATE INDEX idx_error_logs_source ON error_logs(source);

CREATE INDEX idx_model_call_logs_model_name ON model_call_logs(model_name);
CREATE INDEX idx_model_call_logs_application_id ON model_call_logs(application_id);
CREATE INDEX idx_model_call_logs_project_id ON model_call_logs(project_id);

CREATE INDEX idx_system_metrics_metric_name ON system_metrics(metric_name);
CREATE INDEX idx_system_metrics_created_at ON system_metrics(created_at); 