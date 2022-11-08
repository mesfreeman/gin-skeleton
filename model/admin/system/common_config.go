package system

import (
	"encoding/json"
	"errors"
	"gin-skeleton/helper"
	"gin-skeleton/model"

	"gorm.io/gorm"
)

const (
	CommonConfigModuleFileStorage = "fileStorage" // 文件存储
	CommonConfigModuleEmailServer = "emailServer" // 邮件服务
)

// CommonConfig  公共配置表
type CommonConfig struct {
	model.BaseModel
	Module  string `json:"module"`  // 模块
	Keyword string `json:"keyword"` // 关键词
	Value   string `json:"value"`   // 配置值
	Remark  string `json:"remark"`  // 备注
}

// NewCommonConfig 初始化对象
func NewCommonConfig() *CommonConfig {
	return &CommonConfig{}
}

// GetConfigs 获取指定模块下的所有配置
func (cc *CommonConfig) GetConfigs(module string) (res map[string]string, err error) {
	var configs []struct {
		Keyword string
		Value   string
	}

	res = make(map[string]string, 0)
	err = helper.GormDefaultDb.Model(NewCommonConfig()).Select("keyword", "value").Where("module = ?", module).Find(&configs).Error
	if err != nil {
		return
	}

	for _, item := range configs {
		res[item.Keyword] = item.Value
	}
	return
}

// FindConfigValue 查找指定模块下指定关键词的配置值
func (cc *CommonConfig) FindConfigValue(module, keyword string) (value string, err error) {
	err = helper.GormDefaultDb.Model(NewCommonConfig()).Where("module = ? and keyword = ?", module, keyword).Pluck("value", &value).Error
	return
}

// FindConfigValueTo 查找指定模块下指定关键词的配置值并将数据转化为指定的格式（注：仅适用于配置值为JSON字符串的情况）
func (cc *CommonConfig) FindConfigValueTo(module, keyword string, to any) (err error) {
	value, err := cc.FindConfigValue(module, keyword)
	if err != nil {
		return
	}

	if value == "" {
		value = "{}"
	}

	err = json.Unmarshal([]byte(value), &to)
	return
}

// FindConfigValueItem 查找指定模块下指定关键词配置值下的指定元素的值（注：仅适用于配置值为JSON字符串的情况）
func (cc *CommonConfig) FindConfigValueItem(module, keyword, itemCode string) (itemValue any, err error) {
	itemMap := make(map[string]any)
	err = cc.FindConfigValueTo(module, keyword, &itemMap)
	if err != nil {
		return
	}

	itemValue, _ = itemMap[itemCode]
	return
}

// CreateOrUpdateConfig 创建或更新配置
func (cc *CommonConfig) CreateOrUpdateConfig(params *CommonConfig) (id int64, err error) {
	var dbConfig CommonConfig
	err = helper.GormDefaultDb.Model(NewCommonConfig()).Where("module = ? and keyword = ?", params.Module, params.Keyword).First(&dbConfig).Error
	if err != nil && !errors.Is(gorm.ErrRecordNotFound, err) {
		return
	}

	// 配置不存在，创建操作
	if dbConfig.ID == 0 {
		newConfig := CommonConfig{
			Module:  params.Module,
			Keyword: params.Keyword,
			Value:   params.Value,
			Remark:  params.Remark,
		}

		if err = helper.GormDefaultDb.Create(&newConfig).Error; err != nil {
			return
		}
		id = newConfig.ID
		return
	}

	// 否则更新操作
	dbConfig.Value = params.Value
	dbConfig.Remark = params.Remark
	if err = helper.GormDefaultDb.Save(&dbConfig).Error; err != nil {
		return
	}
	id = dbConfig.ID
	return
}
