package conversation

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/agent-platform/server/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	ErrConversationNotFound = errors.New("对话不存在")
	ErrMessageNotFound      = errors.New("消息不存在")
	ErrAgentNotFound        = errors.New("智能体不存在")
	ErrUnauthorized         = errors.New("无权访问此资源")
)

// Service 提供对话相关功能
type Service struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewService 创建新的对话服务
func NewService(db *gorm.DB) *Service {
	return &Service{
		db:     db,
		logger: zap.L().With(zap.String("service", "conversation")),
	}
}

// GetConversationsByAgentID 获取智能体的所有对话
func (s *Service) GetConversationsByAgentID(agentID uuid.UUID, userID uuid.UUID) ([]models.ConversationResponse, error) {
	// 检查智能体是否存在
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", agentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAgentNotFound
		}
		s.logger.Error("Failed to find agent", zap.Error(err))
		return nil, err
	}

	// 获取智能体所有对话
	var conversations []models.Conversation
	if err := s.db.Preload("Agent").Where("agent_id = ? AND user_id = ?", agentID, userID).Find(&conversations).Error; err != nil {
		s.logger.Error("Failed to find conversations", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	conversationResponses := make([]models.ConversationResponse, len(conversations))
	for i, conv := range conversations {
		conversationResponses[i] = conv.ToResponse()
	}

	return conversationResponses, nil
}

// GetConversationByID 通过ID获取对话
func (s *Service) GetConversationByID(id uuid.UUID, userID uuid.UUID) (*models.ConversationResponse, error) {
	var conversation models.Conversation
	if err := s.db.Preload("Agent").First(&conversation, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		s.logger.Error("Failed to find conversation", zap.Error(err))
		return nil, err
	}

	// 检查用户权限
	if conversation.UserID != userID {
		return nil, ErrUnauthorized
	}

	response := conversation.ToResponse()
	return &response, nil
}

// CreateConversation 创建新的对话
func (s *Service) CreateConversation(req models.CreateConversationRequest, userID uuid.UUID) (*models.ConversationResponse, error) {
	// 检查智能体是否存在
	var agent models.Agent
	if err := s.db.First(&agent, "id = ?", req.AgentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAgentNotFound
		}
		s.logger.Error("Failed to find agent", zap.Error(err))
		return nil, err
	}

	// 创建对话
	conversation := models.Conversation{
		Title:    req.Title,
		AgentID:  req.AgentID,
		UserID:   userID,
		Status:   "active",
		Metadata: req.Metadata,
	}

	// 如果未提供标题，使用时间戳生成一个
	if conversation.Title == "" {
		conversation.Title = "对话 " + time.Now().Format("2006-01-02 15:04:05")
	}

	if err := s.db.Create(&conversation).Error; err != nil {
		s.logger.Error("Failed to create conversation", zap.Error(err))
		return nil, err
	}

	// 重新加载对话以获取关联
	if err := s.db.Preload("Agent").First(&conversation, "id = ?", conversation.ID).Error; err != nil {
		s.logger.Error("Failed to reload conversation", zap.Error(err))
		return nil, err
	}

	response := conversation.ToResponse()
	return &response, nil
}

// GetMessagesByConversationID 获取对话的所有消息
func (s *Service) GetMessagesByConversationID(conversationID uuid.UUID, userID uuid.UUID) ([]models.MessageResponse, error) {
	// 检查对话是否存在及用户权限
	var conversation models.Conversation
	if err := s.db.First(&conversation, "id = ?", conversationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		s.logger.Error("Failed to find conversation", zap.Error(err))
		return nil, err
	}

	// 检查用户权限
	if conversation.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 获取对话的所有消息
	var messages []models.Message
	if err := s.db.Order("created_at asc").Where("conversation_id = ?", conversationID).Find(&messages).Error; err != nil {
		s.logger.Error("Failed to find messages", zap.Error(err))
		return nil, err
	}

	// 转换为响应格式
	messageResponses := make([]models.MessageResponse, len(messages))
	for i, msg := range messages {
		messageResponses[i] = msg.ToResponse()
	}

	return messageResponses, nil
}

