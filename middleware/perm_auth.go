package middleware

import (
	"gin-skeleton/helper"
	"gin-skeleton/helper/response"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
)

// PermAuth 权限授权
func PermAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		aid := GetTokenAuthInfo(c).ID
		username := GetTokenAuthInfo(c).Username
		if aid == 0 || username == "" {
			response.InvalidArgumentJSON("登录信息为空", c)
			c.Abort()
			return
		}

		// 超级账号，忽略权限检查
		if helper.IsSuperAccount(username) {
			c.Next()
			return
		}

		// 权限检查
		hasPermission, err := system.NewAccount().HasPermission(aid, username, c.FullPath())
		if err != nil {
			response.LogicExceptionJSON("系统出错了："+err.Error(), c)
			c.Abort()
			return
		}
		if !hasPermission {
			response.ForbiddenJSON("你没有权限", c)
			c.Abort()
			return
		}
		c.Next()
		return
	}
}
