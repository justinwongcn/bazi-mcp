package main

import (
	"context"
	
	"fmt"
	"log"

	application "github.com/justinwongcn/bazi-mcp/internal/application"
	baziDomain "github.com/justinwongcn/bazi-mcp/internal/domain/bazi"
	
	baziInfra "github.com/justinwongcn/bazi-mcp/internal/infrastructure/bazi"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

// 2. 将工具名称定义为领域常量（提升领域概念内聚性）
const (
	BaziToolName = "bazi_paipan" // 领域工具名称常量定义
)

// Init 初始化并启动八字排盘MCP服务器。
func Init() error {
	// 1. 初始化依赖
	baziAppService := setupDependencies()

	// 2. 创建并配置服务器
	mcpServer, err := createAndConfigureServer(baziAppService)
	if err != nil {
		return err
	}

	// 3. 运行服务器
	return runServer(mcpServer)
}

// setupDependencies 初始化应用依赖
func setupDependencies() *application.BaziAppService {
	apiClient := baziInfra.NewAPIClient()
	var baziDomainService baziDomain.Service = apiClient
	return application.NewBaziAppService(baziDomainService)
}

// createAndConfigureServer 创建并配置MCP服务器
func createAndConfigureServer(baziAppService *application.BaziAppService) (*server.Server, error) {
	transportServer := transport.NewStdioServerTransport()
	mcpServer, err := server.NewServer(transportServer)
	if err != nil {
		return nil, fmt.Errorf("创建MCP服务器失败: %w", err)
	}

	if err := registerAllResources(mcpServer, baziAppService); err != nil {
		return nil, fmt.Errorf("服务器配置失败: %w", err)
	}

	return mcpServer, nil
}

// registerAllResources 注册所有资源
func registerAllResources(mcpServer *server.Server, baziAppService *application.BaziAppService) error {
	

	registerBaziTool(mcpServer, baziAppService)
	registerPrompts(mcpServer)
	return nil
}

// runServer 启动服务器运行
func runServer(mcpServer *server.Server) error {
	log.Printf("八字排盘MCP服务器已启动，使用 stdio 模式进行通信\n")
	if err := mcpServer.Run(); err != nil {
		return fmt.Errorf("服务器运行失败: %w", err)
	}
	return nil
}

func main() {
	// 调用 Init 方法启动程序
	if err := Init(); err != nil {
		log.Fatalf("程序启动失败: %v", err)
	}
}



// registerBaziTool 注册八字排盘工具及其处理程序
func registerBaziTool(mcpServer *server.Server, baziAppService *application.BaziAppService) {
	tool, err := protocol.NewTool(BaziToolName, "根据生辰八字信息获取排盘结果", baziDomain.Request{}) // 使用领域常量
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

// 3. 迁移提示词内容到领域层（保持领域知识内聚，明确MCP协议层职责）
func registerPrompts(mcpServer *server.Server) {
	// 创建八字排盘提示词
	baziPrompt := &protocol.Prompt{
		Name:        "bazi_prompt",
		Description: baziDomain.PromptDescription,    // 引用领域层定义的常量
		Arguments:   baziDomain.GetPromptArguments(), // 使用领域层构造方法
	}

	mcpServer.RegisterPrompt(baziPrompt, func(ctx context.Context, request *protocol.GetPromptRequest) (*protocol.GetPromptResult, error) {
		return baziDomain.GeneratePromptContent() // 将具体实现移交领域层
	})
}
