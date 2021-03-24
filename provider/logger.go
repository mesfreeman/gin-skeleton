package provider

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 日志配置
func InitLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)
	logout, _ := os.OpenFile("./storage/logs/gin-skeleton.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModeAppend)
	logrus.SetOutput(logout)

	// 根据环境设置日志等级
	if viper.GetString("app.mode") == "release" {
		logrus.SetLevel(logrus.WarnLevel)
	}
}
