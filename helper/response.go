package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 常用的状态码
const (
	successCode         = 0     // 成功
	invalidAuthCode     = 40100 // 无效的授权
	invalidArgumentCode = 40400 // 无效的参数
	logicExceptionCode  = 50100 // 逻辑异常
)

// SuccessJSON 成功时返回
func SuccessJSON(data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": successCode,
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

// InvalidAuthJSON 授权不合法
func InvalidAuthJSON(msg string, c *gin.Context) {
	FailJSON(invalidAuthCode, msg, c)
}

// InvalidArgumentJSON 参数不合法
func InvalidArgumentJSON(msg string, c *gin.Context) {
	FailJSON(invalidArgumentCode, msg, c)
}

// LogicExceptionJSON 逻辑异常
func LogicExceptionJSON(msg string, c *gin.Context) {
	FailJSON(logicExceptionCode, msg, c)
}
