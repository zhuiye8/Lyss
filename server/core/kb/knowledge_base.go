package kb

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	"github.com/google/uuid"
)

// KnowledgeBase 知识库
type KnowledgeBase struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DocumentCount int       `json:"document_count"`
	ChunkCount   int       `json:"chunk_count"`
	EmbeddingModel string    `json:"embedding_model"`
	Metadata     interface{} `json:"metadata,omitempty"`
}

// KnowledgeBaseManager 知识库管理器
type KnowledgeBaseManager struct {
	vectorDB  VectorDatabase
	embedding *EmbeddingManager
	processor *DocumentProcessorRegistry
	knowledgeBases map[string]*KnowledgeBase // id -> KnowledgeBase
	documents     map[string][]*Document     // knowledgeBaseId -> []Document
	mu            sync.RWMutex
}

// NewKnowledgeBaseManager 创建知识库管理器
func NewKnowledgeBaseManager(vectorDB VectorDatabase, embedding *EmbeddingManager, processor *DocumentProcessorRegistry) *KnowledgeBaseManager {
	return &KnowledgeBaseManager{
		vectorDB:       vectorDB,
		embedding:      embedding,
		processor:      processor,
		knowledgeBases: make(map[string]*KnowledgeBase),
		documents:      make(map[string][]*Document),
	}
}

// CreateKnowledgeBase 创建知识库
func (m *KnowledgeBaseManager) CreateKnowledgeBase(ctx context.Context, name, description, embeddingModel string) (*KnowledgeBase, error) {
	// 验证向量模型是否存在
	if _, err := m.embedding.GetModel(embeddingModel); err != nil {
		return nil, fmt.Errorf("invalid embedding model: %w", err)
	}
	
	now := time.Now()
	kbID := uuid.New().String()
	
	// 创建知识库对象
	kb := &KnowledgeBase{
		ID:             kbID,
		Name:           name,
		Description:    description,
		CreatedAt:      now,
		UpdatedAt:      now,
		DocumentCount:  0,
		ChunkCount:     0,
		EmbeddingModel: embeddingModel,
	}
	
	// 创建向量数据库集合
	model, _ := m.embedding.GetModel(embeddingModel)
	dimension := model.Dimensions()
	if err := m.vectorDB.CreateCollection(ctx, kbID, dimension); err != nil {
		return nil, fmt.Errorf("failed to create vector collection: %w", err)
	}
	
	// 保存知识库
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.knowledgeBases[kbID] = kb
	m.documents[kbID] = []*Document{}
	
	return kb, nil
}

// GetKnowledgeBase 获取知识库
func (m *KnowledgeBaseManager) GetKnowledgeBase(id string) (*KnowledgeBase, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	kb, exists := m.knowledgeBases[id]
	if !exists {
		return nil, errors.New("knowledge base not found")
	}
	
	return kb, nil
}

// ListKnowledgeBases 列出所有知识库
func (m *KnowledgeBaseManager) ListKnowledgeBases() []*KnowledgeBase {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	kbs := make([]*KnowledgeBase, 0, len(m.knowledgeBases))
	for _, kb := range m.knowledgeBases {
		kbs = append(kbs, kb)
	}
	
	return kbs
}

// UpdateKnowledgeBase 更新知识库
func (m *KnowledgeBaseManager) UpdateKnowledgeBase(id, name, description string) (*KnowledgeBase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	kb, exists := m.knowledgeBases[id]
	if !exists {
		return nil, errors.New("knowledge base not found")
	}
	
	if name != "" {
		kb.Name = name
	}
	
	if description != "" {
		kb.Description = description
	}
	
	kb.UpdatedAt = time.Now()
	
	return kb, nil
}

// DeleteKnowledgeBase 删除知识库
func (m *KnowledgeBaseManager) DeleteKnowledgeBase(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if _, exists := m.knowledgeBases[id]; !exists {
		return errors.New("knowledge base not found")
	}
	
	// 删除向量数据库集合
	if err := m.vectorDB.DropCollection(ctx, id); err != nil {
		return fmt.Errorf("failed to drop vector collection: %w", err)
	}
	
	// 删除知识库和相关文档
	delete(m.knowledgeBases, id)
	delete(m.documents, id)
	
	return nil
}

// AddDocument 添加文档到知识库
func (m *KnowledgeBaseManager) AddDocument(ctx context.Context, knowledgeBaseID string, file *multipart.FileHeader) (*Document, error) {
	m.mu.Lock()
	kb, exists := m.knowledgeBases[knowledgeBaseID]
	if !exists {
		m.mu.Unlock()
		return nil, errors.New("knowledge base not found")
	}
	m.mu.Unlock()
	
	// 创建文档
	doc, err := NewDocument(knowledgeBaseID, file)
	if err != nil {
		return nil, fmt.Errorf("failed to create document: %w", err)
	}
	
	// 处理文档
	chunks, err := m.processor.ProcessDocument(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}
	
	if len(chunks) == 0 {
		return nil, errors.New("document processing resulted in no chunks")
	}
	
	// 转换文档块为文本
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}
	
	// 生成向量嵌入
	resp, err := m.embedding.Embed(ctx, EmbeddingRequest{
		Texts: texts,
		Model: kb.EmbeddingModel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}
	
	// 准备向量插入数据
	ids := make([]string, len(chunks))
	vectors := resp.Embeddings
	metadata := make([]map[string]interface{}, len(chunks))
	
	for i, chunk := range chunks {
		ids[i] = chunk.ID
		chunk.Vector = vectors[i]
		
		// 设置元数据
		metadata[i] = map[string]interface{}{
			"content":      chunk.Content,
			"document_id":  chunk.DocumentID,
			"chunk_index":  chunk.Metadata.ChunkIndex,
			"page_number":  chunk.Metadata.PageNumber,
			"source":       chunk.Metadata.Source,
		}
	}
	
	// 插入向量
	if err := m.vectorDB.InsertVectors(ctx, knowledgeBaseID, ids, vectors, metadata); err != nil {
		return nil, fmt.Errorf("failed to insert vectors: %w", err)
	}
	
	// 更新知识库信息
	m.mu.Lock()
	defer m.mu.Unlock()
	
	kb.DocumentCount++
	kb.ChunkCount += len(chunks)
	kb.UpdatedAt = time.Now()
	
	// 保存文档
	m.documents[knowledgeBaseID] = append(m.documents[knowledgeBaseID], doc)
	
	return doc, nil
}

