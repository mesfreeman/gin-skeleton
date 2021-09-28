package helper

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// GormDefaultDb mysql默认连接
	GormDefaultDb *gorm.DB

	// RedisDefaultDb redis默认连接
	RedisDefaultDb *redis.Client

	// TimeFormatYmdhis 时间格式 yyyy-mm-dd hh:ii:ss
	TimeFormatYmdhis = "2006-01-02 15:04:05"

	// TimeFormatYmd 时间格式 yyyy-mm-dd
	TimeFormatYmd = "2006-01-02"
)
