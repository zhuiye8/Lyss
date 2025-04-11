# 智能体构建平台数据模型设计

## 1. 数据模型概览

本文档定义了智能体构建平台的核心数据模型，包括实体关系、属性定义和数据存储策略。系统采用关系型数据库(PostgreSQL)存储结构化数据，向量数据库(Milvus/Weaviate)存储向量数据，对象存储(MinIO)存储文件数据。

## 2. 实体关系图

```
┌───────────┐       ┌───────────┐       ┌───────────┐
│    用户    │◄──────┤   组织    │◄──────┤   团队    │
└───────────┘       └───────────┘       └───────────┘
      △                   △                   △
      │                   │                   │
      │                   │                   │
      │                   │                   │
      │              ┌────┴────┐              │
      └──────────────┤  应用   ├──────────────┘
                     └────┬────┘
                          │
         ┌────────────────┼────────────────┐
         │                │                │
         ▼                ▼                ▼
   ┌──────────┐     ┌──────────┐     ┌──────────┐
   │  智能体   │     │  知识库   │     │  模型配置 │
   └────┬─────┘     └────┬─────┘     └──────────┘
        │                │
        ▼                ▼
  ┌──────────┐     ┌──────────┐
  │   对话    │     │   文档    │
  └──────────┘     └──────────┘
```

## 3. 核心数据实体定义

### 3.1 用户 (User)

**描述**：代表系统的用户账户，存储用户的身份验证和个人信息。

**属性**：
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(64) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(128),
    avatar_url TEXT,
    bio TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个用户可以属于多个组织
- 一个用户可以创建多个应用
- 一个用户可以属于多个团队

### 3.2 组织 (Organization)

**描述**：代表一个组织或公司，可以包含多个用户和团队。

