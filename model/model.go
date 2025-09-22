package chatmodel

import (
	"context"
	"github.com/cloudwego/eino-ext/components/model/ark"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"os"
)

var SystemPrompt = `
一、身份定位与核心使命
你是可自定义名称的全能型智能助手，核心使命是为用户提供专业、高效、人性化的服务，覆盖学习、工作、生活全场景，兼顾专业性与亲和力，成为用户信赖的一站式伙伴。
二、核心能力模块
1. 信息查询与知识服务
基础信息（人物、事件、数据等）快速检索，优先引用权威来源；实时数据需标注更新时间。
专业领域问题可通俗科普，复杂问题提示 “建议咨询执业人员”，避免误导。
按逻辑（时间 / 因果 / 层级）结构化整合信息，用分点、表格等清晰呈现。
2. 内容创作与编辑
按需创作文案（营销、邮件、演讲稿等）、创意内容（故事、脚本等），先确认主题、受众等核心信息。
提供文本纠错、润色、缩写 / 扩写服务，优化后标注修改说明。
协助基础代码编写 / 调试、文档格式处理，复杂问题提示参考官方文档。
3. 生活与学习辅助
制定日程、旅行、学习等计划，结合用户优先级和时间节点提供建议。
基于预算、偏好提供产品推荐、景点选择等参考，注明 “需结合个人情况”。
辅助学科答疑（提供思路而非标准答案）、语言学习（单词、语法、翻译），翻译需提示 “正式场景建议专业服务”。
三、交互准则
需求模糊时主动引导确认（如 “请问你需要写什么类型的文案？”）。
默认采用 “亲切礼貌 + 专业” 语气，可根据用户语气灵活调整。
记住上下文对话，避免重复提问；回应需具体，拒绝空话。
四、安全边界
绝对拒绝：违法犯罪、违背公序良俗、侵犯知识产权、危害安全的需求，说明拒绝原因。
谨慎回应：医疗诊断、法律诉讼、金融投资等需求，提示 “建议咨询专业人士”，仅提供基础常识。
不询问 / 存储用户隐私；对知识盲区坦诚说明，不编造信息。
五、特殊场景处理
用户负面情绪：先共情，再给建设性建议，必要时提示寻求亲友 / 专业帮助。
自身错误：主动道歉并修正；功能不支持时说明并建议替代工具。`

// Model 结构体
type Model struct {
	arkModel *ark.ChatModel
}

func NewModel(ctx context.Context) (*Model, error) {
	apiKey := os.Getenv("ARK_API_KEY")
	modelName := os.Getenv("ARK_MODEL")

	arkModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		APIKey: apiKey,
		Model:  modelName,
	})
	if err != nil {
		return nil, err
	}

	return &Model{arkModel: arkModel}, nil
}

func (m *Model) Generate(ctx context.Context, messages []*schema.Message) (string, error) {
	responseMessage, err := m.arkModel.Generate(ctx, messages)
	if err != nil {
		return "", err
	}
	return responseMessage.Content, nil
}
func (m *Model) Stream(ctx context.Context, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	responseMessage, err := m.arkModel.Stream(ctx, messages)
	if err != nil {
		log.Printf("Stream error: %v", err)
		return nil, err
	}
	return responseMessage, nil
}
func NewReactAgent(ctx context.Context) (*react.Agent, error) {
	model, err := NewModel(ctx)
	if err != nil {
		panic(err)
	}
	tools := GetMCPTool(ctx)
	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: model.arkModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: tools,
		},
	})
	if err != nil {
		panic(err)
	}
	return agent, nil
}
func GetMCPTool(ctx context.Context) []tool.BaseTool {
	cli, err := client.NewSSEMCPClient("http://localhost:9090/sse")
	if err != nil {
		log.Fatal(err)
	}
	err = cli.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "example-client",
		Version: "1.0.0",
	}

	_, err = cli.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatal(err)
	}

	tools, err := mcpp.GetTools(ctx, &mcpp.Config{Cli: cli})
	if err != nil {
		log.Fatal(err)
	}

	return tools
}
