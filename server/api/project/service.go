package project

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"gorm.io/gorm"
)

var (
	ErrProjectNotFound = errors.New("项目不存在")
	ErrNoPermission    = errors.New("没有操作权限")
)

// Service 提供项目管理功能
type Service struct {
	db *gorm.DB
}

// NewService 创建新的项目服务
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// CreateProjectRequest 包含创建项目所需的数据
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
	Config      string `json:"config"`
}

// CreateProject 创建新项目
func (s *Service) CreateProject(req CreateProjectRequest, userID uuid.UUID) (*models.Project, error) {
	project := models.Project{
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userID,
		Public:      req.Public,
		Status:      "active",
		Config:      req.Config,
	}

	if err := s.db.Create(&project).Error; err != nil {
		return nil, fmt.Errorf("创建项目失败: %w", err)
	}

	return &project, nil
}

// GetProjectByID 根据ID获取项目
func (s *Service) GetProjectByID(id uuid.UUID, userID uuid.UUID) (*models.Project, error) {
	var project models.Project
	
	// 查询项目，同时检查访问权限（项目所有者或公开项目）
	err := s.db.Where("id = ? AND (owner_id = ? OR public = ?)", id, userID, true).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("获取项目失败: %w", err)
	}
	
	return &project, nil
}

// UpdateProjectRequest 包含更新项目所需的数据
type UpdateProjectRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description"`
	Public      *bool  `json:"public"`
	Status      string `json:"status" binding:"omitempty,oneof=active archived"`
	Config      string `json:"config"`
}

// UpdateProject 更新项目信息
func (s *Service) UpdateProject(id uuid.UUID, req UpdateProjectRequest, userID uuid.UUID) (*models.Project, error) {
	// 首先检查项目是否存在，以及用户是否有权限
	var project models.Project
	if err := s.db.Where("id = ?", id).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, fmt.Errorf("获取项目失败: %w", err)
	}

	// 只有项目所有者可以更新项目
	if project.OwnerID != userID {
		return nil, ErrNoPermission
	}

	// 更新项目字段
	updates := map[string]interface{}{}
	
	if req.Name != "" {
		updates["name"] = req.Name
	}
	
	if req.Description != "" || req.Description == "" {
		updates["description"] = req.Description
	}
	
	if req.Public != nil {
		updates["public"] = *req.Public
	}
	
	if req.Status != "" {
		updates["status"] = req.Status
	}
	
	if req.Config != "" {
		updates["config"] = req.Config
	}
	
	if len(updates) > 0 {
		if err := s.db.Model(&project).Updates(updates).Error; err != nil {
			return nil, fmt.Errorf("更新项目失败: %w", err)
		}
		
		// 重新获取完整的项目信息
		if err := s.db.Where("id = ?", id).First(&project).Error; err != nil {
			return nil, fmt.Errorf("获取更新后的项目失败: %w", err)
		}
	}
	
	return &project, nil
}

// DeleteProject 删除项目
func (s *Service) DeleteProject(id uuid.UUID, userID uuid.UUID) error {
	// 首先检查项目是否存在，以及用户是否有权限
	var project models.Project
	if err := s.db.Where("id = ?", id).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return fmt.Errorf("获取项目失败: %w", err)
	}

	// 只有项目所有者可以删除项目
	if project.OwnerID != userID {
		return ErrNoPermission
	}

	// 软删除项目
	if err := s.db.Delete(&project).Error; err != nil {
		return fmt.Errorf("删除项目失败: %w", err)
	}

	return nil
}

// GetProjects 获取用户的项目列表
func (s *Service) GetProjects(userID uuid.UUID, includeArchived bool) ([]models.ProjectResponse, error) {
	var projects []models.Project
	query := s.db.Where("owner_id = ?", userID)
	
	if !includeArchived {
		query = query.Where("status = ?", "active")
	}
	
	if err := query.Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("获取项目列表失败: %w", err)
	}
	
	// 转换为响应格式
	responses := make([]models.ProjectResponse, len(projects))
	for i, project := range projects {
		// 获取每个项目的应用数量
		var appsCount int64
		s.db.Model(&models.Application{}).Where("project_id = ?", project.ID).Count(&appsCount)
		
		resp := project.ToResponse()
		resp.AppsCount = appsCount
		responses[i] = resp
	}
	
	return responses, nil
}

// GetPublicProjects 获取公开项目列表
func (s *Service) GetPublicProjects() ([]models.ProjectResponse, error) {
	var projects []models.Project
	
	if err := s.db.Where("public = ? AND status = ?", true, "active").
		Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("获取公开项目列表失败: %w", err)
	}
	
	// 转换为响应格式
	responses := make([]models.ProjectResponse, len(projects))
	for i, project := range projects {
		// 获取每个项目的应用数量
		var appsCount int64
		s.db.Model(&models.Application{}).Where("project_id = ?", project.ID).Count(&appsCount)
		
		resp := project.ToResponse()
		resp.AppsCount = appsCount
		responses[i] = resp
	}
	
	return responses, nil
} 
