package provider

import (
	"gin-skeleton/helper"
	"gin-skeleton/router"
	"gin-skeleton/router/admin"

	"github.com/gin-gonic/gin"
)

// Routers 路由
func Routers() *gin.Engine {
	if helper.IsProductEnv() {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}

	var Router = gin.Default()

	// 公共路由
	PublicGroup := Router.Group("/")
	{
		router.InitTestRouter(PublicGroup)
	}

	// 后台路由
	AdminGroup := Router.Group("/admin")
	{
		admin.InintCommonRouter(AdminGroup)
		admin.InintSystemRouter(AdminGroup)
	}

	return Router
}
