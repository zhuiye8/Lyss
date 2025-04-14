package response

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// StandardResponse 标准响应结构
type StandardResponse struct {
	Success   bool        `json:"success"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	RequestId string      `json:"requestId"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// PaginationData 分页数据结构
type PaginationData struct {
	List       interface{} `json:"list"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination 分页信息
type Pagination struct {
	Current  int64 `json:"current"`
	PageSize int64 `json:"pageSize"`
	Total    int64 `json:"total"`
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, StandardResponse{
		Success:   true,
		Code:      "200",
		Message:   "操作成功",
		RequestId: GenRequestId(c),
		Data:      data,
		Timestamp: GetTimestamp(),
	})
}

// SuccessWithPagination 返回带分页的成功响应
func SuccessWithPagination(c *gin.Context, list interface{}, current, pageSize, total int64) {
	paginationData := PaginationData{
		List: list,
		Pagination: Pagination{
			Current:  current,
			PageSize: pageSize,
			Total:    total,
		},
	}

	c.JSON(200, StandardResponse{
		Success:   true,
		Code:      "200",
		Message:   "操作成功",
		RequestId: GenRequestId(c),
		Data:      paginationData,
		Timestamp: GetTimestamp(),
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code string, message string, data ...interface{}) {
	var errorData interface{}
	if len(data) > 0 {
		errorData = data[0]
	}

	c.JSON(getHTTPStatus(code), StandardResponse{
		Success:   false,
		Code:      code,
		Message:   message,
		RequestId: GenRequestId(c),
		Data:      errorData,
		Timestamp: GetTimestamp(),
	})
}

// GenRequestId 获取请求ID
func GenRequestId(c *gin.Context) string {
	// 尝试从上下文获取请求ID
	if requestId, exists := c.Get("requestId"); exists {
		return requestId.(string)
	}
	
	// 如果不存在，生成新的
	requestId := uuid.New().String()
	c.Set("requestId", requestId)
	return requestId
}

// GetTimestamp 获取当前时间戳（毫秒）
func GetTimestamp() int64 {
	return time.Now().UnixMilli()
}

// 根据业务码获取HTTP状态码
func getHTTPStatus(code string) int {
	if len(code) < 3 {
		return 500
	}
	
	switch code[:3] {
	case "200":
		return 200
	case "400":
		return 400
	case "401":
		return 401
	case "403":
		return 403
	case "404":
		return 404
	case "429":
		return 429
	case "500":
		return 500
	default:
		return 200 // 默认200，让业务错误在正常HTTP响应中返回
	}
} 