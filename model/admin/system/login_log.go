package system

import (
	"gin-skeleton/helper"
	"gin-skeleton/model"
)

const (
	LoginLogOperTypeSuccess = 1 // 操作类型 - 登录成功
	LoginLogOperTypeFail    = 2 // 操作类型 - 登录失败
	LoginLogOperTypeLogout  = 3 // 操作类型 - 退出登录
)

// LoginLog 登录日志表
type LoginLog struct {
	model.BaseModel
	Username string `json:"username"` // 用户名
	Nickname string `json:"nickname"` // 昵称
	Ip       string `json:"ip"`       // IP地址
	Device   string `json:"device"`   // 设备型号
	Os       string `json:"os"`       // 操作系统
	Browser  string `json:"browser"`  // 浏览器
	OperType int    `json:"operType"` // 操作类型：1-登录成功，2-登录失败，3-退出登录
	Remark   string `json:"remark"`   // 备注
}

// NewLoginLog 初始化对象
func NewLoginLog() *LoginLog {
	return &LoginLog{}
}

// GetLoginLogList 获取登录日志列表
func (ll *LoginLog) GetLoginLogList(name, ip string, operType int, createdDate []string, pageInfo model.BasePageParams) (pr *model.BasePageResult[LoginLog], err error) {
	loginLogModel := helper.GormDefaultDb.Model(NewLoginLog())
	if name != "" {
		loginLogModel.Where("username like ? or nickname like ?", "%"+name+"%", "%"+name+"%")
	}
	if ip != "" {
		loginLogModel.Where("ip like ?", "%"+ip+"%")
	}
	if len(createdDate) == 2 && createdDate[0] != "" && createdDate[1] != "" {
		loginLogModel.Where("created_at >= ? and created_at <= ?", createdDate[0]+" 00:00:00", createdDate[1]+" 23:59:59")
	}
	if operType > 0 {
		loginLogModel.Where("oper_type = ?", operType)
	}

	pr = &model.BasePageResult[LoginLog]{Items: make([]*LoginLog, 0), Total: 0}
	err = loginLogModel.Scopes(model.Paginate(pageInfo)).Find(&pr.Items).Scopes(model.Count).Count(&pr.Total).Error
	return
}
