package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ModelProvider 表示模型提供者类型
type ModelProvider string

// 模型提供者常量
const (
	ModelProviderOpenAI    ModelProvider = "openai"
	ModelProviderAnthropic ModelProvider = "anthropic"
	ModelProviderBaidu     ModelProvider = "baidu"
	ModelProviderAli       ModelProvider = "aliyun"
	ModelProviderLocal     ModelProvider = "local"
	ModelProviderCustom    ModelProvider = "custom"
)

// String 返回模型提供者的字符串表示
func (m ModelProvider) String() string {
	return string(m)
}

// Scan 实现 sql.Scanner 接口
func (m *ModelProvider) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			*m = ModelProvider(v)
			return nil
		}
	}
	return errors.New("无法将数据库值转换为ModelProvider")
}

// Value 实现 driver.Valuer 接口
func (m ModelProvider) Value() (driver.Value, error) {
	return string(m), nil
}

// ModelStatus 表示模型状态
type ModelStatus string

// 模型状态常量
const (
	ModelStatusActive   ModelStatus = "active"
	ModelStatusInactive ModelStatus = "inactive"
	ModelStatusError    ModelStatus = "error"
)

// String 返回模型状态的字符串表示
func (m ModelStatus) String() string {
	return string(m)
}

// Scan 实现 sql.Scanner 接口
func (m *ModelStatus) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			*m = ModelStatus(v)
			return nil
		}
	}
	return errors.New("无法将数据库值转换为ModelStatus")
}

// Value 实现 driver.Valuer 接口
func (m ModelStatus) Value() (driver.Value, error) {
	return string(m), nil
}

// ModelType 表示模型类型
type ModelType string

// 模型类型常量
const (
	ModelTypeText          ModelType = "text"
	ModelTypeEmbedding     ModelType = "embedding"
	ModelTypeMultimodal    ModelType = "multimodal"
	ModelTypeFineTuned     ModelType = "fine-tuned"
)

// String 返回模型类型的字符串表示
func (m ModelType) String() string {
	return string(m)
}

// Scan 实现 sql.Scanner 接口
func (m *ModelType) Scan(value interface{}) error {
	if value == nil {
		*m = ""
		return nil
	}
	if sv, err := driver.String.ConvertValue(value); err == nil {
		if v, ok := sv.(string); ok {
			*m = ModelType(v)
			return nil
		}
	}
	return errors.New("无法将数据库值转换为ModelType")
}

// Value 实现 driver.Valuer 接口
func (m ModelType) Value() (driver.Value, error) {
	return string(m), nil
}

// ModelParameters 表示模型参数配置
type ModelParameters struct {
	Temperature      *float32 `json:"temperature,omitempty"`       // 温度参数
	TopP             *float32 `json:"top_p,omitempty"`             // 核采样参数
	TopK             *int     `json:"top_k,omitempty"`             // Top-K采样参数
	MaxTokens        *int     `json:"max_tokens,omitempty"`        // 生成的最大token数
	PresencePenalty  *float32 `json:"presence_penalty,omitempty"`  // 存在惩罚因子
	FrequencyPenalty *float32 `json:"frequency_penalty,omitempty"` // 频率惩罚因子
	Stop             []string `json:"stop,omitempty"`              // 停止生成的标记
}

// Scan 实现 sql.Scanner 接口
func (p *ModelParameters) Scan(value interface{}) error {
	if value == nil {
		*p = ModelParameters{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, p)
	case string:
		return json.Unmarshal([]byte(v), p)
	default:
		return errors.New("无法将数据库值转换为ModelParameters")
	}
}

// Value 实现 driver.Valuer 接口
func (p ModelParameters) Value() (driver.Value, error) {
	// 检查结构体是否为零值
	if p.Temperature == nil && p.TopP == nil && p.TopK == nil && 
	   p.MaxTokens == nil && p.PresencePenalty == nil && 
	   p.FrequencyPenalty == nil && len(p.Stop) == 0 {
		return nil, nil
	}
	return json.Marshal(p)
}

// ModelUsageMetrics 模型使用指标
type ModelUsageMetrics struct {
	TotalCalls      int64   `json:"total_calls"`       // 总调用次数
	SuccessfulCalls int64   `json:"successful_calls"`  // 成功调用次数
	FailedCalls     int64   `json:"failed_calls"`      // 失败调用次数
	TotalTokens     int64   `json:"total_tokens"`      // 消耗的总token数
	AverageLatency  float64 `json:"average_latency"`   // 平均延迟(ms)
	Cost            float64 `json:"cost"`              // 总费用(USD)
	LastUsedAt      string  `json:"last_used_at"`      // 最近一次使用时间
}

// ModelProviderConfig 各提供者的特定配置
type ModelProviderConfig struct {
	ApiKey     string `json:"api_key,omitempty"`      // API密钥(加密存储)
	ApiSecret  string `json:"api_secret,omitempty"`   // API密钥(加密存储)
	BaseURL    string `json:"base_url,omitempty"`     // API基础URL
	OrgID      string `json:"org_id,omitempty"`       // 组织ID
	AppID      string `json:"app_id,omitempty"`       // 应用ID(如百度需要)
	Version    string `json:"version,omitempty"`      // API版本
	Deployment string `json:"deployment,omitempty"`   // 部署ID(如Azure OpenAI)
	Region     string `json:"region,omitempty"`       // 区域设置
	ProxyURL   string `json:"proxy_url,omitempty"`    // 代理地址
	Extra      string `json:"extra,omitempty"`        // 其他配置(JSON)
}

