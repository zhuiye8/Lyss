package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

// AgentType 表示智能体类型
type AgentType string

const (
	// ConversationalAgent 对话型智能体
	ConversationalAgent AgentType = "conversational"
	// RAGAgent 检索增强型智能体
	RAGAgent AgentType = "rag"
	// WorkflowAgent 工作流智能体
	WorkflowAgent AgentType = "workflow"
	// CustomAgent 自定义智能体
	CustomAgent AgentType = "custom"
)

// AgentTemplate 智能体模板
type AgentTemplate struct {
	Type         AgentType              `json:"type"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	SystemPrompt string                 `json:"system_prompt"`
	DefaultTools []string               `json:"default_tools,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
}

// AgentFactory 智能体工厂，负责创建智能体实例
type AgentFactory struct {
	templates     map[string]AgentTemplate
	toolRegistry  *ToolRegistry
	modelProvider ModelProvider
	mu            sync.RWMutex
}

// ModelProvider 模型提供商接口
type ModelProvider interface {
	GetAPIKey(provider string) (string, error)
	GetBaseURL(provider string) (string, error)
}

// NewAgentFactory 创建新的智能体工厂
func NewAgentFactory(toolRegistry *ToolRegistry, modelProvider ModelProvider) *AgentFactory {
	factory := &AgentFactory{
		templates:     make(map[string]AgentTemplate),
		toolRegistry:  toolRegistry,
		modelProvider: modelProvider,
	}
	
	// 注册默认模板
	factory.registerDefaultTemplates()
	
	return factory
}

// 注册默认模板
func (f *AgentFactory) registerDefaultTemplates() {
	// 对话型智能体模板
	f.RegisterTemplate("default_conversation", AgentTemplate{
		Type:         ConversationalAgent,
		Name:         "对话智能体",
		Description:  "通用对话型智能体，适合日常交流和信息提供",
		SystemPrompt: "你是一个有帮助的AI助手。回答用户的问题时要准确、有见地且富有帮助。",
		DefaultTools: []string{},
		Config: map[string]interface{}{
			"temperature": 0.7,
			"max_tokens":  1000,
		},
	})
	
	// RAG型智能体模板
	f.RegisterTemplate("default_rag", AgentTemplate{
		Type:         RAGAgent,
		Name:         "知识库智能体",
		Description:  "基于知识库的智能体，可以回答特定领域的问题",
		SystemPrompt: "你是一个专业的知识库助手。根据提供的知识回答用户问题，如果知识库中没有相关信息，请明确说明。",
		DefaultTools: []string{"knowledge_search"},
		Config: map[string]interface{}{
			"temperature": 0.5,
			"max_tokens":  1500,
		},
	})
	
	// 工作流智能体模板
	f.RegisterTemplate("default_workflow", AgentTemplate{
		Type:         WorkflowAgent,
		Name:         "工作流智能体",
		Description:  "可以执行多步骤任务的智能体",
		SystemPrompt: "你是一个工作流自动化助手。你可以协助用户完成多步骤任务，按照逻辑顺序使用工具完成目标。",
		DefaultTools: []string{"web_search", "calculator"},
		Config: map[string]interface{}{
			"temperature": 0.3,
			"max_tokens":  2000,
		},
	})
}

// RegisterTemplate 注册新的智能体模板
func (f *AgentFactory) RegisterTemplate(templateID string, template AgentTemplate) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if templateID == "" {
		return errors.New("template ID cannot be empty")
	}
	
	if template.Name == "" || template.Type == "" {
		return errors.New("template must have a name and type")
	}
	
	f.templates[templateID] = template
	return nil
}

// GetTemplate 获取指定ID的模板
func (f *AgentFactory) GetTemplate(templateID string) (AgentTemplate, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	template, exists := f.templates[templateID]
	if !exists {
		return AgentTemplate{}, fmt.Errorf("template '%s' not found", templateID)
	}
	
	return template, nil
}

