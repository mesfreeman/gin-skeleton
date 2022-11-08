package system

import (
	"encoding/json"
	"errors"

	"gin-skeleton/helper"
	"gin-skeleton/helper/response"
	"gin-skeleton/helper/tool"
	"gin-skeleton/model"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FileList 获取文件列表
func FileList(c *gin.Context) {
	var params struct {
		model.BasePageParams
		Keyword     string   `json:"keyword,omitempty" remark:"关键词"`
		Uploader    string   `json:"uploader,omitempty" remark:"上传者"`
		CreatedDate []string `json:"createdDate,omitempty" remark:"创建时间"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	fileList, err := system.NewFile().GetFileList(params.Keyword, params.Uploader, params.CreatedDate, params.BasePageParams)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(fileList, "", c)
}

// FileDelete 删除文件
func FileDelete(c *gin.Context) {
	var params model.BaseIdParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断文件是否存在
	file, err := system.NewFile().FindFileInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了"+err.Error(), c)
		return
	}
	if file.ID == 0 {
		response.InvalidArgumentJSON("文件已被删除", c)
		return
	}

	// 判断文件配置是否存在
	var fileConfig system.FileConfig
	err = system.NewCommonConfig().FindConfigValueTo(system.CommonConfigModuleFileStorage, "", &fileConfig)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if fileConfig.Provider == "" {
		response.RuntimeExceptionJSON("未配置文件提供商", c)
		return
	}

	// 调用三方服务删除文件
	err = tool.NewStorage(fileConfig).DeleteFile(file.FileUrl)
	if err != nil {
		response.ThirdExceptionJSON("删除三方服务文件异常："+err.Error(), c)
		return
	}

	// 删除数据表的数据
	if err := helper.GormDefaultDb.Delete(&file).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(model.BaseIdReuslt{ID: file.ID}, "删除成功", c)
}

// FileViewConfig 查看文件配置
func FileViewConfig(c *gin.Context) {
	var fileConfig system.FileConfig
	err := system.NewCommonConfig().FindConfigValueTo(system.CommonConfigModuleFileStorage, "", &fileConfig)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(fileConfig, "", c)
}

// FileSaveConfig 保存文件配置
func FileSaveConfig(c *gin.Context) {
	var params struct {
		Provider   string   `json:"provider" remark:"提供商" binding:"oneof=qiniu ali tencent"`
		Bucket     string   `json:"bucket" remark:"空间" binding:"required"`
		AccessKey  string   `json:"accessKey" remark:"密匙" binding:"required"`
		SecretKey  string   `json:"secretKey" remark:"密钥" binding:"required"`
		Domain     string   `json:"domain" remark:"域名" binding:"required"`
		ThumbConf  string   `json:"thumbConf,omitempty" remark:"缩略图配置"`
		AllowTypes []string `json:"allowTypes" remark:"类型"`
		Remark     string   `json:"remark,omitempty" remark:"备注"`
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
		Module:  system.CommonConfigModuleFileStorage,
		Keyword: "",
		Value:   string(value),
		Remark:  "文件存储相关配置",
	}
	id, err := system.NewCommonConfig().CreateOrUpdateConfig(&commonConfig)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(model.BaseIdReuslt{ID: id}, "", c)
}