// Scan 实现 sql.Scanner 接口
func (c *ModelProviderConfig) Scan(value interface{}) error {
	if value == nil {
		*c = ModelProviderConfig{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, c)
	case string:
		return json.Unmarshal([]byte(v), c)
	default:
		return errors.New("无法将数据库值转换为ModelProviderConfig")
	}
}

// Value 实现 driver.Valuer 接口
func (c ModelProviderConfig) Value() (driver.Value, error) {
	if c == (ModelProviderConfig{}) {
		return nil, nil
	}
	return json.Marshal(c)
}

// Model 表示模型实体
type Model struct {
	ID              uuid.UUID         `gorm:"type:uuid;primary_key" json:"id"`
	Name            string            `gorm:"type:varchar(100);not null;unique" json:"name"`
	Provider        ModelProvider     `gorm:"type:varchar(50);not null" json:"provider"`
	ModelID         string            `gorm:"type:varchar(100);not null" json:"model_id"`
	Type            ModelType         `gorm:"type:varchar(20);not null" json:"type"`
	Description     string            `gorm:"type:text" json:"description"`
	Capabilities    []string          `gorm:"type:text[]" json:"capabilities"`
	Parameters      ModelParameters   `gorm:"type:jsonb" json:"parameters"`
	MaxTokens       int               `gorm:"not null" json:"max_tokens"`
	TokenCostPrompt float64           `gorm:"not null" json:"token_cost_prompt"`
	TokenCostCompl   float64           `gorm:"not null" json:"token_cost_completion"`
	Status          ModelStatus       `gorm:"type:varchar(20);not null;default:'active'" json:"status"`
	ProviderConfig  ModelProviderConfig `gorm:"type:jsonb" json:"provider_config"`
	IsSystem        bool              `gorm:"not null;default:false" json:"is_system"`
	IsCustom        bool              `gorm:"not null;default:false" json:"is_custom"`
	OrganizationID  *uuid.UUID        `gorm:"type:uuid" json:"organization_id,omitempty"`
	CreatedBy       *uuid.UUID        `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	DeletedAt       gorm.DeletedAt    `gorm:"index" json:"-"`
}

// BeforeCreate 在创建模型前生成UUID
func (m *Model) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}

// ModelConfig 表示用户或应用的模型配置
type ModelConfig struct {
	ID              uuid.UUID           `json:"id" gorm:"type:uuid;primary_key;"`
	Name            string              `json:"name" gorm:"size:128;not null;"`
	Description     string              `json:"description" gorm:"type:text;"`
	ModelID         uuid.UUID           `json:"model_id" gorm:"type:uuid;not null;"`
	Model           *Model              `json:"model" gorm:"foreignKey:ModelID;"`
	Parameters      ModelParameters     `json:"parameters" gorm:"type:jsonb;default:'{}'::jsonb;"`
	ProviderConfig  ModelProviderConfig `json:"provider_config" gorm:"type:jsonb;default:'{}'::jsonb;"`
	IsShared        bool                `json:"is_shared" gorm:"default:false;"`
	UsageMetrics    ModelUsageMetrics   `json:"usage_metrics" gorm:"-"`
	UsageMetricsStr string              `json:"-" gorm:"column:usage_metrics;type:text;"`
	OrganizationID  uuid.UUID           `json:"organization_id" gorm:"type:uuid;not null;"`
	CreatedBy       uuid.UUID           `json:"created_by" gorm:"type:uuid;not null;"`
	CreatedAt       time.Time           `json:"created_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP;"`
	UpdatedAt       time.Time           `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:CURRENT_TIMESTAMP;"`
}

// BeforeSave GORM钩子，保存前处理
func (c *ModelConfig) BeforeSave() error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	
	// 将使用指标转为JSON字符串
	if c.UsageMetrics != (ModelUsageMetrics{}) {
		bytes, err := json.Marshal(c.UsageMetrics)
		if err != nil {
			return err
		}
		c.UsageMetricsStr = string(bytes)
	}
	
	return nil
}

// AfterFind GORM钩子，查询后处理
func (c *ModelConfig) AfterFind() error {
	// 将JSON字符串转为使用指标
	if c.UsageMetricsStr != "" {
		return json.Unmarshal([]byte(c.UsageMetricsStr), &c.UsageMetrics)
	}
	return nil
}

// ModelResponse 模型API响应
type ModelResponse struct {
	ID          uuid.UUID     `json:"id"`
	Name        string        `json:"name"`
	Provider    ModelProvider `json:"provider"`
	ModelID     string        `json:"model_id"`
	Type        ModelType     `json:"type"`
	Description string        `json:"description"`
	Capabilities []string     `json:"capabilities"`
	Parameters  ModelParameters `json:"parameters"`
	MaxTokens   int           `json:"max_tokens"`
	TokenCost   struct {
		Prompt     float64 `json:"prompt"`
		Completion float64 `json:"completion"`
	} `json:"token_cost"`
	Status      ModelStatus   `json:"status"`
	IsSystem    bool          `json:"is_system"`
	CreatedAt   time.Time     `json:"created_at"`
}

// ModelConfigResponse 是返回给客户端的模型配置结构
type ModelConfigResponse struct {
	ID             uuid.UUID      `json:"id"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Model          ModelResponse  `json:"model"`
	Parameters     ModelParameters `json:"parameters"`
	IsShared       bool           `json:"is_shared"`
	UsageMetrics   ModelUsageMetrics `json:"usage_metrics"`
	OrganizationID uuid.UUID      `json:"organization_id"`
	CreatedBy      uuid.UUID      `json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
} 
