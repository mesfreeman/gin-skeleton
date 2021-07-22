package middleware

import (
	"gin-skeleton/helper"
	"gin-skeleton/helper/response"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	// 防重放攻击最大间隔，单位：秒
	restMaxTime = 600
)

// 密钥映射
var secretKeys = map[int]string{
	10000: "yghi6vnwpc35kmj1tdxbea7zq02o8lf4",
}

// 签名授权
func SignAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params struct {
			AccessKey int    `form:"accessKey" json:"accessKey" uri:"accessKey"`
			Signature string `form:"signature" json:"signature" uri:"signature"`
			Timestamp int    `form:"timestamp" json:"timestamp" uri:"timestamp"`
			Nonce     string `form:"noce" json:"noce" uri:"noce"`
		}
		c.ShouldBind(&params)
		logger := helper.GetLogger("sign").WithFields(logrus.Fields{"params": c.Request.Form})

		// 参数校验
		if params.AccessKey == 0 || params.Signature == "" || params.Timestamp == 0 || params.Nonce == "" {
			logger.Warnln("签名参数缺失")
			response.InvalidArgumentJSON("签名参数缺失", c)
			c.Abort()
			return
		}

		// nonce唯一性检查，10分钟内唯一，防止重放攻击
		redisKey := "signAuthNonce:" + strconv.Itoa(params.AccessKey) + ":" + params.Nonce
		if helper.RedisDefaultDb.Exists(helper.RedisDefaultDb.Context(), redisKey).Val() > 0 {
			logger.Warnln("重复请求")
			response.InvalidAuthJSON("重复请求", c)
			c.Abort()
			return
		}
		if math.Abs(float64(time.Now().Unix()-int64(params.Timestamp))) > restMaxTime {
			logger.Warnln("签名已过期")
			response.InvalidAuthJSON("签名已过期", c)
			c.Abort()
			return
		}

		secretKey, ok := secretKeys[params.AccessKey]
		if !ok {
			logger.Warnln("密钥不存在")
			response.InvalidAuthJSON("密钥不存在", c)
			c.Abort()
			return
		}

		// 签名校验
		localSignature := helper.GetMD5("accessKey=" + strconv.Itoa(params.AccessKey) + "&secretKey=" + secretKey + "&timestamp=" + strconv.Itoa(params.Timestamp) + "&nonce=" + params.Nonce)
		if params.Signature != localSignature {
			logger.Warnln("签名验证失败", params.Signature, localSignature)
			response.InvalidAuthJSON("签名验证失败", c)
			c.Abort()
			return
		}

		helper.RedisDefaultDb.SetEX(helper.RedisDefaultDb.Context(), redisKey, 1, restMaxTime*time.Second)
		c.Next()
	}
}
