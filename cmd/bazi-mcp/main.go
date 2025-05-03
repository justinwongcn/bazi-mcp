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
	registerLocationResources(mcpServer)

	// 4. 注册八字排盘工具
	registerBaziTool(mcpServer, baziAppService)

	// 5. 运行服务器
	log.Println("八字排盘MCP服务器已启动，使用 Stdio 输入/输出进行通信")
	if err = mcpServer.Run(); err != nil {
		log.Fatalf("服务器运行失败: %v", err)
	}
}

// registerLocationResources 注册省份和城市相关的 MCP 资源。
func registerLocationResources(mcpServer *server.Server) {
	// 注册省份资源
	mcpServer.RegisterResource(&protocol.Resource{
		Name:        "可用省份列表",
		URI:         "data://provinces",
		Description: "八字排盘可用的省份列表",
		MimeType:    "application/json",
	}, server.ResourceHandlerFunc(func(ctx context.Context, req *protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
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
	}))

	// 注册城市资源模板
	mcpServer.RegisterResourceTemplate(&protocol.ResourceTemplate{
		Name:        "省份城市查询",
		URITemplate: "data://cities/{province}",
		Description: "根据省份查询城市列表",
	}, server.ResourceHandlerFunc(func(ctx context.Context, req *protocol.ReadResourceRequest) (*protocol.ReadResourceResult, error) {
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
	}))
}

// registerBaziTool 注册八字排盘工具及其处理程序。
func registerBaziTool(mcpServer *server.Server, baziAppService *app.BaziAppService) {
	tool, err := protocol.NewTool("bazi_paipan", "根据生辰八字信息获取排盘结果", baziDomain.BaziRequest{})
	if err != nil {
		log.Fatalf("创建工具失败: %v", err)
	}

	mcpServer.RegisterTool(tool, server.ToolHandlerFunc(func(ctx context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
		// 1. 解析参数
		var baziReq baziDomain.BaziRequest
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
	}))
}
