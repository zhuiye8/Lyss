package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"gorm.io/gorm"
)

var (
	ErrLogNotFound      = errors.New("日志记录不存�?)
	ErrInvalidLogType   = errors.New("无效的日志类�?)
	ErrNoPermission     = errors.New("没有操作权限")
	ErrInvalidLogFormat = errors.New("无效的日志格�?)
)

// Service 提供日志查询和管理功�?
type Service struct {
	db *gorm.DB
}

// NewService 创建新的日志服务
func NewService(db *gorm.DB) *Service {
	return &Service{
		db: db,
	}
}

// LogQueryParams 日志查询参数
type LogQueryParams struct {
	Level       string    `form:"level"`
	Category    string    `form:"category"`
	StartTime   time.Time `form:"start_time"`
	EndTime     time.Time `form:"end_time"`
	UserID      string    `form:"user_id"`
	RequestID   string    `form:"request_id"`
	Method      string    `form:"method"`
	Path        string    `form:"path"`
	StatusCode  int       `form:"status_code"`
	MinDuration int64     `form:"min_duration"`
	MaxDuration int64     `form:"max_duration"`
	ErrorCode   string    `form:"error_code"`
	ModelName   string    `form:"model_name"`
	ProjectID   string    `form:"project_id"`
	AppID       string    `form:"app_id"`
	Page        int       `form:"page"`
	PageSize    int       `form:"page_size"`
	SortBy      string    `form:"sort_by"`
	SortOrder   string    `form:"sort_order"`
}

// LogType 表示查询的日志类�?
type LogType string

const (
	LogTypeAPI        LogType = "api"
	LogTypeError      LogType = "error"
	LogTypeModelCall  LogType = "model_call"
	LogTypeAll        LogType = "all"
)

// GetLogs 根据查询参数获取日志列表
func (s *Service) GetLogs(params LogQueryParams, logType LogType) ([]models.LogResponse, int64, error) {
	var responses []models.LogResponse
	var totalCount int64
	
	// 构建查询，根据日志类型选择�?
	var query *gorm.DB
	switch logType {
	case LogTypeAPI:
		query = s.db.Model(&models.APILog{})
	case LogTypeError:
		query = s.db.Model(&models.ErrorLog{})
	case LogTypeModelCall:
		query = s.db.Model(&models.ModelCallLog{})
	case LogTypeAll:
		query = s.db.Model(&models.Log{})
	default:
		return nil, 0, ErrInvalidLogType
	}
	
	// 添加过滤条件
	if params.Level != "" {
		query = query.Where("level = ?", params.Level)
	}
	
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}
	
	if !params.StartTime.IsZero() {
		query = query.Where("created_at >= ?", params.StartTime)
	}
	
	if !params.EndTime.IsZero() {
		query = query.Where("created_at <= ?", params.EndTime)
	}
	
	if params.UserID != "" {
		userID, err := uuid.Parse(params.UserID)
		if err == nil {
			query = query.Where("user_id = ?", userID)
		}
	}
	
	// 添加特定日志类型的过滤条�?
	if logType == LogTypeAPI || logType == LogTypeAll {
		if params.RequestID != "" {
			query = query.Where("request_id = ?", params.RequestID)
		}
		
		if params.Method != "" {
			query = query.Where("method = ?", params.Method)
		}
		
		if params.Path != "" {
			query = query.Where("path LIKE ?", fmt.Sprintf("%%%s%%", params.Path))
		}
		
		if params.StatusCode > 0 {
			query = query.Where("status_code = ?", params.StatusCode)
		}
		
		if params.MinDuration > 0 {
			query = query.Where("duration >= ?", params.MinDuration)
		}
		
		if params.MaxDuration > 0 {
			query = query.Where("duration <= ?", params.MaxDuration)
		}
	}
	
	if logType == LogTypeError || logType == LogTypeAll {
		if params.ErrorCode != "" {
			query = query.Where("error_code = ?", params.ErrorCode)
		}
	}
	
	if logType == LogTypeModelCall || logType == LogTypeAll {
		if params.ModelName != "" {
			query = query.Where("model_name = ?", params.ModelName)
		}
		
		if params.ProjectID != "" {
			projectID, err := uuid.Parse(params.ProjectID)
			if err == nil {
				query = query.Where("project_id = ?", projectID)
			}
		}
		
		if params.AppID != "" {
			appID, err := uuid.Parse(params.AppID)
			if err == nil {
				query = query.Where("application_id = ?", appID)
			}
		}
	}
	
	// 计算总数
	err := query.Count(&totalCount).Error
	if err != nil {
		return nil, 0, err
	}
	
	// 分页设置
	if params.Page <= 0 {
		params.Page = 1
	}
	
	if params.PageSize <= 0 {
		params.PageSize = 20
	}
	
	offset := (params.Page - 1) * params.PageSize
	
	// 排序
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}
	
	// 添加排序和分�?
	query = query.Order(fmt.Sprintf("%s %s", params.SortBy, params.SortOrder))
	query = query.Offset(offset).Limit(params.PageSize)
	
	// 基于日志类型获取数据
	switch logType {
	case LogTypeAPI:
		var logs []models.APILog
		if err := query.Find(&logs).Error; err != nil {
			return nil, 0, err
		}
		
		// 转换为响应格�?
		responses = make([]models.LogResponse, len(logs))
		for i, log := range logs {
			var metadata interface{}
			if log.Metadata != "" {
				if err := json.Unmarshal([]byte(log.Metadata), &metadata); err == nil {
					// 如果无法解析，保留原始字符串
					metadata = log.Metadata
				}
			}
			
			responses[i] = models.LogResponse{
				ID:         log.ID,
				Level:      log.Level,
				Category:   log.Category,
				Message:    log.Message,
				UserID:     log.UserID,
				Metadata:   metadata,
				CreatedAt:  log.CreatedAt,
				Method:     log.Method,
				Path:       log.Path,
				StatusCode: log.StatusCode,
				Duration:   log.Duration,
			}
		}
		
	case LogTypeError:
		var logs []models.ErrorLog
		if err := query.Find(&logs).Error; err != nil {
			return nil, 0, err
		}
		
		// 转换为响应格�?
		responses = make([]models.LogResponse, len(logs))
		for i, log := range logs {
			var metadata interface{}
			if log.Metadata != "" {
				if err := json.Unmarshal([]byte(log.Metadata), &metadata); err == nil {
					// 如果无法解析，保留原始字符串
					metadata = log.Metadata
				}
			}
			
			responses[i] = models.LogResponse{
				ID:         log.ID,
				Level:      log.Level,
				Category:   log.Category,
				Message:    log.Message,
				UserID:     log.UserID,
				Metadata:   metadata,
				CreatedAt:  log.CreatedAt,
				ErrorCode:  log.ErrorCode,
				ResolvedAt: log.ResolvedAt,
			}
		}
		
	case LogTypeModelCall:
		var logs []models.ModelCallLog
		if err := query.Find(&logs).Error; err != nil {
			return nil, 0, err
		}
		
		// 转换为响应格�?
		responses = make([]models.LogResponse, len(logs))
		for i, log := range logs {
			var metadata interface{}
			if log.Metadata != "" {
				if err := json.Unmarshal([]byte(log.Metadata), &metadata); err == nil {
					// 如果无法解析，保留原始字符串
					metadata = log.Metadata
				}
			}
			
			responses[i] = models.LogResponse{
				ID:         log.ID,
				Level:      log.Level,
				Category:   log.Category,
				Message:    log.Message,
				UserID:     log.UserID,
				Metadata:   metadata,
				CreatedAt:  log.CreatedAt,
				ModelName:  log.ModelName,
				TotalTokens: log.TotalTokens,
				Duration:   log.Duration,
			}
		}
		
	case LogTypeAll:
		var logs []models.Log
		if err := query.Find(&logs).Error; err != nil {
			return nil, 0, err
		}
		
		// 转换为响应格�?
		responses = make([]models.LogResponse, len(logs))
		for i, log := range logs {
			var metadata interface{}
			if log.Metadata != "" {
				if err := json.Unmarshal([]byte(log.Metadata), &metadata); err == nil {
					// 如果无法解析，保留原始字符串
					metadata = log.Metadata
				}
			}
			
			responses[i] = models.LogResponse{
				ID:        log.ID,
				Level:     log.Level,
				Category:  log.Category,
				Message:   log.Message,
				UserID:    log.UserID,
				Metadata:  metadata,
				CreatedAt: log.CreatedAt,
			}
		}
	}
	
	return responses, totalCount, nil
}

// GetLogByID 根据ID获取日志详情
func (s *Service) GetLogByID(id uuid.UUID) (*models.LogResponse, error) {
	// 首先尝试在基础日志表中查找
	var baseLog models.Log
	if err := s.db.Where("id = ?", id).First(&baseLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLogNotFound
		}
		return nil, err
	}
	
	// 解析元数�?
	var metadata interface{}
	if baseLog.Metadata != "" {
		if err := json.Unmarshal([]byte(baseLog.Metadata), &metadata); err == nil {
			// 如果无法解析，保留原始字符串
			metadata = baseLog.Metadata
		}
	}
	
	// 创建基本响应
	response := &models.LogResponse{
		ID:        baseLog.ID,
		Level:     baseLog.Level,
		Category:  baseLog.Category,
		Message:   baseLog.Message,
		UserID:    baseLog.UserID,
		Metadata:  metadata,
		CreatedAt: baseLog.CreatedAt,
	}
	
	// 根据类别获取更多详细信息
	switch baseLog.Category {
	case models.LogCategoryAPI:
		var apiLog models.APILog
		if err := s.db.Where("id = ?", id).First(&apiLog).Error; err == nil {
			response.Method = apiLog.Method
			response.Path = apiLog.Path
			response.StatusCode = apiLog.StatusCode
			response.Duration = apiLog.Duration
		}
		
	case models.LogCategoryAuth:
		// 可以添加认证日志特有字段
		
	case models.LogCategoryModel:
		var modelLog models.ModelCallLog
		if err := s.db.Where("id = ?", id).First(&modelLog).Error; err == nil {
			response.ModelName = modelLog.ModelName
			response.TotalTokens = modelLog.TotalTokens
			response.Duration = modelLog.Duration
		}
	}
	
	return response, nil
}

