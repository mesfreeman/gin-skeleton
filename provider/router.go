package provider

import (
	"gin-skeleton/router"

	"github.com/gin-gonic/gin"
)

// Routers 路由
func Routers() *gin.Engine {
	var Router = gin.Default()

	// 公共路由
	PublicGroup := Router.Group("")
	{
		router.InitTestRouter(PublicGroup)
	}

	return Router
}
