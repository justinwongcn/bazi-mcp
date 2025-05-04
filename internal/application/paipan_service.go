package application

import (
	"context"
	"fmt"
	"strconv" // 添加 strconv 包
	"strings"

	"github.com/john/bazi-mcp/internal/domain/bazi"
	"github.com/john/bazi-mcp/internal/domain/location"
)

var pillars = []string{"年", "月", "日", "时"}

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
func (s *BaziAppService) GetBaziPaipan(ctx context.Context, req bazi.Request) (string, bool, error) {
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
func (s *BaziAppService) validateInput(req bazi.Request) (string, bool) {
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
func (s *BaziAppService) handleAPIResponse(req bazi.Request, resp *bazi.PaipanResponse) (string, bool, error) {
	// 处理 API 返回的业务错误
	if resp.ErrCode != 0 {
		promptText := s.formatErrorPrompt(req, resp)
		// 检查是否因为缺少省市信息导致真太阳时计算失败
		if req.Zhen == 1 && (req.Province == "" || req.City == "") {
			promptText += "\n⚠️ 注意：当选择考虑真太阳时(zhen=1)时，必须提供省份和城市信息！\n"
		}

		// 将 API 返回的原始数据附加到错误信息后

		return promptText + "\n【遇到错误】\n" + resp.ErrMsg, true, nil
	}

	// 格式化成功结果
	promptText := s.formatSuccessPrompt(req)
	// 使用formatDetailedText格式化详细排盘数据
	// 从 bazi.Result 中提取数据
	// 并使用 formatDetailedText 格式化
	formattedText := s.formatDetailedText(resp)
	promptText += formattedText

	return promptText, false, nil
}

// formatSuccessPrompt 生成成功的提示信息。
func (s *BaziAppService) formatSuccessPrompt(req bazi.Request) string {
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
func (s *BaziAppService) formatErrorPrompt(req bazi.Request, resp *bazi.PaipanResponse) string {
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

// writeBaseInfo 输出基本信息
func (s *BaziAppService) writeBaseInfo(builder *strings.Builder, baseInfo *bazi.BaseInfo) {
	builder.WriteString("\n【基本信息】\n")
	// 预分配足够大的缓冲区
	builder.Grow(512)

	// 使用更高效的字符串构建方式
	infoFields := []struct {
		label string
		value string
	}{
		{"姓名", baseInfo.Name},
		{"性别", fmt.Sprintf("%s（乾造为男，坤造为女）", baseInfo.Sex)},
		{"公历", baseInfo.Gongli},
		{"农历", baseInfo.Nongli},
		{"起运时间", baseInfo.Qiyun},
		{"交运", baseInfo.Jiaoyun},
		{"八字正格", baseInfo.Zhengge},
	}

	for _, field := range infoFields {
		builder.WriteString(field.label)
		builder.WriteString("：")
		builder.WriteString(field.value)
		builder.WriteByte('\n')
	}
}

// writeBaziInfo 输出八字排盘信息
func (s *BaziAppService) writeBaziInfo(builder *strings.Builder, baziInfo *bazi.BaziInfo) {
	// 预分配足够大的缓冲区
	builder.Grow(2048)

	if len(baziInfo.Bazi) != 4 {
		builder.WriteString("错误：未获取到有效的八字排盘数据")
		return
	}

	builder.WriteString("\n【八字排盘】\n")
	builder.WriteString("空亡位置：")
	builder.WriteString(baziInfo.Kw)
	builder.WriteString("\n")

	for i := range pillars {
		builder.WriteString(pillars[i])
		builder.WriteString("柱如下：\n")
		builder.WriteString("  干支八字：")
		builder.WriteString(baziInfo.Bazi[i])
		builder.WriteString("\n")
		builder.WriteString("  天干十神：")
		builder.WriteString(baziInfo.TgCgGod[i])
		builder.WriteString("\n")
		builder.WriteString("  地支藏干：")
		builder.WriteString(baziInfo.DzCg[i])
		builder.WriteString("\n")
		builder.WriteString("  地支藏干十神：")
		builder.WriteString(baziInfo.DzCgGod[i])
		builder.WriteString("\n")
		builder.WriteString("  十二长生衰亡：")
		builder.WriteString(baziInfo.DayCs[i])
		builder.WriteString("\n")
		builder.WriteString("  纳音：")
		builder.WriteString(baziInfo.NaYin[i])
		builder.WriteString("\n\n")
	}
}

// writeDayunInfo 输出大运信息
func (s *BaziAppService) writeDayunInfo(builder *strings.Builder, dayunInfo *bazi.DayunInfo) {
	builder.WriteString("\n【大运信息】\n")

	for i := range dayunInfo.Big {
		// 大运基本信息
		builder.WriteString("\n第")
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString("个大运： ")
		builder.WriteString(strconv.Itoa(dayunInfo.XuSui[i]))
		builder.WriteString("-")
		builder.WriteString(strconv.Itoa(dayunInfo.XuSui[i] + 10))
		builder.WriteString("岁（")
		builder.WriteString(strconv.Itoa(dayunInfo.BigStartYear[i]))
		builder.WriteString("年-")
		builder.WriteString(strconv.Itoa(dayunInfo.BigEndYear[i]))
		builder.WriteString("年）：\n")

		// 大运详细信息
		builder.WriteString("  大运干支：")
		builder.WriteString(dayunInfo.Big[i])
		builder.WriteString("\n")
		builder.WriteString("  大运天干十神：")
		builder.WriteString(dayunInfo.BigGod[i])
		builder.WriteString("\n")
		builder.WriteString("  大运长生衰旺：")
		builder.WriteString(dayunInfo.BigCs[i])
		builder.WriteString("\n")
		// 流年信息
		type yearOffsetType struct {
			yearChar   string
			yearOffset int
		}

		yearsInfoFields := []yearOffsetType{
			{dayunInfo.YearsInfo0[i].YearChar, -1},
			{dayunInfo.YearsInfo1[i].YearChar, 0},
			{dayunInfo.YearsInfo2[i].YearChar, 1},
			{dayunInfo.YearsInfo3[i].YearChar, 2},
			{dayunInfo.YearsInfo4[i].YearChar, 3},
			{dayunInfo.YearsInfo5[i].YearChar, 4},
			{dayunInfo.YearsInfo6[i].YearChar, 5},
			{dayunInfo.YearsInfo7[i].YearChar, 6},
			{dayunInfo.YearsInfo8[i].YearChar, 7},
			{dayunInfo.YearsInfo9[i].YearChar, 8},
		}

		for _, info := range yearsInfoFields {
			// 使用 strings.Builder 和 strconv.Itoa 优化字符串拼接
			builder.WriteString("  流年年柱：")
			builder.WriteString(info.yearChar)
			builder.WriteString(" ")
			builder.WriteString(strconv.Itoa(dayunInfo.BigStartYear[i] + info.yearOffset))
			builder.WriteString("年(虚岁")
			builder.WriteString(strconv.Itoa(dayunInfo.XuSui[i] + info.yearOffset))
			builder.WriteString(")\n")
		}
	}
}

// writeDetailInfo 输出详细信息，对应bazi.DetailInfo结构，按四柱组织
func (s *BaziAppService) writeDetailInfo(builder *strings.Builder, detailInfo bazi.DetailInfo) {
	// 预分配足够大的缓冲区
	builder.Grow(4096)

	builder.WriteString("\n【四柱详细信息】\n")

	for i := range pillars {
		builder.WriteString("\n")
		builder.WriteString(pillars[i])
		builder.WriteString("柱详细信息：\n")

		var tg, dz, zhuxing, xingyun, zizuo, kongwang, nayin, shensha, cangganStr, fuxingStr string
		var canggan, fuxing []string

		switch i {
		case 0: // 年柱
			zhuxing = detailInfo.Zhuxing.Year
			xingyun = detailInfo.Xingyun.Year
			zizuo = detailInfo.Zizuo.Year
			kongwang = detailInfo.Kongwang.Year
			nayin = detailInfo.Nayin.Year
			shensha = detailInfo.Shensha.Year
			canggan = detailInfo.Canggan.Year
			fuxing = detailInfo.Fuxing.Year
			tg = detailInfo.Sizhu.Year.Tg
			dz = detailInfo.Sizhu.Year.Dz
		case 1: // 月柱
			zhuxing = detailInfo.Zhuxing.Month
			xingyun = detailInfo.Xingyun.Month
			zizuo = detailInfo.Zizuo.Month
			kongwang = detailInfo.Kongwang.Month
			nayin = detailInfo.Nayin.Month
			shensha = detailInfo.Shensha.Month
			canggan = detailInfo.Canggan.Month
			fuxing = detailInfo.Fuxing.Month
			tg = detailInfo.Sizhu.Month.Tg
			dz = detailInfo.Sizhu.Month.Dz
		case 2: // 日柱
			zhuxing = detailInfo.Zhuxing.Day
			xingyun = detailInfo.Xingyun.Day
			zizuo = detailInfo.Zizuo.Day
			kongwang = detailInfo.Kongwang.Day
			nayin = detailInfo.Nayin.Day
			shensha = detailInfo.Shensha.Day
			canggan = detailInfo.Canggan.Day
			fuxing = detailInfo.Fuxing.Day
			tg = detailInfo.Sizhu.Day.Tg
			dz = detailInfo.Sizhu.Day.Dz
		case 3: // 时柱
			zhuxing = detailInfo.Zhuxing.Hour
			xingyun = detailInfo.Xingyun.Hour
			zizuo = detailInfo.Zizuo.Hour
			kongwang = detailInfo.Kongwang.Hour
			nayin = detailInfo.Nayin.Hour
			shensha = detailInfo.Shensha.Hour
			canggan = detailInfo.Canggan.Hour
			fuxing = detailInfo.Fuxing.Hour
			tg = detailInfo.Sizhu.Hour.Tg
			dz = detailInfo.Sizhu.Hour.Dz
		}

		// 处理地支藏干和藏干十神字符串
		if len(canggan) > 0 {
			cangganStr = strings.Join(canggan, "|")
		} else {
			cangganStr = "无"
		}
		if len(fuxing) > 0 {
			fuxingStr = strings.Join(fuxing, "|")
		} else {
			fuxingStr = "无"
		}

		// 使用结构体切片统一输出
		detailFields := []struct {
			label string
			value string
		}{
			{"天干", tg},
			{"地支", dz},
			{"天干透出十神", zhuxing},
			{"星运信息", xingyun},
			{"自坐特性", zizuo},
			{"空亡方位", kongwang},
			{"纳音五行", nayin},
			{"神煞组合", shensha},
			{"地支藏干", cangganStr},
			{"藏干十神", fuxingStr},
		}

		for _, field := range detailFields {
			builder.WriteString("  ")
			builder.WriteString(field.label)
			builder.WriteString("：")
			builder.WriteString(field.value)
			builder.WriteByte('\n')
		}
	}

	// 大运神煞 (非柱位相关，单独列出)
	builder.WriteString("\n【大运神煞】\n")
	for i, ds := range detailInfo.Dayunshensha {
		builder.WriteString("  第")
		builder.WriteString(strconv.Itoa(i + 1))
		builder.WriteString("个大运天干地支：")
		builder.WriteString(ds.Tgdz)
		builder.WriteString("，对应神煞：")
		builder.WriteString(ds.Shensha)
		builder.WriteString("\n")
	}
}

// writeStartInfo 输出起运信息
func (s *BaziAppService) writeStartInfo(builder *strings.Builder, startInfo *bazi.StartInfo) {
	// 预分配足够大的缓冲区
	builder.Grow(512)

	builder.WriteString("\n【起运信息】\n")
	// 按照年月日时顺序输出吉神
	if len(startInfo.Jishen) >= 4 {
		builder.WriteString("吉神：\n")
		for i, pillar := range pillars {
			builder.WriteString("  ")
			builder.WriteString(pillar)
			builder.WriteString("柱：")
			builder.WriteString(startInfo.Jishen[i])
			builder.WriteString("\n")
		}
	} else {
		builder.WriteString("吉神：")
		builder.WriteString(strings.Join(startInfo.Jishen, "、"))
		builder.WriteString("\n")
	}
	builder.WriteString("星座：")
	builder.WriteString(startInfo.Xz)
	builder.WriteString("\n")
	builder.WriteString("生肖：")
	builder.WriteString(startInfo.Sx)
	builder.WriteString("\n")
}

// writeZhenSolarTimeInfo 输出真太阳时信息
func (s *BaziAppService) writeZhenSolarTimeInfo(builder *strings.Builder, zhen *bazi.ZhenInfo) {
	if zhen == nil {
		return
	}

	// 预分配足够大的缓冲区
	builder.Grow(256)

	builder.WriteString("\n【真太阳时信息】\n")
	builder.WriteString("省份：")
	builder.WriteString(zhen.Province)
	builder.WriteString("\n城市：")
	builder.WriteString(zhen.City)
	builder.WriteString("\n经度：")
	builder.WriteString(zhen.Jingdu)
	builder.WriteString("\n纬度：")
	builder.WriteString(zhen.Weidu)
	builder.WriteString("\n时差：")
	builder.WriteString(zhen.Shicha)
	builder.WriteString("\n")
}

// formatDetailedText 格式化详细文本
func (s *BaziAppService) formatDetailedText(data *bazi.PaipanResponse) string {
	var builder strings.Builder
	resData := data.Data

	// 输出基本信息
	s.writeBaseInfo(&builder, &resData.BaseInfo)

	// 输出真太阳时信息
	s.writeZhenSolarTimeInfo(&builder, resData.BaseInfo.Zhen)

	// 输出八字排盘信息
	s.writeBaziInfo(&builder, &resData.BaziInfo)

	// 输出大运信息
	s.writeDayunInfo(&builder, &resData.DayunInfo)

	// 起运信息
	s.writeStartInfo(&builder, &resData.StartInfo)

	s.writeDetailInfo(&builder, resData.DetailInfo)

	return builder.String()
}
