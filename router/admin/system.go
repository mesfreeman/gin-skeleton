package admin

import (
	"gin-skeleton/controller/admin/system"
	"gin-skeleton/middleware"

	"github.com/gin-gonic/gin"
)

// InitSystemRouter 初始化系统管理相关路由
func InitSystemRouter(Router *gin.RouterGroup) {
	SystemRouter := Router.Group("/system").Use(middleware.TokenAuth(), middleware.PermAuth(), middleware.OperateLog())
	{
		// 账号管理
		SystemRouter.POST("/account/list", system.AccountList)
		SystemRouter.POST("/account/add", system.AccountAdd)
		SystemRouter.POST("/account/modify", system.AccountModify)
		SystemRouter.POST("/account/modifyPwd", system.AccountModifyPwd)
		SystemRouter.POST("/account/delete", system.AccountDelete)

		// 角色管理
		SystemRouter.POST("/role/list", system.RoleList)
		SystemRouter.POST("/role/add", system.RoleAdd)
		SystemRouter.POST("/role/modify", system.RoleModify)
		SystemRouter.POST("/role/delete", system.RoleDelete)

		// 菜单管理
		SystemRouter.POST("/menu/list", system.MenuList)
		SystemRouter.POST("/menu/add", system.MenuAdd)
		SystemRouter.POST("/menu/modify", system.MenuModify)
		SystemRouter.POST("/menu/delete", system.MenuDelete)

		// 文件管理
		SystemRouter.POST("/file/list", system.FileList)
		SystemRouter.POST("/file/delete", system.FileDelete)
		SystemRouter.POST("/file/viewConfig", system.FileViewConfig)
		SystemRouter.POST("/file/saveConfig", system.FileSaveConfig)

		// 邮件配置
		SystemRouter.POST("/email/viewConfig", system.EmailViewConfig)
		SystemRouter.POST("/email/saveConfig", system.EmailSaveConfig)
		SystemRouter.POST("/email/templateList", system.EmailTemplateList)
		SystemRouter.POST("/email/templateAdd", system.EmailTemplateAdd)
		SystemRouter.POST("/email/templateModify", system.EmailTemplateModify)
		SystemRouter.POST("/email/templateDelete", system.EmailTemplateDelete)
		SystemRouter.POST("/email/batchSend", system.EmailBatchSend)

		// 登录日志
		SystemRouter.POST("/loginLog/list", system.LoginLogList)

		// 操作日志
		SystemRouter.POST("/operateLog/list", system.OperateLogList)
	}
}
