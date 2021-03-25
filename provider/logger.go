package provider

import (
	"gin-skeleton/util"
	"log"
	"os"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 日志配置
func InitLogger() {
	// 记录文件名和行号
	if viper.GetString("app.logLevel") == "debug" {
		logrus.SetReportCaller(true)
	}

	// 设置日志输出格式：json, text
	switch viper.GetString("app.logOutFormat") {
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logrus.SetFormatter(&logrus.TextFormatter{})
	}

	// 设置日志等级：debug, info, error, warn, panic, fatal
	switch viper.GetString("app.logLevel") {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.DebugLevel)
	}

	// 设置日志输出路径：console, file
	switch viper.GetString("app.logOutPath") {
	case "file":
		// 日志打印到指定的目录
		logFileName := path.Join(util.GetRootPath(), "storage", "logs", viper.GetString("app.name"), ".log")
		logOut, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
		if err != nil {
			log.Fatal("Open log file fail: ", err)
			break
		}
		logrus.SetOutput(logOut)

		// 创建日志输出对象
		var logMaxSaveDay = time.Duration(viper.GetInt("app.logMaxSaveDay"))
		logWriter, err := rotatelogs.New(
			logFileName+".%Y-%m-%d.log",                              // 日志切割名称
			rotatelogs.WithLinkName(logFileName),                     // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(logMaxSaveDay*24*time.Hour),        // 文件最大保存时间
			rotatelogs.WithRotationTime(time.Duration(24)*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			log.Fatal("Create rotatelogs object fail: ", err)
		}

		// 为不同级别设置不同的输出目的
		writeMap := lfshook.WriterMap{
			logrus.InfoLevel:  logWriter,
			logrus.FatalLevel: logWriter,
			logrus.DebugLevel: logWriter,
			logrus.WarnLevel:  logWriter,
			logrus.ErrorLevel: logWriter,
			logrus.PanicLevel: logWriter,
		}

		// 创建logrus的本地文件系统钩子
		lfHook := lfshook.NewHook(writeMap, &logrus.JSONFormatter{})
		logrus.AddHook(lfHook)
	default:
		logrus.SetOutput(os.Stdout)
	}
}
