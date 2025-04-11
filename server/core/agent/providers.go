package agent

// providers.go 提供自定义的LLM提供者实现

import (
	"context"

	"github.com/cloudwego/eino"
)

// AliyunProvider 是阿里云通义千问的LLM提供者实现
type AliyunProvider struct {
	apiKey  string
	baseURL string
}

// NewAliyunProvider 创建阿里云通义千问LLM提供者
func NewAliyunProvider(apiKey, baseURL string) *AliyunProvider {
	if baseURL == "" {
		baseURL = "https://dashscope.aliyuncs.com/api/v1"
	}
	return &AliyunProvider{
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

// BaiduProvider 是百度文心一言的LLM提供者实现
type BaiduProvider struct {
	apiKey  string
	appID   string
	baseURL string
}

// NewBaiduProvider 创建百度文心一言LLM提供者
func NewBaiduProvider(apiKey, appID, baseURL string) *BaiduProvider {
	if baseURL == "" {
		baseURL = "https://aip.baidubce.com/rpc/2.0/ai_custom/v1/wenxinworkshop"
	}
	return &BaiduProvider{
		apiKey:  apiKey,
		appID:   appID,
		baseURL: baseURL,
	}
} 
