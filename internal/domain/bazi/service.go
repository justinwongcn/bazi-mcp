package bazi

import "context"

// Service 定义了八字排盘领域服务的接口。
type Service interface {
	// GetPaipanResult 根据请求获取八字排盘结果。
	GetPaipanResult(ctx context.Context, req BaziRequest) (*BaziResponse, error)
}
