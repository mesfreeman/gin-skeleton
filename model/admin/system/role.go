package system

import (
	"errors"
	"gin-skeleton/helper"
	"gin-skeleton/model"

	"gorm.io/gorm"
)

// Role  角色表
type Role struct {
	model.BaseModel
	Name   string  `json:"name"`          // 名称，中文名
	Status int     `json:"status"`        // 状态：1-禁用，2-启用
	Weight int     `json:"weight"`        // 权重，值越大越靠前
	Remark string  `json:"remark"`        // 备注
	Mids   []int64 `json:"mids" gorm:"-"` // 菜单ID数组
}

// NewRole 初始化角色
func NewRole() *Role {
	return &Role{}
}

// GetRoleList 获取角色列表
func (r *Role) GetRoleList(name string, status int, pageInfo model.BasePageParams) (pr *model.BasePageResult[Role], err error) {
	roleModel := helper.GormDefaultDb.Model(NewRole())
	if name != "" {
		roleModel.Where("name like ?", "%"+name+"%")
	}
	if status > 0 {
		roleModel.Where("status = ?", status)
	}

	pr = &model.BasePageResult[Role]{Items: make([]*Role, 0), Total: 0}
	err = roleModel.Scopes(model.Paginate(pageInfo)).Find(&pr.Items).Scopes(model.CancelPaginate).Count(&pr.Total).Error
	if err != nil || len(pr.Items) == 0 {
		return
	}

	// 获取角色相关菜单ID
	for _, role := range pr.Items {
		mids, err2 := NewAuthRelation().GetRoleMids(role.ID)
		if err2 != nil {
			err = err2
			return
		}

		role.Mids = mids
	}
	return
}

// FindBasicInfo 查询基本的角色信息
func (r *Role) FindBasicInfo(id int64) (role *Role, err error) {
	err = helper.GormDefaultDb.First(&role, id).Error
	return
}

// FindFullInfo 查询完整的角色信息
func (r *Role) FindFullInfo(id int64) (role *Role, err error) {
	role, err = r.FindBasicInfo(id)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	// 获取权限菜单
	mids, err := NewAuthRelation().GetRoleMids(role.ID)
	if err != nil {
		return
	}
	role.Mids = mids
	return
}

// FindByName 基于角色名查找角色
func (r *Role) FindByName(name string) (role *Role, err error) {
	err = helper.GormDefaultDb.Where("name = ?", name).First(&role).Error
	return
}
