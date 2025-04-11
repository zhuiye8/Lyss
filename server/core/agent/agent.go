package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/bytedance/eino"
	"github.com/bytedance/eino/provider"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AgentRuntime 表示智能体运行时的状态和选项
type AgentRuntime struct {
	Streaming bool                   // 是否启用流式响应
	Callbacks []AgentRuntimeCallback // 回调函数列表
}

// AgentRuntimeCallback 是运行时事件的回调函数类型
type AgentRuntimeCallback func(ctx context.Context, event AgentRuntimeEvent)

// AgentRuntimeEventType 表示运行时事件类型
type AgentRuntimeEventType string

const (
	// 事件类型常量
	EventStart       AgentRuntimeEventType = "start"       // 开始处理请求
	EventThinking    AgentRuntimeEventType = "thinking"    // 思考中（LLM生成中）
	EventToolCall    AgentRuntimeEventType = "tool_call"   // 调用工具
	EventToolResult  AgentRuntimeEventType = "tool_result" // 工具结果返回
	EventToken       AgentRuntimeEventType = "token"       // 流式响应的单个token
	EventComplete    AgentRuntimeEventType = "complete"    // 完成响应
	EventError       AgentRuntimeEventType = "error"       // 发生错误
)

// AgentRuntimeEvent 表示运行时事件
type AgentRuntimeEvent struct {
	Type      AgentRuntimeEventType `json:"type"`
	AgentID   string                `json:"agent_id"`
	Timestamp int64                 `json:"timestamp"`
	Data      interface{}           `json:"data,omitempty"`
}

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
	SystemPrompt string                 `json:"system_prompt"`
	einoAgent   *eino.Agent
	Runtime     *AgentRuntime           `json:"-"`
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
		SystemPrompt: "",
		Runtime:     &AgentRuntime{
			Streaming: false,
			Callbacks: []AgentRuntimeCallback{},
		},
	}

	// 默认创建一个简单内存
	agent.Memory = NewSimpleMemory(100)

	return agent, nil
}

// SetSystemPrompt 设置系统提示词
func (a *Agent) SetSystemPrompt(prompt string) {
	a.SystemPrompt = prompt
}

// SetMemory 设置自定义内存
func (a *Agent) SetMemory(memory Memory) {
	a.Memory = memory
}

// EnableStreaming 启用流式响应
func (a *Agent) EnableStreaming(enabled bool) {
	if a.Runtime != nil {
		a.Runtime.Streaming = enabled
	}
}

// AddCallback 添加运行时回调函数
func (a *Agent) AddCallback(callback AgentRuntimeCallback) {
	if a.Runtime != nil {
		a.Runtime.Callbacks = append(a.Runtime.Callbacks, callback)
	}
}

// emitEvent 发送运行时事件
func (a *Agent) emitEvent(ctx context.Context, eventType AgentRuntimeEventType, data interface{}) {
	if a.Runtime == nil || len(a.Runtime.Callbacks) == 0 {
		return
	}

	event := AgentRuntimeEvent{
		Type:      eventType,
		AgentID:   a.ID,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Data:      data,
	}

	for _, callback := range a.Runtime.Callbacks {
		callback(ctx, event)
	}
}

