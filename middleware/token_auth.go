package middleware

import (
	"gin-skeleton/helper/response"
	"gin-skeleton/helper/tool"

	"github.com/gin-gonic/gin"
)

// TokenAuth Token授权
func TokenAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := getAccessToken(c)
		if accessToken == "" {
			response.InvalidArgumentJSON("请先登录", c)
			c.Abort()
			return
		}

		// 解析Token
		id, err := tool.ParseJwtToken(accessToken)
		if err != nil {
			response.InvalidAuthJSON("请重新登录", c)
			c.Abort()
			return
		}

		if id <= 0 {
			response.InvalidAuthJSON("令牌不合法", c)
			c.Abort()
			return
		}

		// 将ID合并到请求中（注意：请求参数不要使用到该参数）
		c.Set("identity", id)
		c.Next()
	}
}

// 获取令牌字符串 header > cookie > post > get
func getAccessToken(c *gin.Context) string {
	accessToken := c.GetHeader("accessToken")
	if accessToken != "" {
		return accessToken
	}

	accessToken, _ = c.Cookie("accessToken")
	if accessToken != "" {
		return accessToken
	}

	accessToken = c.PostForm("accessToken")
	if accessToken != "" {
		return accessToken
	}

	accessToken = c.Query("accessToken")
	if accessToken != "" {
		return accessToken
	}

	return ""
}
