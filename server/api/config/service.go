package config

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
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

// UpdateSystemSettingsRequest 系统设置更新请求
type UpdateSystemSettingsRequest struct {
	SiteName          string `json:"siteName"`
	LogoURL           string `json:"logoUrl"`
	APIRateLimit      int    `json:"apiRateLimit"`
	AllowRegistration bool   `json:"allowRegistration"`
	DefaultLanguage   string `json:"defaultLanguage"`
	DefaultModel      string `json:"defaultModel"`
	StorageProvider   string `json:"storageProvider"`
	S3Config          struct {
		Bucket    string `json:"bucket"`
		Region    string `json:"region"`
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
	} `json:"s3Config"`
	EmailSettings struct {
		SMTPServer   string `json:"smtpServer"`
		SMTPPort     int    `json:"smtpPort"`
		SMTPUser     string `json:"smtpUser"`
		SMTPPassword string `json:"smtpPassword"`
		SenderEmail  string `json:"senderEmail"`
	} `json:"emailSettings"`
}

// SystemSettings 系统设置响应
type SystemSettings struct {
	SiteName          string `json:"siteName"`
	LogoURL           string `json:"logoUrl"`
	APIRateLimit      int    `json:"apiRateLimit"`
	AllowRegistration bool   `json:"allowRegistration"`
	DefaultLanguage   string `json:"defaultLanguage"`
	DefaultModel      string `json:"defaultModel"`
	StorageProvider   string `json:"storageProvider"`
	S3Config          struct {
		Bucket    string `json:"bucket"`
		Region    string `json:"region"`
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
	} `json:"s3Config"`
	EmailSettings struct {
		SMTPServer   string `json:"smtpServer"`
		SMTPPort     int    `json:"smtpPort"`
		SMTPUser     string `json:"smtpUser"`
		SMTPPassword string `json:"smtpPassword,omitempty"`
		SenderEmail  string `json:"senderEmail"`
	} `json:"emailSettings"`
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

// GetSystemSettings 获取系统设置
func (s *Service) GetSystemSettings() (*SystemSettings, error) {
	// 从系统配置中获取设置
	configs, err := s.GetSystemConfigs()
	if err != nil {
		return nil, err
	}
	
	// 构造系统设置对象
	settings := &SystemSettings{
		SiteName:          "智能体构建平台",
		LogoURL:           "/logo.png",
		APIRateLimit:      100,
		AllowRegistration: true,
		DefaultLanguage:   "zh-CN",
		DefaultModel:      "",
		StorageProvider:   "local",
	}
	
	// 填充从配置中获取的值
	for _, config := range configs {
		switch config.Key {
		case "site.name":
			settings.SiteName = config.Value
		case "site.logo_url":
			settings.LogoURL = config.Value
		case "api.rate_limit":
			if limit, err := strconv.Atoi(config.Value); err == nil {
				settings.APIRateLimit = limit
			}
		case "user.allow_registration":
			settings.AllowRegistration = config.Value == "true"
		case "site.default_language":
			settings.DefaultLanguage = config.Value
		case "model.default":
			settings.DefaultModel = config.Value
		case "storage.provider":
			settings.StorageProvider = config.Value
		case "storage.s3.bucket":
			settings.S3Config.Bucket = config.Value
		case "storage.s3.region":
			settings.S3Config.Region = config.Value
		case "storage.s3.access_key":
			settings.S3Config.AccessKey = config.Value
		case "email.smtp_server":
			settings.EmailSettings.SMTPServer = config.Value
		case "email.smtp_port":
			if port, err := strconv.Atoi(config.Value); err == nil {
				settings.EmailSettings.SMTPPort = port
			}
		case "email.smtp_user":
			settings.EmailSettings.SMTPUser = config.Value
		case "email.sender":
			settings.EmailSettings.SenderEmail = config.Value
		}
	}
	
	return settings, nil
}

// UpdateSystemSettings 更新系统设置
func (s *Service) UpdateSystemSettings(req UpdateSystemSettingsRequest) error {
	// 更新系统设置
	tx := s.db.Begin()
	
	// 更新基本设置
	if req.SiteName != "" {
		if err := s.setConfig(tx, "site.name", req.SiteName, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.LogoURL != "" {
		if err := s.setConfig(tx, "site.logo_url", req.LogoURL, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.APIRateLimit > 0 {
		if err := s.setConfig(tx, "api.rate_limit", strconv.Itoa(req.APIRateLimit), models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if err := s.setConfig(tx, "user.allow_registration", strconv.FormatBool(req.AllowRegistration), models.ScopeSystem, nil); err != nil {
		tx.Rollback()
		return err
	}
	
	if req.DefaultLanguage != "" {
		if err := s.setConfig(tx, "site.default_language", req.DefaultLanguage, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.DefaultModel != "" {
		if err := s.setConfig(tx, "model.default", req.DefaultModel, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.StorageProvider != "" {
		if err := s.setConfig(tx, "storage.provider", req.StorageProvider, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	// 更新S3配置
	if req.StorageProvider == "s3" {
		if req.S3Config.Bucket != "" {
			if err := s.setConfig(tx, "storage.s3.bucket", req.S3Config.Bucket, models.ScopeSystem, nil); err != nil {
				tx.Rollback()
				return err
			}
		}
		
		if req.S3Config.Region != "" {
			if err := s.setConfig(tx, "storage.s3.region", req.S3Config.Region, models.ScopeSystem, nil); err != nil {
				tx.Rollback()
				return err
			}
		}
		
		if req.S3Config.AccessKey != "" {
			if err := s.setConfig(tx, "storage.s3.access_key", req.S3Config.AccessKey, models.ScopeSystem, nil); err != nil {
				tx.Rollback()
				return err
			}
		}
		
		if req.S3Config.SecretKey != "" {
			if err := s.setConfig(tx, "storage.s3.secret_key", req.S3Config.SecretKey, models.ScopeSystem, nil); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	
	// 更新邮件设置
	if req.EmailSettings.SMTPServer != "" {
		if err := s.setConfig(tx, "email.smtp_server", req.EmailSettings.SMTPServer, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.EmailSettings.SMTPPort > 0 {
		if err := s.setConfig(tx, "email.smtp_port", strconv.Itoa(req.EmailSettings.SMTPPort), models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.EmailSettings.SMTPUser != "" {
		if err := s.setConfig(tx, "email.smtp_user", req.EmailSettings.SMTPUser, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.EmailSettings.SMTPPassword != "" {
		if err := s.setConfig(tx, "email.smtp_password", req.EmailSettings.SMTPPassword, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	if req.EmailSettings.SenderEmail != "" {
		if err := s.setConfig(tx, "email.sender", req.EmailSettings.SenderEmail, models.ScopeSystem, nil); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit().Error
}

// setConfig 设置配置项（内部方法，用于批量更新）
func (s *Service) setConfig(tx *gorm.DB, key, value string, scope models.ConfigScope, scopeID *uuid.UUID) error {
	var config models.Config
	
	query := tx.Where("key = ? AND scope = ?", key, scope)
	if scopeID != nil {
		query = query.Where("scope_id = ?", *scopeID)
	} else {
		query = query.Where("scope_id IS NULL")
	}
	
	if err := query.First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 创建新配置
			newConfig := models.Config{
				Key:     key,
				Value:   value,
				Scope:   scope,
				ScopeID: scopeID,
			}
			
			return tx.Create(&newConfig).Error
		}
		return err
	}
	
	// 更新现有配置
	config.Value = value
	return tx.Save(&config).Error
} 
