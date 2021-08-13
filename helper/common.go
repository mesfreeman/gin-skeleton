package helper

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 日志对象集合
var loggerMap sync.Map

// GetLogger 获取日志对象
func GetLogger(fileName string) *logrus.Logger {
	if fileName == "" {
		fileName = viper.GetString("app.name")
	}

	// 日志配置
	logLevel := viper.GetString("app.logLevel")
	logOutFormat := viper.GetString("app.logOutFormat")
	logOutPath := viper.GetString("app.logOutPath")
	logMaxSaveDay := viper.GetInt("app.logMaxSaveDay")

	// 如果存在，直接返回
	key := fmt.Sprintf("%s_%s_%s_%s_%d", fileName, logLevel, logOutFormat, logOutPath, logMaxSaveDay)
	if logger, ok := loggerMap.Load(key); ok {
		return logger.(*logrus.Logger)
	}

	// 创建日志对象
	logger := logrus.New()

	// 设置日志输出格式：json, text
	switch logOutFormat {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	// 设置日志等级：debug, info, error, warn, panic, fatal
	switch logLevel {
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.DebugLevel)
		logger.SetReportCaller(true) // 记录文件名和行号
	}

	// 设置日志输出路径：console, file
	switch logOutPath {
	case "file":
		// 日志打印到指定的目录
		logFileName := path.Join(GetRootPath(), "storage", "logs", (fileName + ".log"))
		logOut, err := os.OpenFile(logFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			log.Fatal("Open log file fail: ", err)
			break
		}
		logger.SetOutput(logOut)

		// 创建日志输出对象
		var logMaxSaveDay = time.Duration(logMaxSaveDay)
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
		logger.AddHook(lfHook)
	default:
		logger.SetOutput(os.Stdout)
	}

	loggerMap.Store(key, logger)
	return logger
}

// GetRootPath 获取项目根目录
func GetRootPath() string {
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("get project root path faild: %s", err))
	}
	return rootPath
}

// GetMD5 MD5 加密
func GetMD5(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
