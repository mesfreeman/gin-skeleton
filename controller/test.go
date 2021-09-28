package controller

import (
	"gin-skeleton/helper"
	"gin-skeleton/helper/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Test 测试动作
func Test(c *gin.Context) {
	// 开发者可以在这里加上自己的任意的测试代码，但是测试代码不应提交到仓库中！

	var result = struct {
		Mode        string `json:"mode"`
		Welcome     string `json:"welcome"`
		ClientIp    string `json:"clientIp"`
		CurrentTime string `json:"currentTime"`
	}{
		Mode:        viper.GetString("Server.Mode"),
		Welcome:     "hello world!",
		ClientIp:    c.ClientIP(),
		CurrentTime: time.Now().Format(helper.TimeFormatYmdhis),
	}
	response.SuccessJSON(result, "success", c)
}
