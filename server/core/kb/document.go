package kb

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DocumentType 表示文档类型
type DocumentType string

const (
	// TypeText 纯文本文档
	TypeText DocumentType = "text"
	// TypeMarkdown Markdown文档
	TypeMarkdown DocumentType = "markdown"
	// TypePDF PDF文档
	TypePDF DocumentType = "pdf"
	// TypeWord Word文档
	TypeWord DocumentType = "docx"
	// TypeExcel Excel文档
	TypeExcel DocumentType = "xlsx"
	// TypeHTML HTML文档
	TypeHTML DocumentType = "html"
)

// Document 表示一个文档
type Document struct {
	ID             string       `json:"id"`
	KnowledgeBaseID string       `json:"knowledge_base_id"`
	Name           string       `json:"name"`
	Type           DocumentType `json:"type"`
	Size           int64        `json:"size"`
	Content        string       `json:"-"` // 原始内容，不在JSON中返回
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
	Metadata       interface{}  `json:"metadata,omitempty"`
}

// Chunk 表示文档的一个分块
type Chunk struct {
	ID         string    `json:"id"`
	DocumentID string    `json:"document_id"`
	Content    string    `json:"content"`
	Vector     []float32 `json:"-"`
	Metadata   struct {
		ChunkIndex int    `json:"chunk_index"`
		PageNumber int    `json:"page_number,omitempty"`
		Source     string `json:"source"`
	} `json:"metadata"`
}

// DocumentProcessor 文档处理器接口
type DocumentProcessor interface {
	// Process 处理文档并返回分块结果
	Process(doc *Document) ([]Chunk, error)
	// SupportsType 检查是否支持特定文档类型
	SupportsType(docType DocumentType) bool
}

// DocumentProcessorRegistry 文档处理器注册表
type DocumentProcessorRegistry struct {
	processors map[DocumentType]DocumentProcessor
}

// NewDocumentProcessorRegistry 创建新的文档处理器注册表
func NewDocumentProcessorRegistry() *DocumentProcessorRegistry {
	return &DocumentProcessorRegistry{
		processors: make(map[DocumentType]DocumentProcessor),
	}
}

// Register 注册文档处理器
func (r *DocumentProcessorRegistry) Register(docType DocumentType, processor DocumentProcessor) {
	r.processors[docType] = processor
}

// GetProcessor 获取指定类型的处理器
func (r *DocumentProcessorRegistry) GetProcessor(docType DocumentType) (DocumentProcessor, error) {
	processor, exists := r.processors[docType]
	if !exists {
		return nil, errors.New("no processor found for document type: " + string(docType))
	}
	return processor, nil
}

// ProcessDocument 处理文档
func (r *DocumentProcessorRegistry) ProcessDocument(doc *Document) ([]Chunk, error) {
	processor, err := r.GetProcessor(doc.Type)
	if err != nil {
		return nil, err
	}
	return processor.Process(doc)
}

// NewDocument 从文件创建新文档
func NewDocument(knowledgeBaseID string, file *multipart.FileHeader) (*Document, error) {
	// 打开文件
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	
	// 读取文件内容
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}
	
	// 获取文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	docType := getDocumentTypeFromExt(ext)
	
	now := time.Now()
	doc := &Document{
		ID:             uuid.New().String(),
		KnowledgeBaseID: knowledgeBaseID,
		Name:           file.Filename,
		Type:           docType,
		Size:           file.Size,
		Content:        buf.String(),
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	
	return doc, nil
}

