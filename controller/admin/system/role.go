package system

import (
	"errors"

	"gin-skeleton/helper"
	"gin-skeleton/helper/response"
	"gin-skeleton/model"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RoleList 获取角色列表
func RoleList(c *gin.Context) {
	var params struct {
		model.BasePageParams
		Name   string `json:"name" remark:"名称" binding:"max=64"`
		Status int    `json:"status" remark:"状态" binding:"oneof=0 1 2"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	roleList, err := system.NewRole().GetRoleList(params.Name, params.Status, params.BasePageParams)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(roleList, "", c)
}

// RoleAdd 添加角色
func RoleAdd(c *gin.Context) {
	var params struct {
		Name   string  `json:"name" remark:"名称" binding:"required,max=64"`
		Status int     `json:"status" remark:"状态" binding:"oneof=1 2"`
		Weight int     `json:"weight" remark:"排序" binding:"gte=0"`
		Remark string  `json:"remark" remark:"备注" binding:"max=255"`
		Mids   []int64 `json:"mids" remark:"授权菜单"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断角色名是否存在
	dbRole, err := system.NewRole().FindByName(params.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if dbRole.ID > 0 {
		response.InvalidArgumentJSON("角色名已存在", c)
		return
	}

	// 事务处理
	var rid int64
	err = helper.GormDefaultDb.Transaction(func(tx *gorm.DB) error {
		// 创建角色
		newRole := system.Role{
			Name:   params.Name,
			Status: params.Status,
			Weight: params.Weight,
			Remark: params.Remark,
		}
		if err := tx.Create(&newRole).Error; err != nil {
			return err
		}

		// 创建角色菜单
		rid = newRole.ID
		if err := system.NewAuthRelation().CreateMids(tx, rid, params.Mids); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdReuslt{ID: rid}, "创建角色成功", c)
}

// RoleModify 修改角色
func RoleModify(c *gin.Context) {
	var params struct {
		model.BaseIdParams
		Name   string  `json:"name" remark:"名称" binding:"required,max=64"`
		Status int     `json:"status" remark:"状态" binding:"oneof=1 2"`
		Weight int     `json:"weight" remark:"排序" binding:"gte=0"`
		Remark string  `json:"remark" remark:"备注" binding:"max=255"`
		Mids   []int64 `json:"mids" remark:"授权菜单"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断角色是否存在
	role, err := system.NewRole().FindFullInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if role.ID == 0 {
		response.InvalidArgumentJSON("角色不存在", c)
		return
	}

	// 判断修改的角色名是否存在
	if params.Name != role.Name {
		dbRole, err := system.NewRole().FindByName(params.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.LogicExceptionJSON("系统出错了："+err.Error(), c)
			return
		}
		if dbRole.ID > 0 {
			response.InvalidArgumentJSON("角色名已存在", c)
			return
		}
	}

	// 事务处理
	err = helper.GormDefaultDb.Transaction(func(tx *gorm.DB) error {
		// 修改角色
		role.Name = params.Name
		role.Status = params.Status
		role.Weight = params.Weight
		role.Remark = params.Remark
		if err := tx.Save(&role).Error; err != nil {
			return err
		}

		// 添加新增权限
		if err := system.NewAuthRelation().CreateMids(tx, role.ID, helper.DiffSilce(params.Mids, role.Mids)); err != nil {
			return err
		}

		// 删除弃用权限
		if err := system.NewAuthRelation().DeleteMids(tx, role.ID, helper.DiffSilce(role.Mids, params.Mids)); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdReuslt{ID: role.ID}, "修改角色成功", c)
}

// RoleDelete 删除角色
func RoleDelete(c *gin.Context) {
	var params model.BaseIdParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断角色是否存在
	role, err := system.NewRole().FindFullInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if role.ID == 0 {
		response.InvalidArgumentJSON("角色已被删除", c)
		return
	}

	// 事务处理
	err = helper.GormDefaultDb.Transaction(func(tx *gorm.DB) error {
		// 删除角色
		if err := tx.Delete(&role).Error; err != nil {
			return err
		}

		// 删除权限菜单
		if err := system.NewAuthRelation().DeleteMids(tx, role.ID, role.Mids); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdReuslt{ID: role.ID}, "删除角色成功", c)
}
