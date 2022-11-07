package main

import (
	"context"
	"gin-skeleton/helper"
	"gin-skeleton/provider"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	provider.InitConfig()
	// provider.InitGormDB()
	// provider.InitRedisDB()

	// 服务配置
	router := provider.Routers()
	server := &http.Server{
		Addr:           ":" + viper.GetString("Server.Port"),
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln("Listen server error: ", err, syscall.Getpid())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	helper.GetLogger("").Warnln("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		helper.GetLogger("").Fatalln("Shutdown server error: ", err)
	}
	helper.GetLogger("").Warnln("Server exiting")
}
