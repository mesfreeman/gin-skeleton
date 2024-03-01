package system

import (
	"gin-skeleton/helper"
	"gin-skeleton/model"
)

// EmailTemplate 邮件模板表
type EmailTemplate struct {
	model.BaseModel
	Subject string `json:"subject"` // 邮件主题
	Content string `json:"content"` // 邮件内容
	Slug    string `json:"slug"`    // 邮件标识
	Remark  string `json:"remark"`  // 备注
}

// NewEmailTemplate 初始化对象
func NewEmailTemplate() *EmailTemplate {
	return &EmailTemplate{}
}

// GetEmailTemplateList 获取邮件模板列表
func (e *EmailTemplate) GetEmailTemplateList(subject, slug string, pageInfo model.BasePageParams) (pr *model.BasePageResult[EmailTemplate], err error) {
	emailTemplateModel := helper.GormDefaultDb.Model(NewEmailTemplate())
	if subject != "" {
		emailTemplateModel.Where("subject like ?", "%"+subject+"%")
	}
	if slug != "" {
		emailTemplateModel.Where("slug like ?", "%"+slug+"%")
	}

	pr = &model.BasePageResult[EmailTemplate]{Items: make([]*EmailTemplate, 0), Total: 0}
	err = emailTemplateModel.Scopes(model.Paginate(pageInfo)).Find(&pr.Items).Scopes(model.CancelPaginate).Error
	return
}

// FindFullInfo 查询完整的邮件模板信息
func (e *EmailTemplate) FindFullInfo(id int64) (emailTemplate *EmailTemplate, err error) {
	err = helper.GormDefaultDb.First(&emailTemplate, id).Error
	return
}

// GetEmailTemplateBySlug 根据slug获取邮件模板
func (e *EmailTemplate) GetEmailTemplateBySlug(slug string) (emailTemplate *EmailTemplate, err error) {
	err = helper.GormDefaultDb.Where("slug = ?", slug).First(&emailTemplate).Error
	return
}
