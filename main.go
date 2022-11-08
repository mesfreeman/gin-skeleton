package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gin-skeleton/helper/log"
	"gin-skeleton/helper/validator"
	"gin-skeleton/provider"

	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	provider.InitConfig()
	provider.InitGormDB()
	provider.InitRedisDB()
	validator.InitTranslator()

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
			panic("Listen server error: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.GetLogger("").Warnln("Shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic("Shutdown server error: " + err.Error())
	}
	log.GetLogger("").Warnln("Server exiting")
}
