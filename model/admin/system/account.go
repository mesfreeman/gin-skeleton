package system

import (
	"errors"
	"fmt"
	"gin-skeleton/helper"
	"gin-skeleton/helper/encrypt"
	"gin-skeleton/model"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Account  账号表
type Account struct {
	model.BaseModel
	Username string          `json:"username"`      // 账号
	Nickname string          `json:"nickname"`      // 昵称
	Password string          `json:"-"`             // 密钥
	Avatar   string          `json:"avatar"`        // 头像
	Email    string          `json:"email"`         // 邮箱
	Phone    string          `json:"phone"`         // 手机号
	Status   int             `json:"status"`        // 状态：1-禁用，2-启用
	Remark   string          `json:"remark"`        // 备注
	LoginAt  model.LocalTime `json:"loginAt"`       // 最后登录时间
	Rids     []int64         `json:"rids" gorm:"-"` // 角色ID
}

// MyInfo 我的账号信息
type MyInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	Status   int    `json:"-"`
}

// NewAccount 初始化账号
func NewAccount() *Account {
	return &Account{}
}

// GetAccountList 获取账号列表
func (a *Account) GetAccountList(name string, status int, pageInfo model.BasePageParams) (pr *model.BasePageResult[Account], err error) {
	accountModel := helper.GormDefaultDb.Model(NewAccount())
	if name != "" {
		accountModel.Where("username like ? or nickname like ?", "%"+name+"%", "%"+name+"%")
	}
	if status > 0 {
		accountModel.Where("status = ?", status)
	}

	pr = &model.BasePageResult[Account]{Items: make([]*Account, 0), Total: 0}
	err = accountModel.Scopes(model.Paginate(pageInfo)).Find(&pr.Items).Scopes(model.CancelPaginate).Count(&pr.Total).Error
	if err != nil || len(pr.Items) == 0 {
		return
	}

	// 获取账号相关的角色ID
	for _, account := range pr.Items {
		rids, err2 := NewAuthRelation().GetAccountRids(account.ID)
		if err2 != nil {
			err = err2
			return
		}
		account.Rids = rids
	}
	return
}

// FindMyInfo 查找指定账号ID的个人信息
func (a *Account) FindMyInfo(id int64) (loginInfo *MyInfo, err error) {
	err = helper.GormDefaultDb.Model(NewAccount()).First(&loginInfo, id).Error
	return
}

// FindBasicInfo 查询基本的账号信息
func (a *Account) FindBasicInfo(id int64) (account *Account, err error) {
	err = helper.GormDefaultDb.First(&account, id).Error
	return
}

// FindFullInfo 查询完整的账号信息
func (a *Account) FindFullInfo(id int64) (account *Account, err error) {
	account, err = a.FindBasicInfo(id)
	if err != nil || errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	// 获取账号角色
	rids, err := NewAuthRelation().GetAccountRids(account.ID)
	if err != nil {
		return
	}
	account.Rids = rids
	return
}

// FindByUsername 基于账号查找账号信息
func (a *Account) FindByUsername(username string) (account *Account, err error) {
	err = helper.GormDefaultDb.Where("username = ?", username).First(&account).Error
	return
}

// FindByEmail 基于邮箱查找账号信息
func (a *Account) FindByEmail(email string) (account *Account, err error) {
	err = helper.GormDefaultDb.Where("email = ?", email).First(&account).Error
	return
}

// EncryptPassword 加密密码
func (a *Account) EncryptPassword(password string) string {
	return encrypt.GetMD5(fmt.Sprintf("%s@@%s", viper.GetString("Server.PwdSalt"), password))
}

// HasPermission 判断指定账号是否授权某个权限
func (a *Account) HasPermission(aid int64, username string, path string) (has bool, err error) {
	myPerms, err := NewMenu().GetMyPerms(aid, username)
	if err != nil {
		return
	}

	has = helper.InSlice(path, myPerms)
	return
}
