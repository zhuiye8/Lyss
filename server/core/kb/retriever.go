package kb

import (
	"context"
	"fmt"

	"github.com/yourusername/myApp/server/core/agent"
)

// QueryRequest 知识库查询请求
type QueryRequest struct {
	KnowledgeBaseID string `json:"knowledge_base_id"`
	Query           string `json:"query"`
	TopK            int    `json:"top_k"`
}

// QueryResponse 知识库查询响应
type QueryResponse struct {
	Results []SearchResult `json:"results"`
	Query   string         `json:"query"`
}

// Retriever 检索器接口
type Retriever interface {
	// Retrieve 根据查询检索相关内容
	Retrieve(ctx context.Context, req QueryRequest) (*QueryResponse, error)
}

// VectorRetriever 基于向量的检索器
type VectorRetriever struct {
	kb        *KnowledgeBaseManager
	embedding *EmbeddingManager
}

// NewVectorRetriever 创建向量检索器
func NewVectorRetriever(kb *KnowledgeBaseManager, embedding *EmbeddingManager) *VectorRetriever {
	return &VectorRetriever{
		kb:        kb,
		embedding: embedding,
	}
}

// Retrieve 实现检索功能
func (r *VectorRetriever) Retrieve(ctx context.Context, req QueryRequest) (*QueryResponse, error) {
	// 获取知识库信息
	kb, err := r.kb.GetKnowledgeBase(req.KnowledgeBaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}
	
	// 对查询进行向量化
	resp, err := r.embedding.Embed(ctx, EmbeddingRequest{
		Texts: []string{req.Query},
		Model: kb.EmbeddingModel,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}
	
	// 设置默认topK
	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}
	
	// 搜索相似向量
	results, err := r.kb.vectorDB.Search(ctx, req.KnowledgeBaseID, resp.Embeddings[0], topK)
	if err != nil {
		return nil, fmt.Errorf("failed to search vectors: %w", err)
	}
	
	return &QueryResponse{
		Results: results,
		Query:   req.Query,
	}, nil
}

// DefaultRetriever 默认的检索器
var DefaultRetriever Retriever

// InitDefaultRetriever 初始化默认检索器
func InitDefaultRetriever() {
	DefaultRetriever = NewVectorRetriever(DefaultKnowledgeBaseManager, DefaultEmbeddingManager)
}

// RAGTool RAG工具
type RAGTool struct {
	retriever Retriever
}

// NewRAGTool 创建RAG工具
func NewRAGTool(retriever Retriever) *RAGTool {
	return &RAGTool{
		retriever: retriever,
	}
}

// RegisterRAGTool 注册RAG工具到工具注册表
func RegisterRAGTool(registry *agent.ToolRegistry, retriever Retriever) error {
	ragTool := agent.Tool{
		Name:        "knowledge_search",
		Description: "从知识库搜索相关信息",
		Parameters: map[string]interface{}{
			"knowledge_base_id": map[string]interface{}{
				"type":        "string",
				"description": "要搜索的知识库ID",
			},
			"query": map[string]interface{}{
				"type":        "string",
				"description": "搜索查询",
			},
			"top_k": map[string]interface{}{
				"type":        "integer",
				"description": "返回的最大结果数",
				"default":     5,
			},
		},
		Handler: func(ctx context.Context, params map[string]interface{}) (interface{}, error) {
			// 提取参数
			kbID, ok := params["knowledge_base_id"].(string)
			if !ok {
				return nil, fmt.Errorf("knowledge_base_id must be a string")
			}
			
			query, ok := params["query"].(string)
			if !ok {
				return nil, fmt.Errorf("query must be a string")
			}
			
			topK := 5
			if topKParam, ok := params["top_k"]; ok {
				if topKFloat, ok := topKParam.(float64); ok {
					topK = int(topKFloat)
				}
			}
			
			// 创建请求
			req := QueryRequest{
				KnowledgeBaseID: kbID,
				Query:           query,
				TopK:            topK,
			}
			
			// 执行检索
			rag := NewRAGTool(retriever)
			return rag.Search(ctx, req)
		},
	}
	
	return registry.RegisterTool(ragTool)
}

// Search 执行知识库搜索
func (r *RAGTool) Search(ctx context.Context, req QueryRequest) (interface{}, error) {
	resp, err := r.retriever.Retrieve(ctx, req)
	if err != nil {
		return nil, err
	}
	
	// 转换为友好的返回格式
	type ResultItem struct {
		Content  string  `json:"content"`
		Source   string  `json:"source"`
		Score    float32 `json:"score"`
	}
	
	items := make([]ResultItem, len(resp.Results))
	for i, result := range resp.Results {
		source := "unknown"
		if meta, ok := result.Metadata.(map[string]interface{}); ok {
			if src, ok := meta["source"].(string); ok {
				source = src
			}
		}
		
		items[i] = ResultItem{
			Content: result.Content,
			Source:  source,
			Score:   result.Score,
		}
	}
	
	return map[string]interface{}{
		"query":   resp.Query,
		"results": items,
	}, nil
}

// GeneratePromptFromResults 根据检索结果生成提示词
func GeneratePromptFromResults(results []SearchResult, userQuery string) string {
	prompt := "以下是与问题相关的信息：\n\n"
	
	for i, result := range results {
		prompt += fmt.Sprintf("信息片段 %d：\n%s\n\n", i+1, result.Content)
	}
	
	prompt += fmt.Sprintf("用户问题：%s\n\n", userQuery)
	prompt += "请根据以上信息回答用户问题。如果提供的信息不足以回答问题，请说明信息不足，不要编造信息。"
	
	return prompt
}

// ApplyRAG 应用RAG到Agent的对话
func ApplyRAG(ctx context.Context, kbID string, query string, agent *agent.Agent) (string, error) {
	// 先检索相关知识
	req := QueryRequest{
		KnowledgeBaseID: kbID,
		Query:           query,
		TopK:            5,
	}
	
	resp, err := DefaultRetriever.Retrieve(ctx, req)
	if err != nil {
		return "", fmt.Errorf("retrieval failed: %w", err)
	}
	
	// 如果没有找到相关内容，直接让Agent处理
	if len(resp.Results) == 0 {
		return agent.Chat(ctx, query)
	}
	
	// 生成增强提示词
	augmentedPrompt := GeneratePromptFromResults(resp.Results, query)
	
	// 让Agent处理增强后的提示词
	return agent.Chat(ctx, augmentedPrompt)
} 