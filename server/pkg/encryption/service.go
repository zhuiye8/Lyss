package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Service 提供加密和解密功能
type Service struct {
	key []byte
}

// NewService 创建加密服务实例
func NewService(secretKey string) (*Service, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("密钥长度不足，需要至少32字节")
	}
	
	// 使用提供的密钥的前32字节
	key := []byte(secretKey)[:32]
	
	return &Service{
		key: key,
	}, nil
}

// Encrypt 加密字符串
func (s *Service) Encrypt(plaintext string) (string, error) {
	plainBytes := []byte(plaintext)
	
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	
	// 创建分组加密模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	// 创建随机数
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	// 加密
	ciphertext := aesGCM.Seal(nonce, nonce, plainBytes, nil)
	
	// 返回Base64编码的密文
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字符串
func (s *Service) Decrypt(encryptedText string) (string, error) {
	// 解码Base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}
	
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}
	
	// 创建分组加密模式
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	// 提取nonce
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("密文长度不足")
	}
	
	nonce, cipherBytes := ciphertext[:nonceSize], ciphertext[nonceSize:]
	
	// 解密
	plainBytes, err := aesGCM.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}
	
	return string(plainBytes), nil
} 
