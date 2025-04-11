package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yourusername/agent-platform/server/models"
)

// OpenAIAdapter 适配OpenAI接口
type OpenAIAdapter struct {
	apiKey     string
	baseURL    string
	orgID      string
	httpClient *http.Client
}

// NewOpenAIAdapter 创建OpenAI适配器
func NewOpenAIAdapter(config models.ModelProviderConfig) *OpenAIAdapter {
	baseURL := "https://api.openai.com/v1"
	if config.BaseURL != "" {
		baseURL = config.BaseURL
	}
	
	return &OpenAIAdapter{
		apiKey:  config.ApiKey,
		baseURL: baseURL,
		orgID:   config.OrgID,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// OpenAIChatRequest OpenAI聊天请求结构
type OpenAIChatRequest struct {
	Model       string             `json:"model"`
	Messages    []OpenAIChatMessage `json:"messages"`
	Functions   []OpenAIFunction    `json:"functions,omitempty"`
	Temperature float32            `json:"temperature,omitempty"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
}

// OpenAIChatMessage OpenAI聊天消息结构
type OpenAIChatMessage struct {
	Role         string              `json:"role"`
	Content      string              `json:"content"`
	Name         string              `json:"name,omitempty"`
	FunctionCall *OpenAIFunctionCall `json:"function_call,omitempty"`
}

// OpenAIFunction OpenAI函数定义
type OpenAIFunction struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Parameters  json.RawMessage `json:"parameters"`
}

// OpenAIFunctionCall OpenAI函数调用
type OpenAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// OpenAIChatResponse OpenAI聊天响应结构
type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int                `json:"index"`
		Message      OpenAIChatMessage  `json:"message"`
		FinishReason string             `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// OpenAIEmbeddingRequest OpenAI嵌入请求结构
type OpenAIEmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// OpenAIEmbeddingResponse OpenAI嵌入响应结构
type OpenAIEmbeddingResponse struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// Chat 实现对话方法
func (a *OpenAIAdapter) Chat(ctx context.Context, request ChatRequest) (*ChatResponse, error) {
	if a.apiKey == "" {
		return nil, ErrAPIKeyRequired
	}
	
	// 转换消息格式
	messages := make([]OpenAIChatMessage, len(request.Messages))
	for i, msg := range request.Messages {
		messages[i] = OpenAIChatMessage{
			Role:    msg.Role,
			Content: msg.Content,
			Name:    msg.Name,
		}
		
		if msg.FuncCall != nil {
			messages[i].FunctionCall = &OpenAIFunctionCall{
				Name:      msg.FuncCall.Name,
				Arguments: msg.FuncCall.Arguments,
			}
		}
	}
	
	// 转换函数定义
	functions := make([]OpenAIFunction, len(request.Functions))
	for i, fn := range request.Functions {
		functions[i] = OpenAIFunction{
			Name:        fn.Name,
			Description: fn.Description,
			Parameters:  json.RawMessage(fn.Parameters),
		}
	}
	
	// 构建请求
	openaiReq := OpenAIChatRequest{
		Model:       "gpt-3.5-turbo", // 默认模型，实际应该使用配置中的模型
		Messages:    messages,
		Temperature: request.Temperature,
		MaxTokens:   request.MaxTokens,
		Stream:      request.Stream,
	}
	
	if len(functions) > 0 {
		openaiReq.Functions = functions
	}
	
	// 序列化请求
	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, err
	}
	
	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	if a.orgID != "" {
		req.Header.Set("OpenAI-Organization", a.orgID)
	}
	
	// 记录开始时间
	startTime := time.Now()
	
	// 发送请求
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 计算延迟
	latency := time.Since(startTime)
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrAPIError, string(bodyBytes))
	}
	
	// 解析响应
	var openaiResp OpenAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return nil, err
	}
	
	// 检查是否返回有效结果
	if len(openaiResp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI返回了空响应")
	}
	
	// 构建我们的响应格式
	response := &ChatResponse{
		ID:               openaiResp.ID,
		PromptTokens:     openaiResp.Usage.PromptTokens,
		CompletionTokens: openaiResp.Usage.CompletionTokens,
		TotalTokens:      openaiResp.Usage.TotalTokens,
		Model:            openaiResp.Model,
		FinishReason:     openaiResp.Choices[0].FinishReason,
		Latency:          latency,
		// 费用将由调用者根据模型定价计算
	}
	
	// 处理消息
	choice := openaiResp.Choices[0]
	response.Message = Message{
		Role:    choice.Message.Role,
		Content: choice.Message.Content,
		Name:    choice.Message.Name,
	}
	
	if choice.Message.FunctionCall != nil {
		response.Message.FuncCall = &struct {
			Name      string `json:"name"`
			Arguments string `json:"arguments"`
		}{
			Name:      choice.Message.FunctionCall.Name,
			Arguments: choice.Message.FunctionCall.Arguments,
		}
	}
	
	return response, nil
}

// Embedding 实现嵌入方法
func (a *OpenAIAdapter) Embedding(ctx context.Context, request EmbeddingRequest) (*EmbeddingResponse, error) {
	if a.apiKey == "" {
		return nil, ErrAPIKeyRequired
	}
	
	// 默认使用embedding模型
	model := "text-embedding-3-small"
	if request.Model != "" {
		model = request.Model
	}
	
	// 构建请求
	openaiReq := OpenAIEmbeddingRequest{
		Model: model,
		Input: request.Texts,
	}
	
	// 序列化请求
	jsonData, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, err
	}
	
	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	if a.orgID != "" {
		req.Header.Set("OpenAI-Organization", a.orgID)
	}
	
	// 记录开始时间
	startTime := time.Now()
	
	// 发送请求
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	// 计算延迟
	latency := time.Since(startTime)
	
	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: %s", ErrAPIError, string(bodyBytes))
	}
	
	// 解析响应
	var openaiResp OpenAIEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&openaiResp); err != nil {
		return nil, err
	}
	
	// 构建我们的响应格式
	embeddings := make([]EmbeddingVector, len(openaiResp.Data))
	for i, data := range openaiResp.Data {
		embeddings[i] = EmbeddingVector{
			Vector: data.Embedding,
			Index:  data.Index,
			Object: data.Object,
		}
	}
	
	response := &EmbeddingResponse{
		Embeddings: embeddings,
		Model:      openaiResp.Model,
		TokenCount: openaiResp.Usage.TotalTokens,
		Latency:    latency,
		// 费用将由调用者根据模型定价计算
	}
	
	return response, nil
}

// TestConnection 测试连接
func (a *OpenAIAdapter) TestConnection(ctx context.Context) error {
	if a.apiKey == "" {
		return ErrAPIKeyRequired
	}
	
	// 使用模型列表API测试连接
	req, err := http.NewRequestWithContext(ctx, "GET", a.baseURL+"/models", nil)
	if err != nil {
		return err
	}
	
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	if a.orgID != "" {
		req.Header.Set("OpenAI-Organization", a.orgID)
	}
	
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: %s", ErrAPIError, string(bodyBytes))
	}
	
	return nil
} 