CREATE TABLE IF NOT EXISTS configs (
    id UUID PRIMARY KEY,
    key VARCHAR(100) NOT NULL,
    value TEXT,
    scope VARCHAR(20) NOT NULL,
    scope_id UUID,
    created_by UUID,
    updated_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- 创建联合唯一索引，确保同一作用域下的配置键是唯一的
CREATE UNIQUE INDEX idx_config_scope_key ON configs (key, scope, COALESCE(scope_id, '00000000-0000-0000-0000-000000000000'));

CREATE INDEX idx_configs_scope ON configs(scope);
CREATE INDEX idx_configs_scope_id ON configs(scope_id);
CREATE INDEX idx_configs_deleted_at ON configs(deleted_at);

-- 默认系统配置
INSERT INTO configs (id, key, value, scope) VALUES 
    (uuid_generate_v4(), 'system.name', '智能体构建平台', 'system'),
    (uuid_generate_v4(), 'system.version', '0.1.0', 'system'),
    (uuid_generate_v4(), 'system.max_file_size', '10485760', 'system'),
    (uuid_generate_v4(), 'system.allowed_file_types', 'pdf,txt,doc,docx,xls,xlsx,csv,json,md', 'system'),
    (uuid_generate_v4(), 'system.default_model', 'gpt-3.5-turbo', 'system')
ON CONFLICT DO NOTHING; 