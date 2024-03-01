package system

import (
	"gin-skeleton/helper"
	"gin-skeleton/model"

	"gorm.io/gorm"
)

const (
	AuthRelationTypA2r = iota + 1 // 类型 - 账号角色
	AuthRelationTypR2m            // 类型 - 角色菜单
)

// AuthRelation  授权关系表
type AuthRelation struct {
	Typ       int             `json:"typ"`       // 类型：1-账号角色，2-角色菜单
	Aid       int64           `json:"aid"`       // 账号表id
	Rid       int64           `json:"rid"`       // 角色表id
	Mid       int64           `json:"mid"`       // 菜单表id
	CreatedAt model.LocalTime `json:"createdAt"` // 创建时间
}

// NewAuthRelation 初始化授权关系
func NewAuthRelation() *AuthRelation {
	return &AuthRelation{}
}

// GetAccountRids 基于账号获取授权的角色信息
func (ar *AuthRelation) GetAccountRids(aid int64) (rids []int64, err error) {
	err = helper.GormDefaultDb.Model(NewAuthRelation()).Where("typ = ? and aid = ?", AuthRelationTypA2r, aid).Pluck("rid", &rids).Error
	return
}

// GetAccountMids 基于账号获取授权的菜单ID信息
func (ar *AuthRelation) GetAccountMids(aid int64) (mids []int64, err error) {
	ridsModel := helper.GormDefaultDb.Model(NewAuthRelation()).Select("rid").Where("typ = ? and aid = ?", AuthRelationTypA2r, aid)
	err = helper.GormDefaultDb.Model(NewAuthRelation()).Model(NewAuthRelation()).Where("typ = ? and rid in (?)", AuthRelationTypR2m, ridsModel).Distinct().Pluck("mid", &mids).Error
	return
}

// GetRoleMids 基于角色获取授权的菜单ID集合
func (ar *AuthRelation) GetRoleMids(rid int64) (mids []int64, err error) {
	err = helper.GormDefaultDb.Model(NewAuthRelation()).Where("typ = ? and rid = ?", AuthRelationTypR2m, rid).Pluck("mid", &mids).Error
	return
}

// CreateRids 批量创建账号角色
func (ar *AuthRelation) CreateRids(tx *gorm.DB, aid int64, rids []int64) (err error) {
	if len(rids) == 0 {
		return
	}
	authRelations := make([]AuthRelation, 0)
	for _, rid := range rids {
		authRelations = append(authRelations, AuthRelation{
			Typ: AuthRelationTypA2r,
			Aid: aid,
			Rid: rid,
		})
	}
	err = tx.Create(&authRelations).Error
	return
}

// DeleteRids 批量删除账号角色
func (ar *AuthRelation) DeleteRids(tx *gorm.DB, aid int64, rids []int64) (err error) {
	if len(rids) == 0 {
		return
	}
	err = tx.Model(NewAuthRelation()).Where("typ = ? and aid = ? and rid in ?", AuthRelationTypA2r, aid, rids).Delete(&ar).Error
	return
}

// CreateMids 批量创建角色菜单
func (ar *AuthRelation) CreateMids(tx *gorm.DB, rid int64, mids []int64) (err error) {
	if len(mids) == 0 {
		return
	}
	authRelations := make([]AuthRelation, 0)
	for _, mid := range mids {
		authRelations = append(authRelations, AuthRelation{
			Typ: AuthRelationTypR2m,
			Rid: rid,
			Mid: mid,
		})
	}
	err = tx.Create(&authRelations).Error
	return
}

func (ar *AuthRelation) DeleteMids(tx *gorm.DB, rid int64, mids []int64) (err error) {
	if len(mids) == 0 {
		return
	}
	err = tx.Model(NewAuthRelation()).Where("typ = ? and rid = ? and mid in ?", AuthRelationTypR2m, rid, mids).Delete(&ar).Error
	return
}
