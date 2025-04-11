package config

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/agent-platform/server/models"
	"gorm.io/gorm"
)

var (
	ErrConfigNotFound = errors.New("配置不存在")
	ErrNoPermission   = errors.New("没有操作权限")
)

// Service 提供配置管理功能
type Service struct {
	db *gorm.DB
}

// NewService 创建新的配置服务
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// UpsertConfigRequest 包含创建或更新配置所需的数据
type UpsertConfigRequest struct {
	Key     string `json:"key" binding:"required"`
	Value   string `json:"value"`
	Scope   string `json:"scope" binding:"required,oneof=system user project application"`
	ScopeID string `json:"scope_id,omitempty"`
}

// UpsertConfig 创建或更新配置
func (s *Service) UpsertConfig(req UpsertConfigRequest, userID uuid.UUID) (*models.Config, error) {
	var scopeID *uuid.UUID
	if req.ScopeID != "" {
		id, err := uuid.Parse(req.ScopeID)
		if err != nil {
			return nil, errors.New("无效的作用域ID")
		}
		scopeID = &id

		// 检查有权限操作的情况
		if models.ConfigScope(req.Scope) == models.ScopeUser && *scopeID != userID {
			return nil, ErrNoPermission
		}
		
		if models.ConfigScope(req.Scope) == models.ScopeProject {
			// 检查项目是否存在以及用户是否有权限
			var project models.Project
			if err := s.db.Where("id = ?", *scopeID).First(&project).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("项目不存在")
				}
				return nil, fmt.Errorf("获取项目失败: %w", err)
			}
			
			if project.OwnerID != userID {
				return nil, ErrNoPermission
			}
		}
		
		if models.ConfigScope(req.Scope) == models.ScopeApplication {
			// 检查应用是否存在以及用户是否有权限
			var application models.Application
			if err := s.db.Preload("Project").Where("id = ?", *scopeID).First(&application).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, errors.New("应用不存在")
				}
				return nil, fmt.Errorf("获取应用失败: %w", err)
			}
			
			if application.Project.OwnerID != userID {
				return nil, ErrNoPermission
			}
		}
	}

	// 查找是否已存在相同的配置
	var existingConfig models.Config
	query := s.db.Where("key = ? AND scope = ?", req.Key, req.Scope)
	if scopeID != nil {
		query = query.Where("scope_id = ?", *scopeID)
	} else {
		query = query.Where("scope_id IS NULL")
	}
	
	err := query.First(&existingConfig).Error
	
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("查询配置失败: %w", err)
		}
		
		// 创建新配置
		config := models.Config{
			Key:       req.Key,
			Value:     req.Value,
			Scope:     models.ConfigScope(req.Scope),
			ScopeID:   scopeID,
			CreatedBy: &userID,
			UpdatedBy: &userID,
		}
		
		if err := s.db.Create(&config).Error; err != nil {
			return nil, fmt.Errorf("创建配置失败: %w", err)
		}
		
		return &config, nil
	}
	
	// 更新现有配置
	existingConfig.Value = req.Value
	existingConfig.UpdatedBy = &userID
	
	if err := s.db.Save(&existingConfig).Error; err != nil {
		return nil, fmt.Errorf("更新配置失败: %w", err)
	}
	
	return &existingConfig, nil
}

// GetConfigByID 根据ID获取配置
func (s *Service) GetConfigByID(id uuid.UUID) (*models.Config, error) {
	var config models.Config
	if err := s.db.Where("id = ?", id).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConfigNotFound
		}
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}
	
	return &config, nil
}

// GetConfigsByScope 获取指定作用域的配置列表
func (s *Service) GetConfigsByScope(scope string, scopeID *uuid.UUID) ([]models.ConfigResponse, error) {
	var configs []models.Config
	query := s.db.Where("scope = ?", scope)
	
	if scopeID != nil {
		query = query.Where("scope_id = ?", *scopeID)
	} else {
		query = query.Where("scope_id IS NULL")
	}
	
	if err := query.Order("key ASC").Find(&configs).Error; err != nil {
		return nil, fmt.Errorf("获取配置列表失败: %w", err)
	}
	
	// 转换为响应格式
	responses := make([]models.ConfigResponse, len(configs))
	for i, config := range configs {
		responses[i] = config.ToResponse()
	}
	
	return responses, nil
}

// GetSystemConfigs 获取系统配置
func (s *Service) GetSystemConfigs() ([]models.ConfigResponse, error) {
	return s.GetConfigsByScope(string(models.ScopeSystem), nil)
}

// GetUserConfigs 获取用户配置
func (s *Service) GetUserConfigs(userID uuid.UUID) ([]models.ConfigResponse, error) {
	return s.GetConfigsByScope(string(models.ScopeUser), &userID)
}

// GetProjectConfigs 获取项目配置
func (s *Service) GetProjectConfigs(projectID uuid.UUID) ([]models.ConfigResponse, error) {
	return s.GetConfigsByScope(string(models.ScopeProject), &projectID)
}

// GetApplicationConfigs 获取应用配置
func (s *Service) GetApplicationConfigs(applicationID uuid.UUID) ([]models.ConfigResponse, error) {
	return s.GetConfigsByScope(string(models.ScopeApplication), &applicationID)
}

// DeleteConfig 删除配置
func (s *Service) DeleteConfig(id uuid.UUID, userID uuid.UUID) error {
	var config models.Config
	if err := s.db.Where("id = ?", id).First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConfigNotFound
		}
		return fmt.Errorf("获取配置失败: %w", err)
	}
	
	// 检查权限
	if config.Scope == models.ScopeSystem {
		// 系统配置，检查是否为管理员
		var user models.User
		if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
			return fmt.Errorf("获取用户失败: %w", err)
		}
		
		if user.Role != "admin" {
			return ErrNoPermission
		}
	} else if config.Scope == models.ScopeUser {
		// 用户配置，只能删除自己的
		if config.ScopeID == nil || *config.ScopeID != userID {
			return ErrNoPermission
		}
	} else if config.Scope == models.ScopeProject && config.ScopeID != nil {
		// 项目配置，检查是否为项目所有者
		var project models.Project
		if err := s.db.Where("id = ?", *config.ScopeID).First(&project).Error; err != nil {
			return fmt.Errorf("获取项目失败: %w", err)
		}
		
		if project.OwnerID != userID {
			return ErrNoPermission
		}
	} else if config.Scope == models.ScopeApplication && config.ScopeID != nil {
		// 应用配置，检查是否有权限
		var application models.Application
		if err := s.db.Preload("Project").Where("id = ?", *config.ScopeID).First(&application).Error; err != nil {
			return fmt.Errorf("获取应用失败: %w", err)
		}
		
		if application.Project.OwnerID != userID {
			return ErrNoPermission
		}
	}
	
	// 删除配置
	if err := s.db.Delete(&config).Error; err != nil {
		return fmt.Errorf("删除配置失败: %w", err)
	}
	
	return nil
} 