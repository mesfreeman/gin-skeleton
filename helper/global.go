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
)
