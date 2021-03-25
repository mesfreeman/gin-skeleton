package middleware

import (
	"gin-skeleton/helper"
	"gin-skeleton/util"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	restMaxTime = 300 // 防重放攻击最大间隔，单位：秒
)

// 密钥映射
var appSecretMapper map[int]string = map[int]string{
	10000: "yghi6vnwpc35kmj1tdxbea7zq02o8lf4",
}

// 签名授权
func SignAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		appId, _ := strconv.Atoi(c.Query("appId"))
		signature := c.Query("signature")
		timestamp, _ := strconv.Atoi(c.Query("timestamp"))
		params := c.Request.URL.Query()
		logger := logrus.WithFields(logrus.Fields{"params": params})

		// 参数校验
		if appId == 0 || signature == "" || timestamp == 0 {
			logger.Warnln("部分签名参数缺失")
			helper.InvalidArgumentJSON("部分签名参数缺失", c)
			c.Abort()
			return
		}

		appSecret, ok := appSecretMapper[appId]
		if !ok {
			logger.Warnln("appId不存在")
			helper.InvalidAuthJSON("appId不存在", c)
			c.Abort()
			return
		}

		if math.Abs(float64(time.Now().Unix()-int64(timestamp))) > restMaxTime {
			logger.Warnln("签名已过期")
			helper.InvalidAuthJSON("签名已过期", c)
			c.Abort()
			return
		}

		// 签名校验
		signRet := paramSign(appSecret, params)
		if signature != signRet {
			logger.Warnln("签名验证失败", signRet, signature)
			helper.InvalidAuthJSON("签名验证失败", c)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 参数签名
func paramSign(appSecret string, params map[string][]string) string {
	// 去掉部分参数
	delete(params, "signature")

	// 对参数名排序
	paramKeys := make([]string, 0, len(params))
	for key := range params {
		paramKeys = append(paramKeys, key)
	}

	// 兼容PHP的ksort 数字会在字母后面
	sort.Slice(paramKeys, func(i, j int) bool {
		_, err1 := strconv.Atoi(paramKeys[i])
		_, err2 := strconv.Atoi(paramKeys[j])
		if err1 == nil && err2 != nil {
			return false
		}
		if err1 != nil && err2 == nil {
			return true
		}
		return paramKeys[j] > paramKeys[i]
	})

	// 拼接签名字符串
	signStr := ""
	for _, key := range paramKeys {
		// 注意：不支持数组或对象参数
		signStr += key + "=" + params[key][0]
	}
	signStr += appSecret

	// MD5加密
	return util.MD5(signStr)
}
