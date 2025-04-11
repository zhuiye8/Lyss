package kb

// 此文件导出包中的重要类型和函数

import (
	"context"
)

// 初始化知识库模块
func Initialize(ctx context.Context, vectorDBConfig VectorDBConfig) error {
	// 初始化向量数据库
	if err := InitDefaultVectorDB(vectorDBConfig); err != nil {
		return err
	}
	
	// 初始化知识库管理器
	InitDefaultKnowledgeBaseManager(DefaultVectorDB)
	
	// 初始化检索器
	InitDefaultRetriever()
	
	return nil
}

// 获取默认知识库管理器
func GetKnowledgeBaseManager() *KnowledgeBaseManager {
	return DefaultKnowledgeBaseManager
}

// 获取默认检索器
func GetRetriever() Retriever {
	return DefaultRetriever
}

// 获取默认向量数据库
func GetVectorDB() VectorDatabase {
	return DefaultVectorDB
}

// 获取默认向量化管理器
func GetEmbeddingManager() *EmbeddingManager {
	return DefaultEmbeddingManager
}

// 获取默认文档处理器注册表
func GetDocumentProcessorRegistry() *DocumentProcessorRegistry {
	return DefaultProcessorRegistry
}

// CreateInMemoryTestDB 创建内存测试数据库
func CreateInMemoryTestDB() VectorDatabase {
	return NewInMemoryVectorDB()
} 
