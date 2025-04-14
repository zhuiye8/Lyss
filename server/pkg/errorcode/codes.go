package errorcode

// 通用错误码 - 200系列 成功
const (
	Success       = "200"    // 请求成功
	Created       = "201"    // 资源创建成功
	Accepted      = "202"    // 请求已接受，但未处理完成
	NoContent     = "204"    // 请求成功，无返回内容
)

// 通用错误码 - 400系列 请求错误
const (
	BadRequest           = "400000"  // 请求参数错误
	InvalidParameter     = "400001"  // 无效的参数
	MissingParameter     = "400002"  // 缺少必要参数
	InvalidFormat        = "400003"  // 参数格式错误
	DuplicateRequest     = "400004"  // 重复的请求
	TooManyParameters    = "400005"  // 参数过多
	RequestBodyTooLarge  = "400006"  // 请求体过大
)

// 通用错误码 - 401系列 认证错误
const (
	Unauthorized         = "401000"  // 未认证或认证已过期
	InvalidToken         = "401001"  // 无效的令牌
	TokenExpired         = "401002"  // 令牌已过期
	InvalidCredentials   = "401003"  // 无效的凭证
	MfaRequired          = "401004"  // 需要多因素认证
)

// 通用错误码 - 403系列 权限错误
const (
	Forbidden            = "403000"  // 权限不足
	ResourceForbidden    = "403001"  // 禁止访问资源
	OperationForbidden   = "403002"  // 禁止的操作
	RateLimitExceeded    = "403003"  // 请求频率超限
	IpForbidden          = "403004"  // IP地址被禁止
	AccountDisabled      = "403005"  // 账号已禁用
)

// 通用错误码 - 404系列 资源错误
const (
	NotFound             = "404000"  // 资源不存在
	UserNotFound         = "404001"  // 用户不存在
	ResourceNotFound     = "404002"  // 其他资源不存在
	EndpointNotFound     = "404003"  // 接口不存在
)

// 通用错误码 - 429系列 请求过频
const (
	TooManyRequests      = "429000"  // 请求过于频繁
	QuotaExceeded        = "429001"  // 配额已用尽
)

// 通用错误码 - 500系列 服务器错误
const (
	InternalError        = "500000"  // 内部服务器错误
	ServiceUnavailable   = "500001"  // 服务不可用
	DatabaseError        = "500002"  // 数据库错误
	ExternalServiceError = "500003"  // 外部服务错误
	NetworkError         = "500004"  // 网络错误
	ConfigError          = "500005"  // 配置错误
	UnknownError         = "500999"  // 未知错误
)

// 业务错误码 - 600系列 用户模块
const (
	UserAlreadyExists    = "600001"  // 用户已存在
	UserCreateFailed     = "600002"  // 用户创建失败
	UserUpdateFailed     = "600003"  // 用户更新失败
	UserDeleteFailed     = "600004"  // 用户删除失败
	InvalidPassword      = "600005"  // 密码错误
	PasswordTooWeak      = "600006"  // 密码强度不足
	EmailAlreadyExists   = "600007"  // 邮箱已存在
	PhoneAlreadyExists   = "600008"  // 手机号已存在
	UserInactive         = "600009"  // 用户未激活
)

// 业务错误码 - 601系列 认证模块
const (
	LoginFailed          = "601001"  // 登录失败
	RegisterFailed       = "601002"  // 注册失败
	TokenCreateFailed    = "601003"  // 令牌创建失败
	TooManyLoginAttempts = "601004"  // 登录尝试次数过多
	VerificationFailed   = "601005"  // 验证失败
)

// 业务错误码 - 602系列 智能体模块
const (
	AgentNotFound        = "602001"  // 智能体不存在
	AgentCreateFailed    = "602002"  // 智能体创建失败
	AgentUpdateFailed    = "602003"  // 智能体更新失败
	AgentDeleteFailed    = "602004"  // 智能体删除失败
	AgentQueryFailed     = "602005"  // 智能体查询失败
	AgentRunFailed       = "602006"  // 智能体运行失败
	InvalidAgentConfig   = "602007"  // 无效的智能体配置
)

// 业务错误码 - 603系列 对话模块
const (
	ConversationNotFound = "603001"  // 对话不存在
	ConversationCreateFailed = "603002"  // 对话创建失败
	ConversationUpdateFailed = "603003"  // 对话更新失败
	ConversationDeleteFailed = "603004"  // 对话删除失败
	MessageSendFailed    = "603005"  // 消息发送失败
	InvalidMessageFormat = "603006"  // 无效的消息格式
)

// 业务错误码 - 604系列 知识库模块
const (
	KnowledgeBaseNotFound = "604001"  // 知识库不存在
	KnowledgeBaseCreateFailed = "604002"  // 知识库创建失败
	KnowledgeBaseUpdateFailed = "604003"  // 知识库更新失败
	KnowledgeBaseDeleteFailed = "604004"  // 知识库删除失败
	DocumentUploadFailed = "604005"  // 文档上传失败
	DocumentProcessFailed = "604006"  // 文档处理失败
	DocumentDeleteFailed = "604007"  // 文档删除失败
	DocumentNotFound     = "604008"  // 文档不存在
	QueryFailed          = "604009"  // 查询失败
)

