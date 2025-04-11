package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ToolCategory 表示工具类别
type ToolCategory string

const (
	// 工具类别常量
	CategoryUtility    ToolCategory = "utility"    // 实用工具
	CategoryKnowledge  ToolCategory = "knowledge"  // 知识相关
	CategoryMedia      ToolCategory = "media"      // 媒体处理
	CategoryConnector  ToolCategory = "connector"  // 外部连接器
	CategoryDeveloper  ToolCategory = "developer"  // 开发者工具
	CategoryCustom     ToolCategory = "custom"     // 自定义工具
)

// ToolRegistry 管理系统中所有可用的工具
type ToolRegistry struct {
	tools        map[string]Tool
	toolsByCategory map[ToolCategory]map[string]Tool
	mu           sync.RWMutex
	httpClient   *http.Client
}

// NewToolRegistry 创建新的工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:        make(map[string]Tool),
		toolsByCategory: make(map[ToolCategory]map[string]Tool),
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// RegisterTool 注册一个新工具
func (r *ToolRegistry) RegisterTool(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if tool.Name == "" {
		return errors.New("tool name cannot be empty")
	}

	if tool.Handler == nil {
		return errors.New("tool handler cannot be nil")
	}

	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool with name '%s' already exists", tool.Name)
	}

	// 初始化类别
	if tool.Category == "" {
		tool.Category = CategoryUtility
	}

	if _, exists := r.toolsByCategory[tool.Category]; !exists {
		r.toolsByCategory[tool.Category] = make(map[string]Tool)
	}

	r.tools[tool.Name] = tool
	r.toolsByCategory[tool.Category][tool.Name] = tool
	return nil
}

// GetTool 获取指定名称的工具
func (r *ToolRegistry) GetTool(name string) (Tool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	if !exists {
		return Tool{}, fmt.Errorf("tool '%s' not found", name)
	}

	return tool, nil
}

// ListTools 列出所有注册的工具
func (r *ToolRegistry) ListTools() []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}

	return tools
}

// ListToolsByCategory 按类别列出工具
func (r *ToolRegistry) ListToolsByCategory(category ToolCategory) []Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categoryTools, exists := r.toolsByCategory[category]
	if !exists {
		return []Tool{}
	}

	tools := make([]Tool, 0, len(categoryTools))
	for _, tool := range categoryTools {
		tools = append(tools, tool)
	}

	return tools
}

// GetCategories 获取所有可用的工具类别
func (r *ToolRegistry) GetCategories() []ToolCategory {
	r.mu.RLock()
	defer r.mu.RUnlock()

	categories := make([]ToolCategory, 0, len(r.toolsByCategory))
	for category := range r.toolsByCategory {
		categories = append(categories, category)
	}

	return categories
}

// UnregisterTool 注销一个工具
func (r *ToolRegistry) UnregisterTool(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	tool, exists := r.tools[name]
	if !exists {
		return fmt.Errorf("tool '%s' not found", name)
	}

	delete(r.tools, name)
	
	// 如果类别存在，从类别中也删除
	if category, exists := r.toolsByCategory[tool.Category]; exists {
		delete(category, name)
	}
	
	return nil
}

// AddToolsToAgent 向智能体添加多个工具
func (r *ToolRegistry) AddToolsToAgent(agent *Agent, toolNames ...string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, name := range toolNames {
		tool, exists := r.tools[name]
		if !exists {
			return fmt.Errorf("tool '%s' not found", name)
		}
		agent.AddTool(tool)
	}

	return nil
}

// Tool 表示智能体可以使用的工具
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Handler     ToolHandler            `json:"-"`
	Category    ToolCategory           `json:"category"`
	IsBuiltin   bool                   `json:"is_builtin"`
	Version     string                 `json:"version,omitempty"`
}

// RegisterWebSearchTool 注册网络搜索工具
func (r *ToolRegistry) RegisterWebSearchTool(searchFunc func(ctx context.Context, query string) ([]map[string]interface{}, error)) error {
	searchTool := Tool{
		Name:        "web_search",
		Description: "搜索互联网以获取信息",
		Category:    CategoryKnowledge,
		IsBuiltin:   true,
		Version:     "1.0",
		Parameters: map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "搜索查询词",
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			query, ok := params["query"].(string)
			if !ok {
				return nil, errors.New("query parameter must be a string")
			}
			return searchFunc(ctx, query)
		},
	}

	return r.RegisterTool(searchTool)
}

