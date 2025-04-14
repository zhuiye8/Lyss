# 基于阿里云API规范的智能体构建平台API重构指南

## 1. 阿里云API规范概述

阿里云API规范采用RESTful架构风格，具有以下特点：
- 资源导向的URL设计
- 标准HTTP方法语义
- 统一的请求和响应格式
- 规范的错误处理机制
- 完善的安全认证机制

## 2. 标准响应结构

### 2.1 成功响应格式

```json
{
  "success": true,
  "code": "200",
  "message": "操作成功",
  "requestId": "7d66ae12-a1b2-c3d4-e5f6-g7h8i9j0k1l2",
  "data": {
    // 业务数据
  },
  "timestamp": 1682956800000
}
```

### 2.2 分页响应格式

```json
{
  "success": true,
  "code": "200",
  "message": "操作成功",
  "requestId": "7d66ae12-a1b2-c3d4-e5f6-g7h8i9j0k1l2",
  "data": {
    "list": [], // 数据列表
    "pagination": {
      "current": 1,
      "pageSize": 10,
      "total": 100
    }
  },
  "timestamp": 1682956800000
}
```

### 2.3 错误响应格式

```json
{
  "success": false,
  "code": "400001",
  "message": "请求参数错误",
  "requestId": "7d66ae12-a1b2-c3d4-e5f6-g7h8i9j0k1l2",
  "data": {
    "field": "name",
    "value": "",
    "reason": "名称不能为空"
  },
  "timestamp": 1682956800000
}
```

## 3. 状态码规范

### 3.1 HTTP状态码使用

- 200: 成功
- 400: 请求错误
- 401: 未认证
- 403: 权限不足
- 404: 资源不存在
- 429: 请求过多
- 500: 服务器错误

### 3.2 业务状态码命名规则

业务状态码采用字符串形式，便于扩展和阅读：

```
{领域}{模块}{具体错误编号}
```

例如：
- `200`: 通用成功
- `400001`: 通用参数错误
- `401001`: 认证令牌无效
- `403001`: 权限不足
- `500001`: 系统内部错误
- `600101`: 智能体模块特定错误

## 4. 重构待办清单

### 4.1 基础设施更新

- [ ] **统一响应工具包创建**
  - [ ] 在 `server/pkg/` 下创建 `response` 目录
  - [ ] 创建 `response.go` 文件，实现统一响应结构
  - [ ] 实现响应函数:
    - [ ] `Success()`
    - [ ] `SuccessWithPagination()`
    - [ ] `Fail()`
    - [ ] `GenRequestId()`
    - [ ] `GetTimestamp()`

- [ ] **错误码体系建立**
  - [ ] 在 `server/pkg/` 下创建 `errorcode` 目录
  - [ ] 创建 `codes.go` 文件定义错误码常量
  - [ ] 按模块分类创建错误码文件
    - [ ] `auth_errors.go`
    - [ ] `agent_errors.go`
    - [ ] `conversation_errors.go`
    - [ ] 等

- [ ] **中间件更新**
  - [ ] 更新错误处理中间件
  - [ ] 添加请求ID中间件
  - [ ] 添加响应时间统计中间件

### 4.2 认证模块重构

- [ ] **更新认证响应格式**
  - [ ] 修改 `Login` 响应
  - [ ] 修改 `Register` 响应
  - [ ] 修改 `RefreshToken` 响应
  - [ ] 修改认证错误响应

### 4.3 用户模块重构

- [ ] **修改用户API响应格式**
  - [ ] 获取用户信息接口
  - [ ] 更新用户信息接口
  - [ ] 用户列表接口

### 4.4 智能体模块重构

- [ ] **修改智能体API响应格式**
  - [ ] 获取智能体列表接口
  - [ ] 创建智能体接口
  - [ ] 更新智能体接口
  - [ ] 删除智能体接口
  - [ ] 测试智能体接口

### 4.5 对话模块重构

- [ ] **修改对话API响应格式**
  - [ ] 创建对话接口
  - [ ] 获取对话列表接口
  - [ ] 获取对话详情接口
  - [ ] 发送消息接口
  - [ ] 消息流式响应接口

### 4.6 知识库模块重构

- [ ] **修改知识库API响应格式**
  - [ ] 获取知识库列表接口
  - [ ] 创建知识库接口
  - [ ] 更新知识库接口
  - [ ] 文档上传接口
  - [ ] 知识库查询接口

### 4.7 模型管理模块重构

- [ ] **修改模型管理API响应格式**
  - [ ] 获取模型配置列表接口
  - [ ] 创建模型配置接口
  - [ ] 更新模型配置接口
  - [ ] 测试模型连接接口

### 4.8 前端适配重构

- [ ] **创建响应拦截器**
  - [ ] 在 `web/src/services/api.ts` 中添加
  - [ ] 处理新旧格式兼容

- [ ] **更新服务层**
  - [ ] 修改 `authService.ts`
  - [ ] 修改 `agentService.ts`
  - [ ] 修改 `conversationService.ts`
  - [ ] 修改 `knowledgeBaseService.ts`
  - [ ] 修改 `modelService.ts`

### 4.9 WebSocket接口重构

- [ ] **更新WebSocket消息格式**
  - [ ] 客户端发送消息格式
  - [ ] 服务端响应消息格式
  - [ ] 流式响应消息格式

### 4.10 API文档更新

- [ ] **更新OpenAPI/Swagger文档**
  - [ ] 添加新的响应模式
  - [ ] 更新响应示例
  - [ ] 更新错误码描述

