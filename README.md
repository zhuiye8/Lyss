# 智能体构建平台Lyss

<div align="center">

![Lyss Logo](https://via.placeholder.com/200x200?text=Lyss)

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?logo=go)](https://golang.org/)
[![React Version](https://img.shields.io/badge/React-18.2+-61DAFB?logo=react)](https://reactjs.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](./CONTRIBUTING.md)

</div>

一个基于Eino框架和React.js的智能体构建平台，旨在提供比Dify更精简高效的AI应用开发体验。
Lyss灵感来自"Bliss"（愉悦）和"Less"（更少），传递简洁带来的轻松高效感。

> **注意**: 项目处于积极开发阶段，部分功能尚在完善中。欢迎贡献代码或提供反馈！

## 🌟 项目概述

Lyss是一个完整的智能体构建平台，支持用户快速构建、部署和管理AI应用。主要功能包括：

- 🤖 多模型支持（本地和云端LLM接入）
- 🔄 可视化智能体编排与流程设计
- 📚 知识库管理与检索增强生成(RAG)
- 🚀 应用一键部署与实时监控
- 📊 对话数据分析与模型性能优化
- 🔌 丰富的API集成与工具调用能力

<div align="center">
    <img src="https://via.placeholder.com/800x450?text=Lyss+Platform+Screenshot" alt="Lyss Platform Screenshot" />
</div>

## 🛠️ 技术栈

### 后端
- **核心框架**: [Eino](https://github.com/cloudwego/eino) (Go) - 高性能智能体框架
- **数据存储**: PostgreSQL (关系数据), Redis (缓存)
- **向量数据库**: Milvus/Weaviate - 高效语义搜索
- **消息队列**: Kafka - 事件驱动架构
- **对象存储**: MinIO - S3兼容存储

### 前端
- **框架**: React.js + TypeScript
- **UI组件**: Ant Design - 企业级设计系统
- **AI组件**: [ant-design/x](https://x.ant.design/) - 专为AI应用设计的组件库
- **状态管理**: Zustand - 轻量高效的状态管理
- **开发工具**: Vite - 快速的前端构建工具

### 基础设施
- **容器化**: Docker + Docker Compose
- **编排**: Kubernetes - 生产环境部署
- **监控**: Prometheus & Grafana - 全方位可观测性

## ✨ 核心特性

- **直观的智能体构建**
  - 对话式应用构建器
  - 流程式应用编排器
  - 提示词管理与优化

- **强大的知识库管理**
  - 多格式文档处理
  - 自动向量化与检索
  - 上下文增强问答

- **丰富的模型支持**
  - OpenAI, Anthropic, Baidu, Aliyun等云端LLM
  - 本地模型接入与管理
  - 模型性能与成本监控

- **企业级功能**
  - 团队协作与权限管理
  - 应用版本控制
  - 完备的审计日志

## 🚀 快速开始

### 环境要求
- Go 1.20+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 14+
- Redis 6+

### 开发环境设置

1. 克隆仓库
```bash
git clone https://github.com/zhuiye8/Lyss.git
cd Lyss
```

2. 启动依赖服务(PostgreSQL, Redis, Milvus)
```bash
docker-compose up -d postgres redis milvus minio
```

3. 后端设置
```bash
cd server
go mod tidy
go run main.go
```

4. 前端设置
```bash
cd web
npm install
npm run dev
```

5. 访问开发环境
```
前端: http://localhost:5173
API: http://localhost:8080
```

## 📚 项目结构

```
.
├── docs/              # 文档与设计
│   ├── api_spec.md    # API规范
│   ├── architecture.md# 系统架构
│   └── data_models.md # 数据模型
├── deploy/            # 部署配置
│   ├── docker/        # Docker配置
│   └── k8s/           # Kubernetes配置
├── scripts/           # 脚本工具
├── server/            # 后端代码
│   ├── api/           # API接口
│   ├── core/          # 核心业务逻辑
│   ├── models/        # 数据模型
│   ├── pkg/           # 公共包
│   └── config/        # 配置文件
└── web/               # 前端代码
    ├── public/        # 静态资源
    ├── src/           # 源代码
    │   ├── components/# 公共组件
    │   ├── pages/     # 页面组件
    │   ├── services/  # API服务
    │   ├── store/     # 状态管理
    │   └── utils/     # 工具函数
    └── tests/         # 测试
```

## 🤝 参与贡献

我们非常欢迎各种形式的贡献：

- 🐛 提交问题或功能请求
- 🔍 审查代码和文档
- 📝 改进文档
- 💻 贡献代码

请查看 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解详细的贡献指南。

## 📄 许可证

[MIT License](./LICENSE) 

## 🙏 致谢

- [Eino](https://github.com/cloudwego/eino) - Go语言智能体框架
- [Ant Design](https://ant.design/) - 优秀的UI设计系统
- [ant-design/x](https://x.ant.design/) - 专业的AI组件库
- 所有贡献者和社区支持者 

## 项目进展

### 后端核心功能开发

我们已经完成了以下后端核心功能的开发：

1. **基础API实现**
   - 用户认证与授权API - 实现了基于JWT的完整认证系统
   - 智能体管理API - 完整的CRUD操作和相关功能
   - 对话历史API - 实现了对话和消息的管理功能
   - 错误处理与日志 - 在所有服务和处理程序中加入了详细的错误处理和日志记录

2. **数据模型实现**
   - 实现了智能体模型 (Agent)
   - 实现了对话模型 (Conversation)
   - 实现了消息模型 (Message)
   - 实现了智能体知识库关联模型 (AgentKnowledgeBase)

### 接下来的任务

以下是接下来需要完成的任务：

1. **Agent运行时实现**
   - 连接到LLM提供商
   - 实现推理逻辑
   - 实现对话上下文管理

2. **工具调用系统**
   - 实现工具注册机制
   - 实现工具执行环境
   - 实现工具结果解析

3. **多LLM提供商集成**
   - 实现OpenAI集成
   - 实现Anthropic集成
   - 实现国内主流模型集成

## 项目架构

### 后端架构

后端采用了分层架构设计：

- **API层** - 处理HTTP请求和响应
  - Handler - 处理请求参数验证和错误处理
  - Service - 实现业务逻辑
  
- **模型层** - 定义数据模型和数据库交互
  - Entity - 数据库模型定义
  - Repository - 数据访问逻辑
  
- **核心层** - 实现核心功能
  - Agent - 智能体核心逻辑
  - LLM - 大语言模型接口
  - Tools - 工具调用系统

## 如何运行

### 后端服务

1. 确保安装了Go 1.20或更高版本
2. 安装依赖：`go mod tidy`
3. 配置环境变量或配置文件
4. 运行服务：`go run main.go`

### 数据库设置

1. 确保PostgreSQL数据库已经运行
2. 配置数据库连接参数
3. 首次运行时，会自动创建必要的数据表

## API文档

API接口遵循RESTful设计原则，主要包括以下端点：

- 用户认证: `/api/v1/auth/*`
- 智能体管理: `/api/v1/applications/:app_id/agents` 和 `/api/v1/agents/*`
- 对话管理: `/api/v1/agents/:agent_id/conversations` 和 `/api/v1/conversations/*` 