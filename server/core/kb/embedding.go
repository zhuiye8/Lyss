package kb

import (
	"context"
	"errors"
	"sync"
)

// EmbeddingModel 向量化模型接口
type EmbeddingModel interface {
	// Embed 将文本转换为向量表示
	Embed(ctx context.Context, texts []string) ([][]float32, error)
	
	// Dimensions 返回向量维度
	Dimensions() int
	
	// ModelName 返回模型名称
	ModelName() string
}

// EmbeddingRequest 向量化请求
type EmbeddingRequest struct {
	Texts []string
	Model string
}

// EmbeddingResponse 向量化响应
type EmbeddingResponse struct {
	Embeddings [][]float32
	Model      string
	Dimensions int
}

// EmbeddingManager 向量化管理器
type EmbeddingManager struct {
	models map[string]EmbeddingModel
	mu     sync.RWMutex
}

// NewEmbeddingManager 创建新的向量化管理器
func NewEmbeddingManager() *EmbeddingManager {
	return &EmbeddingManager{
		models: make(map[string]EmbeddingModel),
	}
}

// RegisterModel 注册向量模型
func (m *EmbeddingManager) RegisterModel(name string, model EmbeddingModel) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.models[name] = model
}

// GetModel 获取向量模型
func (m *EmbeddingManager) GetModel(name string) (EmbeddingModel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	model, exists := m.models[name]
	if !exists {
		return nil, errors.New("embedding model not found: " + name)
	}
	
	return model, nil
}

// Embed 使用指定模型进行向量化
func (m *EmbeddingManager) Embed(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	if len(req.Texts) == 0 {
		return nil, errors.New("no texts provided for embedding")
	}
	
	model, err := m.GetModel(req.Model)
	if err != nil {
		return nil, err
	}
	
	embeddings, err := model.Embed(ctx, req.Texts)
	if err != nil {
		return nil, err
	}
	
	return &EmbeddingResponse{
		Embeddings: embeddings,
		Model:      req.Model,
		Dimensions: model.Dimensions(),
	}, nil
}

// ListModels 列出所有可用的向量模型
func (m *EmbeddingManager) ListModels() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	models := make([]string, 0, len(m.models))
	for name := range m.models {
		models = append(models, name)
	}
	
	return models
}

// MockEmbeddingModel 模拟的向量模型实现（用于测试）
type MockEmbeddingModel struct {
	dimensions int
	name       string
}

// NewMockEmbeddingModel 创建新的模拟向量模型
func NewMockEmbeddingModel(dimensions int, name string) *MockEmbeddingModel {
	return &MockEmbeddingModel{
		dimensions: dimensions,
		name:       name,
	}
}

// Embed 实现向量化接口
func (m *MockEmbeddingModel) Embed(ctx context.Context, texts []string) ([][]float32, error) {
	result := make([][]float32, len(texts))
	
	// 创建模拟向量（所有值为0.1）
	for i := range texts {
		vector := make([]float32, m.dimensions)
		for j := range vector {
			vector[j] = 0.1
		}
		result[i] = vector
	}
	
	return result, nil
}

// Dimensions 返回向量维度
func (m *MockEmbeddingModel) Dimensions() int {
	return m.dimensions
}

// ModelName 返回模型名称
func (m *MockEmbeddingModel) ModelName() string {
	return m.name
}

// DefaultEmbeddingManager 默认的向量化管理器
var DefaultEmbeddingManager = NewEmbeddingManager()

// 初始化默认向量模型
func init() {
	// 注册一个模拟模型用于测试
	mockModel := NewMockEmbeddingModel(384, "mock-embedding-model")
	DefaultEmbeddingManager.RegisterModel("mock", mockModel)
} 
