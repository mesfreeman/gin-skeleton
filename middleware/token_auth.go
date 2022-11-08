package middleware

import (
	"gin-skeleton/helper/jwt"
	"gin-skeleton/helper/response"

	"strings"

	"github.com/gin-gonic/gin"
)

// TokenAuth Token授权
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		if accessToken == "" {
			response.UnauthorizedJSON("请先登录", c)
			c.Abort()
			return
		}

		// 解析Token
		claims, err := jwt.ParseJwtToken(accessToken)
		if err != nil {
			response.UnauthorizedJSON("请重新登录", c)
			c.Abort()
			return
		}

		if claims.ID <= 0 {
			response.UnauthorizedJSON("令牌不合法", c)
			c.Abort()
			return
		}

		// 将ID合并到请求中（注意：请求参数不要使用到该参数）
		c.Set("identity", claims)
		c.Next()
		return
	}
}

// GetTokenAuthInfo 获取令牌信息
func GetTokenAuthInfo(c *gin.Context) *jwt.Claims {
	claims, exist := c.Get("identity")
	if !exist {
		return &jwt.Claims{}
	}
	return claims.(*jwt.Claims)
}
