package application

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"gorm.io/gorm"
)

var (
	ErrApplicationNotFound = errors.New("应用不存在")
	ErrProjectNotFound     = errors.New("项目不存在")
	ErrNoPermission        = errors.New("没有操作权限")
)

// Service 提供应用管理功能
type Service struct {
	db *gorm.DB
}

// NewService 创建新的应用服务
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// CreateApplicationRequest 包含创建应用所需的数据
type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required,oneof=chat workflow custom"`
	ProjectID   string `json:"project_id" binding:"required,uuid"`
	Config      string `json:"config"`
	ModelConfig string `json:"model_config"`
}

// CreateApplication 创建新应用
func (s *Service) CreateApplication(req CreateApplicationRequest, userID uuid.UUID) (*models.Application, error) {
	projectID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, errors.New("无效的项目ID")
	}

	// 检查项目是否存在以及用户是否有权限操作
	var project models.Project
	if err := s.db.Where("id = ?", projectID).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("获取项目失败: %w", err)
	}

	// 只有项目所有者可以创建应用
	if project.OwnerID != userID {
		return nil, ErrNoPermission
	}

	application := models.Application{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		ProjectID:   projectID,
		Status:      "draft",
		Config:      req.Config,
		ModelConfig: req.ModelConfig,
	}

	if err := s.db.Create(&application).Error; err != nil {
		return nil, fmt.Errorf("创建应用失败: %w", err)
	}

	return &application, nil
}

// GetApplicationByID 根据ID获取应用
func (s *Service) GetApplicationByID(id uuid.UUID, userID uuid.UUID) (*models.Application, error) {
	var application models.Application
	
	// 加载项目信息以检查权限
	if err := s.db.Preload("Project").Where("id = ?", id).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("获取应用失败: %w", err)
	}

	// 检查用户是否有权限访问此应用
	if application.Project.OwnerID != userID && !application.Project.Public {
		return nil, ErrNoPermission
	}

	return &application, nil
}

// UpdateApplicationRequest 包含更新应用所需的数据
type UpdateApplicationRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"omitempty,oneof=chat workflow custom"`
	Status      string `json:"status" binding:"omitempty,oneof=draft published archived"`
	Config      string `json:"config"`
	ModelConfig string `json:"model_config"`
}

// UpdateApplication 更新应用信息
func (s *Service) UpdateApplication(id uuid.UUID, req UpdateApplicationRequest, userID uuid.UUID) (*models.Application, error) {
	var application models.Application
	
	// 加载项目信息以检查权限
	if err := s.db.Preload("Project").Where("id = ?", id).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("获取应用失败: %w", err)
	}

	// 检查用户是否有权限更新此应用
	if application.Project.OwnerID != userID {
		return nil, ErrNoPermission
	}

	// 更新应用字段
	updates := map[string]interface{}{}
	
	if req.Name != "" {
		updates["name"] = req.Name
	}
	
	if req.Description != "" || req.Description == "" {
		updates["description"] = req.Description
	}
	
	if req.Type != "" {
		updates["type"] = req.Type
	}
	
	if req.Status != "" {
		updates["status"] = req.Status
	}
	
	if req.Config != "" {
		updates["config"] = req.Config
	}
	
	if req.ModelConfig != "" {
		updates["model_config"] = req.ModelConfig
	}
	
	if len(updates) > 0 {
		if err := s.db.Model(&application).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("更新应用失败: %w", err)
		}
		
		// 重新获取完整的应用信息
		if err := s.db.Where("id = ?", id).First(&application).Error; err != nil {
			return nil, fmt.Errorf("获取更新后的应用失败: %w", err)
		}
	}
	
	return &application, nil
}

// DeleteApplication 删除应用
func (s *Service) DeleteApplication(id uuid.UUID, userID uuid.UUID) error {
	var application models.Application
	
	// 加载项目信息以检查权限
	if err := s.db.Preload("Project").Where("id = ?", id).First(&application).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrApplicationNotFound
		}
		return fmt.Errorf("获取应用失败: %w", err)
	}

	// 检查用户是否有权限删除此应用
	if application.Project.OwnerID != userID {
		return ErrNoPermission
	}

	// 软删除应用
	if err := s.db.Delete(&application).Error; err != nil {
		return fmt.Errorf("删除应用失败: %w", err)
	}

	return nil
}

// GetApplicationsByProject 获取项目下的应用列表
func (s *Service) GetApplicationsByProject(projectID uuid.UUID, userID uuid.UUID) ([]models.ApplicationResponse, error) {
	// 检查项目是否存在以及用户是否有权限
	var project models.Project
	if err := s.db.Where("id = ?", projectID).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("获取项目失败: %w", err)
	}

	// 检查用户是否有权限访问此项目
	if project.OwnerID != userID && !project.Public {
		return nil, ErrNoPermission
	}

	// 获取项目下的应用列表
	var applications []models.Application
	if err := s.db.Where("project_id = ?", projectID).Order("created_at DESC").Find(&applications).Error; err != nil {
		return nil, fmt.Errorf("获取应用列表失败: %w", err)
	}
	
	// 转换为响应格式
	responses := make([]models.ApplicationResponse, len(applications))
	for i, application := range applications {
		responses[i] = application.ToResponse()
	}
	
	return responses, nil
} 