// ListTemplates 列出所有可用模板
func (f *AgentFactory) ListTemplates() []AgentTemplate {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	templates := make([]AgentTemplate, 0, len(f.templates))
	for _, tmpl := range f.templates {
		templates = append(templates, tmpl)
	}
	
	return templates
}

// CreateAgent 使用模板创建智能体实例
func (f *AgentFactory) CreateAgent(ctx context.Context, templateID, name, description, model, provider string, customConfig map[string]interface{}) (*Agent, error) {
	template, err := f.GetTemplate(templateID)
	if err != nil {
		return nil, err
	}
	
	// 合并配置
	config := make(map[string]interface{})
	for k, v := range template.Config {
		config[k] = v
	}
	if customConfig != nil {
		for k, v := range customConfig {
			config[k] = v
		}
	}
	
	// 创建基础智能体
	agent, err := NewAgent(name, description, model, provider, config)
	if err != nil {
		return nil, err
	}
	
	// 设置系统提示词
	agent.SetSystemPrompt(template.SystemPrompt)
	
	// 添加工具
	if f.toolRegistry != nil && len(template.DefaultTools) > 0 {
		if err := f.toolRegistry.AddToolsToAgent(agent, template.DefaultTools...); err != nil {
			return nil, fmt.Errorf("failed to add tools: %w", err)
		}
	}
	
	// 初始化Eino Agent
	if f.modelProvider != nil {
		apiKey, err := f.modelProvider.GetAPIKey(provider)
		if err != nil {
			return nil, fmt.Errorf("failed to get API key: %w", err)
		}
		
		baseURL, err := f.modelProvider.GetBaseURL(provider)
		if err != nil {
			baseURL = "" // 使用默认URL
		}
		
		if err := agent.InitEinoAgent(ctx, apiKey, baseURL); err != nil {
			return nil, fmt.Errorf("failed to initialize agent: %w", err)
		}
	}
	
	return agent, nil
}

// CreateCustomAgent 创建自定义智能体
func (f *AgentFactory) CreateCustomAgent(ctx context.Context, name, description, model, provider, systemPrompt string, toolNames []string, config map[string]interface{}) (*Agent, error) {
	// 创建基础智能体
	agent, err := NewAgent(name, description, model, provider, config)
	if err != nil {
		return nil, err
	}
	
	// 设置系统提示词
	agent.SetSystemPrompt(systemPrompt)
	
	// 添加工具
	if f.toolRegistry != nil && len(toolNames) > 0 {
		if err := f.toolRegistry.AddToolsToAgent(agent, toolNames...); err != nil {
			return nil, fmt.Errorf("failed to add tools: %w", err)
		}
	}
	
	// 初始化Eino Agent
	if f.modelProvider != nil {
		apiKey, err := f.modelProvider.GetAPIKey(provider)
		if err != nil {
			return nil, fmt.Errorf("failed to get API key: %w", err)
		}
		
		baseURL, err := f.modelProvider.GetBaseURL(provider)
		if err != nil {
			baseURL = "" // 使用默认URL
		}
		
		if err := agent.InitEinoAgent(ctx, apiKey, baseURL); err != nil {
			return nil, fmt.Errorf("failed to initialize agent: %w", err)
		}
	}
	
	return agent, nil
}

// SerializeAgent 将智能体序列化为JSON
func (f *AgentFactory) SerializeAgent(agent *Agent) ([]byte, error) {
	return json.Marshal(agent)
}

// DeserializeAgent 从JSON反序列化智能体
func (f *AgentFactory) DeserializeAgent(data []byte) (*Agent, error) {
	agent, err := FromJSON(data)
	if err != nil {
		return nil, err
	}
	
	return agent, nil
}

// DefaultAgentFactory 默认的智能体工厂实例
var DefaultAgentFactory *AgentFactory

// 初始化默认工厂
func init() {
	// 注意：实际使用时需要提供真实的ModelProvider
	DefaultAgentFactory = NewAgentFactory(DefaultToolRegistry, nil)
} 