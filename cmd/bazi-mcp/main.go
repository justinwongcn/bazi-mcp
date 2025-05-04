package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	app "github.com/john/bazi-mcp/internal/application"
	baziDomain "github.com/john/bazi-mcp/internal/domain/bazi"
	"github.com/john/bazi-mcp/internal/domain/location"
	baziInfra "github.com/john/bazi-mcp/internal/infrastructure/bazi"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

func main() {
	// 1. 初始化依赖
	apiClient := baziInfra.NewAPIClient() // Infrastructure: API Client
	// Domain service is implicitly satisfied by APIClient as it implements the interface
	var baziDomainService baziDomain.Service = apiClient
	baziAppService := app.NewBaziAppService(baziDomainService) // Application Service

	// 2. 创建 MCP 服务器
	transportServer := transport.NewStdioServerTransport()
	mcpServer, err := server.NewServer(transportServer)
	if err != nil {
		log.Fatalf("创建MCP服务器失败: %v", err)
	}

	// 3. 注册资源 (省份和城市)
	err = registerLocationResources(mcpServer)
	if err != nil {
		log.Fatalf("注册资源失败: %v", err)
	}

	// 4. 注册八字排盘工具
	registerBaziTool(mcpServer, baziAppService)

	// 5. 注册提示词
	registerPrompts(mcpServer)

	// 5. 运行服务器
	log.Printf("八字排盘MCP服务器已启动，使用 stdio 模式进行通信\n")
	if err = mcpServer.Run(); err != nil {
		log.Fatalf("服务器运行失败: %v", err)
	}
}

// registerLocationResources 注册省份和城市相关的 MCP 资源。
func registerLocationResources(mcpServer *server.Server) error {
	// 注册省份资源
	mcpServer.RegisterResource(&protocol.Resource{
		Name:        "可用省份列表",
		URI:         "data://provinces",
		Description: "八字排盘可用的省份列表",
		MimeType:    "application/json",
	}, func(ctx context.Context, req *protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
		provincesJSON, _ := json.Marshal(location.Provinces)
		return &protocol.ReadResourceResult{
			Contents: []protocol.ResourceContents{
				protocol.TextResourceContents{
					URI:      "data://provinces",
					MimeType: "application/json",
					Text:     string(provincesJSON),
				},
			},
		}, nil
	})

	// 注册城市资源模板
	err := mcpServer.RegisterResourceTemplate(&protocol.ResourceTemplate{
		Name:        "省份城市查询",
		URITemplate: "data://cities/{province}",
		Description: "根据省份查询城市列表",
	}, func(ctx context.Context, req *protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
		province, ok := req.Arguments["province"].(string)
		if !ok || province == "" {
			return nil, fmt.Errorf("未提供有效的省份参数")
		}

		citiesList, found := location.Cities[province]
		if !found {
			// 即使省份在 Provinces 列表中，也可能没有对应的城市数据（数据不完整）
			// 或者省份本身就不在 Provinces 列表中
			if !location.IsValidProvince(province) {
				return nil, fmt.Errorf("无效的省份: %s", province)
			} else {
				// 省份有效，但无城市数据
				return nil, fmt.Errorf("未找到省份 '%s' 对应的城市列表数据", province)
			}
		}

		citiesJSON, _ := json.Marshal(citiesList)
		return &protocol.ReadResourceResult{
			Contents: []protocol.ResourceContents{
				protocol.TextResourceContents{
					URI:      fmt.Sprintf("data://cities/%s", province),
					MimeType: "application/json",
					Text:     string(citiesJSON),
				},
			},
		}, nil
	})

	return err
}

