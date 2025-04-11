package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LogLevel 日志级别枚举
type LogLevel string

const (
	// LogLevelDebug 调试级别
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo 信息级别
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn 警告级别
	LogLevelWarn LogLevel = "warn"
	// LogLevelError 错误级别
	LogLevelError LogLevel = "error"
	// LogLevelFatal 致命级别
	LogLevelFatal LogLevel = "fatal"
)

// LogCategory 日志类别枚举
type LogCategory string

const (
	// LogCategorySystem 系统日志
	LogCategorySystem LogCategory = "system"
	// LogCategoryAPI API请求日志
	LogCategoryAPI LogCategory = "api"
	// LogCategoryAuth 认证日志
	LogCategoryAuth LogCategory = "auth"
	// LogCategoryModel 模型调用日志
	LogCategoryModel LogCategory = "model"
	// LogCategoryApplication 应用日志
	LogCategoryApplication LogCategory = "application"
)

// Log 基础日志模型
type Log struct {
	ID        uuid.UUID   `gorm:"type:uuid;primary_key" json:"id"`
	Level     LogLevel    `gorm:"type:varchar(10);not null;index" json:"level"`
	Category  LogCategory `gorm:"type:varchar(20);not null;index" json:"category"`
	Message   string      `gorm:"type:text;not null" json:"message"`
	UserID    *uuid.UUID  `gorm:"type:uuid;index" json:"user_id,omitempty"`
	User      *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Metadata  string      `gorm:"type:jsonb" json:"metadata"`
	CreatedAt time.Time   `json:"created_at"`
}

// BeforeCreate 在创建日志前生成UUID
func (l *Log) BeforeCreate(tx *gorm.DB) error {
	l.ID = uuid.New()
	return nil
}

// APILog API日志模型，扩展自Log
type APILog struct {
	Log
	Method     string `gorm:"type:varchar(10);not null" json:"method"`
	Path       string `gorm:"type:varchar(255);not null" json:"path"`
	StatusCode int    `gorm:"not null" json:"status_code"`
	IP         string `gorm:"type:varchar(50)" json:"ip"`
	UserAgent  string `gorm:"type:varchar(255)" json:"user_agent"`
	Duration   int64  `gorm:"not null" json:"duration"` // 毫秒
	RequestID  string `gorm:"type:varchar(36);index" json:"request_id"`
}

// ErrorLog 错误日志模型，扩展自Log
type ErrorLog struct {
	Log
	StackTrace string    `gorm:"type:text" json:"stack_trace"`
	ErrorCode  string    `gorm:"type:varchar(50)" json:"error_code"`
	Source     string    `gorm:"type:varchar(100)" json:"source"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	ResolvedBy *uuid.UUID `gorm:"type:uuid" json:"resolved_by,omitempty"`
}

// ModelCallLog 模型调用日志，记录LLM调用
type ModelCallLog struct {
	Log
	ModelName    string    `gorm:"type:varchar(100);not null;index" json:"model_name"`
	PromptTokens int       `json:"prompt_tokens"`
	CompTokens   int       `json:"comp_tokens"`
	TotalTokens  int       `json:"total_tokens"`
	Duration     int64     `gorm:"not null" json:"duration"` // 毫秒
	ApplicationID *uuid.UUID `gorm:"type:uuid;index" json:"application_id,omitempty"`
	ProjectID    *uuid.UUID `gorm:"type:uuid;index" json:"project_id,omitempty"`
	Success      bool      `gorm:"not null" json:"success"`
}

// SystemMetric 系统指标记录
type SystemMetric struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	MetricName  string    `gorm:"type:varchar(100);not null;index" json:"metric_name"`
	MetricValue float64   `gorm:"not null" json:"metric_value"`
	Unit        string    `gorm:"type:varchar(20)" json:"unit"`
	Tags        string    `gorm:"type:jsonb" json:"tags"` // JSON格式的标签
	CreatedAt   time.Time `json:"created_at"`
}

// BeforeCreate 在创建系统指标前生成UUID
func (sm *SystemMetric) BeforeCreate(tx *gorm.DB) error {
	sm.ID = uuid.New()
	return nil
}

// LogResponse 是返回给客户端的日志数据结构
type LogResponse struct {
	ID        uuid.UUID   `json:"id"`
	Level     LogLevel    `json:"level"`
	Category  LogCategory `json:"category"`
	Message   string      `json:"message"`
	UserID    *uuid.UUID  `json:"user_id,omitempty"`
	UserName  string      `json:"user_name,omitempty"`
	Metadata  interface{} `json:"metadata,omitempty"`
	CreatedAt time.Time   `json:"created_at"`
	// 扩展字段，根据日志类型填充
	Method      string     `json:"method,omitempty"`
	Path        string     `json:"path,omitempty"`
	StatusCode  int        `json:"status_code,omitempty"`
	Duration    int64      `json:"duration,omitempty"`
	ModelName   string     `json:"model_name,omitempty"`
	TotalTokens int        `json:"total_tokens,omitempty"`
	ErrorCode   string     `json:"error_code,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
} 