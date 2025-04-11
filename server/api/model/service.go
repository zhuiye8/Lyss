package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"github.com/zhuiye8/Lyss/server/pkg/encryption"
	"gorm.io/gorm"
)

var (
	ErrModelNotFound      = errors.New("模型不存在")
	ErrModelConfigNotFound = errors.New("模型配置不存在")
	ErrDuplicateModelName = errors.New("模型名称已存在")
	ErrInvalidProvider    = errors.New("无效的模型提供商")
	ErrNoPermission       = errors.New("没有操作权限")
)

// Service 提供模型管理功能
type Service struct {
	db           *gorm.DB
	encryptor    *encryption.Service
}

// NewService 创建模型服务
func NewService(db *gorm.DB, encryptor *encryption.Service) *Service {
	return &Service{
		db:        db,
		encryptor: encryptor,
	}
}

// GetModels 获取模型列表
func (s *Service) GetModels(params ModelsQueryParams) ([]models.Model, int64, error) {
	var modelList []models.Model
	var totalCount int64
	
	query := s.db.Model(&models.Model{})
	
	// 应用过滤条件
	if params.Provider != "" {
		query = query.Where("provider = ?", params.Provider)
	}
	
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}
	
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	
	// 如果不包含系统模型，则过滤
	if !params.IncludeSystem {
		query = query.Where("is_system = ?", false)
	}
	
	// 组织过滤
	if params.OrganizationID != nil {
		query = query.Where("organization_id = ? OR organization_id IS NULL", params.OrganizationID)
	}
	
	// 名称搜索
	if params.Search != "" {
		query = query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", params.Search))
	}
	
	// 计算总数
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	
	// 排序
	if params.SortBy != "" {
		direction := "ASC"
		if params.SortDesc {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", params.SortBy, direction))
	} else {
		query = query.Order("name ASC")
	}
	
	// 分页
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}
	
	// 执行查询
	if err := query.Find(&modelList).Error; err != nil {
		return nil, 0, err
	}
	
	return modelList, totalCount, nil
}

// GetModelByID 通过ID获取模型
func (s *Service) GetModelByID(id uuid.UUID) (*models.Model, error) {
	var model models.Model
	if err := s.db.First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelNotFound
		}
		return nil, err
	}
	return &model, nil
}

// CreateModel 创建新模型
func (s *Service) CreateModel(model *models.Model) error {
	// 检查是否已存在同名模型
	var count int64
	if err := s.db.Model(&models.Model{}).Where("name = ? AND (organization_id = ? OR organization_id IS NULL)", 
		model.Name, model.OrganizationID).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return ErrDuplicateModelName
	}
	
	// 加密敏感信息
	if err := s.encryptProviderConfig(&model.ProviderConfig); err != nil {
		return err
	}
	
	// 设置创建时间
	now := time.Now()
	model.CreatedAt = now
	model.UpdatedAt = now
	
	// 创建模型
	return s.db.Create(model).Error
}

// UpdateModel 更新模型
func (s *Service) UpdateModel(id uuid.UUID, updateData *models.Model) error {
	model, err := s.GetModelByID(id)
	if err != nil {
		return err
	}
	
	// 检查权限
	if model.IsSystem && !updateData.IsSystem {
		return ErrNoPermission
	}
	
	// 检查是否存在同名模型
	var count int64
	if err := s.db.Model(&models.Model{}).Where("name = ? AND id != ? AND (organization_id = ? OR organization_id IS NULL)", 
		updateData.Name, id, model.OrganizationID).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return ErrDuplicateModelName
	}
	
	// 准备更新数据
	updateMap := map[string]interface{}{
		"name":             updateData.Name,
		"model_id":         updateData.ModelID,
		"description":      updateData.Description,
		"parameters":       updateData.Parameters,
		"max_tokens":       updateData.MaxTokens,
		"token_cost_prompt": updateData.TokenCostPrompt,
		"token_cost_completion": updateData.TokenCostCompl,
		"status":           updateData.Status,
		"updated_at":       time.Now(),
	}
	
	// 如果提供了能力数组
	if len(updateData.Capabilities) > 0 {
		updateMap["capabilities"] = updateData.Capabilities
	}
	
	// 加密并更新提供者配置
	if updateData.ProviderConfig != (models.ModelProviderConfig{}) {
		if err := s.encryptProviderConfig(&updateData.ProviderConfig); err != nil {
			return err
		}
		updateMap["provider_config"] = updateData.ProviderConfig
	}
	
	// 执行更新
	if err := s.db.Model(&models.Model{}).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return err
	}
	
	return nil
}

