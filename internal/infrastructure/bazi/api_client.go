package bazi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/justinwongcn/bazi-mcp/internal/domain/bazi"
)

const (
	APIEndpoint = "https://api.yuanfenju.com/index.php/v1/Bazi"
)

// APIClient 实现了 bazi.Service 接口，通过调用外部 API 获取八字排盘结果。
type APIClient struct {
	apiKey string
}

// NewAPIClient 创建一个新的 APIClient 实例。
func NewAPIClient() *APIClient {
	return &APIClient{
		apiKey: os.Getenv("API_KEY"),
	}
}

// GetPaipanResult 调用外部 API 获取八字排盘结果。
func (c *APIClient) GetPaipanResult(ctx context.Context, req bazi.Request) (*bazi.PaipanResponse, error) {
	formData := c.buildFormData(req)

	body, err := c.makeAPIRequest(ctx, "/paipan", formData)
	if err != nil {
		return nil, err
	}

	var baziResp bazi.PaipanResponse
	if err := json.Unmarshal(body, &baziResp); err != nil {
		// 如果解析失败，仍然尝试返回原始 body 以便调试
		return nil, fmt.Errorf("解析API响应失败: %w, 原始响应: %s", err, string(body))
	}

	return &baziResp, nil
}

// makeAPIRequest 执行API请求
func (c *APIClient) makeAPIRequest(ctx context.Context, path string, formData url.Values) ([]byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", APIEndpoint+path, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应内容失败: %w", err)
	}
	return body, nil
}

// buildFormData 构建API请求的表单数据
func (c *APIClient) buildFormData(req bazi.Request) url.Values {
	formData := url.Values{}
	formData.Set("api_key", c.apiKey)
	formData.Set("name", req.Name)
	formData.Set("sex", strconv.Itoa(req.Sex))
	formData.Set("type", strconv.Itoa(req.Type))
	formData.Set("year", strconv.Itoa(req.Year))
	formData.Set("month", strconv.Itoa(req.Month))
	formData.Set("day", strconv.Itoa(req.Day))
	formData.Set("hours", strconv.Itoa(req.Hours))
	formData.Set("minute", strconv.Itoa(req.Minute))
	formData.Set("zhen", strconv.Itoa(req.Zhen))
	formData.Set("sect", strconv.Itoa(req.Sect))

	if req.Name != "" {
		formData.Set("name", req.Name)
	} else {
		formData.Set("name", "求测者")
	}

	if req.Province != "" {
		formData.Set("province", req.Province)
	}
	if req.City != "" {
		formData.Set("city", req.City)
	}
	if req.Lang != "" {
		formData.Set("lang", req.Lang)
	}

	return formData
}
