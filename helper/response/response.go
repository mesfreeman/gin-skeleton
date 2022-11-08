package response

import (
	"net/http"

	"gin-skeleton/helper/validator"

	"github.com/gin-gonic/gin"
)

// 业务异常状态码
const (
	success         = 0     // 响应成功
	invalidArgument = 40000 // 无效参数
	logicException  = 40020 // 逻辑异常
	thirdException  = 40040 // 三方异常
)

// 业务码默认消息提示
var businessCodeMessages = map[int]string{
	success:         "成功",
	invalidArgument: "参数错误",
	logicException:  "系统出错了",
	thirdException:  "服务出错了",
}

// SuccessJSON 成功时返回
func SuccessJSON(r interface{}, m string, c *gin.Context) {
	responseJson(http.StatusOK, success, r, m, c)
}

// ValidatorFailedJson 校验器未通过
func ValidatorFailedJson(err error, c *gin.Context) {
	responseJson(http.StatusOK, invalidArgument, "", validator.Translate(err), c)
}

// InvalidArgumentJSON 无效的请求参数
func InvalidArgumentJSON(m string, c *gin.Context) {
	responseJson(http.StatusOK, invalidArgument, "", m, c)
}

// LogicExceptionJSON 逻辑处理异常
func LogicExceptionJSON(m string, c *gin.Context) {
	responseJson(http.StatusOK, logicException, "", m, c)
}

// ThirdExceptionJSON 三方服务异常
func ThirdExceptionJSON(m string, c *gin.Context) {
	responseJson(http.StatusOK, thirdException, "", m, c)
}

// UnauthorizedJSON 未经授权的请求
func UnauthorizedJSON(m string, c *gin.Context) {
	responseJson(http.StatusUnauthorized, success, "", m, c)
}

// ForbiddenJSON 禁止访问的资源
func ForbiddenJSON(m string, c *gin.Context) {
	responseJson(http.StatusForbidden, success, "", m, c)
}

// RuntimeExceptionJSON 运行时异常(一般指缺少必要配置)
func RuntimeExceptionJSON(m string, c *gin.Context) {
	responseJson(http.StatusInternalServerError, success, "", m, c)
}

// responseJson 响应JSON结构
func responseJson(httpCode int, businessCode int, r interface{}, m string, c *gin.Context) {
	// 提示消息为空时，走默认提示
	msg, ok := businessCodeMessages[businessCode]
	if m == "" && ok {
		m = msg
	}

	// 返回JSON响应体
	c.JSON(httpCode, gin.H{
		"code":    businessCode,
		"result":  r,
		"message": m,
	})
}