// DeleteModel 删除模型
func (s *Service) DeleteModel(id uuid.UUID) error {
	model, err := s.GetModelByID(id)
	if err != nil {
		return err
	}
	
	// 检查权限- 系统模型不能删除
	if model.IsSystem {
		return ErrNoPermission
	}
	
	// 检查是否有模型配置在使用这个模型
	var count int64
	if err := s.db.Model(&models.ModelConfig{}).Where("model_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		// 如果有配置在使用，则只将状态设为不活跃
		return s.db.Model(&models.Model{}).Where("id = ?", id).
			Updates(map[string]interface{}{
				"status": models.ModelStatusInactive,
				"updated_at": time.Now(),
			}).Error
	}
	
	// 如果没有配置使用，则可以直接删除
	return s.db.Delete(&models.Model{}, "id = ?", id).Error
}

// GetModelConfigs 获取模型配置列表
func (s *Service) GetModelConfigs(params ConfigsQueryParams) ([]models.ModelConfig, int64, error) {
	var configs []models.ModelConfig
	var totalCount int64
	
	query := s.db.Model(&models.ModelConfig{})
	
	// 应用过滤条件
	if params.ModelID != nil {
		query = query.Where("model_id = ?", params.ModelID)
	}
	
	if params.OrganizationID != nil {
		query = query.Where("organization_id = ?", params.OrganizationID)
	}
	
	if params.CreatedBy != nil {
		query = query.Where("created_by = ?", params.CreatedBy)
	}
	
	// 共享过滤
	if params.Shared {
		query = query.Where("is_shared = ?", true)
	}
	
	// 名称搜索
	if params.Search != "" {
		query = query.Where("name ILIKE ?", fmt.Sprintf("%%%s%%", params.Search))
	}
	
	// 计算总数
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	
	// 排序
	if params.SortBy != "" {
		direction := "ASC"
		if params.SortDesc {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", params.SortBy, direction))
	} else {
		query = query.Order("name ASC")
	}
	
	// 分页
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}
	
	// 预加载模型
	query = query.Preload("Model")
	
	// 执行查询
	if err := query.Find(&configs).Error; err != nil {
		return nil, 0, err
	}
	
	return configs, totalCount, nil
}

// GetModelConfigByID 通过ID获取模型配置
func (s *Service) GetModelConfigByID(id uuid.UUID) (*models.ModelConfig, error) {
	var config models.ModelConfig
	if err := s.db.Preload("Model").First(&config, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelConfigNotFound
		}
		return nil, err
	}
	return &config, nil
}

// CreateModelConfig 创建新的模型配置
func (s *Service) CreateModelConfig(config *models.ModelConfig) error {
	// 检查模型是否存在
	if _, err := s.GetModelByID(config.ModelID); err != nil {
		return err
	}
	
	// 检查是否已存在同名配置
	var count int64
	if err := s.db.Model(&models.ModelConfig{}).Where("name = ? AND organization_id = ?", 
		config.Name, config.OrganizationID).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("配置名称已存在")
	}
	
	// 加密敏感信息
	if err := s.encryptProviderConfig(&config.ProviderConfig); err != nil {
		return err
	}
	
	// 初始化使用指标
	config.UsageMetrics = models.ModelUsageMetrics{
		TotalCalls:      0,
		SuccessfulCalls: 0,
		FailedCalls:     0,
		TotalTokens:     0,
		AverageLatency:  0,
		Cost:            0,
		LastUsedAt:      time.Now().Format(time.RFC3339),
	}
	
	// 设置创建时间
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now
	
	// 创建配置
	return s.db.Create(config).Error
}

