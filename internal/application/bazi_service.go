package application

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/john/bazi-mcp/internal/domain/bazi"
	"github.com/john/bazi-mcp/internal/domain/location"
)

// BaziAppService 定义了八字排盘的应用服务。
type BaziAppService struct {
	BaziDomainService bazi.Service
}

// NewBaziAppService 创建一个新的 BaziAppService 实例。
func NewBaziAppService(baziDomainService bazi.Service) *BaziAppService {
	return &BaziAppService{
		BaziDomainService: baziDomainService,
	}
}

// GetBaziPaipan 处理获取八字排盘结果的请求。
func (s *BaziAppService) GetBaziPaipan(ctx context.Context, req bazi.BaziRequest) (string, bool, error) {
	// 输入验证和默认值设置
	if errMsg, hasError := s.validateInput(req); hasError {
		return errMsg, true, nil
	}
	// 其他默认值在 APIClient 或请求构建时处理，这里主要处理业务逻辑相关的默认值或校验

	// 3. 调用领域服务获取结果
	baziResp, err := s.BaziDomainService.GetPaipanResult(ctx, req)
	if err != nil {
		// 底层错误，直接返回
		return "", true, fmt.Errorf("获取八字结果失败: %w", err)
	}

	// 处理API响应
	return s.handleAPIResponse(req, baziResp)
}

// validateInput 验证输入参数并设置默认值
func (s *BaziAppService) validateInput(req bazi.BaziRequest) (string, bool) {
	// 1. 输入验证 (省份和城市有效性)
	if req.Province != "" && !location.IsValidProvince(req.Province) {
		return fmt.Sprintf("无效省份: %s\n 一般最后面需要带上“省市区”等 例：\"北京市\"", req.Province), true
	}
	if req.City != "" && !location.IsValidCity(req.Province, req.City) {
		return fmt.Sprintf("无效城市: %s\n 最后面一般不带上“县市区”等（除非带上后只有两个字）", req.City), true
	}

	// 2. 设置默认值 (如果请求中未提供)
	if req.Name == "" {
		req.Name = "求测者"
	}

	return "", false
}

// handleAPIResponse 处理API响应
func (s *BaziAppService) handleAPIResponse(req bazi.BaziRequest, resp *bazi.BaziResponse) (string, bool, error) {
	// 处理 API 返回的业务错误
	if resp.ErrCode != 0 {
		promptText := s.formatErrorPrompt(req, resp)
		// 检查是否因为缺少省市信息导致真太阳时计算失败
		if req.Zhen == 1 && (req.Province == "" || req.City == "") {
			promptText += "\n⚠️ 注意：当选择考虑真太阳时(zhen=1)时，必须提供省份和城市信息！\n"
		}
		// 将 API 返回的原始数据附加到错误信息后
		formattedJSON, formatErr := s.formatJSONResponse(resp.Data) // 尝试格式化 Data 部分
		if formatErr != nil {
			formattedJSON = string(resp.Data) // 格式化失败则使用原始 Data
		}
		return promptText + "\n【原始响应数据】\n" + formattedJSON, true, nil // 返回格式化的错误消息，标记为错误
	}

	// 格式化成功结果
	promptText := s.formatSuccessPrompt(req)
	formattedJSON, formatErr := s.formatJSONResponse(resp.Data) // 格式化 Data 部分
	if formatErr != nil {
		// 如果格式化失败，可以记录日志，但仍尝试返回未格式化的数据
		promptText += "\n【详细排盘数据 (格式化失败)】\n" + string(resp.Data)
	} else {
		promptText += "\n【详细排盘数据】\n" + formattedJSON
	}

	return promptText, false, nil
}

// formatSuccessPrompt 生成成功的提示信息。
func (s *BaziAppService) formatSuccessPrompt(req bazi.BaziRequest) string {
	sexText := "男"
	if req.Sex == 1 {
		sexText = "女"
	}
	calendarText := "农历"
	if req.Type == 1 {
		calendarText = "公历"
	}
	trueTimeText := "是"
	if req.Zhen == 2 {
		trueTimeText = "否"
	}

	promptText := fmt.Sprintf(
		"✅ 成功获取 %s 的八字排盘结果！\n\n【基本信息】\n"+
			"姓名：%s\n性别：%s\n"+
			"出生时间：%d年%d月%d日%d时%d分\n"+
			"历法类型：%s\n是否考虑真太阳时：%s\n",
		req.Name, req.Name, sexText,
		req.Year, req.Month, req.Day, req.Hours, req.Minute,
		calendarText, trueTimeText)

	if req.Province != "" && req.City != "" {
		promptText += fmt.Sprintf("出生地点：%s %s\n", req.Province, req.City)
	}

	promptText += "\n【八字排盘结果解读指南】\n" +
		"1. 四柱结构：年柱(祖业)、月柱(父母)、日柱(自己)、时柱(子女)\n" +
		"2. 天干地支：每个柱位的天干地支组合构成命盘基础\n" +
		"3. 藏干十神：地支中隐藏的天干及其十神关系（比肩/正印等）\n" +
		"4. 五行纳音：年份对应的五行属性（如大林木、天河水等）\n" +
		"5. 大运走势：每十年大运对应的天干地支及运势变化\n" +
		"6. 神煞吉凶：包含禄神、太极、空亡等重要神煞说明\n" +
		"7. 真太阳信息：当考虑真太阳时显示经纬度及时差数据\n" +
		"8. 起运交运：标志人生重要阶段开始的关键时间点"

	return promptText
}

// formatErrorPrompt 生成失败的提示信息。
func (s *BaziAppService) formatErrorPrompt(req bazi.BaziRequest, resp *bazi.BaziResponse) string {
	sexText := "男"
	if req.Sex == 1 {
		sexText = "女"
	}
	calendarText := "农历"
	if req.Type == 1 {
		calendarText = "公历"
	}
	trueTimeText := "是"
	if req.Zhen == 2 {
		trueTimeText = "否"
	}

	return fmt.Sprintf(
		"❌ 获取八字排盘结果失败！错误码：%d，错误信息：%s\n\n"+
			"请检查您的输入参数是否正确：\n"+
			"姓名：%s\n性别：%s\n"+
			"出生时间：%d年%d月%d日%d时%d分\n"+
			"历法类型：%s\n是否考虑真太阳时：%s",
		resp.ErrCode, resp.ErrMsg,
		req.Name, sexText,
		req.Year, req.Month, req.Day, req.Hours, req.Minute,
		calendarText, trueTimeText)
}

// formatJSONResponse 格式化 JSON 数据。
func (s *BaziAppService) formatJSONResponse(data json.RawMessage) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return "", fmt.Errorf("格式化JSON失败: %w", err)
	}
	return prettyJSON.String(), nil
}
