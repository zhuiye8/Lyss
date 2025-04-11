package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project 模型表示用户创建的项目
type Project struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	Owner       User      `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Status      string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active, archived
	Public      bool      `gorm:"default:false" json:"public"`
	Config      string    `gorm:"type:jsonb" json:"config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在创建项目前生成UUID
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New()
	return nil
}

// ProjectResponse 是返回给客户端的项目数据结构
type ProjectResponse struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	OwnerID     uuid.UUID      `json:"owner_id"`
	Status      string         `json:"status"`
	Public      bool           `json:"public"`
	Config      string         `json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	AppsCount   int64          `json:"apps_count,omitempty"`
}

// ToResponse 将完整项目模型转换为对外响应
func (p *Project) ToResponse() ProjectResponse {
	return ProjectResponse{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OwnerID:     p.OwnerID,
		Status:      p.Status,
		Public:      p.Public,
		Config:      p.Config,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
} 