// MarkErrorAsResolved 将错误日志标记为已解�?
func (s *Service) MarkErrorAsResolved(id uuid.UUID, userID uuid.UUID) error {
	var errorLog models.ErrorLog
	if err := s.db.Where("id = ?", id).First(&errorLog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLogNotFound
		}
		return err
	}
	
	now := time.Now()
	errorLog.ResolvedAt = &now
	errorLog.ResolvedBy = &userID
	
	if err := s.db.Save(&errorLog).Error; err != nil {
		return err
	}
	
	return nil
}

// AddSystemLog 添加一条系统日�?
func (s *Service) AddSystemLog(level models.LogLevel, message string, metadata map[string]interface{}) error {
	metadataJSON, err := json.Marshal(metadata)
	if err != nil {
		return ErrInvalidLogFormat
	}
	
	log := models.Log{
		Level:     level,
		Category:  models.LogCategorySystem,
		Message:   message,
		Metadata:  string(metadataJSON),
		CreatedAt: time.Now(),
	}
	
	if err := s.db.Create(&log).Error; err != nil {
		return err
	}
	
	return nil
}

// GetLogStats 获取日志统计信息
func (s *Service) GetLogStats(startTime, endTime time.Time) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// 确保时间范围有效
	if endTime.IsZero() {
		endTime = time.Now()
	}
	
	if startTime.IsZero() {
		startTime = endTime.Add(-24 * time.Hour)
	}
	
	// 各级别日志数�?
	var levelCounts []struct {
		Level string `json:"level"`
		Count int64  `json:"count"`
	}
	if err := s.db.Model(&models.Log{}).
		Select("level, count(*) as count").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("level").
		Find(&levelCounts).Error; err != nil {
		return nil, err
	}
	stats["level_counts"] = levelCounts
	
	// 各类别日志数�?
	var categoryCounts []struct {
		Category string `json:"category"`
		Count    int64  `json:"count"`
	}
	if err := s.db.Model(&models.Log{}).
		Select("category, count(*) as count").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Group("category").
		Find(&categoryCounts).Error; err != nil {
		return nil, err
	}
	stats["category_counts"] = categoryCounts
	
	// API请求统计
	var apiStats struct {
		TotalRequests int64   `json:"total_requests"`
		ErrorCount    int64   `json:"error_count"`
		AvgDuration   float64 `json:"avg_duration"`
	}
	if err := s.db.Model(&models.APILog{}).
		Select("count(*) as total_requests, sum(case when status_code >= 400 then 1 else 0 end) as error_count, avg(duration) as avg_duration").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Find(&apiStats).Error; err != nil {
		return nil, err
	}
	stats["api_stats"] = apiStats
	
	// 错误日志统计
	var errorStats struct {
		TotalErrors        int64 `json:"total_errors"`
		UnresolvedErrors   int64 `json:"unresolved_errors"`
		ResolvedErrors     int64 `json:"resolved_errors"`
	}
	if err := s.db.Model(&models.ErrorLog{}).
		Select("count(*) as total_errors, sum(case when resolved_at is null then 1 else 0 end) as unresolved_errors, sum(case when resolved_at is not null then 1 else 0 end) as resolved_errors").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Find(&errorStats).Error; err != nil {
		return nil, err
	}
	stats["error_stats"] = errorStats
	
	// 模型调用统计
	var modelStats struct {
		TotalCalls     int64   `json:"total_calls"`
		SuccessRate    float64 `json:"success_rate"`
		AvgDuration    float64 `json:"avg_duration"`
		TotalTokens    int64   `json:"total_tokens"`
		AvgTokens      float64 `json:"avg_tokens"`
	}
	if err := s.db.Model(&models.ModelCallLog{}).
		Select("count(*) as total_calls, avg(case when success then 1 else 0 end) * 100 as success_rate, avg(duration) as avg_duration, sum(total_tokens) as total_tokens, avg(total_tokens) as avg_tokens").
		Where("created_at BETWEEN ? AND ?", startTime, endTime).
		Find(&modelStats).Error; err != nil {
		return nil, err
	}
	stats["model_stats"] = modelStats
	
	return stats, nil
} 
