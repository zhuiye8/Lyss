package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bytedance/eino"
	"github.com/bytedance/eino/provider"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Agent 表示一个智能体实例
type Agent struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Model       string                 `json:"model"`
	Provider    string                 `json:"provider"`
	Tools       []Tool                 `json:"tools"`
	Memory      Memory                 `json:"memory"`
	Config      map[string]interface{} `json:"config"`
	einoAgent   *eino.Agent
}

// Tool 表示智能体可以使用的工具
type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
	Handler     ToolHandler            `json:"-"`
}

// ToolHandler 是处理工具调用的函数类型
type ToolHandler func(ctx context.Context, params map[string]interface{}) (interface{}, error)

// Memory 智能体的记忆接口
type Memory interface {
	AddMessage(msg eino.Message) error
	GetMessages() ([]eino.Message, error)
	Clear() error
}

// NewAgent 创建一个新的智能体实例
func NewAgent(name, description, model, providerName string, config map[string]interface{}) (*Agent, error) {
	if name == "" || model == "" || providerName == "" {
		return nil, errors.New("name, model and provider are required")
	}

	agent := &Agent{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Model:       model,
		Provider:    providerName,
		Tools:       []Tool{},
		Config:      config,
	}

	return agent, nil
}

// InitEinoAgent 初始化Eino Agent
func (a *Agent) InitEinoAgent(ctx context.Context, apiKey, baseURL string) error {
	// 创建模型提供者
	var p provider.Provider

	switch a.Provider {
	case "openai":
		p = provider.NewOpenAI(apiKey, baseURL)
	case "anthropic":
		p = provider.NewAnthropic(apiKey, baseURL)
	default:
		return fmt.Errorf("unsupported provider: %s", a.Provider)
	}

	// 转换工具为Eino工具格式
	tools := make([]eino.Tool, 0, len(a.Tools))
	for _, t := range a.Tools {
		tools = append(tools, eino.Tool{
			Name:        t.Name,
			Description: t.Description,
			Parameters:  t.Parameters,
		})
	}

	// 创建Eino Agent配置
	agentConfig := eino.AgentConfig{
		Model:    a.Model,
		Provider: p,
		Tools:    tools,
	}

	// 应用自定义配置
	if a.Config != nil {
		if temperature, ok := a.Config["temperature"].(float64); ok {
			agentConfig.Temperature = temperature
		}
		if maxTokens, ok := a.Config["max_tokens"].(int); ok {
			agentConfig.MaxTokens = maxTokens
		}
	}

	// 创建Agent
	einoAgent, err := eino.NewAgent(&agentConfig)
	if err != nil {
		return err
	}

	a.einoAgent = einoAgent
	return nil
}

// AddTool 为智能体添加工具
func (a *Agent) AddTool(tool Tool) {
	a.Tools = append(a.Tools, tool)
}

// Chat 与智能体进行对话
func (a *Agent) Chat(ctx context.Context, userMessage string) (string, error) {
	if a.einoAgent == nil {
		return "", errors.New("agent not initialized, call InitEinoAgent first")
	}

	// 创建用户消息
	msg := eino.Message{
		Role:    eino.RoleUser,
		Content: userMessage,
	}

	// 发送消息并获取响应
	resp, err := a.einoAgent.Chat(ctx, msg)
	if err != nil {
		zap.L().Error("Failed to chat with agent", zap.Error(err))
		return "", err
	}

	// 如果有工具调用，处理它们
	if len(resp.ToolCalls) > 0 {
		return a.handleToolCalls(ctx, resp)
	}

	return resp.Content, nil
}

// 处理工具调用
func (a *Agent) handleToolCalls(ctx context.Context, resp eino.Response) (string, error) {
	var results []eino.ToolResult

	// 处理每个工具调用
	for _, call := range resp.ToolCalls {
		var params map[string]interface{}
		if err := json.Unmarshal([]byte(call.Arguments), &params); err != nil {
			zap.L().Error("Failed to unmarshal tool arguments", zap.Error(err))
			results = append(results, eino.ToolResult{
				ToolCallID: call.ID,
				Error:      fmt.Sprintf("Invalid arguments: %v", err),
			})
			continue
		}

		// 查找匹配的工具
		var tool *Tool
		for i := range a.Tools {
			if a.Tools[i].Name == call.Name {
				tool = &a.Tools[i]
				break
			}
		}

		if tool == nil {
			results = append(results, eino.ToolResult{
				ToolCallID: call.ID,
				Error:      fmt.Sprintf("Tool not found: %s", call.Name),
			})
			continue
		}

		// 执行工具
		if tool.Handler != nil {
			result, err := tool.Handler(ctx, params)
			if err != nil {
				results = append(results, eino.ToolResult{
					ToolCallID: call.ID,
					Error:      err.Error(),
				})
			} else {
				resultJSON, _ := json.Marshal(result)
				results = append(results, eino.ToolResult{
					ToolCallID: call.ID,
					Result:     string(resultJSON),
				})
			}
		} else {
			results = append(results, eino.ToolResult{
				ToolCallID: call.ID,
				Error:      "Tool handler not implemented",
			})
		}
	}

	// 将结果发送回智能体
	finalResp, err := a.einoAgent.SubmitToolResults(ctx, results)
	if err != nil {
		return "", err
	}

	// 如果仍有工具调用，递归处理
	if len(finalResp.ToolCalls) > 0 {
		return a.handleToolCalls(ctx, finalResp)
	}

	return finalResp.Content, nil
}

// ToJSON 将智能体转换为JSON
func (a *Agent) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

// FromJSON 从JSON创建智能体
func FromJSON(data []byte) (*Agent, error) {
	var agent Agent
	if err := json.Unmarshal(data, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
} 