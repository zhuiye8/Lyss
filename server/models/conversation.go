package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Conversation 对话模型
type Conversation struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key" json:"id"`
	Title     string         `gorm:"type:varchar(255)" json:"title"`
	AgentID   uuid.UUID      `gorm:"type:uuid;not null" json:"agent_id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Status    string         `gorm:"type:varchar(16);default:'active'" json:"status"` // active, archived
	Metadata  JSONMap        `gorm:"type:jsonb" json:"metadata"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Agent    *Agent    `gorm:"foreignKey:AgentID" json:"-"`
	User     *User     `gorm:"foreignKey:UserID" json:"-"`
	Messages []Message `gorm:"foreignKey:ConversationID" json:"-"`
}

// Message 消息模型
type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null" json:"conversation_id"`
	Role           string    `gorm:"type:varchar(16);not null" json:"role"` // user, assistant, system
	Content        string    `gorm:"type:text;not null" json:"content"`
	Tokens         int       `json:"tokens"`
	Feedback       string    `gorm:"type:varchar(16)" json:"feedback"` // positive, negative, null
	Metadata       JSONMap   `gorm:"type:jsonb" json:"metadata"`
	CreatedAt      time.Time `json:"created_at"`

	// 关联
	Conversation *Conversation `gorm:"foreignKey:ConversationID" json:"-"`
}

// BeforeCreate 在创建对话前生成UUID
func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}

// BeforeCreate 在创建消息前生成UUID
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}

// ConversationResponse 是返回给客户端的对话数据结构
type ConversationResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	AgentID   uuid.UUID `json:"agent_id"`
	AgentName string    `json:"agent_name,omitempty"`
	UserID    uuid.UUID `json:"user_id"`
	Status    string    `json:"status"`
	Metadata  JSONMap   `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MessageResponse 是返回给客户端的消息数据结构
type MessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Role           string    `json:"role"`
	Content        string    `json:"content"`
	Tokens         int       `json:"tokens"`
	Feedback       string    `json:"feedback,omitempty"`
	Metadata       JSONMap   `json:"metadata,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

// ToResponse 将完整对话模型转换为对外响应
func (c *Conversation) ToResponse() ConversationResponse {
	var agentName string
	if c.Agent != nil {
		agentName = c.Agent.Name
	}

	return ConversationResponse{
		ID:        c.ID,
		Title:     c.Title,
		AgentID:   c.AgentID,
		AgentName: agentName,
		UserID:    c.UserID,
		Status:    c.Status,
		Metadata:  c.Metadata,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// ToResponse 将完整消息模型转换为对外响应
func (m *Message) ToResponse() MessageResponse {
	return MessageResponse{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		Role:           m.Role,
		Content:        m.Content,
		Tokens:         m.Tokens,
		Feedback:       m.Feedback,
		Metadata:       m.Metadata,
		CreatedAt:      m.CreatedAt,
	}
}

// CreateConversationRequest 创建对话请求
type CreateConversationRequest struct {
	AgentID  uuid.UUID `json:"agent_id" binding:"required"`
	Title    string    `json:"title"`
	Metadata JSONMap   `json:"metadata"`
}

// SendMessageRequest 发送消息请求
type SendMessageRequest struct {
	Content string  `json:"content" binding:"required"`
	Stream  bool    `json:"stream"`
}

// RegenerateRequest 重新生成回复请求
type RegenerateRequest struct {
	Stream bool `json:"stream"`
}

// MessageFeedbackRequest 提供消息反馈请求
type MessageFeedbackRequest struct {
	Feedback string `json:"feedback" binding:"required,oneof=positive negative"`
} 