// NewTextDocument 从文本创建新文档
func NewTextDocument(knowledgeBaseID, name, content string, docType DocumentType) *Document {
	now := time.Now()
	return &Document{
		ID:             uuid.New().String(),
		KnowledgeBaseID: knowledgeBaseID,
		Name:           name,
		Type:           docType,
		Size:           int64(len(content)),
		Content:        content,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// 从文件扩展名获取文档类型
func getDocumentTypeFromExt(ext string) DocumentType {
	switch ext {
	case ".txt":
		return TypeText
	case ".md", ".markdown":
		return TypeMarkdown
	case ".pdf":
		return TypePDF
	case ".doc", ".docx":
		return TypeWord
	case ".xls", ".xlsx":
		return TypeExcel
	case ".html", ".htm":
		return TypeHTML
	default:
		return TypeText
	}
}

// DefaultChunkSize 默认的文本分块大小
const DefaultChunkSize = 1000

// DefaultChunkOverlap 默认的文本分块重叠大小
const DefaultChunkOverlap = 200

// BasicTextProcessor 基础文本处理器
type BasicTextProcessor struct {
	ChunkSize    int
	ChunkOverlap int
}

// NewBasicTextProcessor 创建基础文本处理器
func NewBasicTextProcessor(chunkSize, chunkOverlap int) *BasicTextProcessor {
	if chunkSize <= 0 {
		chunkSize = DefaultChunkSize
	}
	if chunkOverlap <= 0 {
		chunkOverlap = DefaultChunkOverlap
	}
	return &BasicTextProcessor{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
	}
}

// Process 处理文档并返回分块
func (p *BasicTextProcessor) Process(doc *Document) ([]Chunk, error) {
	content := doc.Content
	
	// 如果文档为空，返回空结果
	if len(content) == 0 {
		return []Chunk{}, nil
	}
	
	// 简单按照句子或段落分割
	paragraphs := strings.Split(content, "\n\n")
	var chunks []Chunk
	
	currentChunk := ""
	chunkIndex := 0
	
	for _, para := range paragraphs {
		if len(currentChunk)+len(para) > p.ChunkSize {
			// 当前块已经足够大，创建一个新块
			if len(currentChunk) > 0 {
				chunks = append(chunks, Chunk{
					ID:         uuid.New().String(),
					DocumentID: doc.ID,
					Content:    currentChunk,
					Metadata: struct {
						ChunkIndex int    `json:"chunk_index"`
						PageNumber int    `json:"page_number,omitempty"`
						Source     string `json:"source"`
					}{
						ChunkIndex: chunkIndex,
						Source:     doc.Name,
					},
				})
				chunkIndex++
				
				// 使用重叠部分创建下一个块
				lastWords := getLastWords(currentChunk, p.ChunkOverlap)
				currentChunk = lastWords
			}
		}
		
		if len(currentChunk) > 0 {
			currentChunk += "\n\n"
		}
		currentChunk += para
	}
	
	// 添加最后一个块
	if len(currentChunk) > 0 {
		chunks = append(chunks, Chunk{
			ID:         uuid.New().String(),
			DocumentID: doc.ID,
			Content:    currentChunk,
			Metadata: struct {
				ChunkIndex int    `json:"chunk_index"`
				PageNumber int    `json:"page_number,omitempty"`
				Source     string `json:"source"`
			}{
				ChunkIndex: chunkIndex,
				Source:     doc.Name,
			},
		})
	}
	
	return chunks, nil
}

// SupportsType 检查是否支持特定文档类型
func (p *BasicTextProcessor) SupportsType(docType DocumentType) bool {
	return docType == TypeText || docType == TypeMarkdown
}

// 获取文本的最后n个字符，确保不会拆分单词
func getLastWords(text string, n int) string {
	if len(text) <= n {
		return text
	}
	
	// 查找适当的断点
	cutIndex := len(text) - n
	for i := cutIndex; i < len(text); i++ {
		if text[i] == ' ' || text[i] == '\n' {
			cutIndex = i + 1
			break
		}
	}
	
	return text[cutIndex:]
}

// DefaultProcessorRegistry 默认的文档处理器注册表
var DefaultProcessorRegistry = NewDocumentProcessorRegistry()

// 初始化默认处理器
func init() {
	basicProcessor := NewBasicTextProcessor(DefaultChunkSize, DefaultChunkOverlap)
	DefaultProcessorRegistry.Register(TypeText, basicProcessor)
	DefaultProcessorRegistry.Register(TypeMarkdown, basicProcessor)
} 