// registerBaziTool 注册八字排盘工具及其处理程序。
func registerBaziTool(mcpServer *server.Server, baziAppService *app.BaziAppService) {
	tool, err := protocol.NewTool("bazi_paipan", "根据生辰八字信息获取排盘结果", baziDomain.Request{})
	if err != nil {
		log.Fatalf("创建工具失败: %v", err)
	}

	mcpServer.RegisterTool(tool, func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		// 1. 解析参数
		var baziReq baziDomain.Request
		if err := protocol.VerifyAndUnmarshal(req.RawArguments, &baziReq); err != nil {
			// 参数格式本身错误
			return protocol.NewCallToolResult([]protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("参数格式错误: %v\n请检查您的输入是否符合工具要求。\n您可以通过'data://provinces'查看省份列表，并通过'data://cities/{省份}'查询城市。", err),
				},
			}, true), nil // 标记为错误
		}

		// 2. 调用应用服务处理请求
		resultText, isAppError, appErr := baziAppService.GetBaziPaipan(ctx, baziReq)

		// 3. 处理应用服务返回的结果
		if appErr != nil {
			// 应用层或更底层发生内部错误
			log.Printf("处理八字排盘请求时发生内部错误: %v", appErr) // 记录详细错误
			return protocol.NewCallToolResult([]protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: "处理请求时发生内部错误，请稍后再试或联系管理员。",
				},
			}, true), nil // 标记为错误
		}

		// 应用服务正常处理完成，根据 isAppError 判断是成功还是业务错误
		return protocol.NewCallToolResult([]protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: resultText, // 应用服务已格式化好文本
			},
		}, isAppError), nil // 根据应用服务的结果设置 IsError 标志
	})
}

// registerPrompts 注册所有提示词
func registerPrompts(mcpServer *server.Server) {
	// 创建八字排盘提示词
	baziPrompt := &protocol.Prompt{
		Name:        "bazi_prompt",
		Description: "八字排盘系统提示词",
		Arguments: []protocol.PromptArgument{
			{
				Name:        "birth_time",
				Description: "出生时间，格式：YYYY-MM-DD HH:MM",
				Required:    true,
			},
		},
	}

	// 注册提示词处理函数
	mcpServer.RegisterPrompt(baziPrompt, func(ctx context.Context, request *protocol.GetPromptRequest) (*protocol.GetPromptResult, error) {
		return &protocol.GetPromptResult{
			Description: "八字排盘系统提示",
			Messages: []protocol.PromptMessage{
				{
					Role: protocol.RoleAssistant,
					Content: protocol.TextContent{
						Type: "text",
						Text: `你是一位资深命理师，精通子平八字推命术。请按照以下步骤分析：
融合子平八字与认知心理学的数字命理师，遵循「数据验证→格局分析→心理映射→解决方案」四步流程

【核心模块】
===时空校验===

时辰模糊处理：
23:00-24:59出生者标注「夜子时争议」并询问用户是否要区分早晚子时
调用NASA太阳时差数据库自动校正（精度±3分钟）
地域五行匹配： 使用《地理辨正》原理转换：北京（坎宫属水）→ 匹配亥子丑年能量场
===合冲刑害分析===

天干作用模型： 甲己合土＞乙庚合金＞... ｜ 优先度：合化成功＞合绊＞相克
地支关系矩阵：
三刑触发条件：寅巳申需同时出现两个以上
六害化解方案：子未害建议佩戴水晶（土水通关）
===日主喜忌系统===
调用《穷通宝鉴》算法：
甲木日主：
必要元素：庚（斧斤）丙（阳光）
禁忌组合：乙木透干+地支水旺
庚金日主：
优化路径：丁火炼金→壬水淬锋→甲木生火
风险预警：辛金透干引发比劫争财

===解决方案引擎===

五行调节： 「缺庚」→ 申时佩戴钛钢饰品 「丁弱」→ 参加金属锻造体验课
心理干预： 比劫过旺→MBTI矫正：ENTP→ISTJ平衡训练
【交互协议】

风险控制：
所有结论标注置信度（例：财运85%±5）
极端案例触发建议转接人工心理咨询师
验证机制：
提供以往年份重大事件反推验证
支持生成PDF版《命理分析溯源报告》
【示例指令逻辑】
当输入「甲戌 丙子 庚寅 辛巳」：

日主庚金识别 → 检查丁火（时支巳中藏丁）
寅巳申三刑预警 → 标注2016丙申年危机事件概率
生成建议：
五行：每日巳时（9-11点）日光浴补火
心理：每周三次HIIT训练泄金气 ===
身体健康：八字中藏丁，宜补火，忌水旺`,
					},
				},
				{
					Role: protocol.RoleUser,
					Content: protocol.TextContent{
						Type: "text",
						Text: "请结合以下MCP数据源进行深度分析：\n1. 使用《三命通会》的神煞数据验证命局\n2. 参考《渊海子平》的纳音五行表匹配日主特性\n3. 对比MCP时辰数据库中的地区真太阳时差值",
					},
				},
			},
		}, nil
	})
}
