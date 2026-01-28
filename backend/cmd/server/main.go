package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"online-mall/internal/api/routes"
	"online-mall/internal/config"
	"online-mall/internal/utils"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := utils.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer utils.CloseDB()

	// 初始化Redis
	if err := utils.InitRedis(); err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}
	defer utils.CloseRedis()

	// 设置路由
	r := routes.SetupRoutes()

	// 获取配置
	cfg := config.GlobalConfig

	// 创建服务器
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 启动服务器
	log.Printf("Server starting on port %d...", cfg.App.Port)
	log.Printf("Environment: %s", cfg.App.Name)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server shutting down...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
