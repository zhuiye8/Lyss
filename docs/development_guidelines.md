# 智能体构建平台开发流程与代码规范

## 1. 开发流程

### 1.1 Git工作流

我们采用GitFlow工作流模型，具体如下：

- **master**: 主分支，只存放稳定版本，受保护不可直接提交
- **develop**: 开发分支，开发的主要分支，也受保护
- **feature/xxx**: 功能分支，从develop分支创建，完成后合并回develop
- **bugfix/xxx**: 修复分支，从develop分支创建，修复后合并回develop
- **release/x.y.z**: 发布分支，从develop分支创建，测试无误后合并到master和develop
- **hotfix/xxx**: 热修复分支，从master分支创建，修复生产环境紧急问题，修复后合并到master和develop

#### 分支命名规范

- 功能分支: `feature/模块-功能`，例如: `feature/user-login`
- 修复分支: `bugfix/模块-问题简述`，例如: `bugfix/agent-memory-leak`
- 发布分支: `release/x.y.z`，例如: `release/1.0.0`
- 热修复分支: `hotfix/问题简述`，例如: `hotfix/api-crash`

### 1.2 版本控制

采用语义化版本（Semantic Versioning）命名规则：`主版本号.次版本号.修订号`

- **主版本号**: 不兼容的API修改
- **次版本号**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

### 1.3 提交规范

采用Angular提交规范，格式如下：

```
<type>(<scope>): <subject>

<body>

<footer>
```

其中:

- **type**: 提交类型
  - feat: 新功能
  - fix: 修复bug
  - docs: 文档变更
  - style: 代码风格变更(不影响代码运行的变动)
  - refactor: 重构(既不修复bug也不添加特性)
  - perf: 性能优化
  - test: 增加测试
  - chore: 构建过程或辅助工具的变动
  - ci: CI相关变更

- **scope**: 变更范围，例如 user, agent, model 等

- **subject**: 变更的简短描述，不超过50个字符

- **body**: 详细描述，说明代码变更的动机，以及和以前行为的对比

- **footer**: 关闭Issue或Breaking Changes说明

### 1.4 开发迭代流程

1. **需求分析**：分析并明确需求，创建任务或Issue
2. **设计**：必要时进行技术设计，复杂功能需要编写设计文档
3. **开发**：基于设计进行编码实现
4. **自测**：开发人员进行单元测试和集成测试
5. **代码审查**：提交Pull Request，至少一名其他开发人员审查
6. **修改**：根据审查意见进行修改
7. **合并**：代码通过审查后合并到develop分支
8. **持续集成**：自动化测试确保代码质量
9. **发布**：代码经测试无误后发布

### 1.5 项目管理

使用GitHub Projects或Linear进行项目管理：

- **Todo**: 待开始的任务
- **In Progress**: 进行中的任务
- **Review**: 待审查的任务
- **Done**: 已完成的任务

任务卡片包含：
- 任务描述
- 优先级
- 执行人
- 截止日期
- 相关链接(如PR, 设计文档等)

## 2. 代码规范

### 2.1 Go语言规范

#### 目录结构

```
server/
├── api/                  # API定义和接口
│   ├── controllers/      # 控制器
│   ├── middlewares/      # 中间件
│   ├── routes/           # 路由定义
│   └── validators/       # 请求验证
├── core/                 # 核心业务逻辑
│   ├── agent/            # 智能体相关
│   ├── conversation/     # 对话管理
│   ├── knowledge/        # 知识库处理
│   └── model/            # 模型管理
├── models/               # 数据模型
│   ├── entities/         # 实体定义
│   ├── repositories/     # 数据仓库
│   └── migrations/       # 数据库迁移
├── pkg/                  # 公共包
│   ├── auth/             # 认证相关
│   ├── cache/            # 缓存
│   ├── config/           # 配置管理
│   ├── database/         # 数据库连接
│   ├── logger/           # 日志工具
│   ├── storage/          # 存储工具
│   └── utils/            # 工具函数
├── config/               # 配置文件
├── scripts/              # 脚本工具
├── main.go               # 入口文件
└── go.mod                # 依赖管理
```

#### 命名规范

- **文件名**: 小写，下划线分隔，如 `user_service.go`
- **包名**: 小写，单个单词，如 `models`, `controllers`
- **变量名**: 驼峰式，如 `userID`, `errorMessage`
- **常量名**: 大写，下划线分隔，如 `MAX_RETRY_COUNT`
- **接口名**: 驼峰式，通常以 "er" 结尾，如 `Reader`, `Writer`
- **结构体名**: 大驼峰式，如 `UserService`, `ApiController`

#### 代码格式

