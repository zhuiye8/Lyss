package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Application 模型表示项目中的智能体应用
type Application struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"type:varchar(50);not null" json:"type"` // chat, workflow, custom
	ProjectID   uuid.UUID `gorm:"type:uuid;not null" json:"project_id"`
	Project     Project   `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Status      string    `gorm:"type:varchar(20);default:'draft'" json:"status"` // draft, published, archived
	Config      string    `gorm:"type:jsonb" json:"config"`
	ModelConfig string    `gorm:"type:jsonb" json:"model_config"`
	CreatedBy   uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在创建应用前生成UUID
func (a *Application) BeforeCreate(tx *gorm.DB) error {
	a.ID = uuid.New()
	return nil
}

// ApplicationResponse 是返回给客户端的应用数据结构
type ApplicationResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	ProjectID   uuid.UUID `json:"project_id"`
	Status      string    `json:"status"`
	Config      string    `json:"config"`
	ModelConfig string    `json:"model_config"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ToResponse 将完整应用模型转换为对外响应
func (a *Application) ToResponse() ApplicationResponse {
	return ApplicationResponse{
		ID:          a.ID,
		Name:        a.Name,
		Description: a.Description,
		Type:        a.Type,
		ProjectID:   a.ProjectID,
		Status:      a.Status,
		Config:      a.Config,
		ModelConfig: a.ModelConfig,
		CreatedBy:   a.CreatedBy,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
} 
