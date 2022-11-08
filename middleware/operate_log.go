package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-skeleton/helper"
	"gin-skeleton/helper/log"
	"gin-skeleton/model/admin/system"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

// 操作日志黑名单
var blackPaths = []string{
	"/admin/system/operateLog/list",
}

// OperateLog 操作日志
func OperateLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		startAt := time.Now()
		ip := c.ClientIP()
		path := c.FullPath()
		userInfo := GetTokenAuthInfo(c)
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		params, remark := getParams(c)

		// 判断该操作是否忽略记录
		if helper.InSilce(path, blackPaths) {
			c.Next()
			return
		}

		// 组合响应体
		crw := &customerResponseWriter{c.Writer, bytes.NewBufferString("")}
		c.Writer = crw
		c.Next()

		// 查找接口名称
		function, err := system.NewMenu().FindNameByPath(path)
		if err != nil {
			remark += fmt.Sprintf("查询接口名称出错了：%s;", err.Error())
			log.GetLogger("operate-log").Errorln(remark)
		}

		// 如果是上传接口，将忽略请求参数
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			params = ""
			remark += "该操作请求参数忽略记录"
		}

		// 判断是否有异常输出，追加到备注字段里
		sysErr := c.Errors.ByType(gin.ErrorTypePrivate).String()
		if sysErr != "" {
			remark += fmt.Sprintf("系统出错了：%s;", sysErr)
		}

		// 业务码
		code := "-"
		body := crw.body.String()
		if gjson.Get(body, "code").Exists() {
			code = gjson.Get(body, "code").String()
		}

		// 创建操作日志
		operateLog := system.OperateLog{
			Username:  userInfo.Username,
			Nickname:  userInfo.Nickname,
			Ip:        ip,
			Function:  function,
			Uri:       path,
			Method:    method,
			Params:    helper.SubStr(params, 1000),
			Status:    c.Writer.Status(),
			Code:      code,
			SpendTime: time.Since(startAt).Milliseconds(),
			Result:    helper.SubStr(body, 5000),
			UserAgent: helper.SubStr(userAgent, 500),
			Remark:    helper.SubStr(remark, 500),
		}
		if err = helper.GormDefaultDb.Create(&operateLog).Error; err != nil {
			log.GetLogger("operate-log").Errorln("创建操作日志出错了：" + err.Error())
		}
	}
}

// 获取请求参数
func getParams(c *gin.Context) (params, remark string) {
	// GET请求
	if c.Request.Method == http.MethodGet {
		// 获取请求参数并转码
		query, err := url.QueryUnescape(c.Request.URL.RawQuery)
		if err != nil {
			remark = fmt.Sprintf("获取GET请求参数转码出错了：%s;", err.Error())
			log.GetLogger("operate-log").Errorln(remark)
			return
		}

		// 格式化数据
		maps := make(map[string]string)
		for _, val := range strings.Split(query, "&") {
			kv := strings.Split(val, "=")
			if len(kv) == 2 {
				maps[kv[0]] = kv[1]
			}
		}

		// 转码为JSON字符串
		bt, _ := json.Marshal(&maps)
		params = string(bt)
		return
	}

	// 非GET请求（需要将原数据重新放回去，因为io.ReadAll会清空c.Request.Body中的数据）
	body, err := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		remark = fmt.Sprintf("获取POST请求参数出错了：%s;", err.Error())
		log.GetLogger("operate-log").Errorln(remark)
		return
	}

	params = string(body)
	return
}

// 自定义响应体组合
type customerResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写写入方法
func (crw customerResponseWriter) Write(bt []byte) (int, error) {
	crw.body.Write(bt)
	return crw.ResponseWriter.Write(bt)
}
