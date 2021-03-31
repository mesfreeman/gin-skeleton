package provider

import (
	"fmt"
	"gin-skeleton/helper"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GormDB初始化
func InitGormDB() {
	// 扩展其它库 ...
	helper.GormDefaultDb = gormMysql("default")
}

// MYSQL驱动
func gormMysql(connection string) *gorm.DB {
	host := viper.GetString("gorm." + connection + ".host")
	database := viper.GetString("gorm." + connection + ".database")
	port := viper.GetInt("gorm." + connection + ".port")
	username := viper.GetString("gorm." + connection + ".username")
	password := viper.GetString("gorm." + connection + ".password")
	charset := viper.GetString("gorm." + connection + ".charset")

	// 拼接mysql相关配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", username, password, host, port, database, charset)
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}

	// 打开链接
	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{"connection": connection, "err": err}).Fatalln("Gorm mysql start err")
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
