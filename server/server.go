package server

import (
	chatmodel "chat-app/model"
	"chat-app/rag"
	"chat-app/tool"
	"chat-app/user"
	"context"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type ChatRequest struct {
	ID      string `json:"id" binding:"required"`      // 消息ID，必需字段
	Content string `json:"content" binding:"required"` // 消息内容，必需字段
}

type ChatResponse struct {
	ID       string `json:"id"`       // 返回相同的消息ID
	Response string `json:"response"` // 模型生成的回复内容
}
type Server struct {
	router      *gin.Engine
	model       *chatmodel.Model
	rag         *rag.Rag
	toolManager *tool.ToolManager
	userManager *user.UserManager
	userHandler *user.UserHandler
}

func NewServer(model *chatmodel.Model, rag *rag.Rag, tm *tool.ToolManager, redisClient *redis.Client) *Server {
	router := gin.Default()
	// 添加中间件
	router.Use(gin.Logger())   // 记录请求日志
	router.Use(gin.Recovery()) // 从 panic 中恢复
	// 添加 CORS 中间件，允许所有来源（生产环境应限制为特定来源）
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 创建用户管理器
	userManager := user.NewUserManager(redisClient)
	userHandler := user.NewUserHandler(userManager)

	s := &Server{
		router:      router,
		model:       model,
		rag:         rag,
		toolManager: tm,
		userManager: userManager,
		userHandler: userHandler,
	}

	// 聊天相关路由
	router.POST("/AIchat/li", s.handleChat)
	router.POST("/AIchat/upload", s.handleUpload)

	// 用户相关路由
	router.POST("/user/register", userHandler.Register)
	router.POST("/user/login", userHandler.Login)
	router.POST("/user/logout", userHandler.Logout)
	router.GET("/user/current", userHandler.GetCurrentUser)

	return s
}
func (s *Server) Start(port string) error {
	return s.router.Run(port)
}

// 处理聊天请求
func (s *Server) handleChat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Content is required"})
		return
	}

	// 创建带有超时的上下文，防止请求处理时间过长
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var enhancedContent string
	var err error
	if s.toolManager != nil {
		enhancedContent, err = s.toolManager.ProcessWithTools(ctx, req.Content)
		if err != nil {
			log.Printf("Tool processing failed: %v", err)
			enhancedContent = req.Content
		}
	} else {
		enhancedContent = req.Content
	}

	log.Printf("Processed content: %s", enhancedContent)
	// 构建发送给 Ark 模型的消息
	messages := []*schema.Message{
		schema.SystemMessage(chatmodel.SystemPrompt), // 系统消息设定助手行为
		schema.UserMessage(enhancedContent),          // 用户消息来自 HTTP 请求
	}
	// 调用 Ark 模型生成回复
	//responseMessage, err := s.model.Generate(ctx, messages)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"error":   "Failed to generate response from AI model",
	//		"details": err.Error(),
	//	})
	//	return
	//}
	//response := ChatResponse{
	//	ID:       req.ID,
	//	Response: responseMessage,
	//}
	//
	//c.JSON(http.StatusOK, response)

	responseMessage, err := s.model.Stream(ctx, messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate response from AI model",
			"details": err.Error(),
		})
		return
	}
	defer responseMessage.Close()
	c.Stream(func(w io.Writer) bool {
		chunk, err := responseMessage.Recv()
		if err == io.EOF {
			return false
		}
		if err != nil {
			return false
		}
		//c.Render(-1, sse{Data: chunk.Content})
		c.Render(-1, ChatResponse{ID: req.ID, Response: chunk.Content})
		return true
	})

}

// 处理知识库文件上传并入库
func (s *Server) handleUpload(c *gin.Context) {
	// 限制表单大小可在上层中间件设置，这里直接读取
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// 读取文件内容
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()
	data, err := io.ReadAll(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	// 直接入库
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if s.rag != nil {
		if err := s.rag.StoreDoc(ctx, filepath.Base(file.Filename), string(data)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rag store failed", "details": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "uploaded and indexed", "file": file.Filename})
}
func (r ChatResponse) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
}

func (r ChatResponse) Render(w http.ResponseWriter) error {
	_, err := w.Write([]byte("data: " + r.Response + "\n\n"))
	return err
}

// SSE 结构体用于流式输出
type sse struct {
	Data string
}

func (r sse) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
}

func (r sse) Render(w http.ResponseWriter) error {
	_, err := w.Write([]byte("data: " + r.Data + "\n\n"))
	return err
}
