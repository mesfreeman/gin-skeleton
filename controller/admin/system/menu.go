package system

import (
	"errors"
	"gin-skeleton/model"
	"strings"

	"gin-skeleton/helper"
	"gin-skeleton/helper/response"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MenuList 菜单列表
func MenuList(c *gin.Context) {
	var params struct {
		Name   string `json:"name" remark:"名称" binding:"max=32"`
		Status int    `json:"status" remark:"状态" binding:"oneof=0 1 2"`
		Type   []int  `json:"type" remark:"类型"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	sysMenus, err := system.NewMenu().GetSysMenus(params.Name, params.Status, params.Type)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(sysMenus, "", c)
}

// MenuAdd 添加菜单
func MenuAdd(c *gin.Context) {
	var params struct {
		Pid       int64  `json:"pid" remark:"PID"`
		Pname     any    `json:"pname" remark:"上级"`
		Name      string `json:"name" remark:"名称" binding:"required,min=1,max=32"`
		Icon      string `json:"icon" remark:"图标" binding:"max=128"`
		Path      string `json:"path" remark:"地址" binding:"required,max=255"`
		Component string `json:"component" remark:"组件" binding:"max=255"`
		Type      int    `json:"type" remark:"类型" binding:"oneof=1 2 3"`
		Mode      int    `json:"mode" remark:"模式" binding:"oneof=1 2 3"`
		Weight    int    `json:"weight" remark:"排序" binding:"gte=0"`
		IsShow    int    `json:"isShow" remark:"显示" binding:"oneof=1 2"`
		Status    int    `json:"status" remark:"状态" binding:"oneof=1 2"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 关联性检查
	if params.Mode != system.MenuModeComponent && !helper.IsLegalUrl(params.Path) {
		response.InvalidArgumentJSON("模式为内链或外链时，地址必须是合法的链接地址", c)
		return
	}

	// 前端框架的原因，暂时这么处理
	if val, ok := params.Pname.(float64); ok {
		params.Pid = int64(val)
	}
	if params.Pname == nil {
		params.Pid = 0
	}

	// 判断父级菜单是否存在
	level := 1
	if params.Pid > 0 {
		parentMenuInfo, err := system.NewMenu().FindMenuInfo(params.Pid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.LogicExceptionJSON("系统出错了："+err.Error(), c)
			return
		}
		if parentMenuInfo.ID == 0 {
			response.InvalidArgumentJSON("上级菜单不存在", c)
			return
		}
		level = parentMenuInfo.Level + 1
	}

	// 判断同级菜单名称是否存在
	dbMenu, err := system.NewMenu().FindByNameLevel(params.Name, level)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if dbMenu.ID > 0 {
		response.InvalidArgumentJSON("该目录下有同名的菜单", c)
		return
	}

	// 数据处理
	if params.Type == system.MenuTypeDir {
		params.Component = "layout"
	}
	if params.Type == system.MenuTypeMenu && params.Mode != system.MenuModeComponent {
		params.Component = "iframe"
	}
	if params.Type == system.MenuTypeApi {
		params.Component = ""
	}

	// 组件地址，必须以前缀"/"开头
	if params.Mode == system.MenuModeComponent {
		params.Path = "/" + strings.TrimPrefix(params.Path, "/")
	}

	// 创建菜单
	newMenu := system.Menu{
		Pid:       params.Pid,
		Name:      params.Name,
		Icon:      params.Icon,
		Path:      params.Path,
		Component: params.Component,
		Type:      params.Type,
		Mode:      params.Mode,
		Weight:    params.Weight,
		Level:     level,
		IsShow:    params.IsShow,
		Status:    params.Status,
	}
	if err = helper.GormDefaultDb.Create(&newMenu).Error; err != nil {
		response.LogicExceptionJSON(err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdResult{ID: newMenu.ID}, "创建菜单成功", c)
}

// MenuModify 修改菜单
func MenuModify(c *gin.Context) {
	var params struct {
		model.BaseIdParams
		Pid       int64  `json:"pid" remark:"PID"`
		Pname     any    `json:"pname" remark:"上级"`
		Name      string `json:"name" remark:"名称" binding:"required,min=1,max=32"`
		Icon      string `json:"icon" remark:"图标" binding:"max=128"`
		Path      string `json:"path" remark:"地址" binding:"required,max=255"`
		Component string `json:"component" remark:"组件" binding:"max=255"`
		Type      int    `json:"type" remark:"类型" binding:"oneof=1 2 3"`
		Mode      int    `json:"mode" remark:"模式" binding:"oneof=1 2 3"`
		Weight    int    `json:"weight" remark:"排序" binding:"gte=0"`
		IsShow    int    `json:"isShow" remark:"显示" binding:"oneof=1 2"`
		Status    int    `json:"status" remark:"状态" binding:"oneof=1 2"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 关联性检查
	if params.Mode != system.MenuModeComponent && !helper.IsLegalUrl(params.Path) {
		response.InvalidArgumentJSON("模式为内链或外链时，地址必须是合法的链接地址", c)
		return
	}

	// 判断要修改的菜单是否存在
	menu, err := system.NewMenu().FindMenuInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if menu.ID == 0 {
		response.InvalidArgumentJSON("菜单不存在", c)
		return
	}

	// 前端框架的原因，暂时这么处理
	if val, ok := params.Pname.(float64); ok {
		params.Pid = int64(val)
	}
	if params.Pname == nil {
		params.Pid = 0
	}

	// 判断父级菜单是否存在
	level := 1
	if params.Pid > 0 {
		parentMenu, err := system.NewMenu().FindMenuInfo(params.Pid)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.LogicExceptionJSON("系统出错了："+err.Error(), c)
			return
		}
		if parentMenu.ID == 0 {
			response.InvalidArgumentJSON("上级菜单不存在", c)
			return
		}
		level = parentMenu.Level + 1
	}

	// 判断同级菜单名称是否存在
	if params.Name != menu.Name {
		dbMenu, err := system.NewMenu().FindByNameLevel(params.Name, level)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.LogicExceptionJSON("系统出错了："+err.Error(), c)
			return
		}
		if dbMenu.ID > 0 {
			response.InvalidArgumentJSON("该目录下有同名的菜单", c)
			return
		}
	}

	// 数据处理
	if params.Type == system.MenuTypeDir {
		params.Component = "layout"
	}
	if params.Type == system.MenuTypeMenu && params.Mode != system.MenuModeComponent {
		params.Component = "iframe"
	}
	if params.Type == system.MenuTypeApi {
		params.Component = ""
	}

	// 组件地址，必须以前缀"/"开头
	if params.Mode == system.MenuModeComponent {
		params.Path = "/" + strings.TrimPrefix(params.Path, "/")
	}

	// 更新菜单
	menu.Pid = params.Pid
	menu.Name = params.Name
	menu.Icon = params.Icon
	menu.Path = params.Path
	menu.Component = params.Component
	menu.Type = params.Type
	menu.Mode = params.Mode
	menu.Weight = params.Weight
	menu.Level = level
	menu.IsShow = params.IsShow
	menu.Status = params.Status
	err = helper.GormDefaultDb.Save(&menu).Error
	if err != nil {
		response.LogicExceptionJSON(err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdResult{ID: menu.ID}, "修改菜单成功", c)
}

// MenuDelete 删除菜单
func MenuDelete(c *gin.Context) {
	var params model.BaseIdParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断要删除的菜单是否存在
	menu, err := system.NewMenu().FindMenuInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if menu.ID == 0 {
		response.InvalidArgumentJSON("菜单不存在", c)
		return
	}

	// 删除菜单
	if err = helper.GormDefaultDb.Delete(&menu).Error; err != nil {
		response.LogicExceptionJSON(err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdResult{ID: menu.ID}, "删除菜单成功", c)
}