// InitEinoAgent 初始化Eino Agent
func (a *Agent) InitEinoAgent(ctx context.Context, apiKey, baseURL string) error {
	// 发送开始事件
	a.emitEvent(ctx, EventStart, map[string]interface{}{
		"model":    a.Model,
		"provider": a.Provider,
	})

	// 创建模型提供者
	var p provider.Provider

	switch a.Provider {
	case "openai":
		p = provider.NewOpenAI(apiKey, baseURL)
	case "anthropic":
		p = provider.NewAnthropic(apiKey, baseURL)
	case "baidu":
		// 百度文心一言需要额外的AppID
		appID := ""
		if a.Config != nil {
			if val, ok := a.Config["app_id"].(string); ok {
				appID = val
			}
		}
		p = NewBaiduProvider(apiKey, appID, baseURL)
	case "aliyun":
		// 阿里通义千问
		p = NewAliyunProvider(apiKey, baseURL)
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
		a.emitEvent(ctx, EventError, map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	a.einoAgent = einoAgent
	return nil
}

// AddTool 为智能体添加工具
func (a *Agent) AddTool(tool Tool) {
	a.Tools = append(a.Tools, tool)
}

// ChatStream 与智能体进行流式对话，返回一个可以读取流式响应的reader
func (a *Agent) ChatStream(ctx context.Context, userMessage string) (io.Reader, error) {
	if a.einoAgent == nil {
		return nil, errors.New("agent not initialized, call InitEinoAgent first")
	}

	// 创建用户消息并添加到记忆
	userMsg := eino.Message{
		Role:    eino.RoleUser,
		Content: userMessage,
	}

	if a.Memory != nil {
		if err := a.Memory.AddMessage(userMsg); err != nil {
			zap.L().Warn("Failed to add message to memory", zap.Error(err))
		}
	}

	var messages []eino.Message

	// 如果有系统提示词，添加它
	if a.SystemPrompt != "" {
		messages = append(messages, eino.Message{
			Role:    eino.RoleSystem,
			Content: a.SystemPrompt,
		})
	}
	
	// 从记忆中获取历史消息
	if a.Memory != nil {
		historyMsgs, err := a.Memory.GetMessages()
		if err == nil {
			messages = append(messages, historyMsgs...)
		} else {
			zap.L().Warn("Failed to get messages from memory", zap.Error(err))
			// 如果获取历史失败，至少包含当前用户消息
			messages = append(messages, userMsg)
		}
	} else {
		// 如果没有记忆，至少包含当前用户消息
		messages = append(messages, userMsg)
	}

	// 发送思考事件
	a.emitEvent(ctx, EventThinking, nil)

	// 使用stream方法获取流式响应
	stream, err := a.einoAgent.ChatMultipleStream(ctx, messages)
	if err != nil {
		a.emitEvent(ctx, EventError, map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	
	// 创建一个管道，用于将stream转换为io.Reader
	pr, pw := io.Pipe()
	
	// 在goroutine中处理stream
	go func() {
		defer pw.Close()
		
		var fullContent string
		
		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			}
			
			if err != nil {
				zap.L().Error("Error receiving stream chunk", zap.Error(err))
				a.emitEvent(ctx, EventError, map[string]interface{}{
					"error": err.Error(),
				})
				pw.CloseWithError(err)
				return
			}
			
			// 发送token事件
			a.emitEvent(ctx, EventToken, map[string]interface{}{
				"content": chunk.Content,
			})
			
			// 写入管道
			_, err = pw.Write([]byte(chunk.Content))
			if err != nil {
				zap.L().Error("Error writing to pipe", zap.Error(err))
				pw.CloseWithError(err)
				return
			}
			
			fullContent += chunk.Content
		}
		
		// 如果有工具调用，处理它们
		// 注意：这里需要在流式响应完成后检查工具调用
		// Eino的流式实现可能需要修改以支持工具调用与流式响应的组合
		
		// 将助手回复添加到记忆
		if a.Memory != nil {
			assistantMsg := eino.Message{
				Role:    eino.RoleAssistant,
				Content: fullContent,
			}
			if err := a.Memory.AddMessage(assistantMsg); err != nil {
				zap.L().Warn("Failed to add assistant message to memory", zap.Error(err))
			}
		}
		
		// 发送完成事件
		a.emitEvent(ctx, EventComplete, map[string]interface{}{
			"content": fullContent,
		})
	}()
	
	return pr, nil
}

// Chat 与智能体进行对话
func (a *Agent) Chat(ctx context.Context, userMessage string) (string, error) {
	// 如果启用了流式响应，使用不同的处理方式
	if a.Runtime != nil && a.Runtime.Streaming {
		reader, err := a.ChatStream(ctx, userMessage)
		if err != nil {
			return "", err
		}
		
		// 读取整个响应
		contentBytes, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		
		return string(contentBytes), nil
	}

	if a.einoAgent == nil {
		return "", errors.New("agent not initialized, call InitEinoAgent first")
	}

	// 创建用户消息并添加到记忆
	userMsg := eino.Message{
		Role:    eino.RoleUser,
		Content: userMessage,
	}

	if a.Memory != nil {
		if err := a.Memory.AddMessage(userMsg); err != nil {
			zap.L().Warn("Failed to add message to memory", zap.Error(err))
		}
	}

	var messages []eino.Message

	// 如果有系统提示词，添加它
	if a.SystemPrompt != "" {
		messages = append(messages, eino.Message{
			Role:    eino.RoleSystem,
			Content: a.SystemPrompt,
		})
	}
	
	// 从记忆中获取历史消息
	if a.Memory != nil {
		historyMsgs, err := a.Memory.GetMessages()
		if err == nil {
			messages = append(messages, historyMsgs...)
		} else {
			zap.L().Warn("Failed to get messages from memory", zap.Error(err))
			// 如果获取历史失败，至少包含当前用户消息
			messages = append(messages, userMsg)
		}
	} else {
		// 如果没有记忆，至少包含当前用户消息
		messages = append(messages, userMsg)
	}

	// 发送思考事件
	a.emitEvent(ctx, EventThinking, nil)

	// 发送消息并获取响应
	resp, err := a.einoAgent.ChatMultiple(ctx, messages)
	if err != nil {
		zap.L().Error("Failed to chat with agent", zap.Error(err))
		a.emitEvent(ctx, EventError, map[string]interface{}{
			"error": err.Error(),
		})
		return "", err
	}

	// 将助手回复添加到记忆
	if a.Memory != nil {
		assistantMsg := eino.Message{
			Role:    eino.RoleAssistant,
			Content: resp.Content,
		}
		if err := a.Memory.AddMessage(assistantMsg); err != nil {
			zap.L().Warn("Failed to add assistant message to memory", zap.Error(err))
		}
	}

	// 发送完成事件
	a.emitEvent(ctx, EventComplete, map[string]interface{}{
		"content": resp.Content,
	})

	// 如果有工具调用，处理它们
	if len(resp.ToolCalls) > 0 {
		return a.handleToolCalls(ctx, resp)
	}

	return resp.Content, nil
}

// 处理工具调用
func (a *Agent) handleToolCalls(ctx context.Context, resp eino.Response) (string, error) {
	var results []eino.ToolResult
	toolResultsMap := make(map[string]interface{})

	// 处理每个工具调用
	for _, call := range resp.ToolCalls {
		// 发送工具调用事件
		a.emitEvent(ctx, EventToolCall, map[string]interface{}{
			"tool_id":   call.ID,
			"tool_name": call.Name,
			"arguments": call.Arguments,
		})

		var params map[string]interface{}
		if err := json.Unmarshal([]byte(call.Arguments), &params); err != nil {
			zap.L().Error("Failed to unmarshal tool arguments", zap.Error(err))
			results = append(results, eino.ToolResult{
				ToolCallID: call.ID,
				Error:      fmt.Sprintf("Invalid arguments: %v", err),
			})
			toolResultsMap[call.Name] = fmt.Sprintf("Error: Invalid arguments - %v", err)
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
			toolResultsMap[call.Name] = fmt.Sprintf("Error: Tool not found - %s", call.Name)
			continue
		}

		// 执行工具
		if tool.Handler != nil {
			result, err := tool.Handler(ctx, params)
			if err != nil {
				errMsg := err.Error()
				results = append(results, eino.ToolResult{
					ToolCallID: call.ID,
					Error:      errMsg,
				})
				toolResultsMap[call.Name] = fmt.Sprintf("Error: %s", errMsg)
				
				// 发送工具结果事件（错误）
				a.emitEvent(ctx, EventToolResult, map[string]interface{}{
					"tool_id":   call.ID,
					"tool_name": call.Name,
					"error":     errMsg,
				})
				continue
			}
			
			// 序列化结果
			resultJSON, err := json.Marshal(result)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to serialize tool result: %v", err)
				results = append(results, eino.ToolResult{
					ToolCallID: call.ID,
					Error:      errMsg,
				})
				toolResultsMap[call.Name] = fmt.Sprintf("Error: %s", errMsg)
				
				// 发送工具结果事件（错误）
				a.emitEvent(ctx, EventToolResult, map[string]interface{}{
					"tool_id":   call.ID,
					"tool_name": call.Name,
					"error":     errMsg,
				})
				continue
			}
			
			results = append(results, eino.ToolResult{
				ToolCallID: call.ID,
				Content:    string(resultJSON),
			})
			toolResultsMap[call.Name] = result
			
			// 发送工具结果事件（成功）
			a.emitEvent(ctx, EventToolResult, map[string]interface{}{
				"tool_id":   call.ID,
				"tool_name": call.Name,
				"result":    result,
			})
		} else {
			errMsg := "Tool handler not implemented"
			results = append(results, eino.ToolResult{
				ToolCallID: call.ID,
				Error:      errMsg,
			})
			toolResultsMap[call.Name] = fmt.Sprintf("Error: %s", errMsg)
			
			// 发送工具结果事件（错误）
			a.emitEvent(ctx, EventToolResult, map[string]interface{}{
				"tool_id":   call.ID,
				"tool_name": call.Name,
				"error":     errMsg,
			})
		}
	}

	// 如果有工具结果，发送回模型以获取最终回复
	if len(results) > 0 {
		// 发送思考事件
		a.emitEvent(ctx, EventThinking, nil)

		var messages []eino.Message

		// 如果有系统提示词，添加它
		if a.SystemPrompt != "" {
			messages = append(messages, eino.Message{
				Role:    eino.RoleSystem,
				Content: a.SystemPrompt,
			})
		}

		// 获取历史消息
		if a.Memory != nil {
			historyMsgs, err := a.Memory.GetMessages()
			if err == nil {
				messages = append(messages, historyMsgs...)
			} else {
				zap.L().Warn("Failed to get messages from memory", zap.Error(err))
			}
		}

		// 添加助手消息和工具结果
		messages = append(messages, eino.Message{
			Role:    eino.RoleAssistant,
			Content: resp.Content,
			ToolCalls: resp.ToolCalls,
		})

		// 添加工具结果消息
		for _, result := range results {
			messages = append(messages, eino.Message{
				Role:        eino.RoleTool,
				ToolCallID:  result.ToolCallID,
				Content:     result.Content,
				Error:       result.Error,
			})
		}

		// 获取新的回复
		newResp, err := a.einoAgent.ChatMultiple(ctx, messages)
		if err != nil {
			zap.L().Error("Failed to chat with agent after tool call", zap.Error(err))
			a.emitEvent(ctx, EventError, map[string]interface{}{
				"error": err.Error(),
			})
			return "", err
		}

		// 将新的助手回复添加到记忆
		if a.Memory != nil {
			assistantMsg := eino.Message{
				Role:    eino.RoleAssistant,
				Content: newResp.Content,
			}
			if err := a.Memory.AddMessage(assistantMsg); err != nil {
				zap.L().Warn("Failed to add assistant message to memory", zap.Error(err))
			}
		}

		// 发送完成事件
		a.emitEvent(ctx, EventComplete, map[string]interface{}{
			"content": newResp.Content,
		})

		// 如果新的回复中又有工具调用，递归处理
		if len(newResp.ToolCalls) > 0 {
			return a.handleToolCalls(ctx, newResp)
		}

		return newResp.Content, nil
	}

	return resp.Content, nil
}

// ClearMemory 清除智能体记忆
func (a *Agent) ClearMemory() error {
	if a.Memory != nil {
		return a.Memory.Clear()
	}
	return nil
}

// ToJSON 将智能体序列化为JSON
func (a *Agent) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

// FromJSON 从JSON反序列化智能体
func FromJSON(data []byte) (*Agent, error) {
	var agent Agent
	if err := json.Unmarshal(data, &agent); err != nil {
		return nil, err
	}
	return &agent, nil
}

// ProcessMultimodalContent 处理多模态内容
func (a *Agent) ProcessMultimodalContent(ctx context.Context, userMessage string, images [][]byte) (string, error) {
	if a.einoAgent == nil {
		return "", errors.New("agent not initialized, call InitEinoAgent first")
	}

	// 检查模型是否支持多模态
	if a.Provider != "openai" && a.Provider != "anthropic" && a.Provider != "aliyun" {
		return "", errors.New("multimodal content is only supported by OpenAI, Anthropic, and Aliyun providers")
	}

	// 创建多模态消息
	var content []eino.Content
	
	// 添加文本内容
	if userMessage != "" {
		content = append(content, eino.Content{
			Type: eino.ContentTypeText,
			Text: userMessage,
		})
	}
	
	// 添加图像内容
	for _, imgData := range images {
		content = append(content, eino.Content{
			Type: eino.ContentTypeImage, 
			ImageData: &eino.ImageData{
				Data: imgData,
			},
		})
	}
	
	// 创建用户消息
	userMsg := eino.Message{
		Role:    eino.RoleUser,
		Content: "", // 将被忽略
		Contents: content,
	}
	
	// 在记忆中添加用户消息
	// 注意：可能需要修改Memory接口以支持多模态内容
	if a.Memory != nil {
		if err := a.Memory.AddMessage(userMsg); err != nil {
			zap.L().Warn("Failed to add multimodal message to memory", zap.Error(err))
		}
	}
	
	var messages []eino.Message
	
	// 如果有系统提示词，添加它
	if a.SystemPrompt != "" {
		messages = append(messages, eino.Message{
			Role:    eino.RoleSystem,
			Content: a.SystemPrompt,
		})
	}
	
	// 添加当前用户消息
	messages = append(messages, userMsg)
	
	// 发送思考事件
	a.emitEvent(ctx, EventThinking, nil)
	
	// 发送消息并获取响应
	resp, err := a.einoAgent.ChatMultiple(ctx, messages)
	if err != nil {
		zap.L().Error("Failed to chat with agent with multimodal content", zap.Error(err))
		a.emitEvent(ctx, EventError, map[string]interface{}{
			"error": err.Error(),
		})
		return "", err
	}
	
	// 将助手回复添加到记忆
	if a.Memory != nil {
		assistantMsg := eino.Message{
			Role:    eino.RoleAssistant,
			Content: resp.Content,
		}
		if err := a.Memory.AddMessage(assistantMsg); err != nil {
			zap.L().Warn("Failed to add assistant message to memory", zap.Error(err))
		}
	}
	
	// 发送完成事件
	a.emitEvent(ctx, EventComplete, map[string]interface{}{
		"content": resp.Content,
	})
	
	// 如果有工具调用，处理它们
	if len(resp.ToolCalls) > 0 {
		return a.handleToolCalls(ctx, resp)
	}
	
	return resp.Content, nil
} 