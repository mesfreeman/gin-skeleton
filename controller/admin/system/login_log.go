package system

import (
	"gin-skeleton/helper/response"
	"gin-skeleton/model"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
)

// LoginLogList 获取登录日志列表
func LoginLogList(c *gin.Context) {
	var params struct {
		model.BasePageParams
		Name        string   `json:"name,omitempty" remark:"账号"`
		Ip          string   `json:"ip,omitempty" remark:"IP地址"`
		OperType    int      `json:"operType,omitempty" remark:"操作类型" binding:"oneof=0 1 2 3"`
		CreatedDate []string `json:"createdDate,omitempty" remark:"创建时间"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	loginLogList, err := system.NewLoginLog().GetLoginLogList(params.Name, params.Ip, params.OperType, params.CreatedDate, params.BasePageParams)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(loginLogList, "", c)
}