- 使用 `gofmt` 或 `goimports` 格式化代码
- 缩进使用制表符(Tab)
- 行长度控制在120个字符以内
- 函数长度控制在50行以内，过长应考虑拆分

#### 注释规范

- 所有导出的函数、类型、变量和常量都应该有注释
- 包应该有包级注释，位于package子句之前
- 注释应该是完整的句子，句首字母大写，句末有标点
- 遵循 GoDoc 规范，使用正确的注释格式

```go
// User represents a system user.
type User struct {
  ID       string
  Username string
  Email    string
}

// FindByID retrieves a user by its ID.
// Returns nil if the user is not found.
func FindByID(id string) (*User, error) {
  // ...
}
```

#### 错误处理

- 总是检查错误，不要使用 `_` 忽略错误
- 错误应该包含足够的上下文信息
- 使用自定义错误类型或增加错误上下文，如使用 `github.com/pkg/errors`
- 函数应将错误作为最后一个返回值

```go
func ProcessFile(path string) ([]byte, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, errors.Wrap(err, "failed to open file")
  }
  defer file.Close()
  
  // ...
}
```

#### 并发控制

- 使用 context 进行超时控制和取消
- 谨慎使用 goroutine，确保正确管理生命周期
- 使用适当的同步原语，如 Mutex, WaitGroup 等
- 避免 goroutine 泄漏，确保所有 goroutine 都能正常退出

### 2.2 JavaScript/TypeScript规范

#### 目录结构

```
web/
├── public/              # 静态资源
├── src/                 # 源代码
│   ├── components/      # 共享组件
│   │   ├── layout/      # 布局组件
│   │   ├── common/      # 通用组件
│   │   └── ...          # 其他组件
│   ├── pages/           # 页面组件
│   ├── services/        # API服务
│   ├── store/           # 状态管理
│   ├── utils/           # 工具函数
│   ├── hooks/           # 自定义Hooks
│   ├── styles/          # 全局样式
│   ├── types/           # TypeScript类型定义
│   ├── App.tsx          # 根组件
│   └── index.tsx        # 入口文件
├── tests/               # 测试文件
├── .eslintrc.js         # ESLint配置
├── .prettierrc          # Prettier配置
├── package.json         # 项目配置
└── tsconfig.json        # TypeScript配置
```

#### 命名规范

- **文件名**: 
  - 组件: 大驼峰式, 如 `Button.tsx`, `UserForm.tsx`
  - 非组件: 小驼峰式, 如 `apiService.ts`, `useAuth.ts`
- **变量名**: 小驼峰式, 如 `userData`, `isLoading`
- **常量名**: 全大写下划线分隔, 如 `MAX_RETRY_COUNT`
- **组件名**: 大驼峰式, 如 `Button`, `UserProfile`
- **接口名**: 大驼峰式，前缀 I, 如 `IUser`, `IApiResponse`
- **类型名**: 大驼峰式，前缀 T, 如 `TUserData`, `TConfig`

#### 代码格式

- 使用 2 空格缩进
- 行长度控制在 100 个字符以内
- 文件末尾保留一个空行
- 使用分号结束语句
- 使用单引号 `'` 而非双引号 `"`
- 对象、数组最后一个元素后要加逗号

#### React组件编写

- 使用函数组件和Hooks，避免类组件
- 组件文件应只包含一个主要组件，并默认导出
- 小型工具组件可以放在同一个文件中，并命名导出
- 使用React.FC类型标注组件
- 组件属性使用接口定义类型

```tsx
import React from 'react';

interface ButtonProps {
  text: string;
  onClick: () => void;
  disabled?: boolean;
}

const Button: React.FC<ButtonProps> = ({ text, onClick, disabled = false }) => {
  return (
    <button 
      className="button" 
      onClick={onClick} 
      disabled={disabled}
    >
      {text}
    </button>
  );
};

export default Button;
```

#### 状态管理

- 使用Zustand进行全局状态管理
- 组件内部状态使用useState或useReducer
- 复杂表单状态可考虑使用Formik等库
- 避免过度使用全局状态，优先考虑组件内部状态

#### API调用

- 统一使用Axios进行API调用
- API请求封装在services目录中
- 使用拦截器统一处理错误和认证
- 使用React Query管理数据获取、缓存和状态

#### 样式管理

- 使用Less作为CSS预处理器
- 组件样式使用CSS Modules
- 全局样式放在styles目录
- 使用主题变量和设计标记

### 2.3 数据库规范

#### 表命名

- 表名使用小写，下划线分隔，复数形式，如 `users`, `knowledge_bases`
- 关联表名使用单数形式，下划线连接，如 `user_organization`
- 表名应清晰表达表的内容，避免缩写

#### 字段命名

