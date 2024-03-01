package system

import (
	"gin-skeleton/helper"
	"gin-skeleton/model"
)

const (
	LoginLogTypeFail    = iota + 1 // 日志类型 - 登录失败
	LoginLogTypeSuccess            // 日志类型 - 登录成功
	LoginLogTypeLogout             // 日志类型 - 退出登录
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
	Type     int    `json:"type"`     // 日志类型：1-登录失败，2-登录成功，3-退出登录
	Remark   string `json:"remark"`   // 备注
}

// NewLoginLog 初始化对象
func NewLoginLog() *LoginLog {
	return &LoginLog{}
}

// GetLoginLogList 获取登录日志列表
func (ll *LoginLog) GetLoginLogList(name, ip string, logType int, createdDate []string, pageInfo model.BasePageParams) (pr *model.BasePageResult[LoginLog], err error) {
	loginLogModel := helper.GormDefaultDb.Model(NewLoginLog()).Scopes(model.FilterByDate(createdDate, "created_at"))
	if name != "" {
		loginLogModel.Where("username like ? or nickname like ?", "%"+name+"%", "%"+name+"%")
	}
	if ip != "" {
		loginLogModel.Where("ip like ?", "%"+ip+"%")
	}
	if logType > 0 {
		loginLogModel.Where("type = ?", logType)
	}

	pr = &model.BasePageResult[LoginLog]{Items: make([]*LoginLog, 0), Total: 0}
	err = loginLogModel.Scopes(model.Paginate(pageInfo)).Find(&pr.Items).Scopes(model.CancelPaginate).Count(&pr.Total).Error
	return
}
