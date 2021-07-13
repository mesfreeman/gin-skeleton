package provider

import (
	"gin-skeleton/helper"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitRedisDB() {
	// 扩展其它库 ...
	helper.RedisDefaultDb = getRedisDb("default")
}

func getRedisDb(connection string) *redis.Client {
	host := viper.GetString("redis." + connection + ".host")
	port := viper.GetInt("redis." + connection + ".port")
	database := viper.GetInt("redis." + connection + ".database")
	password := viper.GetString("redis." + connection + ".password")

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + strconv.Itoa(port),
		Password: password,
		DB:       database,
	})

	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		helper.GetLogger("").WithFields(logrus.Fields{"connection": connection, "err": err}).Fatalln("Redis start err")
		return nil
	}
	return rdb
}
