package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"mall/ioc"
)

func main() {
	initViper()

	router := ioc.InitGin()

	server := &http.Server{
		Addr:    "0.0.0.0:9000",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	fmt.Println("Server is running on http://localhost:9000")

	// 创建通道监听信号
	quit := make(chan os.Signal, 1)

	// 监听信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞直到收到信号
	<-quit
	fmt.Println("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	// 优雅地关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced shutting down:", err)
	}

	fmt.Println("Server exited gracefully")
}

func initViper() {
	viper.SetConfigFile("config\\dev.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
