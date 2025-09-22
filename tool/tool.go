package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cloudwego/eino-ext/components/tool/bingsearch"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo"
	"github.com/cloudwego/eino-ext/components/tool/duckduckgo/ddgsearch"
	mcpp "github.com/cloudwego/eino-ext/components/tool/mcp"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ToolManager 工具管理器，用于管理所有工具
type ToolManager struct {
	bingSearch *BingSearchTool
	duckSearch *DuckSearchTool
	sumTool    tool.InvokableTool
}

func NewToolManager(ctx context.Context) (*ToolManager, error) {
	// 初始化本地加法工具
	sumTl, err := createSumTool()
	if err != nil {
		log.Printf("sum tool initialization failed: %v", err)
		sumTl = nil
	}
	//bingSearch, err := NewBingSearchTool(ctx)
	//if err != nil {
	//	log.Printf("Bing search initialization failed: %v", err)
	//	bingSearch = nil
	//}

	//duckSearch, err := NewDuckSearchTool(ctx)
	//if err != nil {
	//	log.Printf("DuckDuckGo search initialization failed: %v", err)
	//	duckSearch = nil
	//}

	return &ToolManager{
		bingSearch: nil,
		duckSearch: nil,
		sumTool:    sumTl,
	}, nil
}

func (tm *ToolManager) ProcessWithTools(ctx context.Context, userInput string) (string, error) {

	// 识别并处理加法表达式
	if tm.sumTool != nil {
		if a, b, ok := extractTwoNumbersForAddition(userInput); ok {
			args := fmt.Sprintf("{\"a\": %s, \"b\": %s}", formatNumber(a), formatNumber(b))
			resp, err := tm.sumTool.InvokableRun(ctx, args)
			if err != nil {
				return "", err
			}
			var out struct {
				Result float64 `json:"result"`
			}
			if err := json.Unmarshal([]byte(resp), &out); err != nil {
				return "", err
			}
			return formatNumber(out.Result), nil
		}
	}

	// 无需处理则原样返回
	return userInput, nil
}

// BingSearchTool 封装Bing搜索工具
type BingSearchTool struct {
	tool tool.InvokableTool
}

func NewBingSearchTool(ctx context.Context) (*BingSearchTool, error) {
	apiKey := os.Getenv("BING_API_KEY")
	if apiKey == "" {
		return nil, errors.New("BING_API_KEY not found")
	}

	searchTool, err := bingsearch.NewTool(ctx, &bingsearch.Config{
		APIKey: apiKey,
		Cache:  5 * time.Minute,
	})
	if err != nil {
		return nil, err
	}

	return &BingSearchTool{
		tool: searchTool,
	}, nil
}

func (b *BingSearchTool) Search(ctx context.Context, query string) (string, error) {
	result, err := b.tool.InvokableRun(ctx, query)
	if err != nil {
		return "", err
	}
	return result, nil
}

// DuckSearchTool 封装Duck搜索工具
// DuckDuckGo 是一个注重隐私的搜索引擎，不会追踪用户的搜索行为，重点是无需 api key 鉴权即可直接使用
type DuckSearchTool struct {
	tool tool.InvokableTool
}

func NewDuckSearchTool(ctx context.Context) (*DuckSearchTool, error) {
	searchTool, err := duckduckgo.NewTool(ctx, &duckduckgo.Config{
		ToolName:   "duckduckgo_search",                        // 工具名称
		ToolDesc:   "search web for information by duckduckgo", // 工具描述
		MaxResults: 2,
		Region:     ddgsearch.RegionCN,
		DDGConfig: &ddgsearch.Config{
			Timeout:    15 * time.Second,
			Cache:      true,
			MaxRetries: 2,
		},
	})
	if err != nil {
		return nil, err
	}
	return &DuckSearchTool{
		tool: searchTool,
	}, nil
}
func (d *DuckSearchTool) Search(ctx context.Context, query string) (string, error) {
	searchReq := &duckduckgo.SearchRequest{
		Query: query,
		Page:  1,
	}
	jsonReq, err := json.Marshal(searchReq)
	if err != nil {
		log.Fatalf("Marshal of search request failed, err=%v", err)
	}
	searchResp, err := d.tool.InvokableRun(ctx, string(jsonReq))
	if err != nil {
		return "", err
	}
	//var searchResp duckduckgo.SearchResponse
	//if err := json.Unmarshal([]byte(resp), &searchResp); err != nil {
	//	log.Fatalf("Unmarshal of search response failed, err=%v", err)
	//}
	return searchResp, nil
}
func GetMCPTool(ctx context.Context) []tool.BaseTool {
	cli, err := client.NewSSEMCPClient("http://localhost:12345/sse")
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

// ===================== Sum Tool Implementation =====================

type sumInput struct {
	A float64 `json:"a" jsonschema:"required,description=the first addend"`
	B float64 `json:"b" jsonschema:"required,description=the second addend"`
}

type sumOutput struct {
	Result float64 `json:"result"`
}

func sumLocal(_ context.Context, in *sumInput) (*sumOutput, error) {
	return &sumOutput{Result: in.A + in.B}, nil
}

func createSumTool() (tool.InvokableTool, error) {
	return utils.InferTool("sum_two_numbers", "计算两数之和", sumLocal)
}

func extractTwoNumbersForAddition(text string) (float64, float64, bool) {
	s := strings.TrimSpace(text)
	patterns := []string{
		`(?i)\s*(-?[0-9]+(?:\.[0-9]+)?)\s*\+\s*(-?[0-9]+(?:\.[0-9]+)?)\s*`,
		`(?i)\s*(-?[0-9]+(?:\.[0-9]+)?)\s*加上?\s*(-?[0-9]+(?:\.[0-9]+)?)\s*`,
		`(?i)相加\s*(-?[0-9]+(?:\.[0-9]+)?)\s*和\s*(-?[0-9]+(?:\.[0-9]+)?)`,
		`(?i)求和\s*(-?[0-9]+(?:\.[0-9]+)?)\s*和\s*(-?[0-9]+(?:\.[0-9]+)?)`,
	}
	for _, p := range patterns {
		re := regexp.MustCompile(p)
		m := re.FindStringSubmatch(s)
		if len(m) == 3 {
			a, err1 := strconv.ParseFloat(m[1], 64)
			b, err2 := strconv.ParseFloat(m[2], 64)
			if err1 == nil && err2 == nil {
				return a, b, true
			}
		}
	}
	return 0, 0, false
}

func formatNumber(f float64) string {
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	return strconv.FormatFloat(f, 'f', -1, 64)
}