// AddTextDocument 添加文本到知识库
func (m *KnowledgeBaseManager) AddTextDocument(ctx context.Context, knowledgeBaseID, name, content string, docType DocumentType) (*Document, error) {
	m.mu.Lock()
	kb, exists := m.knowledgeBases[knowledgeBaseID]
	if !exists {
		m.mu.Unlock()
		return nil, errors.New("knowledge base not found")
	}
	m.mu.Unlock()
	
	// 创建文档
	doc := NewTextDocument(knowledgeBaseID, name, content, docType)
	
	// 处理文档
	chunks, err := m.processor.ProcessDocument(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to process document: %w", err)
	}
	
	if len(chunks) == 0 {
		return nil, errors.New("document processing resulted in no chunks")
	}
	
	// 转换文档块为文本
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Content
	}
	
	// 生成向量嵌入
	resp, err := m.embedding.Embed(ctx, EmbeddingRequest{
		Texts: texts,
		Model: kb.EmbeddingModel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}
	
	// 准备向量插入数据
	ids := make([]string, len(chunks))
	vectors := resp.Embeddings
	metadata := make([]map[string]interface{}, len(chunks))
	
	for i, chunk := range chunks {
		ids[i] = chunk.ID
		chunk.Vector = vectors[i]
		
		// 设置元数据
		metadata[i] = map[string]interface{}{
			"content":      chunk.Content,
			"document_id":  chunk.DocumentID,
			"chunk_index":  chunk.Metadata.ChunkIndex,
			"page_number":  chunk.Metadata.PageNumber,
			"source":       chunk.Metadata.Source,
		}
	}
	
	// 插入向量
	if err := m.vectorDB.InsertVectors(ctx, knowledgeBaseID, ids, vectors, metadata); err != nil {
		return nil, fmt.Errorf("failed to insert vectors: %w", err)
	}
	
	// 更新知识库信息
	m.mu.Lock()
	defer m.mu.Unlock()
	
	kb.DocumentCount++
	kb.ChunkCount += len(chunks)
	kb.UpdatedAt = time.Now()
	
	// 保存文档
	m.documents[knowledgeBaseID] = append(m.documents[knowledgeBaseID], doc)
	
	return doc, nil
}

// GetDocuments 获取知识库的所有文档
func (m *KnowledgeBaseManager) GetDocuments(knowledgeBaseID string) ([]*Document, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if _, exists := m.knowledgeBases[knowledgeBaseID]; !exists {
		return nil, errors.New("knowledge base not found")
	}
	
	docs, exists := m.documents[knowledgeBaseID]
	if !exists {
		return []*Document{}, nil
	}
	
	return docs, nil
}

// DeleteDocument 从知识库删除文档
func (m *KnowledgeBaseManager) DeleteDocument(ctx context.Context, knowledgeBaseID, documentID string) error {
	m.mu.Lock()
	kb, exists := m.knowledgeBases[knowledgeBaseID]
	if !exists {
		m.mu.Unlock()
		return errors.New("knowledge base not found")
	}
	
	docs, exists := m.documents[knowledgeBaseID]
	if !exists {
		m.mu.Unlock()
		return errors.New("no documents in knowledge base")
	}
	
	// 查找文档
	var docToDelete *Document
	var newDocs []*Document
	var docIndex int = -1
	
	for i, doc := range docs {
		if doc.ID == documentID {
			docToDelete = doc
			docIndex = i
		} else {
			newDocs = append(newDocs, doc)
		}
	}
	
	if docToDelete == nil {
		m.mu.Unlock()
		return errors.New("document not found")
	}
	
	// 更新文档列表
	m.documents[knowledgeBaseID] = newDocs
	m.mu.Unlock()
	
	// 构建删除表达式（所有与此文档关联的向量）
	// 注意：这里应该是向量数据库中特定的实现
	// 在实际应用中可能需要先查询所有文档相关的向量ID
	expr := fmt.Sprintf("document_id == '%s'", documentID)
	
	// 删除向量数据
	// TODO: 实现根据表达式删除向量的功能
	// 目前简化处理，假设我们已知所有向量ID
	
	// 更新知识库信息
	m.mu.Lock()
	kb.DocumentCount--
	// 注意：这里需要知道删除了多少个chunk
	// kb.ChunkCount -= deletedChunksCount
	kb.UpdatedAt = time.Now()
	m.mu.Unlock()
	
	return nil
}

// DefaultKnowledgeBaseManager 默认的知识库管理器
var DefaultKnowledgeBaseManager *KnowledgeBaseManager

// 初始化默认知识库管理器
func InitDefaultKnowledgeBaseManager(vectorDB VectorDatabase) {
	DefaultKnowledgeBaseManager = NewKnowledgeBaseManager(
		vectorDB,
		DefaultEmbeddingManager,
		DefaultProcessorRegistry,
	)
} 
