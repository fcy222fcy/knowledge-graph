package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"software_engineering/internal/api"
	"software_engineering/internal/middleware"
	"software_engineering/internal/repository/seed"
	"software_engineering/internal/service"
	"software_engineering/pkg/config"
	"software_engineering/pkg/database"
)

type App struct {
	router     *gin.Engine
	httpServer *http.Server
}

func New() *App {
	return &App{}
}

func (a *App) Initialize() {
	// 1. 加载配置
	config.Load()

	// 2. 初始化 AI 客户端（需要在 config.Load 之后）
	service.InitAIClient()

	// 3. 连接数据库
	database.Connect()
	database.ConnectNeo4j()

	// 4. 数据库迁移
	database.AutoMigrate()

	// 5. 初始化种子数据
	seed.SeedAll()

	// 6. 初始化路由
	a.router = gin.New()
	a.router.Use(middleware.Logger())
	a.router.Use(middleware.Recovery())
	a.router.Use(middleware.CORSMiddleware())
	api.SetupRoutes(a.router)

	// 7. 创建 HTTP 服务器
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: a.router,
	}
}

func (a *App) Run() {
	// 启动服务器
	go func() {
		log.Printf("server starting on %s", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.Shutdown()
}

func (a *App) Shutdown() {
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	database.CloseNeo4j()
	log.Println("server exited")
}
