package helper

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// mysql默认连接
	GormDefaultDb *gorm.DB

	// redis默认连接
	RedisDefaultDb *redis.Client
)