- [ ] **创建错误码文档**
  - [ ] 列出所有业务错误码
  - [ ] 提供错误码详细说明

## 5. 具体实施路径

### 5.1 基础更改

1. **第一阶段：创建响应工具包**
   ```
   server/pkg/response/response.go
   server/pkg/response/pagination.go
   server/pkg/errorcode/codes.go
   ```

2. **第二阶段：添加中间件**
   ```
   server/middleware/request_id.go
   server/middleware/error_handler.go
   ```

### 5.2 逐步迁移API

1. **从认证模块开始**
   ```
   server/api/auth/handler.go
   ```

2. **然后是用户模块**
   ```
   server/api/user/handler.go
   ```

3. **依次更新其他模块**

### 5.3 前端适配

1. **创建响应拦截器**
   ```
   web/src/services/api.ts
   ```

2. **逐步更新服务**

## 6. 示例代码

### 6.1 响应工具包示例

```go
// server/pkg/response/response.go
package response

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StandardResponse 标准响应结构
type StandardResponse struct {
	Success   bool        `json:"success"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"requestId"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, StandardResponse{
		Success:   true,
		Code:      "200",
		Message:   "操作成功",
		RequestId: GenRequestId(c),
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code string, message string, data ...interface{}) {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}

	c.JSON(getHTTPStatus(code), StandardResponse{
		Success:   false,
		Code:      code,
		Message:   message,
		RequestId: GenRequestId(c),
		Data:      errorData,
		Timestamp: time.Now().UnixMilli(),
	})
}

// GenRequestId 获取请求ID
func GenRequestId(c *gin.Context) string {
	// 尝试从上下文获取请求ID
	if requestId, exists := c.Get("requestId"); exists {
		return requestId.(string)
	}
	
	// 如果不存在，生成新的
	requestId := uuid.New().String()
	c.Set("requestId", requestId)
	return requestId
}

// 私有函数，根据业务码获取HTTP状态码
func getHTTPStatus(code string) int {
	switch code[:3] {
	case "400":
		return 400
	case "401":
		return 401
	case "403":
		return 403
	case "404":
		return 404
	case "429":
		return 429
	case "500":
		return 500
	default:
		return 500
	}
}
```

### 6.2 控制器示例

```go
// server/api/agent/handler.go
func (h *Handler) GetAgents(c *gin.Context) {
	// 从上下文获取用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, "401001", "未认证")
		return
	}

	// 获取分页参数
	page := 1
	pageSize := 10
	if pageStr := c.Query("page"); pageStr != "" {
		if val, err := strconv.Atoi(pageStr); err == nil && val > 0 {
			page = val
		}
	}

	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		if val, err := strconv.Atoi(pageSizeStr); err == nil && val > 0 {
			pageSize = val
		}
	}

	searchQuery := c.Query("searchQuery")

	// 调用服务获取智能体列表
	agents, total, err := h.service.GetAgents(userID.(uuid.UUID), page, pageSize, searchQuery)
	if err != nil {
		h.logger.Error("Failed to get agents", zap.Error(err))
		response.Fail(c, "600101", "获取智能体列表失败")
		return
	}

	// 使用分页响应格式
	response.SuccessWithPagination(c, agents, page, pageSize, total)
}
```

### 6.3 前端拦截器示例

```typescript
// web/src/services/api.ts
// 响应拦截器
api.interceptors.response.use(
  (response) => {
    // 新格式响应处理
    if (response.data.hasOwnProperty('success')) {
      if (response.data.success) {
        // 成功响应
        return response.data;
      } else {
        // 错误响应
        const error = new Error(response.data.message || '请求失败');
        error.code = response.data.code;
        error.requestId = response.data.requestId;
        error.data = response.data.data;
        throw error;
      }
    }
    
    // 兼容旧格式
    return response;
  },
  (error) => {
    // 处理网络错误等
    return Promise.reject(error);
  }
);
```

## 7. 单元测试范例

为了确保API响应格式的正确性，应为每个重构的端点编写单元测试。

```go
// server/api/agent/handler_test.go
func TestGetAgents(t *testing.T) {
	// 初始化测试环境...
	
	// 模拟请求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/agents?page=1&pageSize=10", nil)
	
	// 设置认证信息
	// ...
	
	// 执行请求
	router.ServeHTTP(w, req)
	
	// 验证响应
	assert.Equal(t, 200, w.Code)
	
	var response response.StandardResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// 验证响应字段
	assert.True(t, response.Success)
	assert.Equal(t, "200", response.Code)
	assert.NotEmpty(t, response.RequestId)
	assert.NotZero(t, response.Timestamp)
	
	// 验证分页数据
	data, ok := response.Data.(map[string]interface{})
	assert.True(t, ok)
	
	list, ok := data["list"].([]interface{})
	assert.True(t, ok)
	
	pagination, ok := data["pagination"].(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(1), pagination["current"])
	assert.Equal(t, float64(10), pagination["pageSize"])
	assert.GreaterOrEqual(t, pagination["total"], float64(0))
}
```

## 8. 渐进式迁移策略

为确保系统稳定性，建议采用渐进式迁移策略：

1. 首先实现基础工具包和中间件
2. 从单个API端点开始验证
3. 添加兼容层支持新旧格式
4. 逐步迁移所有API端点
5. 统一前端调用方式
6. 最后移除兼容层

## 9. 文档和监控

- [ ] 更新API文档反映新的响应格式
- [ ] 添加请求ID到日志系统便于问题追踪
- [ ] 在监控系统中添加API响应时间和错误码统计 