# 智能体构建平台API规范

## 1. API概述

本文档定义了智能体构建平台的RESTful API接口规范。所有API使用HTTPS协议，采用JSON格式交换数据，并通过JWT令牌进行身份认证。

## 2. API版本控制

API使用URL路径方式进行版本控制：`https://api.example.com/v1/resource`

当前API版本：`v1`

## 3. 认证与授权

### 3.1 认证方式

系统使用JWT (JSON Web Token) 进行身份认证：

```
Authorization: Bearer <token>
```

### 3.2 认证相关接口

#### 登录

```
POST /v1/auth/login
```

请求体:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

响应:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2023-12-31T23:59:59Z",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "user123",
    "email": "user@example.com"
  }
}
```

#### 注册

```
POST /v1/auth/register
```

请求体:
```json
{
  "username": "user123",
  "email": "user@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

#### 刷新令牌

```
POST /v1/auth/refresh
```

## 4. 错误处理

所有API错误返回统一格式：

```json
{
  "error": {
    "code": "invalid_request",
    "message": "描述性错误信息",
    "details": {
      "field": "错误字段",
      "reason": "具体原因"
    }
  }
}
```

HTTP状态码:
- 200: 成功
- 400: 请求错误
- 401: 未认证
- 403: 权限不足
- 404: 资源不存在
- 429: 请求过多
- 500: 服务器错误

## 5. 核心API端点

### 5.1 用户管理

```
GET    /v1/users                # 获取用户列表（管理员）
GET    /v1/users/{id}           # 获取用户信息
PUT    /v1/users/{id}           # 更新用户信息
DELETE /v1/users/{id}           # 删除用户（管理员）
GET    /v1/users/me             # 获取当前用户信息
PUT    /v1/users/me/password    # 修改密码
```

### 5.2 组织管理

```
GET    /v1/organizations                    # 获取组织列表
POST   /v1/organizations                    # 创建组织
GET    /v1/organizations/{id}               # 获取组织详情
PUT    /v1/organizations/{id}               # 更新组织
DELETE /v1/organizations/{id}               # 删除组织
GET    /v1/organizations/{id}/members       # 获取组织成员
POST   /v1/organizations/{id}/members       # 添加组织成员
DELETE /v1/organizations/{id}/members/{uid} # 移除组织成员
```

### 5.3 团队管理

```
GET    /v1/organizations/{org_id}/teams           # 获取团队列表
POST   /v1/organizations/{org_id}/teams           # 创建团队
GET    /v1/teams/{id}                             # 获取团队详情
PUT    /v1/teams/{id}                             # 更新团队
DELETE /v1/teams/{id}                             # 删除团队
GET    /v1/teams/{id}/members                     # 获取团队成员
POST   /v1/teams/{id}/members                     # 添加团队成员
DELETE /v1/teams/{id}/members/{uid}               # 移除团队成员
```

### 5.4 应用管理

```
GET    /v1/applications                    # 获取应用列表
POST   /v1/applications                    # 创建应用
GET    /v1/applications/{id}               # 获取应用详情
PUT    /v1/applications/{id}               # 更新应用
DELETE /v1/applications/{id}               # 删除应用
POST   /v1/applications/{id}/publish       # 发布应用
POST   /v1/applications/{id}/unpublish     # 取消发布
POST   /v1/applications/{id}/duplicate     # 复制应用
```

### 5.5 智能体管理

```
GET    /v1/applications/{app_id}/agents         # 获取智能体列表
POST   /v1/applications/{app_id}/agents         # 创建智能体
GET    /v1/agents/{id}                          # 获取智能体详情
PUT    /v1/agents/{id}                          # 更新智能体
DELETE /v1/agents/{id}                          # 删除智能体
PUT    /v1/agents/{id}/system-prompt            # 更新系统提示词
PUT    /v1/agents/{id}/tools                    # 更新工具配置
POST   /v1/agents/{id}/test                     # 测试智能体
```

### 5.6 对话管理

```
GET    /v1/agents/{agent_id}/conversations           # 获取对话列表
POST   /v1/agents/{agent_id}/conversations           # 创建新对话
GET    /v1/conversations/{id}                        # 获取对话详情
DELETE /v1/conversations/{id}                        # 删除对话
GET    /v1/conversations/{conv_id}/messages          # 获取消息列表
POST   /v1/conversations/{conv_id}/messages          # 发送消息
POST   /v1/conversations/{conv_id}/regenerate        # 重新生成回复
POST   /v1/messages/{id}/feedback                    # 提供消息反馈
```

### 5.7 知识库管理

```
GET    /v1/knowledge-bases                     # 获取知识库列表
POST   /v1/knowledge-bases                     # 创建知识库
GET    /v1/knowledge-bases/{id}                # 获取知识库详情
PUT    /v1/knowledge-bases/{id}                # 更新知识库
DELETE /v1/knowledge-bases/{id}                # 删除知识库
POST   /v1/knowledge-bases/{id}/documents      # 上传文档
GET    /v1/knowledge-bases/{id}/documents      # 获取文档列表
DELETE /v1/documents/{id}                      # 删除文档
POST   /v1/knowledge-bases/{id}/query          # 查询知识库
```

### 5.8 模型管理

```
GET    /v1/model-configs                    # 获取模型配置列表
POST   /v1/model-configs                    # 创建模型配置
GET    /v1/model-configs/{id}               # 获取模型配置详情
PUT    /v1/model-configs/{id}               # 更新模型配置
DELETE /v1/model-configs/{id}               # 删除模型配置
POST   /v1/model-configs/{id}/test          # 测试模型连接
GET    /v1/model-providers                  # 获取模型提供商列表
GET    /v1/model-providers/{id}/models      # 获取提供商支持的模型
```

### 5.9 工具管理

```
GET    /v1/tools                            # 获取可用工具列表
GET    /v1/tools/{id}                       # 获取工具详情
POST   /v1/organizations/{org_id}/tools     # 创建自定义工具
PUT    /v1/tools/{id}                       # 更新自定义工具
DELETE /v1/tools/{id}                       # 删除自定义工具
```

### 5.10 API访问

```
POST   /v1/api/chat                         # 公开API: 聊天
POST   /v1/api/query                        # 公开API: 知识库查询
GET    /v1/api/status                       # API状态检查
```

## 6. WebSocket接口

为支持实时聊天，系统提供WebSocket接口，端点为:

```
WSS /v1/ws/conversations/{conversation_id}
```

连接时需要在URL参数中提供令牌:
```
?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 6.1 消息类型

客户端发送:
```json
{
  "type": "message",
  "content": "用户消息内容"
}
```

服务器推送:
```json
{
  "type": "message",
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "role": "assistant",
  "content": "智能体回复内容",
  "created_at": "2023-12-31T23:59:59Z"
}
```

流式响应:
```json
{
  "type": "stream",
  "message_id": "550e8400-e29b-41d4-a716-446655440000",
  "chunk": "部分",
  "done": false
}
```

## 7. API限流策略

API实施以下限流策略:

1. 基本限制: 每IP 60次请求/分钟
2. 认证用户: 每用户 300次请求/分钟
3. 聊天API: 每用户 30次请求/分钟
4. 文档上传: 每用户 10次请求/小时

响应头包含限流信息:
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1640995200
```

## 8. API示例

### 8.1 创建智能体

请求:
```
POST /v1/applications/550e8400-e29b-41d4-a716-446655440000/agents
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

```json
{
  "name": "客服助手",
  "description": "帮助回答用户问题的智能客服",
  "model_config_id": "550e8400-e29b-41d4-a716-446655440001",
  "system_prompt": "你是一个专业的客服助手，负责回答用户的产品问题。",
  "tools": [
    {
      "name": "search_knowledge_base",
      "description": "搜索知识库获取信息",
      "parameters": {
        "query": {
          "type": "string",
          "description": "搜索查询"
        },
        "top_k": {
          "type": "number",
          "description": "返回结果数量",
          "default": 3
        }
      }
    }
  ],
  "knowledge_base_ids": [
    "550e8400-e29b-41d4-a716-446655440002"
  ]
}
```

响应:
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440003",
  "name": "客服助手",
  "description": "帮助回答用户问题的智能客服",
  "application_id": "550e8400-e29b-41d4-a716-446655440000",
  "model_config_id": "550e8400-e29b-41d4-a716-446655440001",
  "system_prompt": "你是一个专业的客服助手，负责回答用户的产品问题。",
  "tools": [...],
  "knowledge_base_ids": ["550e8400-e29b-41d4-a716-446655440002"],
  "created_at": "2023-12-31T23:59:59Z",
  "updated_at": "2023-12-31T23:59:59Z"
}
```

### 8.2 发送消息

请求:
```
POST /v1/conversations/550e8400-e29b-41d4-a716-446655440004/messages
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

```json
{
  "content": "你能告诉我如何退款吗？",
  "stream": true
}
```

流式响应:
```json
{
  "type": "stream",
  "message_id": "550e8400-e29b-41d4-a716-446655440005",
  "chunk": "要",
  "done": false
}
{
  "type": "stream",
  "message_id": "550e8400-e29b-41d4-a716-446655440005",
  "chunk": "进行退款",
  "done": false
}
...
{
  "type": "stream",
  "message_id": "550e8400-e29b-41d4-a716-446655440005",
  "chunk": "",
  "done": true
}
```

## 9. SDK支持

平台提供多种语言的SDK:

- JavaScript/TypeScript: `@agent-platform/js-sdk`
- Python: `agent-platform-python`
- Go: `github.com/agent-platform/go-sdk`

SDK用法示例:

```javascript
// JavaScript示例
const client = new AgentPlatform.Client('YOUR_API_KEY');

// 发送消息
const response = await client.conversations
  .send('550e8400-e29b-41d4-a716-446655440004', {
    content: '你能告诉我如何退款吗？'
  });

console.log(response.content);
```

## 10. API文档与开发者资源

- 交互式API文档: `/docs/api`
- OpenAPI/Swagger规范: `/api/v1/openapi.json`
- API状态: `/api/v1/status`
- SDK下载: `/docs/sdk`

## 11. 生产环境与测试环境

- 生产环境: `https://api.agent-platform.com/v1`
- 测试环境: `https://api.staging.agent-platform.com/v1`
- 本地开发: `http://localhost:8080/v1` 