package provider

import (
	"fmt"
	"strings"

	"gin-skeleton/helper"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// InitRedisDB redis初始化
func InitRedisDB() {
	// 扩展其它库 ...
	helper.RedisDefaultDb = getRedisDb("default")
}

func getRedisDb(connection string) *redis.Client {
	connection = strings.ToUpper(connection)
	host := viper.GetString("Redis." + connection + ".Host")
	port := viper.GetInt("Redis." + connection + ".Port")
	database := viper.GetInt("Redis." + connection + ".Database")
	password := viper.GetString("Redis." + connection + ".Password")

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: password,
		DB:       database,
	})

	if err := rdb.Ping(rdb.Context()).Err(); err != nil {
		fmt.Println("Redis start err: ", err, connection)
		return nil
	}
	return rdb
}
