package router

import (
	"gin-skeleton/controller"

	"github.com/gin-gonic/gin"
)

// InitTestRouter 测试接口相关
func InitTestRouter(Router *gin.RouterGroup) {
	TestRouter := Router.Group("/test")
	{
		TestRouter.GET("/test", controller.Test)
		TestRouter.POST("/test", controller.Test)
	}
}
