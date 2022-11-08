package system

import (
	"gin-skeleton/helper/response"
	"gin-skeleton/model"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
)

// OperateLogList 获取操作日志列表
func OperateLogList(c *gin.Context) {
	var params struct {
		model.BasePageParams
		Name        string   `json:"name,omitempty" remark:"操作账号"`
		Function    string   `json:"function,omitempty" remark:"操作功能"`
		CreatedDate []string `json:"createdDate,omitempty" remark:"创建时间"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	operateLogList, err := system.NewOerateLog().GetOperateLogList(params.Name, params.Function, params.CreatedDate, params.BasePageParams)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(operateLogList, "", c)
}
