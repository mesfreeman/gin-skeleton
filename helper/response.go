package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 常用的状态码
const (
	SuccessCode         = 0     // 成功
	InvalidAuthCode     = 40100 // 无效的授权
	InvalidArgumentCode = 40400 // 无效的参数
	LogicExceptionCode  = 50100 // 逻辑异常
)

// SuccessJSON 成功时返回
func SuccessJSON(data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": SuccessCode,
		"data": data,
		"msg":  msg,
	})
}

// FailJSON 失败时返回
func FailJSON(code int, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"data": "",
		"msg":  msg,
	})
}

// InvalidAuthErrorJSON 授权不合法
func InvalidAuthErrorJSON(msg string, c *gin.Context) {
	FailJSON(InvalidAuthCode, msg, c)
}

// InvalidArgumentErrorJSON 参数不合法
func InvalidArgumentErrorJSON(msg string, c *gin.Context) {
	FailJSON(InvalidArgumentCode, msg, c)
}

// LogicExceptionJSON 逻辑异常
func LogicExceptionJSON(msg string, c *gin.Context) {
	FailJSON(LogicExceptionCode, msg, c)
}
