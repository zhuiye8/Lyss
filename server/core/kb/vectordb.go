package kb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// VectorDBConfig 向量数据库配置
type VectorDBConfig struct {
	// Milvus配置
	MilvusHost     string
	MilvusPort     int
	MilvusUsername string
	MilvusPassword string
}

// SearchResult 搜索结果
type SearchResult struct {
	ChunkID  string    `json:"chunk_id"`
	Content  string    `json:"content"`
	Score    float32   `json:"score"`
	Metadata interface{} `json:"metadata,omitempty"`
}

// VectorDatabase 向量数据库接口
type VectorDatabase interface {
	// Connect 连接数据库
	Connect(ctx context.Context) error
	
	// Disconnect 断开连接
	Disconnect(ctx context.Context) error
	
	// CreateCollection 创建集合
	CreateCollection(ctx context.Context, name string, dimension int) error
	
	// DropCollection 删除集合
	DropCollection(ctx context.Context, name string) error
	
	// InsertVectors 插入向量
	InsertVectors(ctx context.Context, collectionName string, ids []string, vectors [][]float32, metadata []map[string]interface{}) error
	
	// Search 搜索向量
	Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]SearchResult, error)
	
	// DeleteVectors 删除向量
	DeleteVectors(ctx context.Context, collectionName string, ids []string) error
}

// MilvusDB Milvus向量数据库实现
type MilvusDB struct {
	client     client.Client
	config     VectorDBConfig
	collections map[string]bool
	mu          sync.RWMutex
	connected   bool
}

// NewMilvusDB 创建Milvus数据库实例
func NewMilvusDB(config VectorDBConfig) *MilvusDB {
	return &MilvusDB{
		config:      config,
		collections: make(map[string]bool),
		connected:   false,
	}
}

// Connect 连接Milvus
func (m *MilvusDB) Connect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if m.connected {
		return nil
	}
	
	addr := fmt.Sprintf("%s:%d", m.config.MilvusHost, m.config.MilvusPort)
	
	var err error
	
	// 根据是否提供认证信息选择连接方式
	if m.config.MilvusUsername != "" && m.config.MilvusPassword != "" {
		m.client, err = client.NewGrpcClient(ctx, addr, 
			client.WithUsername(m.config.MilvusUsername), 
			client.WithPassword(m.config.MilvusPassword))
	} else {
		m.client, err = client.NewGrpcClient(ctx, addr)
	}
	
	if err != nil {
		return fmt.Errorf("failed to connect to Milvus: %w", err)
	}
	
	m.connected = true
	return nil
}

// Disconnect 断开Milvus连接
func (m *MilvusDB) Disconnect(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if !m.connected {
		return nil
	}
	
	if err := m.client.Close(); err != nil {
		return fmt.Errorf("failed to disconnect from Milvus: %w", err)
	}
	
	m.connected = false
	return nil
}

// CreateCollection 创建Milvus集合
func (m *MilvusDB) CreateCollection(ctx context.Context, name string, dimension int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if !m.connected {
		return errors.New("not connected to Milvus")
	}
	
	// 检查集合是否已存在
	has, err := m.client.HasCollection(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}
	
	if has {
		m.collections[name] = true
		return nil // 集合已存在，直接返回
	}
	
	// 定义集合模式
	schema := &entity.Schema{
		CollectionName: name,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				AutoID:     false,
				MaxLength:  100,
			},
			{
				Name:     "vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": fmt.Sprintf("%d", dimension),
				},
			},
			{
				Name:       "content",
				DataType:   entity.FieldTypeVarChar,
				MaxLength:  65535, // 最大文本长度
			},
			{
				Name:       "metadata",
				DataType:   entity.FieldTypeJSON,
			},
		},
	}
	
	// 创建集合
	err = m.client.CreateCollection(ctx, schema, entity.DefaultShardNumber)
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	
	// 创建索引
	idx, err := entity.NewIndexIvfFlat(entity.L2, 1024)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	
	err = m.client.CreateIndex(ctx, name, "vector", idx, false)
	if err != nil {
		return fmt.Errorf("failed to create index on collection: %w", err)
	}
	
	// 加载集合到内存
	err = m.client.LoadCollection(ctx, name, false)
	if err != nil {
		return fmt.Errorf("failed to load collection: %w", err)
	}
	
	m.collections[name] = true
	return nil
}

// DropCollection 删除Milvus集合
func (m *MilvusDB) DropCollection(ctx context.Context, name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	if !m.connected {
		return errors.New("not connected to Milvus")
	}
	
	// 检查集合是否存在
	has, err := m.client.HasCollection(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check collection existence: %w", err)
	}
	
	if !has {
		delete(m.collections, name)
		return nil // 集合不存在，直接返回
	}
	
	// 释放集合
	err = m.client.ReleaseCollection(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to release collection: %w", err)
	}
	
	// 删除集合
	err = m.client.DropCollection(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to drop collection: %w", err)
	}
	
	delete(m.collections, name)
	return nil
}

