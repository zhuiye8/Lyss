package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 模型表示系统用户
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex" json:"username"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password  string    `gorm:"type:varchar(100)" json:"-"`
	FullName  string    `gorm:"type:varchar(100)" json:"full_name"`
	AvatarURL string    `gorm:"type:varchar(255)" json:"avatar_url,omitempty"`
	Role      string    `gorm:"type:varchar(20);default:'user'" json:"role"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate 在创建用户前生成UUID并加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	
	// 加密密码
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	
	return nil
}

// CheckPassword 验证密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// SetPassword 设置用户密码
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// UserResponse 是返回给客户端的用户数据结构（不包含敏感信息）
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// ToResponse 将完整用户模型转换为对外响应
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FullName:  u.FullName,
		AvatarURL: u.AvatarURL,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
} 