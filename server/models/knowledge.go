package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// KnowledgeBase 模型表示知识库
type KnowledgeBase struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type"` // file, web, database
	Config      string    `gorm:"type:jsonb" json:"config"`
	Status      string    `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在创建知识库前生成UUID
func (kb *KnowledgeBase) BeforeCreate(tx *gorm.DB) error {
	kb.ID = uuid.New()
	return nil
}

// KnowledgeBaseResponse 是返回给客户端的知识库数据结构
type KnowledgeBaseResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Config      string    `json:"config"`
	Status      string    `json:"status"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse 将完整知识库模型转换为对外响应
func (kb *KnowledgeBase) ToResponse() KnowledgeBaseResponse {
	return KnowledgeBaseResponse{
		ID:          kb.ID,
		Name:        kb.Name,
		Description: kb.Description,
		Type:        kb.Type,
		Config:      kb.Config,
		Status:      kb.Status,
		CreatedBy:   kb.CreatedBy,
		CreatedAt:   kb.CreatedAt,
		UpdatedAt:   kb.UpdatedAt,
	}
} 