// InsertVectors 插入向量数据
func (m *MilvusDB) InsertVectors(ctx context.Context, collectionName string, ids []string, vectors [][]float32, metadata []map[string]interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if !m.connected {
		return errors.New("not connected to Milvus")
	}
	
	if _, exists := m.collections[collectionName]; !exists {
		return fmt.Errorf("collection %s does not exist", collectionName)
	}
	
	if len(ids) != len(vectors) || len(ids) != len(metadata) {
		return errors.New("ids, vectors, and metadata must have the same length")
	}
	
	idColumn := entity.NewColumnVarChar("id", ids)
	vectorColumn := entity.NewColumnFloatVector("vector", int64(len(vectors[0])), vectors)
	
	// 提取内容字段
	contents := make([]string, len(metadata))
	for i, meta := range metadata {
		if content, ok := meta["content"].(string); ok {
			contents[i] = content
			delete(meta, "content") // 从metadata中移除content
		}
	}
	contentColumn := entity.NewColumnVarChar("content", contents)
	
	// 构造JSON元数据
	metadataJson := make([][]byte, len(metadata))
	for i, meta := range metadata {
		jsonBytes, err := json.Marshal(meta)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataJson[i] = jsonBytes
	}
	metadataColumn := entity.NewColumnJSONBytes("metadata", metadataJson)
	
	// 插入数据
	_, err := m.client.Insert(ctx, collectionName, "", idColumn, vectorColumn, contentColumn, metadataColumn)
	if err != nil {
		return fmt.Errorf("failed to insert vectors: %w", err)
	}
	
	return nil
}

// Search 执行向量搜索
func (m *MilvusDB) Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]SearchResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if !m.connected {
		return nil, errors.New("not connected to Milvus")
	}
	
	if _, exists := m.collections[collectionName]; !exists {
		return nil, fmt.Errorf("collection %s does not exist", collectionName)
	}
	
	// 准备搜索参数
	vectors := []entity.Vector{entity.FloatVector(vector)}
	searchParam, err := entity.NewIndexIvfFlatSearchParam(10)
	if err != nil {
		return nil, fmt.Errorf("failed to create search param: %w", err)
	}
	
	// 执行搜索
	searchResult, err := m.client.Search(
		ctx,
		collectionName,
		"",
		"",
		[]string{"id", "content", "metadata"},
		vectors,
		"vector",
		entity.L2,
		topK,
		searchParam,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute search: %w", err)
	}
	
	results := make([]SearchResult, 0, topK)
	
	// 解析搜索结果
	for i := 0; i < len(searchResult); i++ {
		ids := searchResult[i].IDs.(*entity.ColumnVarChar).Data()
		distances := searchResult[i].Scores
		
		// 获取内容和元数据
		fieldData := searchResult[i].Fields
		contents := fieldData[0].(*entity.ColumnVarChar).Data()
		metadataBytes := fieldData[1].(*entity.ColumnJSONBytes).Data()
		
		for j := 0; j < len(ids); j++ {
			// 解析元数据
			var metadata map[string]interface{}
			if err := json.Unmarshal(metadataBytes[j], &metadata); err != nil {
				// 忽略解析错误，使用空元数据
				metadata = make(map[string]interface{})
			}
			
			results = append(results, SearchResult{
				ChunkID:  ids[j],
				Content:  contents[j],
				Score:    distances[j],
				Metadata: metadata,
			})
		}
	}
	
	return results, nil
}

// DeleteVectors 删除向量
func (m *MilvusDB) DeleteVectors(ctx context.Context, collectionName string, ids []string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if !m.connected {
		return errors.New("not connected to Milvus")
	}
	
	if _, exists := m.collections[collectionName]; !exists {
		return fmt.Errorf("collection %s does not exist", collectionName)
	}
	
	// 构建删除表达式
	var expr string
	if len(ids) == 1 {
		expr = fmt.Sprintf("id == '%s'", ids[0])
	} else {
		expr = "id in ['" + strings.Join(ids, "','") + "']"
	}
	
	// 执行删除
	_, err := m.client.Delete(ctx, collectionName, "", expr)
	if err != nil {
		return fmt.Errorf("failed to delete vectors: %w", err)
	}
	
	return nil
}

// DefaultVectorDB 默认的向量数据库实例
var DefaultVectorDB VectorDatabase

// 初始化默认向量数据库
func InitDefaultVectorDB(config VectorDBConfig) error {
	DefaultVectorDB = NewMilvusDB(config)
	return nil
}

