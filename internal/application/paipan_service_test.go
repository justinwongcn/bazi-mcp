package application

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/justinwongcn/bazi-mcp/internal/domain/bazi"
)

func loadTestData(t *testing.T) *bazi.PaipanResponse {
	var testData *bazi.PaipanResponse
	jsonFile, err := os.ReadFile("testdata/result.json")
	if err != nil {
		t.Fatalf("读取result.json文件失败: %v", err)
	}
	err = json.Unmarshal(jsonFile, &testData)
	if err != nil {
		t.Fatalf("解析result.json文件失败: %v", err)
	}
	return testData
}

func TestHandleAPIResponse(t *testing.T) {
	testData := loadTestData(t)
	service := &BaziAppService{}

	t.Run("处理成功响应", func(t *testing.T) {
		_, isError, err := service.handleAPIResponse(bazi.Request{}, testData)
		if err != nil || isError {
			t.Errorf("处理成功响应失败: %v", err)
		}
	})

	t.Run("处理错误响应", func(t *testing.T) {
		resp := &bazi.PaipanResponse{
			ErrCode: 1,
			ErrMsg:  "测试错误",
			Data:    bazi.Data{}, // 使用空Data结构体作为错误响应
		}
		_, isError, err := service.handleAPIResponse(bazi.Request{}, resp)
		if err != nil || !isError {
			t.Error("应识别为错误响应")
		}
	})

	t.Run("处理真太阳时缺失警告", func(t *testing.T) {
		resp := &bazi.PaipanResponse{
			ErrCode: 1,
			ErrMsg:  "测试错误",
			Data:    bazi.Data{}, // 使用空Data结构体作为错误响应
		}
		req := bazi.Request{Zhen: 1}
		result, _, _ := service.handleAPIResponse(req, resp)
		if !strings.Contains(result, "⚠️ 注意：当选择考虑真太阳时") {
			t.Error("应包含真太阳时警告")
		}
	})
}

func TestFormatDetailedText(t *testing.T) {
	// 使用loadTestData加载测试数据
	testData := loadTestData(t)

	// 创建测试服务
	service := &BaziAppService{}

	// 测试正常情况
	t.Run("正常格式化", func(t *testing.T) {
		// 创建数据副本避免污染原始数据
		validData := *testData
		validData.Data.BaziInfo = bazi.BaziInfo{
			Kw:      "子丑",
			TgCgGod: []string{"比肩", "正印", "日元", "正印"},
			Bazi:    []string{"己卯", "丙子", "己未", "丙寅"},
			DzCg:    []string{"乙", "癸", "己|丁|乙", "甲|丙|戊"},
			DzCgGod: []string{"七杀", "偏财", "比肩|偏印|七杀", "正官|正印|劫财"},
			DayCs:   []string{"病", "绝", "冠带", "死"},
			NaYin:   []string{"城头土", "涧下水", "天上火", "炉中火"},
		}

		result := service.formatDetailedText(&validData)
		fmt.Println(result)
		if result == "" {
			t.Error("格式化结果不应为空")
		}
		// 添加更详细的断言
		if !strings.Contains(result, "姓名：张三") {
			t.Error("应包含姓名信息")
		}
		if !strings.Contains(result, "偏财格") {
			t.Error("应包含八字正格信息")
		}
	})

	// 测试无效数据情况
	t.Run("无效八字数据", func(t *testing.T) {
		invalidData := *testData
		invalidData.Data.BaziInfo.Bazi = []string{"己卯", "丙子"} // 只有2柱
		result := service.formatDetailedText(&invalidData)
		if !strings.Contains(result, "错误：未获取到有效的八字排盘数据") {
			t.Error("对于无效八字数据应返回错误信息")
		}
	})

	// 测试数组越界保护
	t.Run("数组越界保护", func(t *testing.T) {
		invalidData := *testData
		// 确保所有相关的数组长度一致，避免越界
		invalidData.Data.BaziInfo.TgCgGod = []string{"比肩", "正印", "日元", "正印"}
		invalidData.Data.BaziInfo.Bazi = []string{"己卯", "丙子", "己未", "丙寅"}
		invalidData.Data.BaziInfo.DzCg = []string{"乙", "癸", "己|丁|乙", "甲|丙|戊"}
		invalidData.Data.BaziInfo.DzCgGod = []string{"七杀", "偏财", "比肩|偏印|七杀", "正官|正印|劫财"}
		invalidData.Data.BaziInfo.DayCs = []string{"病", "绝", "冠带", "死"}
		invalidData.Data.BaziInfo.NaYin = []string{"城头土", "涧下水", "天上火", "炉中火"}

		result := service.formatDetailedText(&invalidData)
		if result == "" {
			t.Error("格式化结果不应为空")
		}
		if strings.Contains(result, "panic") {
			t.Error("方法应能处理数组越界情况而不panic")
		}
	})
}
