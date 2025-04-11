# 智能体构建平台服务端

智能体构建平台的后端服务，基于Go语言和Eino框架开发。

## 目录结构

```
server/
├── api/        # API路由与处理程序
├── config/     # 配置文件和设置
├── core/       # 核心业务逻辑
│   ├── agent/  # 智能体引擎
│   ├── auth/   # 认证与授权
│   ├── model/  # 模型管理
│   └── kb/     # 知识库管理
├── db/         # 数据库交互层
├── middleware/ # 中间件组件
├── models/     # 数据模型定义
├── pkg/        # 公共工具包
└── server.go   # 主程序入口
```

## 主要功能模块

### 智能体引擎模块

智能体引擎是平台的核心组件，基于字节跳动的Eino框架实现。该模块提供了智能体的运行时环境、对话管理、工具调用和编排管理等功能。

#### 核心组件

1. **Agent（智能体）**
   - 基本智能体实现，集成Eino框架
   - 支持系统提示词配置
   - 实现了对话历史记录管理
   - 工具调用与结果处理

2. **Memory（内存）**
   - 对话历史记录存储
   - 上下文窗口管理
   - 支持自定义内存实现

3. **ToolRegistry（工具注册中心）**
   - 管理可用工具集合
   - 工具注册与注销
   - 提供默认工具实现

4. **ConversationManager（对话管理器）**
   - 对话会话管理
   - 对话历史记录存储
   - 消息追踪与反馈收集

5. **AgentFactory（智能体工厂）**
   - 智能体模板管理
   - 支持多种类型智能体创建
   - 智能体序列化与反序列化

### 知识库与数据处理模块

知识库模块提供了文档处理、向量化、检索和RAG实现的功能，支持智能体使用外部知识进行交互。

#### 核心组件

1. **Document（文档）**
   - 文档模型定义
   - 多种文档类型支持
   - 文档分块处理

2. **Embedding（向量化）**
   - 文本向量化服务
   - 多模型支持
   - 向量表示管理

3. **VectorDatabase（向量数据库）**
   - 通用向量数据库接口
   - Milvus集成实现
   - 高效相似度搜索

4. **KnowledgeBaseManager（知识库管理器）**
   - 知识库创建与管理
   - 文档上传与处理
   - 内容索引与检索

5. **Retriever（检索器）**
   - 基于语义的相似内容检索
   - RAG工具实现
   - 提示词增强与优化

#### 使用示例

```go
// 创建知识库
kb, err := DefaultKnowledgeBaseManager.CreateKnowledgeBase(
    ctx,
    "产品知识库",
    "包含产品说明文档",
    "mock", // 使用的向量模型
)

// 添加文档
doc, err := DefaultKnowledgeBaseManager.AddTextDocument(
    ctx,
    kb.ID,
    "使用手册.txt",
    "这是产品使用说明...",
    TypeText,
)

// 使用RAG进行回答
agent, _ := agentFactory.CreateAgent(ctx, "default_rag", "助手", "帮助回答问题", "gpt-3.5-turbo", "openai", nil)
response, err := kb.ApplyRAG(ctx, kb.ID, "产品怎么使用？", agent)
```

## 环境要求

- Go 1.20+
- PostgreSQL 14+
- Redis 6+

## 开发指南

详见 [docs/development_guidelines.md](../docs/development_guidelines.md) 