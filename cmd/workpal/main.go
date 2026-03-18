package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aK1r4z/workpal/internal/auth"
	"github.com/aK1r4z/workpal/internal/store/postgres"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// 初始化资源
	ctx := context.TODO()

	// 获取环境变量
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// 创建 Postgres 数据库连接
	connString := os.Getenv("CONNECTION_STRING")

	db, err := postgres.New(ctx, connString)
	if err != nil {
		panic(err)
	}

	// 创建 HTTP 请求处理器
	e := echo.New()
	e.Use(
		middleware.RequestLogger(),
	)

	// 创建用户认证服务
	auth.Config.Load()

	authService := auth.NewService(db.UserStore())
	authHandler := auth.NewHandler(authService)

	authHandler.RegisterRoutes(e)

	// 创建 HTTP 服务器
	s := &http.Server{
		Handler: e,
		Addr:    ":8080",
		// ReadTimeout: 30 * time.Second,
		// WriteTimeout: 30 * time.Second,
	}

	// 在协程中启动服务器
	go func() {
		if err := s.ListenAndServe(); err != http.ErrServerClosed {
			e.Logger.Error("failed to start server", "error", err)
		}
	}()

	log.Println("server must be on now.")

	// 优雅终结流程
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	<-shutdown

	log.Println("gracefully shutting down the server...")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		e.Logger.Error("failed to shutdown server", "error", err)
	}

	log.Println("server closed.")
}