**属性**：
```sql
CREATE TABLE organizations (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    slug VARCHAR(64) NOT NULL UNIQUE,
    description TEXT,
    logo_url TEXT,
    website_url TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个组织可以有多个用户
- 一个组织可以有多个团队
- 一个组织可以有多个应用

### 3.3 团队 (Team)

**描述**：代表组织内的团队，用于组织用户和管理权限。

**属性**：
```sql
CREATE TABLE teams (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个团队属于一个组织
- 一个团队可以有多个用户
- 一个团队可以有多个应用

### 3.4 应用 (Application)

**描述**：代表用户创建的AI应用，是平台的核心实体。

**属性**：
```sql
CREATE TABLE applications (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    type VARCHAR(32) NOT NULL, -- 'chat', 'workflow', 'api'
    icon_url TEXT,
    api_key UUID UNIQUE,
    status VARCHAR(16) NOT NULL DEFAULT 'draft', -- 'draft', 'published', 'archived'
    visibility VARCHAR(16) NOT NULL DEFAULT 'private', -- 'private', 'team', 'public'
    organization_id UUID NOT NULL REFERENCES organizations(id),
    team_id UUID REFERENCES teams(id),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个应用属于一个组织
- 一个应用可以属于一个团队
- 一个应用由一个用户创建
- 一个应用可以有多个智能体
- 一个应用可以关联多个知识库
- 一个应用可以使用多个模型配置

### 3.5 智能体 (Agent)

**描述**：代表应用中的AI智能体，负责处理用户交互和执行任务。

**属性**：
```sql
CREATE TABLE agents (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    application_id UUID NOT NULL REFERENCES applications(id),
    model_config_id UUID NOT NULL REFERENCES model_configs(id),
    system_prompt TEXT,
    tools JSONB, -- 可用工具列表
    variables JSONB, -- 变量定义
    max_history_length INTEGER DEFAULT 10,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个智能体属于一个应用
- 一个智能体使用一个模型配置
- 一个智能体可以有多个对话
- 一个智能体可以关联多个知识库

### 3.6 模型配置 (ModelConfig)

**描述**：代表LLM模型的配置信息，包括模型提供者、参数等。

**属性**：
```sql
CREATE TABLE model_configs (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    provider VARCHAR(32) NOT NULL, -- 'openai', 'anthropic', 'local', etc.
    model VARCHAR(64) NOT NULL, -- 'gpt-4', 'claude-3-opus', etc.
    parameters JSONB, -- temperature, max_tokens, etc.
    organization_id UUID NOT NULL REFERENCES organizations(id),
    is_shared BOOLEAN DEFAULT FALSE,
    api_key_encrypted TEXT, -- 加密存储的API密钥
    base_url TEXT, -- 模型API的基础URL
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个模型配置属于一个组织
- 一个模型配置可以被多个智能体使用

### 3.7 知识库 (KnowledgeBase)

**描述**：代表智能体可以访问的知识库，用于存储和检索文档。

**属性**：
```sql
CREATE TABLE knowledge_bases (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    description TEXT,
    organization_id UUID NOT NULL REFERENCES organizations(id),
    embedding_model VARCHAR(64), -- 使用的嵌入模型
    chunk_size INTEGER DEFAULT 1000,
    chunk_overlap INTEGER DEFAULT 200,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个知识库属于一个组织
- 一个知识库可以被多个应用使用
- 一个知识库可以包含多个文档

### 3.8 文档 (Document)

**描述**：代表知识库中的文档，可以是各种格式的文件。

**属性**：
```sql
CREATE TABLE documents (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    knowledge_base_id UUID NOT NULL REFERENCES knowledge_bases(id),
    file_path TEXT NOT NULL, -- MinIO中的文件路径
    file_type VARCHAR(32) NOT NULL, -- 'pdf', 'docx', 'txt', etc.
    file_size BIGINT NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'processing', -- 'processing', 'indexed', 'failed'
    metadata JSONB,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个文档属于一个知识库
- 一个文档由一个用户上传

### 3.9 文档块 (DocumentChunk)

**描述**：代表文档分块后的片段，用于向量检索。

**属性**：
```sql
CREATE TABLE document_chunks (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(id),
    content TEXT NOT NULL,
    metadata JSONB,
    embedding_id VARCHAR(64), -- 在向量数据库中的ID
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个文档块属于一个文档
- 文档块的向量表示存储在向量数据库中

### 3.10 对话 (Conversation)

**描述**：代表用户与智能体的对话会话。

**属性**：
```sql
CREATE TABLE conversations (
    id UUID PRIMARY KEY,
    title VARCHAR(255),
    agent_id UUID NOT NULL REFERENCES agents(id),
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(16) NOT NULL DEFAULT 'active', -- 'active', 'archived'
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一个对话属于一个智能体
- 一个对话由一个用户发起
- 一个对话包含多条消息

### 3.11 消息 (Message)

**描述**：代表对话中的单条消息。

**属性**：
```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY,
    conversation_id UUID NOT NULL REFERENCES conversations(id),
    role VARCHAR(16) NOT NULL, -- 'user', 'assistant', 'system'
    content TEXT NOT NULL,
    tokens INTEGER,
    feedback VARCHAR(16), -- 'positive', 'negative', null
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**关系**：
- 一条消息属于一个对话

### 3.12 权限 (Permission)

**描述**：定义用户、团队对资源的访问权限。

**属性**：
```sql
CREATE TABLE permissions (
    id SERIAL PRIMARY KEY,
    resource_type VARCHAR(32) NOT NULL, -- 'application', 'knowledge_base', etc.
    resource_id UUID NOT NULL,
    principal_type VARCHAR(16) NOT NULL, -- 'user', 'team', 'organization'
    principal_id UUID NOT NULL,
    permission VARCHAR(16) NOT NULL, -- 'view', 'edit', 'admin'
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## 4. 多对多关系表

### 4.1 用户-组织关系 (UserOrganization)

```sql
CREATE TABLE user_organizations (
    user_id UUID NOT NULL REFERENCES users(id),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    role VARCHAR(16) NOT NULL DEFAULT 'member', -- 'owner', 'admin', 'member'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, organization_id)
);
```

### 4.2 用户-团队关系 (UserTeam)

```sql
CREATE TABLE user_teams (
    user_id UUID NOT NULL REFERENCES users(id),
    team_id UUID NOT NULL REFERENCES teams(id),
    role VARCHAR(16) NOT NULL DEFAULT 'member', -- 'lead', 'member'
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, team_id)
);
```

### 4.3 应用-知识库关系 (ApplicationKnowledgeBase)

```sql
CREATE TABLE application_knowledge_bases (
    application_id UUID NOT NULL REFERENCES applications(id),
    knowledge_base_id UUID NOT NULL REFERENCES knowledge_bases(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (application_id, knowledge_base_id)
);
```

### 4.4 智能体-知识库关系 (AgentKnowledgeBase)

```sql
CREATE TABLE agent_knowledge_bases (
    agent_id UUID NOT NULL REFERENCES agents(id),
    knowledge_base_id UUID NOT NULL REFERENCES knowledge_bases(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (agent_id, knowledge_base_id)
);
```

## 5. 索引策略

为优化查询性能，我们定义以下索引：

```sql
-- 用户相关索引
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);

-- 组织相关索引
CREATE INDEX idx_organizations_slug ON organizations(slug);
CREATE INDEX idx_organizations_created_by ON organizations(created_by);

-- 团队相关索引
CREATE INDEX idx_teams_organization_id ON teams(organization_id);

-- 应用相关索引
CREATE INDEX idx_applications_organization_id ON applications(organization_id);
CREATE INDEX idx_applications_team_id ON applications(team_id);
CREATE INDEX idx_applications_created_by ON applications(created_by);
CREATE INDEX idx_applications_status ON applications(status);

-- 智能体相关索引
CREATE INDEX idx_agents_application_id ON agents(application_id);
CREATE INDEX idx_agents_model_config_id ON agents(model_config_id);

-- 知识库相关索引
CREATE INDEX idx_knowledge_bases_organization_id ON knowledge_bases(organization_id);

-- 文档相关索引
CREATE INDEX idx_documents_knowledge_base_id ON documents(knowledge_base_id);
CREATE INDEX idx_documents_status ON documents(status);

-- 文档块相关索引
CREATE INDEX idx_document_chunks_document_id ON document_chunks(document_id);

-- 对话相关索引
CREATE INDEX idx_conversations_agent_id ON conversations(agent_id);
CREATE INDEX idx_conversations_user_id ON conversations(user_id);

-- 消息相关索引
CREATE INDEX idx_messages_conversation_id ON messages(conversation_id);
CREATE INDEX idx_messages_created_at ON messages(created_at);

-- 权限相关索引
CREATE INDEX idx_permissions_resource ON permissions(resource_type, resource_id);
CREATE INDEX idx_permissions_principal ON permissions(principal_type, principal_id);
```

## 6. 向量数据存储

向量数据使用Milvus/Weaviate存储，主要存储文档块的嵌入向量表示。

### 6.1 Milvus集合结构

```
Collection: document_embeddings
Fields:
  - id: VARCHAR, primary key
  - document_chunk_id: VARCHAR, reference to document_chunks.id
  - vector: FLOAT_VECTOR(1536), 文档块的向量表示
  - metadata: JSON, 包含文档和文档块的元数据
```

### 6.2 索引配置

```
Index Type: HNSW
Metric Type: IP (内积)
Parameters:
  - M: 16
  - efConstruction: 200
```

## 7. 对象存储

使用MinIO存储文件数据，主要包括：

1. **文档文件**：存储在`documents/`前缀下，按组织ID和知识库ID进行分类
2. **用户头像**：存储在`avatars/`前缀下
3. **应用图标**：存储在`app-icons/`前缀下
4. **导出/导入数据**：存储在`exports/`前缀下

## 8. 缓存策略

使用Redis进行缓存，主要缓存：

1. **用户会话**：缓存用户登录会话和JWT令牌
2. **API速率限制**：记录API调用频率和限制
3. **热点数据**：缓存频繁访问的数据，如应用配置、模型配置等
4. **智能体状态**：缓存对话上下文和临时状态

## 9. 数据迁移与版本控制

系统将使用数据库迁移工具（如golang-migrate）来管理数据库架构变更：

1. 每个变更都有UP和DOWN脚本
2. 迁移版本在数据库中跟踪
3. 新环境部署时自动应用所有迁移

## 10. 数据备份策略

1. **定时备份**：每日执行PostgreSQL完整备份
2. **增量备份**：每小时执行WAL日志备份
3. **对象存储备份**：定期备份MinIO数据
4. **跨区域备份**：关键数据备份到不同地区
5. **备份验证**：定期测试备份恢复流程 