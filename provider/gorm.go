package provider

import (
	"fmt"
	"strings"
	"time"

	"gin-skeleton/helper"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitGormDB GormDB初始化
func InitGormDB() {
	// 扩展其它库 ...
	helper.GormDefaultDb = gormMysql("default")
}

// MYSQL驱动
func gormMysql(connection string) *gorm.DB {
	connection = strings.ToUpper(connection)
	host := viper.GetString("Gorm." + connection + ".Host")
	port := viper.GetInt("Gorm." + connection + ".Port")
	database := viper.GetString("Gorm." + connection + ".Database")
	username := viper.GetString("Gorm." + connection + ".Username")
	password := viper.GetString("Gorm." + connection + ".Password")
	charset := viper.GetString("Gorm." + connection + ".Charset")

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

	// Gorm相关配置
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true, // 禁用默认的事务
	}

	// 开发环境，打印执行的SQL语句
	if helper.IsDevelopmentEnv() {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	// 打开链接
	db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig)
	if err != nil {
		fmt.Println("Gorm mysql start err: ", err, connection)
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db
}
