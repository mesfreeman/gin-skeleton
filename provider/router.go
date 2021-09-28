package provider

import (
	"gin-skeleton/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Routers 路由
func Routers() *gin.Engine {
	if viper.GetString("Server.Mode") == "production" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	var Router = gin.Default()

	// 公共路由
	PublicGroup := Router.Group("/")
	{
		router.InitTestRouter(PublicGroup)
	}

	return Router
}
