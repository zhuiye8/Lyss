package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Agent 智能体模型
type Agent struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Name            string         `gorm:"type:varchar(128);not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	ApplicationID   uuid.UUID      `gorm:"type:uuid;not null" json:"application_id"`
	ModelConfigID   uuid.UUID      `gorm:"type:uuid;not null" json:"model_config_id"`
	SystemPrompt    string         `gorm:"type:text" json:"system_prompt"`
	Tools           JSONMap        `gorm:"type:jsonb" json:"tools"`
	Variables       JSONMap        `gorm:"type:jsonb" json:"variables"`
	MaxHistoryLength int           `gorm:"default:10" json:"max_history_length"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Application   *Application `gorm:"foreignKey:ApplicationID" json:"-"`
	ModelConfig   *ModelConfig `gorm:"foreignKey:ModelConfigID" json:"-"`
	Conversations []Conversation `gorm:"foreignKey:AgentID" json:"-"`
}

// JSONMap 是JSON字段的辅助类型
type JSONMap map[string]interface{}

// AgentKnowledgeBase 智能体与知识库的多对多关联
type AgentKnowledgeBase struct {
	AgentID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"agent_id"`
	KnowledgeBaseID uuid.UUID `gorm:"type:uuid;primaryKey" json:"knowledge_base_id"`
	CreatedAt      time.Time `json:"created_at"`
}

// BeforeCreate 在创建智能体前生成UUID
func (a *Agent) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}

// AgentResponse 是返回给客户端的智能体数据结构
type AgentResponse struct {
	ID               uuid.UUID      `json:"id"`
	Name             string         `json:"name"`
	Description      string         `json:"description"`
	ApplicationID    uuid.UUID      `json:"application_id"`
	ModelConfigID    uuid.UUID      `json:"model_config_id"`
	SystemPrompt     string         `json:"system_prompt"`
	Tools            JSONMap        `json:"tools"`
	Variables        JSONMap        `json:"variables"`
	MaxHistoryLength int            `json:"max_history_length"`
	KnowledgeBaseIDs []uuid.UUID    `json:"knowledge_base_ids,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

// ToResponse 将完整智能体模型转换为对外响应
func (a *Agent) ToResponse(knowledgeBaseIDs []uuid.UUID) AgentResponse {
	return AgentResponse{
		ID:               a.ID,
		Name:             a.Name,
		Description:      a.Description,
		ApplicationID:    a.ApplicationID,
		ModelConfigID:    a.ModelConfigID,
		SystemPrompt:     a.SystemPrompt,
		Tools:            a.Tools,
		Variables:        a.Variables,
		MaxHistoryLength: a.MaxHistoryLength,
		KnowledgeBaseIDs: knowledgeBaseIDs,
		CreatedAt:        a.CreatedAt,
		UpdatedAt:        a.UpdatedAt,
	}
}

// CreateAgentRequest 创建智能体请求
type CreateAgentRequest struct {
	Name             string    `json:"name" binding:"required,min=1,max=128"`
	Description      string    `json:"description"`
	ModelConfigID    uuid.UUID `json:"model_config_id" binding:"required"`
	SystemPrompt     string    `json:"system_prompt"`
	Tools            JSONMap   `json:"tools"`
	Variables        JSONMap   `json:"variables"`
	MaxHistoryLength int       `json:"max_history_length"`
	KnowledgeBaseIDs []uuid.UUID `json:"knowledge_base_ids"`
}

// UpdateAgentRequest 更新智能体请求
type UpdateAgentRequest struct {
	Name             string    `json:"name" binding:"omitempty,min=1,max=128"`
	Description      string    `json:"description"`
	ModelConfigID    uuid.UUID `json:"model_config_id"`
	SystemPrompt     string    `json:"system_prompt"`
	Tools            JSONMap   `json:"tools"`
	Variables        JSONMap   `json:"variables"`
	MaxHistoryLength int       `json:"max_history_length"`
	KnowledgeBaseIDs []uuid.UUID `json:"knowledge_base_ids"`
}

// UpdateSystemPromptRequest 更新系统提示词请求
type UpdateSystemPromptRequest struct {
	SystemPrompt string `json:"system_prompt" binding:"required"`
}

// UpdateToolsRequest 更新工具配置请求
type UpdateToolsRequest struct {
	Tools JSONMap `json:"tools" binding:"required"`
}

// TestAgentRequest 测试智能体请求
type TestAgentRequest struct {
	Message string `json:"message" binding:"required"`
} 
