package bazi

import "github.com/ThinkInAIXYZ/go-mcp/protocol"

// 领域层维护的提示词元数据
const (
	PromptDescription = "八字排盘系统提示词"
)

// GetPromptArguments 提供符合领域规范的参数定义（返回协议依赖需明确职责边界）
func GetPromptArguments() []protocol.PromptArgument {
	return []protocol.PromptArgument{
		{
			Name:        "birth_time",
			Description: "出生时间，格式：YYYY-MM-DD HH:MM",
			Required:    true,
		},
	}
}

// GeneratePromptContent 生成领域特定的提示内容（实现与协议类型的转换，迁移原提示词内容至此）
func GeneratePromptContent() (*protocol.GetPromptResult, error) {
	return &protocol.GetPromptResult{
		Description: "八字排盘系统提示",
		Messages:    []protocol.PromptMessage{
			// 迁移原提示词内容至此...
		},
	}, nil
}
