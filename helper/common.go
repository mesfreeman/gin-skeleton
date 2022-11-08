package helper

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/spf13/viper"
)

// GetRootPath 获取项目根目录
func GetRootPath() string {
	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(fmt.Sprintf("get project root path faild: %s", err))
	}
	return rootPath
}

// GetRandomString 返回指定长度的随机字符串
func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}

// SubStr 截取指定长度的字段串
func SubStr(str string, num int) string {
	tmp := []rune(str)
	if len(tmp) > num {
		return string(tmp[:num])
	}
	return string(tmp)
}

// Calmel2Case 大小驼峰转下划线格式
func Calmel2Case(str string) string {
	bts := make([]byte, 0)
	for idx, val := range str {
		if unicode.IsUpper(val) && idx != 0 {
			bts = append(bts, '_')
		}
		bts = append(bts, byte(unicode.ToLower(val)))
	}
	return string(bts)
}

// InSilce 检查切片中是否存在值
func InSilce[T comparable](needle T, haystack []T) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// DiffSilce 比较切片中差异的元素
func DiffSilce[T comparable](targetObj []T, compareObj []T) []T {
	var res []T
	for _, to := range targetObj {
		if !InSilce(to, compareObj) {
			res = append(res, to)
		}
	}
	return res
}

// IsDevelopmentEnv 检查是否为开发环境
func IsDevelopmentEnv() bool {
	return viper.GetString("Server.Mode") == "development"
}

// IsTestEnv 检查是否为测试环境
func IsTestEnv() bool {
	return viper.GetString("Server.Mode") == "test"
}

// IsProductEnv 检查是否为线上环境
func IsProductEnv() bool {
	return viper.GetString("Server.Mode") == "production"
}

// IsLegalUrl 检查是否为合法的网址(暂时简单判断一下)
func IsLegalUrl(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}

// IsSuperAccount 检查是否为超级账号
func IsSuperAccount(username string) bool {
	return viper.GetString("Server.SuperAccount") == username
}