// 业务错误码 - 605系列 模型模块
const (
	ModelNotFound        = "605001"  // 模型不存在
	ModelCreateFailed    = "605002"  // 模型创建失败
	ModelUpdateFailed    = "605003"  // 模型更新失败
	ModelDeleteFailed    = "605004"  // 模型删除失败
	ModelConnectionFailed = "605005"  // 模型连接失败
	ModelQuotaExceeded   = "605006"  // 模型配额超限
	InvalidModelConfig   = "605007"  // 无效的模型配置
)

// 错误码映射表，用于获取错误信息
var ErrorMessages = map[string]string{
	// 200系列
	Success:              "操作成功",
	Created:              "创建成功",
	Accepted:             "请求已接受",
	NoContent:            "无内容",
	
	// 400系列
	BadRequest:           "请求参数错误",
	InvalidParameter:     "无效的参数",
	MissingParameter:     "缺少必要参数",
	InvalidFormat:        "参数格式错误",
	DuplicateRequest:     "重复的请求",
	TooManyParameters:    "参数过多",
	RequestBodyTooLarge:  "请求体过大",
	
	// 401系列
	Unauthorized:         "未认证或认证已过期",
	InvalidToken:         "无效的令牌",
	TokenExpired:         "令牌已过期",
	InvalidCredentials:   "无效的凭证",
	MfaRequired:          "需要多因素认证",
	
	// 403系列
	Forbidden:            "权限不足",
	ResourceForbidden:    "禁止访问资源",
	OperationForbidden:   "禁止的操作",
	RateLimitExceeded:    "请求频率超限",
	IpForbidden:          "IP地址被禁止",
	AccountDisabled:      "账号已禁用",
	
	// 404系列
	NotFound:             "资源不存在",
	UserNotFound:         "用户不存在",
	ResourceNotFound:     "其他资源不存在",
	EndpointNotFound:     "接口不存在",
	
	// 429系列
	TooManyRequests:      "请求过于频繁",
	QuotaExceeded:        "配额已用尽",
	
	// 500系列
	InternalError:        "内部服务器错误",
	ServiceUnavailable:   "服务不可用",
	DatabaseError:        "数据库错误",
	ExternalServiceError: "外部服务错误",
	NetworkError:         "网络错误",
	ConfigError:          "配置错误",
	UnknownError:         "未知错误",
	
	// 600系列 - 用户模块
	UserAlreadyExists:    "用户已存在",
	UserCreateFailed:     "用户创建失败",
	UserUpdateFailed:     "用户更新失败",
	UserDeleteFailed:     "用户删除失败",
	InvalidPassword:      "密码错误",
	PasswordTooWeak:      "密码强度不足",
	EmailAlreadyExists:   "邮箱已存在",
	PhoneAlreadyExists:   "手机号已存在",
	UserInactive:         "用户未激活",
	
	// 601系列 - 认证模块
	LoginFailed:          "登录失败",
	RegisterFailed:       "注册失败",
	TokenCreateFailed:    "令牌创建失败",
	TooManyLoginAttempts: "登录尝试次数过多",
	VerificationFailed:   "验证失败",
	
	// 602系列 - 智能体模块
	AgentNotFound:        "智能体不存在",
	AgentCreateFailed:    "智能体创建失败",
	AgentUpdateFailed:    "智能体更新失败",
	AgentDeleteFailed:    "智能体删除失败",
	AgentQueryFailed:     "智能体查询失败",
	AgentRunFailed:       "智能体运行失败",
	InvalidAgentConfig:   "无效的智能体配置",
	
	// 603系列 - 对话模块
	ConversationNotFound: "对话不存在",
	ConversationCreateFailed: "对话创建失败",
	ConversationUpdateFailed: "对话更新失败",
	ConversationDeleteFailed: "对话删除失败",
	MessageSendFailed:    "消息发送失败",
	InvalidMessageFormat: "无效的消息格式",
	
	// 604系列 - 知识库模块
	KnowledgeBaseNotFound: "知识库不存在",
	KnowledgeBaseCreateFailed: "知识库创建失败",
	KnowledgeBaseUpdateFailed: "知识库更新失败",
	KnowledgeBaseDeleteFailed: "知识库删除失败",
	DocumentUploadFailed: "文档上传失败",
	DocumentProcessFailed: "文档处理失败", 
	DocumentDeleteFailed: "文档删除失败",
	DocumentNotFound:     "文档不存在",
	QueryFailed:          "查询失败",
	
	// 605系列 - 模型模块
	ModelNotFound:        "模型不存在",
	ModelCreateFailed:    "模型创建失败",
	ModelUpdateFailed:    "模型更新失败",
	ModelDeleteFailed:    "模型删除失败",
	ModelConnectionFailed: "模型连接失败",
	ModelQuotaExceeded:   "模型配额超限",
	InvalidModelConfig:   "无效的模型配置",
}

// GetMessage 根据错误码获取对应的错误信息
func GetMessage(code string) string {
	if msg, exists := ErrorMessages[code]; exists {
		return msg
	}
	return "未知错误"
} 