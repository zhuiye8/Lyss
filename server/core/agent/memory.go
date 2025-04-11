package agent

import (
	"errors"
	"sync"

	"github.com/cloudwego/eino"
)

// SimpleMemory 是一个简单的内存实现，用于存储对话历�?
type SimpleMemory struct {
	messages []eino.Message
	maxSize  int
	mu       sync.Mutex
}

// NewSimpleMemory 创建一个新的简单内存实�?
// maxSize 指定最大消息数量，如果�?则不限制
func NewSimpleMemory(maxSize int) *SimpleMemory {
	return &SimpleMemory{
		messages: make([]eino.Message, 0),
		maxSize:  maxSize,
	}
}

// AddMessage 添加消息到内�?
func (m *SimpleMemory) AddMessage(msg eino.Message) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 如果设置了最大消息数量且已达到上限，则移除最旧的消息
	if m.maxSize > 0 && len(m.messages) >= m.maxSize {
		m.messages = m.messages[1:]
	}

	m.messages = append(m.messages, msg)
	return nil
}

// GetMessages 获取所有消�?
func (m *SimpleMemory) GetMessages() ([]eino.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 创建消息的副本以避免外部修改
	result := make([]eino.Message, len(m.messages))
	copy(result, m.messages)

	return result, nil
}

// Clear 清除所有消�?
func (m *SimpleMemory) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = make([]eino.Message, 0)
	return nil
}

// GetLastUserMessage 获取最后一条用户消�?
func (m *SimpleMemory) GetLastUserMessage() (eino.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 反向遍历消息查找最后一条用户消�?
	for i := len(m.messages) - 1; i >= 0; i-- {
		if m.messages[i].Role == eino.RoleUser {
			return m.messages[i], nil
		}
	}

	return eino.Message{}, errors.New("no user message found")
}

// GetLastAssistantMessage 获取最后一条助手消�?
func (m *SimpleMemory) GetLastAssistantMessage() (eino.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 反向遍历消息查找最后一条助手消�?
	for i := len(m.messages) - 1; i >= 0; i-- {
		if m.messages[i].Role == eino.RoleAssistant {
			return m.messages[i], nil
		}
	}

	return eino.Message{}, errors.New("no assistant message found")
} 