// InMemoryVectorDB 基于内存的向量数据库实现（用于测试）
type InMemoryVectorDB struct {
	collections map[string]*inMemoryCollection
	mu          sync.RWMutex
}

type inMemoryCollection struct {
	dimension int
	vectors   map[string]struct {
		Vector   []float32
		Content  string
		Metadata map[string]interface{}
	}
}

// NewInMemoryVectorDB 创建内存向量数据库
func NewInMemoryVectorDB() *InMemoryVectorDB {
	return &InMemoryVectorDB{
		collections: make(map[string]*inMemoryCollection),
	}
}

// Connect 连接（内存实现为空操作）
func (db *InMemoryVectorDB) Connect(ctx context.Context) error {
	return nil
}

// Disconnect 断开连接（内存实现为空操作）
func (db *InMemoryVectorDB) Disconnect(ctx context.Context) error {
	return nil
}

// CreateCollection 创建集合
func (db *InMemoryVectorDB) CreateCollection(ctx context.Context, name string, dimension int) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	if _, exists := db.collections[name]; exists {
		return nil // 集合已存在
	}
	
	db.collections[name] = &inMemoryCollection{
		dimension: dimension,
		vectors:   make(map[string]struct {
			Vector   []float32
			Content  string
			Metadata map[string]interface{}
		}),
	}
	
	return nil
}

// DropCollection 删除集合
func (db *InMemoryVectorDB) DropCollection(ctx context.Context, name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	delete(db.collections, name)
	return nil
}

// InsertVectors 插入向量
func (db *InMemoryVectorDB) InsertVectors(ctx context.Context, collectionName string, ids []string, vectors [][]float32, metadata []map[string]interface{}) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	collection, exists := db.collections[collectionName]
	if !exists {
		return fmt.Errorf("collection %s does not exist", collectionName)
	}
	
	if len(ids) != len(vectors) || len(ids) != len(metadata) {
		return errors.New("ids, vectors, and metadata must have the same length")
	}
	
	for i, id := range ids {
		content := ""
		if contentVal, ok := metadata[i]["content"]; ok {
			if contentStr, ok := contentVal.(string); ok {
				content = contentStr
				delete(metadata[i], "content")
			}
		}
		
		collection.vectors[id] = struct {
			Vector   []float32
			Content  string
			Metadata map[string]interface{}
		}{
			Vector:   vectors[i],
			Content:  content,
			Metadata: metadata[i],
		}
	}
	
	return nil
}

// Search 搜索向量
func (db *InMemoryVectorDB) Search(ctx context.Context, collectionName string, vector []float32, topK int) ([]SearchResult, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	
	collection, exists := db.collections[collectionName]
	if !exists {
		return nil, fmt.Errorf("collection %s does not exist", collectionName)
	}
	
	type ScoredResult struct {
		ID       string
		Content  string
		Score    float32
		Metadata map[string]interface{}
	}
	
	results := make([]ScoredResult, 0, len(collection.vectors))
	
	// 计算所有向量的距离
	for id, data := range collection.vectors {
		score := calculateCosineSimilarity(vector, data.Vector)
		results = append(results, ScoredResult{
			ID:       id,
			Content:  data.Content,
			Score:    score,
			Metadata: data.Metadata,
		})
	}
	
	// 按相似度排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	// 限制结果数量
	if topK > 0 && len(results) > topK {
		results = results[:topK]
	}
	
	// 转换为SearchResult格式
	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		searchResults[i] = SearchResult{
			ChunkID:  r.ID,
			Content:  r.Content,
			Score:    r.Score,
			Metadata: r.Metadata,
		}
	}
	
	return searchResults, nil
}

// DeleteVectors 删除向量
func (db *InMemoryVectorDB) DeleteVectors(ctx context.Context, collectionName string, ids []string) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	
	collection, exists := db.collections[collectionName]
	if !exists {
		return fmt.Errorf("collection %s does not exist", collectionName)
	}
	
	for _, id := range ids {
		delete(collection.vectors, id)
	}
	
	return nil
}

// 计算余弦相似度（用于内存实现）
func calculateCosineSimilarity(v1, v2 []float32) float32 {
	var dotProduct float32
	var norm1 float32
	var norm2 float32
	
	for i := 0; i < len(v1) && i < len(v2); i++ {
		dotProduct += v1[i] * v2[i]
		norm1 += v1[i] * v1[i]
		norm2 += v2[i] * v2[i]
	}
	
	if norm1 == 0 || norm2 == 0 {
		return 0
	}
	
	return dotProduct / (float32(math.Sqrt(float64(norm1))) * float32(math.Sqrt(float64(norm2))))
} 
