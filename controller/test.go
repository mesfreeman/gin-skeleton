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

	result := struct {
		Mode        string `json:"mode"`
		Version     string `json:"version"`
		Welcome     string `json:"welcome"`
		ClientIp    string `json:"clientIp"`
		BuildTime   string `json:"buildTime"`
		StartTime   string `json:"startTime"`
		CurrentTime string `json:"currentTime"`
	}{
		Mode:        viper.GetString("Server.Mode"),
		Welcome:     "Hello World!",
		Version:     helper.Version,
		ClientIp:    c.ClientIP(),
		BuildTime:   helper.BuildTime,
		StartTime:   helper.StartTime,
		CurrentTime: time.Now().Format(helper.ToDateTimeString),
	}
	response.SuccessJSON(result, "", c)
}