- 字段名使用小写，下划线分隔，如 `first_name`, `email_address`
- 主键统一命名为 `id`
- 外键命名为 `{关联表单数形式}_id`，如 `user_id`, `organization_id`
- 创建和更新时间字段分别命名为 `created_at` 和 `updated_at`
- 状态字段命名为 `status`，类型字段命名为 `type`
- 布尔字段应用 `is_` 或 `has_` 前缀，如 `is_active`, `has_attachment`

#### 索引命名

- 主键索引：`pk_{表名}`
- 唯一索引：`uk_{表名}_{字段名}`
- 普通索引：`idx_{表名}_{字段名}`
- 外键索引：`fk_{表名}_{关联表名}`

#### SQL编写

- 关键字大写，如 `SELECT`, `INSERT`, `UPDATE`
- 表名和字段名使用小写
- 复杂查询应添加注释
- 避免使用 `SELECT *`，明确指定需要的字段
- 对于大型查询，每个子句应换行
- 使用参数化查询，避免SQL注入

### 2.4 文档规范

#### 代码注释

- 关键业务逻辑必须有注释
- 复杂算法需要详细注释和算法描述
- API接口必须有完整的注释，包括参数和返回值说明
- 使用TODO、FIXME等标记注明待改进的地方

#### API文档

- 使用OpenAPI/Swagger记录API接口
- 所有API接口必须在文档中定义
- 文档应包含：接口描述、请求方法、URL、请求参数、响应结果、错误码
- 文档应随代码一起更新

#### 技术文档

- 重要模块应有设计文档
- 文档采用Markdown格式
- 复杂流程应有流程图
- 系统架构应有架构图

## 3. 测试规范

### 3.1 单元测试

- 所有核心业务逻辑必须有单元测试
- 测试文件命名为 `{原文件名}_test.{后缀}`
- Go使用标准库testing包，TypeScript使用Jest
- 单元测试覆盖率目标不低于80%
- 测试应独立，不依赖外部服务

### 3.2 集成测试

- 关键流程必须有集成测试
- API接口必须有集成测试
- 测试应包括正常流程和异常流程
- 可使用Docker容器进行依赖服务的集成测试

### 3.3 端到端测试

- 关键用户场景应有端到端测试
- 前端使用Cypress或Playwright进行E2E测试
- 测试应覆盖主要用户流程

## 4. CI/CD流程

### 4.1 持续集成

每次代码提交后，自动执行以下流程：

1. 代码检查 (Lint)
2. 单元测试
3. 构建检查

当合并到develop分支时，额外执行：

1. 集成测试
2. 构建Docker镜像
3. 部署到开发环境

### 4.2 持续部署

Release分支自动执行：

1. 所有测试
2. 构建生产环境Docker镜像
3. 部署到预生产环境
4. 自动化验收测试

合并到master分支后：

1. 自动创建Release标签
2. 部署到生产环境

### 4.3 环境管理

- **开发环境**：供开发人员使用，自动部署最新develop分支代码
- **测试环境**：供测试人员使用，部署特定版本的代码
- **预生产环境**：与生产环境配置相同，用于最终验证
- **生产环境**：最终用户使用的环境，只部署经过充分测试的版本

## 5. 性能与安全规范

### 5.1 性能规范

- API响应时间应在300ms以内
- 页面首次加载时间控制在2秒以内
- 定期进行性能测试和优化
- 大型查询必须有分页
- 频繁访问的数据应使用缓存

### 5.2 安全规范

- 所有API调用必须进行认证和授权检查
- 敏感数据必须加密存储
- 密码必须使用强哈希算法 (如bcrypt)
- API请求必须进行输入验证
- 避免在日志中输出敏感信息
- 定期进行安全漏洞扫描

## 6. 代码审查清单

代码审查应关注以下几点：

1. **功能性**：代码是否实现了预期功能
2. **可读性**：代码是否易于理解
3. **可维护性**：代码结构是否合理，是否易于修改
4. **性能**：是否有性能问题
5. **安全性**：是否有安全隐患
6. **测试覆盖**：是否有足够的测试
7. **代码风格**：是否符合项目规范
8. **文档**：是否有必要的注释和文档

## 7. 发布流程

1. **准备阶段**
   - 确认所有功能完成并测试通过
   - 更新版本号和CHANGELOG
   - 准备发布说明

2. **预发布阶段**
   - 创建Release分支
   - 部署到预生产环境
   - 进行最终验收测试

3. **发布阶段**
   - 合并到master分支
   - 创建版本标签
   - 部署到生产环境
   - 监控系统状态

4. **发布后**
   - 公布发布说明
   - 收集用户反馈
   - 监控系统性能和异常 