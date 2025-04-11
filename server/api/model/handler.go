package model

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"github.com/zhuiye8/Lyss/server/pkg/auth"
	"github.com/zhuiye8/Lyss/server/pkg/middleware"
)

// ModelResponse 模型API响应
type ModelResponse struct {
	ID          uuid.UUID          `json:"id"`
	Name        string             `json:"name"`
	Provider    models.ModelProvider `json:"provider"`
	ModelID     string             `json:"model_id"`
	Type        models.ModelType    `json:"type"`
	Description string             `json:"description"`
	Capabilities []string           `json:"capabilities"`
	Parameters  models.ModelParameters `json:"parameters"`
	MaxTokens   int                `json:"max_tokens"`
	TokenCost   struct {
		Prompt     float64 `json:"prompt"`
		Completion float64 `json:"completion"`
	} `json:"token_cost"`
	Status      models.ModelStatus `json:"status"`
	IsSystem    bool               `json:"is_system"`
	CreatedAt   time.Time          `json:"created_at"`
}

// ModelConfigResponse 是返回给客户端的模型配置结构
type ModelConfigResponse struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Model          ModelResponse   `json:"model"`
	Parameters     models.ModelParameters `json:"parameters"`
	IsShared       bool            `json:"is_shared"`
	UsageMetrics   models.ModelUsageMetrics `json:"usage_metrics"`
	OrganizationID uuid.UUID       `json:"organization_id"`
	CreatedBy      uuid.UUID       `json:"created_by"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

// Handler 处理模型相关的请求
type Handler struct {
	service        *Service
	authMiddleware *middleware.AuthMiddleware
}

// NewHandler 创建模型处理器
func NewHandler(service *Service, authMiddleware *middleware.AuthMiddleware) *Handler {
	return &Handler{
		service:        service,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes 注册API路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	modelGroup := router.Group("/models")
	{
		// 模型管理路由
		modelGroup.GET("", h.authMiddleware.Authenticate(), h.GetModels)
		modelGroup.GET("/:id", h.authMiddleware.Authenticate(), h.GetModel)
		modelGroup.POST("", h.authMiddleware.Authenticate(), h.CreateModel)
		modelGroup.PUT("/:id", h.authMiddleware.Authenticate(), h.UpdateModel)
		modelGroup.DELETE("/:id", h.authMiddleware.Authenticate(), h.DeleteModel)
		
		// 模型配置路由
		modelGroup.GET("/configs", h.authMiddleware.Authenticate(), h.GetModelConfigs)
		modelGroup.GET("/configs/:id", h.authMiddleware.Authenticate(), h.GetModelConfig)
		modelGroup.POST("/configs", h.authMiddleware.Authenticate(), h.CreateModelConfig)
		modelGroup.PUT("/configs/:id", h.authMiddleware.Authenticate(), h.UpdateModelConfig)
		modelGroup.DELETE("/configs/:id", h.authMiddleware.Authenticate(), h.DeleteModelConfig)
		
		// 模型提供者路由
		modelGroup.GET("/providers", h.authMiddleware.Authenticate(), h.GetProviders)
		
		// 测试模型连接
		modelGroup.POST("/test-connection", h.authMiddleware.Authenticate(), h.TestModelConnection)
	}
}

// GetModels 获取模型列表
func (h *Handler) GetModels(c *gin.Context) {
	// 解析查询参数
	params := ModelsQueryParams{
		Provider:      c.Query("provider"),
		Type:          c.Query("type"),
		Status:        c.Query("status"),
		IncludeSystem: c.Query("include_system") == "true",
		Search:        c.Query("search"),
	}
	
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	params.Page = page
	params.PageSize = pageSize
	
	// 解析排序参数
	params.SortBy = c.DefaultQuery("sort_by", "name")
	params.SortDesc = c.DefaultQuery("sort_desc", "false") == "true"
	
	// 解析组织ID
	orgIDStr := c.Query("organization_id")
	if orgIDStr != "" {
		orgID, err := uuid.Parse(orgIDStr)
		if err == nil {
			params.OrganizationID = &orgID
		}
	}
	
	// 获取模型列表
	models, total, err := h.service.GetModels(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取模型失败: " + err.Error()})
		return
	}
	
	// 构建响应
	var responses []ModelResponse
	for _, model := range models {
		responses = append(responses, ModelResponse{
			ID:           model.ID,
			Name:         model.Name,
			Provider:     model.Provider,
			ModelID:      model.ModelID,
			Type:         model.Type,
			Description:  model.Description,
			Capabilities: model.Capabilities,
			Parameters:   model.Parameters,
			MaxTokens:    model.MaxTokens,
			TokenCost: struct {
				Prompt     float64 `json:"prompt"`
				Completion float64 `json:"completion"`
			}{
				Prompt:     model.TokenCostPrompt,
				Completion: model.TokenCostCompl,
			},
			Status:    model.Status,
			IsSystem:  model.IsSystem,
			CreatedAt: model.CreatedAt,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"meta": gin.H{
			"total":     total,
			"page":      params.Page,
			"page_size": params.PageSize,
		},
	})
}

// GetModel 获取单个模型
func (h *Handler) GetModel(c *gin.Context) {
	// 解析模型ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模型ID"})
		return
	}
	
	// 获取模型
	model, err := h.service.GetModelByID(id)
	if err != nil {
		if err == ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "模型不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取模型失败: " + err.Error()})
		}
		return
	}
	
	// 构建响应
	response := ModelResponse{
		ID:           model.ID,
		Name:         model.Name,
		Provider:     model.Provider,
		ModelID:      model.ModelID,
		Type:         model.Type,
		Description:  model.Description,
		Capabilities: model.Capabilities,
		Parameters:   model.Parameters,
		MaxTokens:    model.MaxTokens,
		TokenCost: struct {
			Prompt     float64 `json:"prompt"`
			Completion float64 `json:"completion"`
		}{
			Prompt:     model.TokenCostPrompt,
			Completion: model.TokenCostCompl,
		},
		Status:    model.Status,
		IsSystem:  model.IsSystem,
		CreatedAt: model.CreatedAt,
	}
	
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// CreateModelRequest 创建模型请求
type CreateModelRequest struct {
	Name            string                  `json:"name" binding:"required"`
	Provider        models.ModelProvider    `json:"provider" binding:"required"`
	ModelID         string                  `json:"model_id" binding:"required"`
	Type            models.ModelType        `json:"type" binding:"required"`
	Description     string                  `json:"description"`
	Capabilities    []string                `json:"capabilities"`
	Parameters      models.ModelParameters  `json:"parameters"`
	MaxTokens       int                     `json:"max_tokens"`
	TokenCostPrompt float64                 `json:"token_cost_prompt"`
	TokenCostCompl  float64                 `json:"token_cost_completion"`
	ProviderConfig  models.ModelProviderConfig `json:"provider_config"`
}

// CreateModel 创建新模型
func (h *Handler) CreateModel(c *gin.Context) {
	// 仅管理员可创建自定义模型
	if !auth.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以创建模型"})
		return
	}
	
	// 解析请求
	var req CreateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	
	// 获取用户信息
	userID := auth.GetUserIDFromContext(c)
	orgID := auth.GetOrgIDFromContext(c)
	
	// 创建模型实体
	model := models.Model{
		Name:            req.Name,
		Provider:        req.Provider,
		ModelID:         req.ModelID,
		Type:            req.Type,
		Description:     req.Description,
		Capabilities:    req.Capabilities,
		Parameters:      req.Parameters,
		MaxTokens:       req.MaxTokens,
		TokenCostPrompt: req.TokenCostPrompt,
		TokenCostCompl:  req.TokenCostCompl,
		Status:          models.ModelStatusActive,
		ProviderConfig:  req.ProviderConfig,
		IsSystem:        false,
		IsCustom:        true,
		OrganizationID:  &orgID,
		CreatedBy:       &userID,
	}
	
	// 创建模型
	if err := h.service.CreateModel(&model); err != nil {
		if err == ErrDuplicateModelName {
			c.JSON(http.StatusConflict, gin.H{"error": "模型名称已存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建模型失败: " + err.Error()})
		}
		return
	}
	
	// 构建响应
	response := ModelResponse{
		ID:           model.ID,
		Name:         model.Name,
		Provider:     model.Provider,
		ModelID:      model.ModelID,
		Type:         model.Type,
		Description:  model.Description,
		Capabilities: model.Capabilities,
		Parameters:   model.Parameters,
		MaxTokens:    model.MaxTokens,
		TokenCost: struct {
			Prompt     float64 `json:"prompt"`
			Completion float64 `json:"completion"`
		}{
			Prompt:     model.TokenCostPrompt,
			Completion: model.TokenCostCompl,
		},
		Status:    model.Status,
		IsSystem:  model.IsSystem,
		CreatedAt: model.CreatedAt,
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// UpdateModelRequest 更新模型请求
type UpdateModelRequest struct {
	Name            string                  `json:"name"`
	ModelID         string                  `json:"model_id"`
	Description     string                  `json:"description"`
	Capabilities    []string                `json:"capabilities"`
	Parameters      models.ModelParameters  `json:"parameters"`
	MaxTokens       int                     `json:"max_tokens"`
	TokenCostPrompt float64                 `json:"token_cost_prompt"`
	TokenCostCompl  float64                 `json:"token_cost_completion"`
	Status          models.ModelStatus      `json:"status"`
	ProviderConfig  models.ModelProviderConfig `json:"provider_config"`
}

// UpdateModel 更新模型
func (h *Handler) UpdateModel(c *gin.Context) {
	// 仅管理员可更新模型
	if !auth.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以更新模型"})
		return
	}
	
	// 解析模型ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模型ID"})
		return
	}
	
	// 解析请求
	var req UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	
	// 创建更新数据
	updateData := models.Model{
		Name:            req.Name,
		ModelID:         req.ModelID,
		Description:     req.Description,
		Capabilities:    req.Capabilities,
		Parameters:      req.Parameters,
		MaxTokens:       req.MaxTokens,
		TokenCostPrompt: req.TokenCostPrompt,
		TokenCostCompl:  req.TokenCostCompl,
		Status:          req.Status,
		ProviderConfig:  req.ProviderConfig,
	}
	
	// 更新模型
	if err := h.service.UpdateModel(id, &updateData); err != nil {
		if err == ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "模型不存在"})
		} else if err == ErrDuplicateModelName {
			c.JSON(http.StatusConflict, gin.H{"error": "模型名称已存在"})
		} else if err == ErrNoPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限更新此模型"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新模型失败: " + err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "模型更新成功"})
}

// DeleteModel 删除模型
func (h *Handler) DeleteModel(c *gin.Context) {
	// 仅管理员可删除模型
	if !auth.IsAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{"error": "只有管理员可以删除模型"})
		return
	}
	
	// 解析模型ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的模型ID"})
		return
	}
	
	// 删除模型
	if err := h.service.DeleteModel(id); err != nil {
		if err == ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "模型不存在"})
		} else if err == ErrNoPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "无法删除系统模型"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除模型失败: " + err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "模型删除成功"})
}

// GetModelConfigs 获取模型配置列表
func (h *Handler) GetModelConfigs(c *gin.Context) {
	// 解析查询参数
	params := ConfigsQueryParams{
		Search: c.Query("search"),
		Shared: c.Query("shared") == "true",
	}
	
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	params.Page = page
	params.PageSize = pageSize
	
	// 解析排序参数
	params.SortBy = c.DefaultQuery("sort_by", "name")
	params.SortDesc = c.DefaultQuery("sort_desc", "false") == "true"
	
	// 解析模型ID
	modelIDStr := c.Query("model_id")
	if modelIDStr != "" {
		modelID, err := uuid.Parse(modelIDStr)
		if err == nil {
			params.ModelID = &modelID
		}
	}
	
	// 获取用户信息
	userID := auth.GetUserIDFromContext(c)
	orgID := auth.GetOrgIDFromContext(c)
	
	// 设置过滤条件
	params.OrganizationID = &orgID
	
	// 如果不是管理员，只能看到自己的配置和共享配置
	if !auth.IsAdmin(c) {
		params.CreatedBy = &userID
	}
	
	// 获取配置列表
	configs, total, err := h.service.GetModelConfigs(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取模型配置失败: " + err.Error()})
		return
	}
	
	// 构建响应
	var responses []ModelConfigResponse
	for _, config := range configs {
		modelResponse := ModelResponse{
			ID:           config.Model.ID,
			Name:         config.Model.Name,
			Provider:     config.Model.Provider,
			ModelID:      config.Model.ModelID,
			Type:         config.Model.Type,
			Description:  config.Model.Description,
			Capabilities: config.Model.Capabilities,
			Parameters:   config.Model.Parameters,
			MaxTokens:    config.Model.MaxTokens,
			TokenCost: struct {
				Prompt     float64 `json:"prompt"`
				Completion float64 `json:"completion"`
			}{
				Prompt:     config.Model.TokenCostPrompt,
				Completion: config.Model.TokenCostCompl,
			},
			Status:    config.Model.Status,
			IsSystem:  config.Model.IsSystem,
			CreatedAt: config.Model.CreatedAt,
		}
		
		responses = append(responses, ModelConfigResponse{
			ID:             config.ID,
			Name:           config.Name,
			Description:    config.Description,
			Model:          modelResponse,
			Parameters:     config.Parameters,
			IsShared:       config.IsShared,
			UsageMetrics:   config.UsageMetrics,
			OrganizationID: config.OrganizationID,
			CreatedBy:      config.CreatedBy,
			CreatedAt:      config.CreatedAt,
			UpdatedAt:      config.UpdatedAt,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": responses,
		"meta": gin.H{
			"total":     total,
			"page":      params.Page,
			"page_size": params.PageSize,
		},
	})
}

// GetModelConfig 获取单个模型配置
func (h *Handler) GetModelConfig(c *gin.Context) {
	// 解析配置ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}
	
	// 获取配置
	config, err := h.service.GetModelConfigByID(id)
	if err != nil {
		if err == ErrModelConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败: " + err.Error()})
		}
		return
	}
	
	// 检查访问权限
	userID := auth.GetUserIDFromContext(c)
	orgID := auth.GetOrgIDFromContext(c)
	
	if config.OrganizationID != orgID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问此配置"})
		return
	}
	
	// 非管理员只能查看自己的或共享的配置
	if !auth.IsAdmin(c) && config.CreatedBy != userID && !config.IsShared {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限访问此配置"})
		return
	}
	
	// 构建响应
	modelResponse := ModelResponse{
		ID:           config.Model.ID,
		Name:         config.Model.Name,
		Provider:     config.Model.Provider,
		ModelID:      config.Model.ModelID,
		Type:         config.Model.Type,
		Description:  config.Model.Description,
		Capabilities: config.Model.Capabilities,
		Parameters:   config.Model.Parameters,
		MaxTokens:    config.Model.MaxTokens,
		TokenCost: struct {
			Prompt     float64 `json:"prompt"`
			Completion float64 `json:"completion"`
		}{
			Prompt:     config.Model.TokenCostPrompt,
			Completion: config.Model.TokenCostCompl,
		},
		Status:    config.Model.Status,
		IsSystem:  config.Model.IsSystem,
		CreatedAt: config.Model.CreatedAt,
	}
	
	response := ModelConfigResponse{
		ID:             config.ID,
		Name:           config.Name,
		Description:    config.Description,
		Model:          modelResponse,
		Parameters:     config.Parameters,
		IsShared:       config.IsShared,
		UsageMetrics:   config.UsageMetrics,
		OrganizationID: config.OrganizationID,
		CreatedBy:      config.CreatedBy,
		CreatedAt:      config.CreatedAt,
		UpdatedAt:      config.UpdatedAt,
	}
	
	c.JSON(http.StatusOK, gin.H{"data": response})
}

// CreateModelConfigRequest 创建模型配置请求
type CreateModelConfigRequest struct {
	Name           string                  `json:"name" binding:"required"`
	Description    string                  `json:"description"`
	ModelID        uuid.UUID               `json:"model_id" binding:"required"`
	Parameters     models.ModelParameters  `json:"parameters"`
	ProviderConfig models.ModelProviderConfig `json:"provider_config"`
	IsShared       bool                    `json:"is_shared"`
}

// CreateModelConfig 创建新的模型配置
func (h *Handler) CreateModelConfig(c *gin.Context) {
	// 解析请求
	var req CreateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	
	// 获取用户信息
	userID := auth.GetUserIDFromContext(c)
	orgID := auth.GetOrgIDFromContext(c)
	
	// 创建配置实体
	config := models.ModelConfig{
		Name:           req.Name,
		Description:    req.Description,
		ModelID:        req.ModelID,
		Parameters:     req.Parameters,
		ProviderConfig: req.ProviderConfig,
		IsShared:       req.IsShared,
		OrganizationID: orgID,
		CreatedBy:      userID,
	}
	
	// 创建配置
	if err := h.service.CreateModelConfig(&config); err != nil {
		if err == ErrModelNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "模型不存在"})
		} else if err.Error() == "配置名称已存在" {
			c.JSON(http.StatusConflict, gin.H{"error": "配置名称已存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建配置失败: " + err.Error()})
		}
		return
	}
	
	// 获取完整配置，包括模型信息
	fullConfig, err := h.service.GetModelConfigByID(config.ID)
	if err != nil {
		c.JSON(http.StatusCreated, gin.H{
			"message": "配置创建成功，但无法获取完整信息",
			"data": gin.H{
				"id": config.ID,
			},
		})
		return
	}
	
	// 构建响应
	modelResponse := ModelResponse{
		ID:           fullConfig.Model.ID,
		Name:         fullConfig.Model.Name,
		Provider:     fullConfig.Model.Provider,
		ModelID:      fullConfig.Model.ModelID,
		Type:         fullConfig.Model.Type,
		Description:  fullConfig.Model.Description,
		Capabilities: fullConfig.Model.Capabilities,
		Parameters:   fullConfig.Model.Parameters,
		MaxTokens:    fullConfig.Model.MaxTokens,
		TokenCost: struct {
			Prompt     float64 `json:"prompt"`
			Completion float64 `json:"completion"`
		}{
			Prompt:     fullConfig.Model.TokenCostPrompt,
			Completion: fullConfig.Model.TokenCostCompl,
		},
		Status:    fullConfig.Model.Status,
		IsSystem:  fullConfig.Model.IsSystem,
		CreatedAt: fullConfig.Model.CreatedAt,
	}
	
	response := ModelConfigResponse{
		ID:             fullConfig.ID,
		Name:           fullConfig.Name,
		Description:    fullConfig.Description,
		Model:          modelResponse,
		Parameters:     fullConfig.Parameters,
		IsShared:       fullConfig.IsShared,
		UsageMetrics:   fullConfig.UsageMetrics,
		OrganizationID: fullConfig.OrganizationID,
		CreatedBy:      fullConfig.CreatedBy,
		CreatedAt:      fullConfig.CreatedAt,
		UpdatedAt:      fullConfig.UpdatedAt,
	}
	
	c.JSON(http.StatusCreated, gin.H{"data": response})
}

// UpdateModelConfigRequest 更新模型配置请求
type UpdateModelConfigRequest struct {
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	ModelID        uuid.UUID               `json:"model_id"`
	Parameters     models.ModelParameters  `json:"parameters"`
	ProviderConfig models.ModelProviderConfig `json:"provider_config"`
	IsShared       bool                    `json:"is_shared"`
}

// UpdateModelConfig 更新模型配置
func (h *Handler) UpdateModelConfig(c *gin.Context) {
	// 解析配置ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}
	
	// 解析请求
	var req UpdateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	
	// 获取用户信息
	userID := auth.GetUserIDFromContext(c)
	
	// 创建更新数据
	updateData := models.ModelConfig{
		Name:           req.Name,
		Description:    req.Description,
		ModelID:        req.ModelID,
		Parameters:     req.Parameters,
		ProviderConfig: req.ProviderConfig,
		IsShared:       req.IsShared,
	}
	
	// 更新配置
	if err := h.service.UpdateModelConfig(id, &updateData, userID); err != nil {
		if err == ErrModelConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		} else if err == ErrNoPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限更新此配置"})
		} else if err.Error() == "配置名称已存在" {
			c.JSON(http.StatusConflict, gin.H{"error": "配置名称已存在"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败: " + err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "配置更新成功"})
}

// DeleteModelConfig 删除模型配置
func (h *Handler) DeleteModelConfig(c *gin.Context) {
	// 解析配置ID
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的配置ID"})
		return
	}
	
	// 获取用户信息
	userID := auth.GetUserIDFromContext(c)
	
	// 删除配置
	if err := h.service.DeleteModelConfig(id, userID); err != nil {
		if err == ErrModelConfigNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "配置不存在"})
		} else if err == ErrNoPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此配置"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除配置失败: " + err.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "配置删除成功"})
}

// ProviderInfo 提供者信息
type ProviderInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Website     string `json:"website"`
	DocURL      string `json:"doc_url"`
}

// GetProviders 获取模型提供者列表
func (h *Handler) GetProviders(c *gin.Context) {
	providers := []ProviderInfo{
		{
			ID:          string(models.ModelProviderOpenAI),
			Name:        "OpenAI",
			Description: "OpenAI提供的GPT系列模型和DALL-E图像模型",
			Website:     "https://openai.com/",
			DocURL:      "https://platform.openai.com/docs/",
		},
		{
			ID:          string(models.ModelProviderAnthropic),
			Name:        "Anthropic",
			Description: "Anthropic提供的Claude系列大语言模型",
			Website:     "https://www.anthropic.com/",
			DocURL:      "https://docs.anthropic.com/",
		},
		{
			ID:          string(models.ModelProviderBaidu),
			Name:        "百度",
			Description: "百度文心一言大语言模型",
			Website:     "https://cloud.baidu.com/",
			DocURL:      "https://cloud.baidu.com/doc/WENXINWORKSHOP/index.html",
		},
		{
			ID:          string(models.ModelProviderAli),
			Name:        "阿里",
			Description: "阿里云通义系列模型",
			Website:     "https://www.aliyun.com/",
			DocURL:      "https://help.aliyun.com/document_detail/2400395.html",
		},
		{
			ID:          string(models.ModelProviderLocal),
			Name:        "本地模型",
			Description: "本地部署的开源模型",
			Website:     "",
			DocURL:      "",
		},
		{
			ID:          string(models.ModelProviderCustom),
			Name:        "自定义API",
			Description: "自定义的模型API接口",
			Website:     "",
			DocURL:      "",
		},
	}
	
	c.JSON(http.StatusOK, gin.H{"data": providers})
}

// TestConnectionRequest 测试连接请求
type TestConnectionRequest struct {
	Provider       models.ModelProvider       `json:"provider" binding:"required"`
	ProviderConfig models.ModelProviderConfig `json:"provider_config" binding:"required"`
}

// TestModelConnection 测试模型连接
func (h *Handler) TestModelConnection(c *gin.Context) {
	// 解析请求
	var req TestConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据: " + err.Error()})
		return
	}
	
	// TODO: 实现实际的连接测试逻辑，通常需要调用各个提供者的API测试连接
	// 这里仅作为示例，始终返回成功
	
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "连接测试成功",
	})
} 
