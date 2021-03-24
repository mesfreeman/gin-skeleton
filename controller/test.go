package controller

import (
	"gin-skeleton/helper"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// Test 测试动作
func Test(c *gin.Context) {
	// 开发者可以在这里加上自己的任意的测试代码，但是测试代码不应提交到仓库中！

	result := map[string]string{
		"welcome":  "hello world!",
		"time":     time.Now().Format("2006-01-02 15:04:05"),
		"clientIp": c.ClientIP(),
		"mode":     viper.GetString("app.mode"),
	}

	helper.SuccessJSON(result, "success", c)
}
