package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ConfigScope 定义配置的作用域
type ConfigScope string

const (
	// ScopeSystem 系统级别配置
	ScopeSystem ConfigScope = "system"
	// ScopeUser 用户级别配置
	ScopeUser ConfigScope = "user"
	// ScopeProject 项目级别配置
	ScopeProject ConfigScope = "project"
	// ScopeApplication 应用级别配置
	ScopeApplication ConfigScope = "application"
)

// Config 模型表示系统配置定义
type Config struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
	Key       string      `gorm:"type:varchar(100);not null;uniqueIndex:idx_config_scope_key" json:"key"`
	Value     string      `gorm:"type:text" json:"value"`
	Scope     ConfigScope `gorm:"type:varchar(20);not null;uniqueIndex:idx_config_scope_key" json:"scope"`
	ScopeID   *uuid.UUID  `gorm:"type:uuid;uniqueIndex:idx_config_scope_key" json:"scope_id"`
	CreatedBy *uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	UpdatedBy *uuid.UUID  `gorm:"type:uuid" json:"updated_by"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在创建配置前生成UUID
func (c *Config) BeforeCreate(tx *gorm.DB) error {
	c.ID = uuid.New()
	return nil
}

// ConfigResponse 是返回给客户端的配置数据结构
type ConfigResponse struct {
	ID        uuid.UUID   `json:"id"`
	Key       string      `json:"key"`
	Value     string      `json:"value"`
	Scope     ConfigScope `json:"scope"`
	ScopeID   *uuid.UUID  `json:"scope_id,omitempty"`
	CreatedBy *uuid.UUID  `json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID  `json:"updated_by,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// ToResponse 将完整配置模型转换为对外响应
func (c *Config) ToResponse() ConfigResponse {
	return ConfigResponse{
		ID:        c.ID,
		Key:       c.Key,
		Value:     c.Value,
		Scope:     c.Scope,
		ScopeID:   c.ScopeID,
		CreatedBy: c.CreatedBy,
		UpdatedBy: c.UpdatedBy,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
} 
