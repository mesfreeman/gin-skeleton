package system

import (
	"encoding/json"
	"errors"
	"sync"
	"sync/atomic"

	"gin-skeleton/helper"
	"gin-skeleton/helper/log"
	"gin-skeleton/helper/response"
	"gin-skeleton/helper/tool"
	"gin-skeleton/model"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

// EmailTemplateList 获取邮件模板列表
func EmailTemplateList(c *gin.Context) {
	var params struct {
		model.BasePageParams
		Subject string `json:"subject" remark:"邮件主题" binding:"max=64"`
		Slug    string `json:"slug" remark:"邮件标识" binding:"max=32"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	emailTemplateList, err := system.NewEmailTemplate().GetEmailTemplateList(params.Subject, params.Slug, params.BasePageParams)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(emailTemplateList, "", c)
}

// EmailTemplateAdd 新增邮件模板
func EmailTemplateAdd(c *gin.Context) {
	var params struct {
		Subject string `json:"subject" remark:"邮件主题" binding:"required,max=64"`
		Content string `json:"content" remark:"邮件内容" binding:"required,max=65535"`
		Slug    string `json:"slug" remark:"邮件标识" binding:"required,max=32"`
		Remark  string `json:"remark" remark:"备注" binding:"max=255"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断邮件标识是否已存在
	dbEmailTemplate, err := system.NewEmailTemplate().GetEmailTemplateBySlug(params.Slug)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if dbEmailTemplate.ID > 0 {
		response.LogicExceptionJSON("邮件标识已存在", c)
		return
	}

	// 新增邮件模板
	emailTemplate := system.EmailTemplate{
		Subject: params.Subject,
		Content: params.Content,
		Slug:    params.Slug,
		Remark:  params.Remark,
	}
	if err := helper.GormDefaultDb.Create(&emailTemplate).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(model.BaseIdResult{ID: emailTemplate.ID}, "新增成功", c)
}

// EmailTemplateModify 修改邮件模板
func EmailTemplateModify(c *gin.Context) {
	var params struct {
		model.BaseIdParams
		Subject string `json:"subject" remark:"邮件主题" binding:"required,max=64"`
		Content string `json:"content" remark:"邮件内容" binding:"required,max=65535"`
		Slug    string `json:"slug" remark:"邮件标识" binding:"required,max=32"`
		Remark  string `json:"remark" remark:"备注" binding:"max=255"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断邮件模板是否存在
	emailTemplate, err := system.NewEmailTemplate().FindFullInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if emailTemplate.ID == 0 {
		response.LogicExceptionJSON("邮件模板不存在", c)
		return
	}

	// 判断邮件标识是否已存在
	if params.Slug != emailTemplate.Slug {
		dbEmailTemplate, err := system.NewEmailTemplate().GetEmailTemplateBySlug(params.Slug)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.LogicExceptionJSON("系统出错了："+err.Error(), c)
			return
		}
		if dbEmailTemplate.ID > 0 {
			response.LogicExceptionJSON("邮件标识已存在", c)
			return
		}
	}

	// 修改邮件模板
	emailTemplate.Subject = params.Subject
	emailTemplate.Content = params.Content
	emailTemplate.Slug = params.Slug
	emailTemplate.Remark = params.Remark
	if err := helper.GormDefaultDb.Save(&emailTemplate).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(model.BaseIdResult{ID: emailTemplate.ID}, "修改成功", c)
}

// EmailTemplateDelete 删除邮件模板
func EmailTemplateDelete(c *gin.Context) {
	var params model.BaseIdParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断邮件模板是否存在
	emailTemplate, err := system.NewEmailTemplate().FindFullInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if emailTemplate.ID == 0 {
		response.LogicExceptionJSON("邮件模板已被删除", c)
		return
	}

	// 删除邮件模板
	if err := helper.GormDefaultDb.Delete(&emailTemplate).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(model.BaseIdResult{ID: emailTemplate.ID}, "删除成功", c)
}

// EmailBatchSend 批量发送邮件
func EmailBatchSend(c *gin.Context) {
	var params struct {
		Subject    string   `json:"subject" remark:"邮件主题" binding:"required,max=64"`
		Content    string   `json:"content" remark:"邮件内容" binding:"required,max=65535"`
		Recipients []string `json:"recipients" remark:"收件人" binding:"required,dive,email"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 读取邮件配置
	email, err := tool.FindEmail()
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if email.Server == "" {
		response.LogicExceptionJSON("邮件服务未配置", c)
		return
	}

	// 使用协程批量发送邮件，最多同时发送10封
	var wg sync.WaitGroup
	var success, fail int32
	var ch = make(chan bool, 10)
	for _, recipient := range params.Recipients {
		ch <- true
		wg.Add(1)

		go func(recipient string) {
			defer wg.Done()
			defer func() { <-ch }()
			if err := email.Send(params.Subject, params.Content, []string{recipient}, []string{}, []string{}); err != nil {
				atomic.AddInt32(&fail, 1)
				log.GetLogger("email").Errorf("邮件发送失败，邮件主题《%s》，收件人：%s，异常信息：%s", params.Subject, recipient, err.Error())
				return
			}
			atomic.AddInt32(&success, 1)
			log.GetLogger("email").Infof("邮件发送成功，邮件主题《%s》，收件人：%s", params.Subject, recipient)
		}(recipient)
	}

	wg.Wait()
	response.SuccessJSON(gin.H{
		"total":   len(params.Recipients),
		"success": success,
		"fail":    fail,
	}, "发送成功", c)
}
