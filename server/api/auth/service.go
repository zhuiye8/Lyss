package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zhuiye8/Lyss/server/models"
	"github.com/zhuiye8/Lyss/server/pkg/auth"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("无效的凭证")
	ErrUserAlreadyExists  = errors.New("用户已存在")
)

// Service 提供认证相关功能
type Service struct {
	db         *gorm.DB
	jwtManager *auth.JWTManager
}

// NewService 创建新的认证服务
func NewService(db *gorm.DB, jwtManager *auth.JWTManager) *Service {
	return &Service{
		db:         db,
		jwtManager: jwtManager,
	}
}

// RegisterRequest 包含注册用户所需的数据
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name" binding:"required"`
}

// RegisterResponse 是注册API的响应
type RegisterResponse struct {
	User  models.UserResponse `json:"user"`
	Token string              `json:"token"`
}

// Register 注册新用户
func (s *Service) Register(req RegisterRequest) (*RegisterResponse, error) {
	// 检查用户名和邮箱是否已存在
	var count int64
	s.db.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		return nil, ErrUserAlreadyExists
	}

	// 创建新用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		Role:     "user",
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// 生成JWT令牌
	token, err := s.jwtManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// 返回响应
	return &RegisterResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

// LoginRequest 包含登录所需的数据
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 是登录API的响应
type LoginResponse struct {
	User         models.UserResponse `json:"user"`
	Token        string              `json:"token"`
	RefreshToken string              `json:"refresh_token"`
	ExpiresAt    time.Time           `json:"expires_at"`
}

// Login 验证用户凭证并生成令牌
func (s *Service) Login(req LoginRequest) (*LoginResponse, error) {
	// 查找用户
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, ErrInvalidCredentials
	}

	// 生成访问令牌
	token, err := s.jwtManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 返回响应
	return &LoginResponse{
		User:         user.ToResponse(),
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour), // 假设令牌有效期为24小时
	}, nil
}

// RefreshRequest 包含更新令牌所需的数据
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse 是刷新令牌API的响应
type RefreshResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// RefreshToken 使用刷新令牌生成新的访问令牌
func (s *Service) RefreshToken(req RefreshRequest) (*RefreshResponse, error) {
	// 解析刷新令牌
	claims, err := s.jwtManager.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	var user models.User
	if err := s.db.First(&user, "id = ?", claims.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// 生成新的访问令牌
	token, err := s.jwtManager.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// 返回响应
	return &RefreshResponse{
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 假设令牌有效期为24小时
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *Service) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
} 