// RegisterCalculatorTool 注册计算器工具
func (r *ToolRegistry) RegisterCalculatorTool() error {
	calculatorTool := Tool{
		Name:        "calculator",
		Description: "执行数学计算",
		Category:    CategoryUtility,
		IsBuiltin:   true,
		Version:     "1.0",
		Parameters: map[string]interface{}{
			"expression": map[string]interface{}{
				"type":        "string",
				"description": "要计算的数学表达式",
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			// 在实际实现中，这里会解析和计算表达式
			// 简化版本，仅返回示例
			expr, ok := params["expression"].(string)
			if !ok {
				return nil, errors.New("expression parameter must be a string")
			}
			return map[string]interface{}{
				"expression": expr,
				"result":     "计算结果示例", // 实际实现会计算真实结果
			}, nil
		},
	}

	return r.RegisterTool(calculatorTool)
}

// RegisterWeatherTool 注册天气查询工具
func (r *ToolRegistry) RegisterWeatherTool() error {
	weatherTool := Tool{
		Name:        "weather",
		Description: "获取指定城市的天气信息",
		Category:    CategoryUtility,
		IsBuiltin:   true,
		Version:     "1.0",
		Parameters: map[string]interface{}{
			"city": map[string]interface{}{
				"type":        "string",
				"description": "城市名称",
			},
			"country": map[string]interface{}{
				"type":        "string",
				"description": "国家名称，可选",
				"required":    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			city, ok := params["city"].(string)
			if !ok || city == "" {
				return nil, errors.New("city parameter must be a non-empty string")
			}

			// 在实际实现中，这里会调用天气API
			// 这里仅提供模拟数据
			weather := map[string]interface{}{
				"city":        city,
				"temperature": 25,
				"condition":   "晴朗",
				"humidity":    60,
				"wind":        "东北风3级",
				"updated_at":  time.Now().Format(time.RFC3339),
			}

			return weather, nil
		},
	}

	return r.RegisterTool(weatherTool)
}

// RegisterTimezoneTool 注册时区转换工具
func (r *ToolRegistry) RegisterTimezoneTool() error {
	timezoneTool := Tool{
		Name:        "timezone_converter",
		Description: "转换不同时区的时间",
		Category:    CategoryUtility,
		IsBuiltin:   true,
		Version:     "1.0",
		Parameters: map[string]interface{}{
			"time": map[string]interface{}{
				"type":        "string",
				"description": "要转换的时间，格式为ISO8601或自然语言描述",
			},
			"from_timezone": map[string]interface{}{
				"type":        "string",
				"description": "源时区，如'Asia/Shanghai'或'UTC+8'",
			},
			"to_timezone": map[string]interface{}{
				"type":        "string",
				"description": "目标时区，如'America/New_York'或'UTC-5'",
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			// 参数验证
			timeStr, ok := params["time"].(string)
			if !ok || timeStr == "" {
				return nil, errors.New("time parameter must be a non-empty string")
			}
			
			fromTz, ok := params["from_timezone"].(string)
			if !ok || fromTz == "" {
				return nil, errors.New("from_timezone must be a non-empty string")
			}
			
			toTz, ok := params["to_timezone"].(string)
			if !ok || toTz == "" {
				return nil, errors.New("to_timezone must be a non-empty string")
			}
			
			// 在实际实现中，这里会解析时间和时区并进行转换
			// 这里仅提供模拟数据
			result := map[string]interface{}{
				"original_time": timeStr,
				"from_timezone": fromTz,
				"to_timezone":   toTz,
				"converted_time": time.Now().Format(time.RFC3339),
			}
			
			return result, nil
		},
	}
	
	return r.RegisterTool(timezoneTool)
}

// RegisterKnowledgeSearchTool 注册知识库搜索工具
func (r *ToolRegistry) RegisterKnowledgeSearchTool(searchFunc func(ctx context.Context, query string, filters map[string]interface{}) ([]map[string]interface{}, error)) error {
	knowledgeTool := Tool{
		Name:        "knowledge_search",
		Description: "搜索知识库以获取相关信息",
		Category:    CategoryKnowledge,
		IsBuiltin:   true,
		Version:     "1.0",
		Parameters: map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "搜索查询词",
			},
			"filters": map[string]interface{}{
				"type":        "object",
				"description": "可选的过滤条件",
				"required":    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			query, ok := params["query"].(string)
			if !ok || query == "" {
				return nil, errors.New("query parameter must be a non-empty string")
			}
			
			var filters map[string]interface{}
			if filtersParam, ok := params["filters"]; ok {
				if filtersMap, ok := filtersParam.(map[string]interface{}); ok {
					filters = filtersMap
				}
			}
			
			return searchFunc(ctx, query, filters)
		},
	}
	
	return r.RegisterTool(knowledgeTool)
}

// RegisterHttpRequestTool 注册HTTP请求工具
func (r *ToolRegistry) RegisterHttpRequestTool() error {
	httpTool := Tool{
		Name:        "http_request",
		Description: "向指定URL发送HTTP请求",
		Category:    CategoryConnector,
		IsBuiltin:   true,
		Version:     "1.0",
		Parameters: map[string]interface{}{
			"url": map[string]interface{}{
				"type":        "string",
				"description": "请求的URL",
			},
			"method": map[string]interface{}{
				"type":        "string",
				"description": "HTTP方法，如GET、POST等",
				"enum":        []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"},
				"default":     "GET",
			},
			"headers": map[string]interface{}{
				"type":        "object",
				"description": "请求头",
				"required":    false,
			},
			"body": map[string]interface{}{
				"type":        "string",
				"description": "请求体，用于POST、PUT等方法",
				"required":    false,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			url, ok := params["url"].(string)
			if !ok || url == "" {
				return nil, errors.New("url parameter must be a non-empty string")
			}
			
			method := "GET"
			if methodParam, ok := params["method"].(string); ok && methodParam != "" {
				method = methodParam
			}
			
			// 创建请求
			req, err := http.NewRequestWithContext(ctx, method, url, nil)
			if err != nil {
				return nil, err
			}
			
			// 添加请求头
			if headersParam, ok := params["headers"].(map[string]interface{}); ok {
				for key, value := range headersParam {
					if strValue, ok := value.(string); ok {
						req.Header.Add(key, strValue)
					}
				}
			}
			
			// 设置请求体
			if bodyParam, ok := params["body"].(string); ok && method != "GET" && method != "HEAD" {
				req.Body = http.NoBody
				// 实际实现应该设置请求体
			}
			
			// 执行请求
			resp, err := r.httpClient.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()
			
			// 读取响应
			// 实际实现应该读取响应体并返回
			return map[string]interface{}{
				"status_code": resp.StatusCode,
				"headers":     resp.Header,
				"body":        "响应内容示例", // 实际实现应该读取真实内容
			}, nil
		},
	}
	
	return r.RegisterTool(httpTool)
}

// RegisterFileReadTool 注册文件读取工具
func (r *ToolRegistry) RegisterFileReadTool(basePath string) error {
	fileTool := Tool{
		Name:        "file_read",
		Description: "读取指定路径的文件内容",
		Category:    CategoryDeveloper,
		IsBuiltin:   true,
		Version:     "1.0",
		Parameters: map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "文件路径，相对于允许的基础路径",
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			path, ok := params["path"].(string)
			if !ok || path == "" {
				return nil, errors.New("path parameter must be a non-empty string")
			}
			
			// 安全检查：确保路径在允许的基础路径内
			fullPath := filepath.Join(basePath, path)
			if !isPathSafe(fullPath, basePath) {
				return nil, errors.New("access denied: path is outside the allowed directory")
			}
			
			// 检查文件是否存在
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				return nil, errors.New("file not found")
			}
			
			// 读取文件内容
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return nil, err
			}
			
			return map[string]interface{}{
				"path":    path,
				"content": string(content),
			}, nil
		},
	}
	
	return r.RegisterTool(fileTool)
}

