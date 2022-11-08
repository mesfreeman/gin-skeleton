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

// AccountList 获取账号列表
func AccountList(c *gin.Context) {
	var params struct {
		model.BasePageParams
		Name   string `json:"name" remark:"账号名" binding:"max=32"`
		Status int    `json:"status" remark:"状态" binding:"oneof=0 1 2"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	accountList, err := system.NewAccount().GetAccountList(params.Name, params.Status, params.BasePageParams)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	response.SuccessJSON(accountList, "", c)
}

// AccountAdd 新增账号
func AccountAdd(c *gin.Context) {
	var params struct {
		Username string  `json:"username" remark:"账号" binding:"required,min=4,max=16"`
		Password string  `json:"password" remark:"密码" binding:"required,min=6,max=15"`
		Nickname string  `json:"nickname" remark:"昵称" binding:"required,min=2,max=32"`
		Avatar   string  `json:"avatar" remark:"头像" binding:"max=255"`
		Email    string  `json:"email" remark:"邮箱" binding:"required,email"`
		Phone    string  `json:"phone" remark:"手机号" binding:"omitempty,len=11"`
		Status   int     `json:"status" remark:"状态" binding:"oneof=1 2"`
		Remark   string  `json:"remark" remark:"备注" binding:"max=255"`
		Rids     []int64 `json:"rids" remark:"角色"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断账号是否存在
	dbAccount, err := system.NewAccount().FindByUsername(params.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if dbAccount.ID > 0 {
		response.InvalidArgumentJSON("账号已存在", c)
		return
	}

	// 事务处理
	var aid int64
	err = helper.GormDefaultDb.Transaction(func(tx *gorm.DB) error {
		// 创建账号
		newAccount := &system.Account{
			Username: params.Username,
			Nickname: params.Nickname,
			Password: system.NewAccount().EncryptPassword(params.Password),
			Avatar:   params.Avatar,
			Email:    params.Email,
			Phone:    params.Phone,
			Status:   system.AccountStatusNormal,
			Remark:   params.Remark,
		}
		if err := tx.Create(&newAccount).Error; err != nil {
			return err
		}

		// 创建角色
		aid = newAccount.ID
		if err := system.NewAuthRelation().CreateRids(tx, aid, params.Rids); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdReuslt{ID: aid}, "创建账号成功", c)
}

// AccountModify 修改账号
func AccountModify(c *gin.Context) {
	var params struct {
		model.BaseIdParams
		Username string  `json:"username" remark:"账号" binding:"required,min=4,max=16"`
		Nickname string  `json:"nickname" remark:"昵称" binding:"required,min=2,max=32"`
		Avatar   string  `json:"avatar" remark:"头像" binding:"max=255"`
		Email    string  `json:"email" remark:"邮箱" binding:"required,email"`
		Phone    string  `json:"phone" remark:"手机号" binding:"omitempty,len=11"`
		Status   int     `json:"status" remark:"状态" binding:"oneof=1 2"`
		Remark   string  `json:"remark" remark:"备注" binding:"max=255"`
		Rids     []int64 `json:"rids" remark:"角色"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断账号是否存在
	account, err := system.NewAccount().FindFullInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if account.ID == 0 {
		response.InvalidArgumentJSON("账号不存在", c)
		return
	}

	// 判断修改的账号是否存在
	if params.Username != account.Username {
		dbAccount, err := system.NewAccount().FindByUsername(params.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			response.LogicExceptionJSON("系统出错了："+err.Error(), c)
			return
		}
		if dbAccount.ID > 0 {
			response.InvalidArgumentJSON("账号已存在", c)
			return
		}
	}

	// 事务处理
	err = helper.GormDefaultDb.Transaction(func(tx *gorm.DB) error {
		// 修改账号
		account.Username = params.Username
		account.Nickname = params.Nickname
		account.Avatar = params.Avatar
		account.Email = params.Email
		account.Phone = params.Phone
		account.Status = params.Status
		account.Remark = params.Remark
		if err := tx.Save(&account).Error; err != nil {
			return err
		}

		// 添加新增角色
		if err := system.NewAuthRelation().CreateRids(tx, account.ID, helper.DiffSilce(params.Rids, account.Rids)); err != nil {
			return err
		}

		// 删除弃用角色
		if err := system.NewAuthRelation().DeleteRids(tx, account.ID, helper.DiffSilce(account.Rids, params.Rids)); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdReuslt{ID: account.ID}, "修改账号成功", c)
}

// AccountModifyPwd 修改账号密码
func AccountModifyPwd(c *gin.Context) {
	var params struct {
		model.BaseIdParams
		NewPwd string `json:"newPwd" remark:"新密码" binding:"min=6,max=15"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断账号是否存在
	account, err := system.NewAccount().FindBasicInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if account.ID == 0 {
		response.InvalidArgumentJSON("账号不存在", c)
		return
	}

	// 修改密码
	account.Password = system.NewAccount().EncryptPassword(params.NewPwd)
	if err := helper.GormDefaultDb.Save(&account).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdReuslt{ID: account.ID}, "修改账号密码成功", c)
}

// AccountDelete 删除账号
func AccountDelete(c *gin.Context) {
	var params model.BaseIdParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断账号是否存在
	account, err := system.NewAccount().FindFullInfo(params.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if account.ID == 0 {
		response.InvalidArgumentJSON("账号已被删除", c)
		return
	}

	// 事务处理
	err = helper.GormDefaultDb.Transaction(func(tx *gorm.DB) error {
		// 删除账号
		if err := tx.Delete(&account).Error; err != nil {
			return err
		}

		// 删除账号角色
		if err := system.NewAuthRelation().DeleteRids(tx, account.ID, account.Rids); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdReuslt{ID: account.ID}, "删除账号成功", c)
}
