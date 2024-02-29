package admin

import (
	"gin-skeleton/controller/admin"
	"gin-skeleton/middleware"

	"github.com/gin-gonic/gin"
)

// InitCommonRouter 初始化公共功能相关路由
func InitCommonRouter(Router *gin.RouterGroup) {
	// 登录登出
	Router.POST("/login/login", admin.Login)
	Router.Use(middleware.TokenAuth()).POST("/login/logout", admin.Logout)

	// 基础功能
	BaseRouter := Router.Group("/base").Use(middleware.TokenAuth())
	{
		BaseRouter.POST("/myInfo", admin.MyInfo)
		BaseRouter.POST("/myMenus", admin.MyMenus)
		BaseRouter.POST("/myPerms", admin.MyPerms)
		BaseRouter.POST("/modifyMyPwd", admin.ModifyMyPwd)
		BaseRouter.POST("/modifyMyInfo", admin.ModifyMyInfo)
		BaseRouter.POST("/liteRoles", admin.LiteRoles)
		BaseRouter.POST("/liteAccounts", admin.LiteAccounts)
		BaseRouter.POST("/uploadFile", admin.UploadFile)
	}
}
