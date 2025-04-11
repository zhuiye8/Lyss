package llm

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
)

var (
	ErrInvalidProvider = errors.New("无效的模型提供�?)
	ErrConfigRequired  = errors.New("需要提供模型配�?)
	ErrAPIKeyRequired  = errors.New("需要提供API密钥")
	ErrAPIError        = errors.New("API调用失败")
)

// Adapter LLM适配器接�?
type Adapter interface {
	// Chat 执行对话请求
	Chat(ctx context.Context, request ChatRequest) (*ChatResponse, error)
	
	// Embedding 生成文本嵌入向量
	Embedding(ctx context.Context, request EmbeddingRequest) (*EmbeddingResponse, error)
	
	// TestConnection 测试API连接
	TestConnection(ctx context.Context) error
}

// Manager 模型管理�?
type Manager struct {
	adapters map[models.ModelProvider]Adapter
}

// NewManager 创建模型管理�?
func NewManager() *Manager {
	return &Manager{
		adapters: make(map[models.ModelProvider]Adapter),
	}
}

// RegisterAdapter 注册适配�?
func (m *Manager) RegisterAdapter(provider models.ModelProvider, adapter Adapter) {
	m.adapters[provider] = adapter
}

// GetAdapter 获取适配�?
func (m *Manager) GetAdapter(config models.ModelConfig) (Adapter, error) {
	provider := config.Model.Provider
	
	adapter, exists := m.adapters[provider]
	if !exists {
		return nil, ErrInvalidProvider
	}
	
	return adapter, nil
}

// Message 表示对话中的一条消�?
type Message struct {
	Role     string `json:"role"`      // system, user, assistant
	Content  string `json:"content"`   // 消息内容
	Name     string `json:"name,omitempty"` // 可选名�?
	FuncCall *struct {
		Name      string `json:"name"`      // 函数�?
		Arguments string `json:"arguments"` // 函数参数 (JSON格式)
	} `json:"function_call,omitempty"`
}

// FunctionDefinition 表示可以调用的函数定�?
type FunctionDefinition struct {
	Name        string `json:"name"`        // 函数名称
	Description string `json:"description"` // 函数描述
	Parameters  string `json:"parameters"`  // 参数定义 (JSON Schema格式)
}

// ChatRequest 对话请求
type ChatRequest struct {
	ConfigID   uuid.UUID            // 使用的模型配置ID
	Messages   []Message            // 对话历史
	Functions  []FunctionDefinition // 可用函数定义
	MaxTokens  int                  // 生成的最大token�?
	Temperature float32             // 温度参数
	Stream     bool                 // 是否使用流式响应
}

// ChatResponse 对话响应
type ChatResponse struct {
	ID               string   `json:"id"`               // 响应ID
	Message          Message  `json:"message"`          // 响应消息
	PromptTokens     int      `json:"prompt_tokens"`    // 提示使用的token�?
	CompletionTokens int      `json:"completion_tokens"`// 生成使用的token�?
	TotalTokens      int      `json:"total_tokens"`     // 总token�?
	Model            string   `json:"model"`            // 使用的模�?
	FinishReason     string   `json:"finish_reason"`    // 结束原因 (stop, length, function_call)
	Latency          time.Duration `json:"latency"`     // 延迟时间
	Cost             float64  `json:"cost"`             // 费用
}

// EmbeddingRequest 嵌入请求
type EmbeddingRequest struct {
	ConfigID uuid.UUID // 使用的模型配置ID
	Texts    []string  // 需要嵌入的文本列表
	Model    string    // 可选，指定模型
}

// EmbeddingVector 表示嵌入向量
type EmbeddingVector struct {
	Vector   []float32 `json:"vector"`   // 嵌入向量
	Index    int       `json:"index"`    // 索引
	Object   string    `json:"object"`   // 对象类型
}

// EmbeddingResponse 嵌入响应
type EmbeddingResponse struct {
	Embeddings []EmbeddingVector `json:"embeddings"` // 嵌入向量列表
	Model      string            `json:"model"`      // 使用的模�?
	TokenCount int               `json:"token_count"`// 使用的token�?
	Latency    time.Duration     `json:"latency"`    // 延迟时间
	Cost       float64           `json:"cost"`       // 费用
}

// CreateAdapter 根据提供者创建适配�?
func CreateAdapter(provider models.ModelProvider, config models.ModelProviderConfig) (Adapter, error) {
	switch provider {
	case models.ModelProviderOpenAI:
		return NewOpenAIAdapter(config), nil
	case models.ModelProviderAnthropic:
		return NewAnthropicAdapter(config), nil
	case models.ModelProviderBaidu:
		return NewBaiduAdapter(config), nil
	case models.ModelProviderAli:
		return NewAliAdapter(config), nil
	case models.ModelProviderLocal:
		return NewLocalAdapter(config), nil
	case models.ModelProviderCustom:
		return NewCustomAdapter(config), nil
	default:
		return nil, ErrInvalidProvider
	}
}

// GetChatCompletionCost 计算对话完成的费�?
func GetChatCompletionCost(model *models.Model, promptTokens, completionTokens int) float64 {
	// 计算提示和完成部分费�?
	promptCost := float64(promptTokens) * model.TokenCostPrompt
	completionCost := float64(completionTokens) * model.TokenCostCompl
	
	// 总费�?
	return promptCost + completionCost
}

// GetEmbeddingCost 计算嵌入的费�?
func GetEmbeddingCost(model *models.Model, tokenCount int) float64 {
	// 嵌入只收取输入token费用
	return float64(tokenCount) * model.TokenCostPrompt
} 
