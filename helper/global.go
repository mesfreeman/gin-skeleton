package helper

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// Version 当前版本号
	Version string

	// BuildTime 编译时间
	BuildTime string

	// StartTime 服务启动时间
	StartTime string

	// GormDefaultDb mysql默认连接
	GormDefaultDb *gorm.DB

	// RedisDefaultDb redis默认连接
	RedisDefaultDb *redis.Client

	// ToDateTimeString 时间格式 yyyy-mm-dd hh:ii:ss
	ToDateTimeString = "2006-01-02 15:04:05"
)
