package agent

import (
	"errors"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrAgentNotFound      = errors.New("智能体不存在")
	ErrApplicationNotFound = errors.New("应用不存在")
	ErrModelConfigNotFound = errors.New("模型配置不存在")
	ErrKnowledgeBaseNotFound = errors.New("知识库不存在")
	ErrUnauthorized       = errors.New("无权访问此资源")
)

// Service 提供智能体相关功能
type Service struct {
	db *gorm.DB
	logger *zap.Logger
}

// NewService 创建新的智能体服务
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
		logger: zap.L().With(zap.String("service", "agent")),
	}
}

// GetAgentsByApplicationID 获取应用下的所有智能体
func (s *Service) GetAgentsByApplicationID(appID uuid.UUID, userID uuid.UUID) ([]models.AgentResponse, error) {
	// 首先检查用户是否有权限访问此应用
	var application models.Application
	if err := s.db.First(&application, "id = ?", appID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		s.logger.Error("Failed to find application", zap.Error(err))
		return nil, err
	}

	// 检查用户权限（简化版，实际应该使用权限系统）
	if application.CreatedBy != userID {
		// 这里需要集成更完善的权限检查
		return nil, ErrUnauthorized
	}

	// 获取应用下的所有智能体
	var agents []models.Agent
	if err := s.db.Where("application_id = ?", appID).Find(&agents).Error; err != nil {
		s.logger.Error("Failed to find agents", zap.Error(err), zap.String("app_id", appID.String()))
		return nil, err
	}

	// 转换为响应格式
	agentResponses := make([]models.AgentResponse, len(agents))
	for i, agent := range agents {
		// 获取智能体关联的知识库
		var knowledgeBaseIDs []uuid.UUID
		if err := s.db.Model(&models.AgentKnowledgeBase{}).
			Where("agent_id = ?", agent.ID).
			Pluck("knowledge_base_id", &knowledgeBaseIDs).Error; err != nil {
			s.logger.Error("Failed to get knowledge base IDs", zap.Error(err))
		}

		agentResponses[i] = agent.ToResponse(knowledgeBaseIDs)
	}

	return agentResponses, nil
}

// GetAgentByID 通过ID获取智能体
func (s *Service) GetAgentByID(id uuid.UUID, userID uuid.UUID) (*models.AgentResponse, error) {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAgentNotFound
		}
		s.logger.Error("Failed to find agent", zap.Error(err), zap.String("agent_id", id.String()))
		return nil, err
	}

	// 获取应用信息以检查权限
	var application models.Application
	if err := s.db.First(&application, "id = ?", agent.ApplicationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		s.logger.Error("Failed to find application", zap.Error(err))
		return nil, err
	}

	// 检查用户权限（简化版，实际应该使用权限系统）
	if application.CreatedBy != userID {
		// 这里需要集成更完善的权限检查
		return nil, ErrUnauthorized
	}

	// 获取智能体关联的知识库
	var knowledgeBaseIDs []uuid.UUID
	if err := s.db.Model(&models.AgentKnowledgeBase{}).
		Where("agent_id = ?", agent.ID).
		Pluck("knowledge_base_id", &knowledgeBaseIDs).Error; err != nil {
		s.logger.Error("Failed to get knowledge base IDs", zap.Error(err))
	}

	response := agent.ToResponse(knowledgeBaseIDs)
	return &response, nil
}

