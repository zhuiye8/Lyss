package agent

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bytedance/eino"
	"github.com/google/uuid"
)

// Conversation 表示一个对话会话
type Conversation struct {
	ID           string      `json:"id"`
	AgentID      string      `json:"agent_id"`
	Title        string      `json:"title"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	Memory       Memory      `json:"-"`
	Agent        *Agent      `json:"-"`
	Metadata     interface{} `json:"metadata,omitempty"`
}

// Message 表示对话中的一条消息
type Message struct {
	ID             string                 `json:"id"`
	ConversationID string                 `json:"conversation_id"`
	Role           string                 `json:"role"` // user, assistant, system
	Content        string                 `json:"content"`
	ToolCalls      []interface{}          `json:"tool_calls,omitempty"`
	ToolResults    map[string]interface{} `json:"tool_results,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	Feedback       *Feedback              `json:"feedback,omitempty"`
}

// Feedback 表示对消息的反馈
type Feedback struct {
	Rating     int       `json:"rating"` // 1-5
	Comment    string    `json:"comment,omitempty"`
	SubmittedAt time.Time `json:"submitted_at"`
}

// ConversationManager 管理对话会话
type ConversationManager struct {
	conversations map[string]*Conversation
	messages      map[string][]Message
	mu            sync.RWMutex
}

// NewConversationManager 创建新的对话管理器
func NewConversationManager() *ConversationManager {
	return &ConversationManager{
		conversations: make(map[string]*Conversation),
		messages:      make(map[string][]Message),
	}
}

// CreateConversation 创建新的对话
func (cm *ConversationManager) CreateConversation(agentID string, title string, agent *Agent) (*Conversation, error) {
	if agent == nil {
		return nil, errors.New("agent cannot be nil")
	}

	now := time.Now()
	convID := uuid.New().String()

	conv := &Conversation{
		ID:        convID,
		AgentID:   agentID,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
		Memory:    NewSimpleMemory(100), // 默认使用简单内存
		Agent:     agent,
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.conversations[convID] = conv
	cm.messages[convID] = []Message{}

	return conv, nil
}

// GetConversation 获取指定ID的对话
func (cm *ConversationManager) GetConversation(id string) (*Conversation, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conv, exists := cm.conversations[id]
	if !exists {
		return nil, errors.New("conversation not found")
	}

	return conv, nil
}

// ListConversations 列出指定智能体的所有对话
func (cm *ConversationManager) ListConversations(agentID string) []*Conversation {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	var conversations []*Conversation
	for _, conv := range cm.conversations {
		if conv.AgentID == agentID {
			conversations = append(conversations, conv)
		}
	}

	return conversations
}

// DeleteConversation 删除指定ID的对话
func (cm *ConversationManager) DeleteConversation(id string) error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if _, exists := cm.conversations[id]; !exists {
		return errors.New("conversation not found")
	}

	delete(cm.conversations, id)
	delete(cm.messages, id)
	return nil
}

// AddMessage 向对话添加消息
func (cm *ConversationManager) AddMessage(conversationID string, role string, content string) (*Message, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	conv, exists := cm.conversations[conversationID]
	if !exists {
		return nil, errors.New("conversation not found")
	}

	msg := Message{
		ID:             uuid.New().String(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	// 更新对话的修改时间
	conv.UpdatedAt = time.Now()

	// 添加消息到eino内存
	if conv.Memory != nil {
		einoMsg := eino.Message{
			Role:    role,
			Content: content,
		}
		if err := conv.Memory.AddMessage(einoMsg); err != nil {
			return nil, err
		}
	}

	// 保存消息到管理器
	cm.messages[conversationID] = append(cm.messages[conversationID], msg)

	return &msg, nil
}

// GetMessages 获取对话的所有消息
func (cm *ConversationManager) GetMessages(conversationID string) ([]Message, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	msgs, exists := cm.messages[conversationID]
	if !exists {
		return nil, errors.New("conversation not found")
	}

	return msgs, nil
}

// SendMessage 发送消息并获取智能体回复
func (cm *ConversationManager) SendMessage(ctx context.Context, conversationID string, content string) (*Message, error) {
	// 添加用户消息
	userMsg, err := cm.AddMessage(conversationID, "user", content)
	if err != nil {
		return nil, err
	}

	// 获取对话
	conv, err := cm.GetConversation(conversationID)
	if err != nil {
		return nil, err
	}

	// 获取回复
	reply, err := conv.Agent.Chat(ctx, content)
	if err != nil {
		return nil, err
	}

	// 添加助手消息
	assistantMsg, err := cm.AddMessage(conversationID, "assistant", reply)
	if err != nil {
		return nil, err
	}

	return assistantMsg, nil
}

// AddFeedback 为消息添加反馈
func (cm *ConversationManager) AddFeedback(messageID, conversationID string, rating int, comment string) error {
	if rating < 1 || rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	cm.mu.Lock()
	defer cm.mu.Unlock()

	msgs, exists := cm.messages[conversationID]
	if !exists {
		return errors.New("conversation not found")
	}

	for i, msg := range msgs {
		if msg.ID == messageID {
			msgs[i].Feedback = &Feedback{
				Rating:     rating,
				Comment:    comment,
				SubmittedAt: time.Now(),
			}
			cm.messages[conversationID] = msgs
			return nil
		}
	}

	return errors.New("message not found")
}

// DefaultConversationManager 全局默认的对话管理器实例
var DefaultConversationManager = NewConversationManager() 