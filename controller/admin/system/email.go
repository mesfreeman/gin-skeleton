package system

import (
	"encoding/json"
	"gin-skeleton/helper/response"
	"gin-skeleton/helper/tool"
	"gin-skeleton/model"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
)

// EmailViewConfig 查看邮件配置
func EmailViewConfig(c *gin.Context) {
	email, err := tool.FindEmail()
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(email, "", c)
}

// EmailSaveConfig 保存邮件配置
func EmailSaveConfig(c *gin.Context) {
	var params struct {
		Server   string `json:"server" remark:"收信服务器" binding:"required"`
		Port     string `json:"port" remark:"端口" binding:"required"`
		Sender   string `json:"sender" remark:"发件人" binding:"required"`
		Account  string `json:"account" remark:"账号" binding:"required"`
		Password string `json:"password" remark:"密码" binding:"required"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 转换为JSON字符串
	value, err := json.Marshal(&params)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	commonConfig := system.CommonConfig{
		Module:  system.CommonConfigModuleEmailServer,
		Keyword: "",
		Value:   string(value),
		Remark:  "邮件服务相关配置",
	}
	id, err := system.NewCommonConfig().CreateOrUpdateConfig(&commonConfig)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(model.BaseIdResult{ID: id}, "", c)
}
