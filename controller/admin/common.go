package admin

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"

	"gin-skeleton/helper"
	"gin-skeleton/helper/jwt"
	"gin-skeleton/helper/log"
	"gin-skeleton/helper/response"
	"gin-skeleton/helper/tool"
	"gin-skeleton/middleware"
	"gin-skeleton/model"
	"gin-skeleton/model/admin/system"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Login 账号登录
func Login(c *gin.Context) {
	var params struct {
		Username string `json:"username" remark:"账号" binding:"required"`
		Password string `json:"password" remark:"密码" binding:"required"`
		Remember bool   `json:"remember"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// @todo 增加IP试错限制

	// 获取账号信息
	account, err := system.NewAccount().FindByUsername(params.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		go addLoginLog(c, params.Username, "", system.LoginLogOperTypeFail, "系统出错了："+err.Error())
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if account.ID == 0 || account.Password != system.NewAccount().EncryptPassword(params.Password) {
		go addLoginLog(c, params.Username, account.Nickname, system.LoginLogOperTypeFail, "账号或密码错误")
		response.InvalidArgumentJSON("账号或密码错误", c)
		return
	}

	// 判断是否允许登录
	if account.Status == system.AccountStatusDisable {
		go addLoginLog(c, params.Username, account.Nickname, system.LoginLogOperTypeFail, "账号已被禁用")
		response.ForbiddenJSON("账号已被禁用", c)
		return
	}

	// 令牌生命周期
	tokenLifeTime := 24
	if params.Remember {
		tokenLifeTime = tokenLifeTime * 7
	}

	// 生成登录令牌
	token, err := jwt.GenerateJwtToken(*account, tokenLifeTime)
	if err != nil {
		go addLoginLog(c, params.Username, account.Nickname, system.LoginLogOperTypeFail, "生成登录令牌出错了："+err.Error())
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	// 更新最后登录时间并返回
	account.LoginAt = model.LocalTime(time.Now())
	helper.GormDefaultDb.Save(account)
	go addLoginLog(c, params.Username, account.Nickname, system.LoginLogOperTypeSuccess, "")
	response.SuccessJSON(gin.H{"token": token}, "登录成功", c)
}

// Logout 账号退出
func Logout(c *gin.Context) {
	go addLoginLog(c, middleware.GetTokenAuthInfo(c).Username, middleware.GetTokenAuthInfo(c).Nickname, system.LoginLogOperTypeLogout, "")
	response.SuccessJSON(gin.H{"id": middleware.GetTokenAuthInfo(c).ID}, "退出成功", c)
}

// MyInfo 我的个人信息
func MyInfo(c *gin.Context) {
	myInfo, err := system.NewAccount().FindMyInfo(middleware.GetTokenAuthInfo(c).ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if myInfo.ID <= 0 {
		response.InvalidArgumentJSON("账户不存在", c)
		return
	}

	// 判断是否允许登录
	if myInfo.Status == system.AccountStatusDisable {
		response.ForbiddenJSON("账号已被禁用", c)
		return
	}
	response.SuccessJSON(myInfo, "登录成功", c)
}

// MyMenus 我的权限菜单
func MyMenus(c *gin.Context) {
	myMenus, err := system.NewMenu().GetMyMenus(middleware.GetTokenAuthInfo(c).ID, middleware.GetTokenAuthInfo(c).Username)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(myMenus, "", c)
}

// MyPerms 我的权限代码
func MyPerms(c *gin.Context) {
	myPerms, err := system.NewMenu().GetMyPerms(middleware.GetTokenAuthInfo(c).ID, middleware.GetTokenAuthInfo(c).Username)
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(myPerms, "success", c)
}

// ModifyMyPwd 修改我的账号密码
func ModifyMyPwd(c *gin.Context) {
	var params struct {
		OldPwd string `json:"oldPwd" remark:"旧密码" binding:"required,min=6,max=15"`
		NewPwd string `json:"newPwd" remark:"新密码" binding:"required,min=6,max=15"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断账号是否存在
	account, err := system.NewAccount().FindBasicInfo(middleware.GetTokenAuthInfo(c).ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if account.ID == 0 {
		response.InvalidArgumentJSON("账号不存在", c)
		return
	}

	// 判断旧密码是否正确
	if system.NewAccount().EncryptPassword(params.OldPwd) != account.Password {
		response.InvalidArgumentJSON("旧密码错误", c)
		return
	}

	// 修改为新密码
	account.Password = system.NewAccount().EncryptPassword(params.NewPwd)
	if err := helper.GormDefaultDb.Save(&account).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdResult{ID: account.ID}, "修改密码成功", c)
}

// ModifyMyInfo 修改我的个人信息
func ModifyMyInfo(c *gin.Context) {
	var params struct {
		Nickname string `json:"nickname" remark:"昵称" binding:"required,min=2,max=32"`
		Email    string `json:"email" remark:"邮箱" binding:"required,email"`
		Phone    string `json:"phone" remark:"手机号" binding:"omitempty,len=11"`
		Avatar   string `json:"avatar" remark:"头像" binding:"max=255"`
	}
	if err := c.ShouldBindJSON(&params); err != nil {
		response.ValidatorFailedJson(err, c)
		return
	}

	// 判断账号是否存在
	account, err := system.NewAccount().FindBasicInfo(middleware.GetTokenAuthInfo(c).ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	if account.ID == 0 {
		response.InvalidArgumentJSON("账号不存在", c)
		return
	}

	// 更新账号信息
	account.Nickname = params.Nickname
	account.Email = params.Email
	account.Phone = params.Phone
	account.Avatar = params.Avatar
	if err := helper.GormDefaultDb.Save(&account).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(model.BaseIdResult{ID: account.ID}, "保存成功", c)
}

// LiteRoles 获取角色简单信息
func LiteRoles(c *gin.Context) {
	var liteRoles []struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	err := helper.GormDefaultDb.Model(system.NewRole()).Select("id, name").Where("status = ?", system.RoleStatusOn).
		Order("weight desc, id asc").Find(&liteRoles).Error
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(liteRoles, "", c)
}

// LiteAccounts 获取账号简单信息
func LiteAccounts(c *gin.Context) {
	var liteAccounts []struct {
		Email   string `json:"email"`
		Account string `json:"account"`
	}
	err := helper.GormDefaultDb.Model(system.NewAccount()).Select("concat(nickname, '(', username, ')') as account, email").Where("status = ?", system.AccountStatusNormal).
		Find(&liteAccounts).Error
	if err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}
	response.SuccessJSON(liteAccounts, "", c)
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.LogicExceptionJSON("获取文件数据失败："+err.Error(), c)
		return
	}

	if file.Size > 10*1024*1024 {
		response.InvalidArgumentJSON("文件大小不能超过10M", c)
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

	// 判断文件类型是否支持
	fileType := strings.TrimPrefix(path.Ext(file.Filename), ".")
	if !helper.InSilce(fileType, fileConfig.AllowTypes) {
		response.InvalidArgumentJSON(fmt.Sprintf("文件类型暂不支持[%s]", fileType), c)
		return
	}

	// 打开文件流
	fileIo, err := file.Open()
	if err != nil {
		response.LogicExceptionJSON("获取文件内容异常："+err.Error(), c)
		return
	}
	defer fileIo.Close()

	// 调用三方服务上传文件
	fileUrl, err := tool.NewStorage(fileConfig).PutFileByIo("admin", file.Filename, file.Size, fileIo)
	if err != nil {
		response.ThirdExceptionJSON("上传文件到三方服务异常："+err.Error(), c)
		return
	}

	// 缩略图
	thumbnail := ""
	if fileConfig.ThumbConf != "" {
		thumbnail = fmt.Sprintf("%s?%s", fileUrl, fileConfig.ThumbConf)
	}

	// 添加文件上传记录到文件管理
	fileInfo := system.File{
		FileName:  file.Filename,
		FileSize:  file.Size,
		FileType:  fileType,
		FileUrl:   fileUrl,
		Thumbnail: thumbnail,
		Provider:  fileConfig.Provider,
		Username:  middleware.GetTokenAuthInfo(c).Username,
		Nickname:  middleware.GetTokenAuthInfo(c).Nickname,
		Remark:    fmt.Sprintf("后台接口上传[%s]", fileConfig.Bucket),
	}
	if err := helper.GormDefaultDb.Create(&fileInfo).Error; err != nil {
		response.LogicExceptionJSON("系统出错了："+err.Error(), c)
		return
	}

	// 返回结果
	result := gin.H{
		"id":       fileInfo.ID,
		"fileName": fileInfo.FileName,
		"fileSize": fileInfo.FileSize,
		"fileType": fileInfo.FileType,
		"fileUrl":  fileInfo.FileUrl,
		"thumbail": fileInfo.Thumbnail,
		"provider": fileInfo.Provider,
	}
	response.SuccessJSON(result, "上传成功", c)
}

// 添加登录日志
func addLoginLog(c *gin.Context, username, nickname string, operType int, remark string) {
	userAgent := user_agent.New(c.Request.UserAgent())
	browserName, browserVersion := userAgent.Browser()

	loginLog := system.LoginLog{
		Username: username,
		Nickname: nickname,
		Ip:       c.ClientIP(),
		Device:   userAgent.Platform(),
		Os:       userAgent.OS(),
		Browser:  browserName + " " + browserVersion,
		OperType: operType,
		Remark:   remark,
	}
	if err := helper.GormDefaultDb.Create(&loginLog).Error; err != nil {
		log.GetLogger("common").WithFields(logrus.Fields{"loginLog": loginLog, "error": err}).Errorln("添加登录日志异常")
		return
	}
	return
}