// UpdateModelConfig 更新模型配置
func (s *Service) UpdateModelConfig(id uuid.UUID, updateData *models.ModelConfig, userID uuid.UUID) error {
	config, err := s.GetModelConfigByID(id)
	if err != nil {
		return err
	}
	
	// 检查权限
	if config.CreatedBy != userID {
		return ErrNoPermission
	}
	
	// 检查是否存在同名配置
	var count int64
	if err := s.db.Model(&models.ModelConfig{}).Where("name = ? AND id != ? AND organization_id = ?", 
		updateData.Name, id, config.OrganizationID).Count(&count).Error; err != nil {
		return err
	}
	
	if count > 0 {
		return errors.New("配置名称已存在")
	}
	
	// 准备更新数据
	updateMap := map[string]interface{}{
		"name":         updateData.Name,
		"description":  updateData.Description,
		"parameters":   updateData.Parameters,
		"is_shared":    updateData.IsShared,
		"updated_at":   time.Now(),
	}
	
	// 检查是否更改了模型
	if updateData.ModelID != uuid.Nil && updateData.ModelID != config.ModelID {
		// 验证新模型存在
		if _, err := s.GetModelByID(updateData.ModelID); err != nil {
			return err
		}
		updateMap["model_id"] = updateData.ModelID
	}
	
	// 加密并更新提供者配置
	if updateData.ProviderConfig != (models.ModelProviderConfig{}) {
		if err := s.encryptProviderConfig(&updateData.ProviderConfig); err != nil {
			return err
		}
		updateMap["provider_config"] = updateData.ProviderConfig
	}
	
	// 执行更新
	if err := s.db.Model(&models.ModelConfig{}).Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return err
	}
	
	return nil
}

// DeleteModelConfig 删除模型配置
func (s *Service) DeleteModelConfig(id uuid.UUID, userID uuid.UUID) error {
	config, err := s.GetModelConfigByID(id)
	if err != nil {
		return err
	}
	
	// 检查权限
	if config.CreatedBy != userID {
		return ErrNoPermission
	}
	
	// 检查是否有应用在使用这个配置
	// TODO: 检查应用和智能体引用
	
	// 删除配置
	return s.db.Delete(&models.ModelConfig{}, "id = ?", id).Error
}

// UpdateModelUsageMetrics 更新模型使用指标
func (s *Service) UpdateModelUsageMetrics(configID uuid.UUID, success bool, tokens int, latency time.Duration, cost float64) error {
	config, err := s.GetModelConfigByID(configID)
	if err != nil {
		return err
	}
	
	// 更新指标
	metrics := config.UsageMetrics
	metrics.TotalCalls++
	
	if success {
		metrics.SuccessfulCalls++
	} else {
		metrics.FailedCalls++
	}
	
	metrics.TotalTokens += int64(tokens)
	
	// 更新平均延迟
	if metrics.TotalCalls > 1 {
		metrics.AverageLatency = (metrics.AverageLatency*float64(metrics.TotalCalls-1) + float64(latency.Milliseconds())) / float64(metrics.TotalCalls)
	} else {
		metrics.AverageLatency = float64(latency.Milliseconds())
	}
	
	metrics.Cost += cost
	metrics.LastUsedAt = time.Now().Format(time.RFC3339)
	
	// 保存指标
	config.UsageMetrics = metrics
	
	// 序列化指标
	if err := config.BeforeSave(); err != nil {
		return err
	}
	
	// 更新数据
	return s.db.Model(&models.ModelConfig{}).Where("id = ?", configID).
		Update("usage_metrics", config.UsageMetricsStr).Error
}

// 辅助方法 - 加密提供者配置中的敏感信息
func (s *Service) encryptProviderConfig(config *models.ModelProviderConfig) error {
	if config.ApiKey != "" {
		encrypted, err := s.encryptor.Encrypt(config.ApiKey)
		if err != nil {
			return err
		}
		config.ApiKey = encrypted
	}
	
	if config.ApiSecret != "" {
		encrypted, err := s.encryptor.Encrypt(config.ApiSecret)
		if err != nil {
			return err
		}
		config.ApiSecret = encrypted
	}
	
	return nil
}

// ModelsQueryParams 模型查询参数
type ModelsQueryParams struct {
	Provider       string
	Type           string
	Status         string
	IncludeSystem  bool
	OrganizationID *uuid.UUID
	Search         string
	Page           int
	PageSize       int
	SortBy         string
	SortDesc       bool
}

// ConfigsQueryParams 模型配置查询参数
type ConfigsQueryParams struct {
	ModelID        *uuid.UUID
	OrganizationID *uuid.UUID
	CreatedBy      *uuid.UUID
	Shared         bool
	Search         string
	Page           int
	PageSize       int
	SortBy         string
	SortDesc       bool
} 
