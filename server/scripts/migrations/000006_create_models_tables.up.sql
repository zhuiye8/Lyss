-- 创建模型表
CREATE TABLE IF NOT EXISTS models (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    provider VARCHAR(32) NOT NULL,
    model_id VARCHAR(64) NOT NULL,
    type VARCHAR(32) NOT NULL,
    description TEXT,
    capabilities TEXT,
    parameters JSONB DEFAULT '{}'::jsonb,
    max_tokens INTEGER DEFAULT 0,
    token_cost_prompt DOUBLE PRECISION DEFAULT 0,
    token_cost_completion DOUBLE PRECISION DEFAULT 0,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    provider_config JSONB DEFAULT '{}'::jsonb,
    is_system BOOLEAN DEFAULT FALSE,
    is_custom BOOLEAN DEFAULT FALSE,
    organization_id UUID,
    created_by UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 创建模型配置表
CREATE TABLE IF NOT EXISTS model_configs (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    model_id UUID NOT NULL REFERENCES models(id),
    parameters JSONB DEFAULT '{}'::jsonb,
    provider_config JSONB DEFAULT '{}'::jsonb,
    is_shared BOOLEAN DEFAULT FALSE,
    usage_metrics TEXT,
    organization_id UUID NOT NULL,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- 索引
CREATE INDEX idx_models_provider ON models(provider);
CREATE INDEX idx_models_type ON models(type);
CREATE INDEX idx_models_status ON models(status);
CREATE INDEX idx_models_org_id ON models(organization_id);

CREATE INDEX idx_model_configs_model_id ON model_configs(model_id);
CREATE INDEX idx_model_configs_org_id ON model_configs(organization_id);
CREATE INDEX idx_model_configs_created_by ON model_configs(created_by);
CREATE INDEX idx_model_configs_shared ON model_configs(is_shared);

-- 插入默认模型
INSERT INTO models (
    id, name, provider, model_id, type, description, capabilities, 
    parameters, max_tokens, token_cost_prompt, token_cost_completion, status,
    is_system, is_custom
) VALUES 
-- OpenAI 模型
(
    gen_random_uuid(), 'GPT-4 Turbo', 'openai', 'gpt-4-0125-preview', 'text', 
    'OpenAI最新的高性能GPT-4大型语言模型，具有更强的推理能力和更广泛的知识',
    '["text-generation", "function-calling", "code-generation", "reasoning"]',
    '{"temperature": 0.7, "top_p": 1, "max_tokens": 4096}',
    128000, 0.00001, 0.00003, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'GPT-4', 'openai', 'gpt-4', 'text',
    'OpenAI的高性能大型语言模型，擅长复杂推理和精确执行任务',
    '["text-generation", "function-calling", "code-generation", "reasoning"]',
    '{"temperature": 0.7, "top_p": 1, "max_tokens": 4096}',
    8192, 0.00003, 0.00006, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'GPT-3.5 Turbo', 'openai', 'gpt-3.5-turbo', 'text',
    'OpenAI的性能与成本平衡的大型语言模型，适合大多数常见任务',
    '["text-generation", "function-calling", "code-generation"]',
    '{"temperature": 0.7, "top_p": 1, "max_tokens": 4096}',
    16385, 0.0000015, 0.000002, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'text-embedding-3-large', 'openai', 'text-embedding-3-large', 'embedding',
    'OpenAI的最新高性能嵌入模型，生成高维度嵌入向量，适合高精度相似度匹配',
    '["embedding-generation"]',
    '{}',
    8191, 0.00000013, 0, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'text-embedding-3-small', 'openai', 'text-embedding-3-small', 'embedding',
    'OpenAI的轻量级嵌入模型，生成低维度嵌入向量，适合大规模处理',
    '["embedding-generation"]',
    '{}',
    8191, 0.00000002, 0, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'DALL-E 3', 'openai', 'dall-e-3', 'multimodal',
    'OpenAI的图像生成模型，可以根据文本描述生成高质量、符合要求的图像',
    '["image-generation"]',
    '{}',
    4096, 0.04, 0, 'active', TRUE, FALSE
),

-- Anthropic 模型
(
    gen_random_uuid(), 'Claude 3 Opus', 'anthropic', 'claude-3-opus-20240229', 'text',
    'Anthropic最先进的大型语言模型，擅长复杂推理、分析和创意任务',
    '["text-generation", "function-calling", "code-generation", "reasoning", "multimodal"]',
    '{"temperature": 0.7, "top_p": 0.9, "max_tokens": 4096}',
    200000, 0.00001, 0.00003, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'Claude 3 Sonnet', 'anthropic', 'claude-3-sonnet-20240229', 'text',
    'Anthropic的平衡性能和成本的中端模型，适合大多数企业应用',
    '["text-generation", "function-calling", "code-generation", "reasoning", "multimodal"]',
    '{"temperature": 0.7, "top_p": 0.9, "max_tokens": 4096}',
    200000, 0.000003, 0.000015, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'Claude 3 Haiku', 'anthropic', 'claude-3-haiku-20240307', 'text',
    'Anthropic的轻量级模型，提供高速响应和成本效益，适合简单对话和大规模集成',
    '["text-generation", "function-calling", "code-generation", "multimodal"]',
    '{"temperature": 0.7, "top_p": 0.9, "max_tokens": 4096}',
    200000, 0.00000025, 0.00000125, 'active', TRUE, FALSE
),

-- 百度文心一言模型
(
    gen_random_uuid(), '文心一言4.0', 'baidu', 'ernie-4.0', 'text',
    '百度最新的中英双语大模型，具有很强的中文理解能力和知识',
    '["text-generation", "function-calling", "reasoning"]',
    '{"temperature": 0.8, "top_p": 0.8, "max_tokens": 2048}',
    3000, 0.000012, 0.000012, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'ERNIE-Speed', 'baidu', 'ernie-speed', 'text',
    '百度最新的轻量化模型，提供高性能、低成本的文本处理能力',
    '["text-generation", "function-calling"]',
    '{"temperature": 0.8, "top_p": 0.8, "max_tokens": 2048}',
    3000, 0.000005, 0.000005, 'active', TRUE, FALSE
),
(
    gen_random_uuid(), 'ERNIE-Embedding', 'baidu', 'ernie-embedding', 'embedding',
    '百度文心的文本嵌入模型，专注中文语义理解',
    '["embedding-generation"]',
    '{}',
    8192, 0.00000002, 0, 'active', TRUE, FALSE
); 