// CreateAgent 创建新的智能体
func (s *Service) CreateAgent(appID uuid.UUID, req models.CreateAgentRequest, userID uuid.UUID) (*models.AgentResponse, error) {
	// 检查应用是否存在及用户权限
	var application models.Application
	if err := s.db.First(&application, "id = ?", appID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApplicationNotFound
		}
		s.logger.Error("Failed to find application", zap.Error(err))
		return nil, err
	}

	// 检查用户权限
	if application.CreatedBy != userID {
		return nil, ErrUnauthorized
	}

	// 检查模型配置是否存在
	var modelConfig models.ModelConfig
	if err := s.db.First(&modelConfig, "id = ?", req.ModelConfigID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrModelConfigNotFound
		}
		s.logger.Error("Failed to find model config", zap.Error(err))
		return nil, err
	}

	// 创建智能体
	agent := models.Agent{
		Name:             req.Name,
		Description:      req.Description,
		ApplicationID:    appID,
		ModelConfigID:    req.ModelConfigID,
		SystemPrompt:     req.SystemPrompt,
		Tools:            req.Tools,
		Variables:        req.Variables,
		MaxHistoryLength: req.MaxHistoryLength,
	}

	// 默认值处理
	if agent.MaxHistoryLength == 0 {
		agent.MaxHistoryLength = 10
	}

	// 开启事务
	tx := s.db.Begin()
	if err := tx.Create(&agent).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create agent", zap.Error(err))
		return nil, err
	}

	// 处理关联的知识库
	if len(req.KnowledgeBaseIDs) > 0 {
		for _, kbID := range req.KnowledgeBaseIDs {
			// 验证知识库存在
			var kb models.KnowledgeBase
			if err := tx.First(&kb, "id = ?", kbID).Error; err != nil {
				tx.Rollback()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrKnowledgeBaseNotFound
				}
				s.logger.Error("Failed to find knowledge base", zap.Error(err))
				return nil, err
			}

			// 创建关联
			agentKb := models.AgentKnowledgeBase{
				AgentID:         agent.ID,
				KnowledgeBaseID: kbID,
			}
			if err := tx.Create(&agentKb).Error; err != nil {
				tx.Rollback()
				s.logger.Error("Failed to create agent-kb relation", zap.Error(err))
				return nil, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	response := agent.ToResponse(req.KnowledgeBaseIDs)
	return &response, nil
}

// UpdateAgent 更新智能体
func (s *Service) UpdateAgent(id uuid.UUID, req models.UpdateAgentRequest, userID uuid.UUID) (*models.AgentResponse, error) {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAgentNotFound
		}
		s.logger.Error("Failed to find agent", zap.Error(err))
		return nil, err
	}

	// 检查用户权限
	var application models.Application
	if err := s.db.First(&application, "id = ?", agent.ApplicationID).Error; err != nil {
		s.logger.Error("Failed to find application", zap.Error(err))
		return nil, err
	}

	if application.CreatedBy != userID {
		return nil, ErrUnauthorized
	}

	// 更新字段
	tx := s.db.Begin()

	// 只更新提供的字段
	updates := make(map[string]interface{})
	
	if req.Name != "" {
		updates["name"] = req.Name
	}
	
	if req.Description != "" {
		updates["description"] = req.Description
	}
	
	if req.ModelConfigID != uuid.Nil {
		// 验证模型配置
		var modelConfig models.ModelConfig
		if err := tx.First(&modelConfig, "id = ?", req.ModelConfigID).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrModelConfigNotFound
			}
			s.logger.Error("Failed to find model config", zap.Error(err))
			return nil, err
		}
		updates["model_config_id"] = req.ModelConfigID
	}
	
	if req.SystemPrompt != "" {
		updates["system_prompt"] = req.SystemPrompt
	}
	
	if req.Tools != nil {
		updates["tools"] = req.Tools
	}
	
	if req.Variables != nil {
		updates["variables"] = req.Variables
	}
	
	if req.MaxHistoryLength > 0 {
		updates["max_history_length"] = req.MaxHistoryLength
	}

	// 更新智能体
	if len(updates) > 0 {
		if err := tx.Model(&agent).Updates(updates).Error; err != nil {
			tx.Rollback()
			s.logger.Error("Failed to update agent", zap.Error(err))
			return nil, err
		}
	}

	// 更新知识库关联
	if req.KnowledgeBaseIDs != nil {
		// 删除现有关联
		if err := tx.Delete(&models.AgentKnowledgeBase{}, "agent_id = ?", agent.ID).Error; err != nil {
			tx.Rollback()
			s.logger.Error("Failed to delete agent-kb relations", zap.Error(err))
			return nil, err
		}

		// 添加新关联
		for _, kbID := range req.KnowledgeBaseIDs {
			// 验证知识库存在
			var kb models.KnowledgeBase
			if err := tx.First(&kb, "id = ?", kbID).Error; err != nil {
				tx.Rollback()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrKnowledgeBaseNotFound
				}
				s.logger.Error("Failed to find knowledge base", zap.Error(err))
				return nil, err
			}

			// 创建关联
			agentKb := models.AgentKnowledgeBase{
				AgentID:         agent.ID,
				KnowledgeBaseID: kbID,
			}
			if err := tx.Create(&agentKb).Error; err != nil {
				tx.Rollback()
				s.logger.Error("Failed to create agent-kb relation", zap.Error(err))
				return nil, err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	// 获取更新后的智能体
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		s.logger.Error("Failed to reload agent", zap.Error(err))
		return nil, err
	}

	// 获取关联的知识库
	var knowledgeBaseIDs []uuid.UUID
	if err := s.db.Model(&models.AgentKnowledgeBase{}).
		Where("agent_id = ?", agent.ID).
		Pluck("knowledge_base_id", &knowledgeBaseIDs).Error; err != nil {
		s.logger.Error("Failed to get knowledge base IDs", zap.Error(err))
	}

	response := agent.ToResponse(knowledgeBaseIDs)
	return &response, nil
}

