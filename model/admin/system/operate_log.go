package system

import (
	"gin-skeleton/helper"
	"gin-skeleton/model"
)

// OperateLog 操作日志表
type OperateLog struct {
	model.BaseModel
	Username  string `json:"username"`  // 用户名
	Nickname  string `json:"nickname"`  // 昵称
	Ip        string `json:"ip"`        // IP地址
	Function  string `json:"function"`  // 操作功能
	Uri       string `json:"uri"`       // 请求地址
	Method    string `json:"method"`    // 请求方式
	Params    string `json:"params"`    // 请求参数
	Status    int    `json:"status"`    // 状态码
	Code      string `json:"code"`      // 业务码
	SpendTime int64  `json:"spendTime"` // 耗时，单位：ms
	Result    string `json:"result"`    // 响应结果
	UserAgent string `json:"userAgent"` // 浏览器信息
	Remark    string `json:"remark"`    // 备注
}

// NewOerateLog 初始化对象
func NewOerateLog() *OperateLog {
	return &OperateLog{}
}

// GetOperateLogList 获取操作日志列表
func (ll *OperateLog) GetOperateLogList(name, function string, createdDate []string, pageInfo model.BasePageParams) (pr *model.BasePageResult[OperateLog], err error) {
	loginLogModel := helper.GormDefaultDb.Model(NewOerateLog())
	if name != "" {
		loginLogModel.Where("username like ? or nickname like ?", "%"+name+"%", "%"+name+"%")
	}
	if function != "" {
		loginLogModel.Where("function like ?", "%"+function+"%")
	}
	if len(createdDate) == 2 && createdDate[0] != "" && createdDate[1] != "" {
		loginLogModel.Where("created_at >= ? and created_at <= ?", createdDate[0]+" 00:00:00", createdDate[1]+" 23:59:59")
	}

	pr = &model.BasePageResult[OperateLog]{Items: make([]*OperateLog, 0), Total: 0}
	err = loginLogModel.Scopes(model.Paginate(pageInfo)).Find(&pr.Items).Scopes(model.CancelPaginate).Count(&pr.Total).Error
	return
}
