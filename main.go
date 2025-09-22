package main

import (
	chatmodel "chat-app/model"
	"chat-app/rag"
	"chat-app/server"
	"chat-app/tool"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	// 创建日志文件
	f, err := os.Create("app.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 设置日志输出到文件
	log.SetOutput(f)
	ctx := context.Background()
	ragInstance, err := rag.NewRAG(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize RAG: %v", err)
	}
	//res := ragInstance.Find(ctx, "滴滴地图的点序列，比如设置出发点，增加途径点")
	//for _, doc := range res {
	//	fmt.Printf("id:%s,ct:%s\n", doc.ID, doc.Content)
	//}

	//err = ragInstance.Store(ctx, "./docs")
	//if err != nil {
	//	log.Fatalf("Failed to store documents: %v", err)
	//}
	//初始化模型
	modelInstance, err := chatmodel.NewModel(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize model: %v", err)
	}
	// 初始化工具管理器
	toolManager, err := tool.NewToolManager(ctx)
	if err != nil {
		log.Printf("Failed to initialize tool manager: %v", err)
		toolManager = nil
	}

	// 初始化Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",            // Redis服务器地址
		Password: os.Getenv("REDIS_PASSWORD"), // Redis密码，如果没有设置则为空
		DB:       0,                           // 使用默认数据库
		Username: "default",
	})

	// 测试Redis连接
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")

	// 创建服务器
	s := server.NewServer(modelInstance, ragInstance, toolManager, redisClient)

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = ":9343"
	}
	fmt.Printf("Server starting on port %s\n", port)
	go func() {
		if err := s.Start(port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	select {}
}