// SendMessage 发送消息到对话
func (s *Service) SendMessage(conversationID uuid.UUID, req models.SendMessageRequest, userID uuid.UUID) (*models.MessageResponse, error) {
	// 检查对话是否存在及用户权限
	var conversation models.Conversation
	if err := s.db.Preload("Agent").First(&conversation, "id = ?", conversationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		s.logger.Error("Failed to find conversation", zap.Error(err))
		return nil, err
	}

	// 检查用户权限
	if conversation.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 创建用户消息
	userMessage := models.Message{
		ConversationID: conversationID,
		Role:           "user",
		Content:        req.Content,
		Metadata:       models.JSONMap{},
	}

	// 开启事务
	tx := s.db.Begin()

	if err := tx.Create(&userMessage).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create user message", zap.Error(err))
		return nil, err
	}

	// 生成AI响应（示例，实际需要调用LLM服务）
	// TODO: 当Agent运行时实现后，这里会调用Agent运行时
	aiResponse := models.Message{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        "这是一个模拟的AI响应。实际功能将在Agent运行时实现后完成。",
		Metadata:       models.JSONMap{},
	}

	if err := tx.Create(&aiResponse).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to create AI response", zap.Error(err))
		return nil, err
	}

	// 更新对话的更新时间
	if err := tx.Model(&conversation).Update("updated_at", time.Now()).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to update conversation time", zap.Error(err))
		return nil, err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	response := aiResponse.ToResponse()
	return &response, nil
}

// RegenerateResponse 重新生成AI回复
func (s *Service) RegenerateResponse(conversationID uuid.UUID, userID uuid.UUID) (*models.MessageResponse, error) {
	// 检查对话是否存在及用户权限
	var conversation models.Conversation
	if err := s.db.Preload("Agent").First(&conversation, "id = ?", conversationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrConversationNotFound
		}
		s.logger.Error("Failed to find conversation", zap.Error(err))
		return nil, err
	}

	// 检查用户权限
	if conversation.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 获取最后一条用户消息
	var lastUserMessage models.Message
	if err := s.db.Where("conversation_id = ? AND role = ?", conversationID, "user").Order("created_at desc").First(&lastUserMessage).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("没有找到用户消息")
		}
		s.logger.Error("Failed to find last user message", zap.Error(err))
		return nil, err
	}

	// 获取最后一条AI回复
	var lastAIMessage models.Message
	if err := s.db.Where("conversation_id = ? AND role = ?", conversationID, "assistant").Order("created_at desc").First(&lastAIMessage).Error; err != nil {
		// 如果没有找到AI消息，可能是第一次对话
		s.logger.Warn("No AI message found, creating new one")
	} else {
		// 删除最后一条AI回复
		if err := s.db.Delete(&lastAIMessage).Error; err != nil {
			s.logger.Error("Failed to delete last AI message", zap.Error(err))
			return nil, err
		}
	}

	// 生成新的AI响应（示例，实际需要调用LLM服务）
	// TODO: 当Agent运行时实现后，这里会调用Agent运行时
	newAIResponse := models.Message{
		ConversationID: conversationID,
		Role:           "assistant",
		Content:        "这是一个重新生成的模拟AI响应。实际功能将在Agent运行时实现后完成。",
		Metadata:       models.JSONMap{},
	}

	if err := s.db.Create(&newAIResponse).Error; err != nil {
		s.logger.Error("Failed to create new AI response", zap.Error(err))
		return nil, err
	}

	// 更新对话的更新时间
	if err := s.db.Model(&conversation).Update("updated_at", time.Now()).Error; err != nil {
		s.logger.Error("Failed to update conversation time", zap.Error(err))
		return nil, err
	}

	response := newAIResponse.ToResponse()
	return &response, nil
}

// ProvideFeedback 提供消息反馈
func (s *Service) ProvideFeedback(messageID uuid.UUID, req models.MessageFeedbackRequest, userID uuid.UUID) error {
	// 查找消息
	var message models.Message
	if err := s.db.Preload("Conversation").First(&message, "id = ?", messageID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMessageNotFound
		}
		s.logger.Error("Failed to find message", zap.Error(err))
		return err
	}

	// 检查用户权限
	if message.Conversation.UserID != userID {
		return ErrUnauthorized
	}

	// 更新反馈
	if err := s.db.Model(&message).Update("feedback", req.Feedback).Error; err != nil {
		s.logger.Error("Failed to update message feedback", zap.Error(err))
		return err
	}

	return nil
}

// DeleteConversation 删除对话
func (s *Service) DeleteConversation(id uuid.UUID, userID uuid.UUID) error {
	// 检查对话是否存在及用户权限
	var conversation models.Conversation
	if err := s.db.First(&conversation, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrConversationNotFound
		}
		s.logger.Error("Failed to find conversation", zap.Error(err))
		return err
	}

	// 检查用户权限
	if conversation.UserID != userID {
		return ErrUnauthorized
	}

	// 开启事务
	tx := s.db.Begin()

	// 删除对话的所有消息
	if err := tx.Where("conversation_id = ?", id).Delete(&models.Message{}).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to delete messages", zap.Error(err))
		return err
	}

	// 删除对话
	if err := tx.Delete(&conversation).Error; err != nil {
		tx.Rollback()
		s.logger.Error("Failed to delete conversation", zap.Error(err))
		return err
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		s.logger.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
} 