// isPathSafe 检查给定路径是否在允许的基础路径内
func isPathSafe(path, basePath string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}
	absBasePath, err := filepath.Abs(basePath)
	if err != nil {
		return false
	}
	return filepath.HasPrefix(absPath, absBasePath)
}

// RegisterCustomTool 注册自定义工具
func (r *ToolRegistry) RegisterCustomTool(name, description string, parameters map[string]interface{}, handler ToolHandler, category ToolCategory) error {
	if category == "" {
		category = CategoryCustom
	}
	
	customTool := Tool{
		Name:        name,
		Description: description,
		Parameters:  parameters,
		Handler:     handler,
		Category:    category,
		IsBuiltin:   false,
		Version:     "1.0",
	}
	
	return r.RegisterTool(customTool)
}

// RegisterBatchTools 批量注册多个内置工具
func (r *ToolRegistry) RegisterBatchTools(toolNames []string) error {
	for _, name := range toolNames {
		var err error
		switch name {
		case "calculator":
			err = r.RegisterCalculatorTool()
		case "weather":
			err = r.RegisterWeatherTool()
		case "timezone_converter":
			err = r.RegisterTimezoneTool()
		case "http_request":
			err = r.RegisterHttpRequestTool()
		default:
			zap.L().Warn("Unknown built-in tool", zap.String("name", name))
			continue
		}
		
		if err != nil {
			return fmt.Errorf("failed to register tool %s: %w", name, err)
		}
	}
	
	return nil
}

// RegisterAllBuiltinTools 注册所有内置工具
func (r *ToolRegistry) RegisterAllBuiltinTools() error {
	// 注册基本工具
	if err := r.RegisterCalculatorTool(); err != nil {
		return err
	}
	if err := r.RegisterWeatherTool(); err != nil {
		return err
	}
	if err := r.RegisterTimezoneTool(); err != nil {
		return err
	}
	if err := r.RegisterHttpRequestTool(); err != nil {
		return err
	}
	
	// 添加更多工具注册...
	
	return nil
}

// ExportToolSpecifications 导出工具规范为JSON格式
func (r *ToolRegistry) ExportToolSpecifications() ([]byte, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	specs := make([]map[string]interface{}, 0, len(r.tools))
	
	for _, tool := range r.tools {
		spec := map[string]interface{}{
			"name":        tool.Name,
			"description": tool.Description,
			"parameters":  tool.Parameters,
			"category":    tool.Category,
			"is_builtin":  tool.IsBuiltin,
		}
		if tool.Version != "" {
			spec["version"] = tool.Version
		}
		specs = append(specs, spec)
	}
	
	return json.Marshal(specs)
}

// DefaultToolRegistry 默认的工具注册表实例
var DefaultToolRegistry = NewToolRegistry() 