// UpdateSystemPrompt 更新智能体系统提示词
func (s *Service) UpdateSystemPrompt(id uuid.UUID, req models.UpdateSystemPromptRequest, userID uuid.UUID) error {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAgentNotFound
		}
		s.logger.Error("Failed to find agent", zap.Error(err))
		return err
	}

	// 检查用户权限
	var application models.Application
	if err := s.db.First(&application, "id = ?", agent.ApplicationID).Error; err != nil {
		s.logger.Error("Failed to find application", zap.Error(err))
		return err
	}

	if application.CreatedBy != userID {
		return ErrUnauthorized
	}

	// 更新系统提示词
	if err := s.db.Model(&agent).Update("system_prompt", req.SystemPrompt).Error; err != nil {
		s.logger.Error("Failed to update system prompt", zap.Error(err))
		return err
	}

	return nil
}

// UpdateTools 更新智能体工具配置
func (s *Service) UpdateTools(id uuid.UUID, req models.UpdateToolsRequest, userID uuid.UUID) error {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAgentNotFound
		}
		s.logger.Error("Failed to find agent", zap.Error(err))
		return err
	}

	// 检查用户权限
	var application models.Application
	if err := s.db.First(&application, "id = ?", agent.ApplicationID).Error; err != nil {
		s.logger.Error("Failed to find application", zap.Error(err))
		return err
	}

	if application.CreatedBy != userID {
		return ErrUnauthorized
	}

	// 更新工具配置
	if err := s.db.Model(&agent).Update("tools", req.Tools).Error; err != nil {
		s.logger.Error("Failed to update tools", zap.Error(err))
		return err
	}

	return nil
}

// DeleteAgent 删除智能体
func (s *Service) DeleteAgent(id uuid.UUID, userID uuid.UUID) error {
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrAgentNotFound
		}
		s.logger.Error("Failed to find agent", zap.Error(err))
		return err
	}

	// 检查用户权限
	var application models.Application
	if err := s.db.First(&application, "id = ?", agent.ApplicationID).Error; err != nil {
		s.logger.Error("Failed to find application", zap.Error(err))
		return err
	}

	if application.CreatedBy != userID {
		return ErrUnauthorized
	}

	// 开启事务
	tx := s.db.Begin()

	// 删除智能体知识库关联
	if err := tx.Delete(&models.AgentKnowledgeBase{}, "agent_id = ?", id).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to delete agent-kb relations", zap.Error(err))
		return err
	}

	// 删除智能体（软删除）
	if err := tx.Delete(&agent).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to delete agent", zap.Error(err))
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
} 
