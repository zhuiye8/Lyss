# 智能体构建平台

一个基于Eino框架和React.js的智能体构建平台，旨在提供比Dify更精简高效的AI应用开发体验。

## 项目概述

本项目是一个完整的智能体构建平台，支持用户快速构建、部署和管理AI应用。主要功能包括：

- 多模型支持（本地和云端）
- 可视化智能体编排
- 知识库管理
- 应用部署与监控
- 数据分析与优化
- API集成与工具调用

## 技术栈

### 后端
- **框架**: Eino (Go)
- **数据库**: PostgreSQL, Redis
- **向量数据库**: Milvus/Weaviate
- **消息队列**: Kafka
- **存储**: MinIO

### 前端
- **框架**: React.js
- **UI库**: Ant Design
- **AI组件**: ant-design/x
- **状态管理**: Redux/Zustand
- **路由**: React Router

### 基础设施
- **容器化**: Docker
- **编排**: Kubernetes
- **CI/CD**: GitHub Actions
- **监控**: Prometheus & Grafana

## 安装指南

### 环境要求
- Go 1.20+
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL 14+
- Redis 6+

### 开发环境设置

1. 克隆仓库
```bash
git clone https://github.com/yourusername/agent-platform.git
cd agent-platform
```

2. 后端设置
```bash
cd server
go mod tidy
go run main.go
```

3. 前端设置
```bash
cd web
npm install
npm run dev
```

4. 使用Docker Compose启动相关服务
```bash
docker-compose up -d
```

## 项目结构

```
.
├── docs/              # 文档
├── deploy/            # 部署配置
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

## 贡献指南

请查看 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解如何贡献代码。

## 许可证

[MIT License](./LICENSE) 