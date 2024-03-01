package system

import (
	"fmt"
	"gin-skeleton/helper"
	"gin-skeleton/model"
)

const (
	MenuTypeDir  = 1 // 类型 - 目录
	MenuTypeMenu = 2 // 类型 - 菜单
	MenuTypeApi  = 3 // 类型 - 接口

	MenuModeComponent = 1 // 模式 - 组件
	MenuModeInnerLink = 2 // 模式 - 内链
	MenuModeOuterLink = 3 // 模式 - 外链

	MenuIsShowNo  = 1 // 是否显示 - 否
	MenuIsShowYes = 2 // 是否显示 - 是

	MenuStatusOff = 1 // 状态 - 禁用
	MenuStatusOn  = 2 // 状态 - 启用
)

// Menu  菜单表
type Menu struct {
	model.BaseModel
	Pid       int64  `json:"pid"`       // 父id
	Name      string `json:"name"`      // 名称
	Icon      string `json:"icon"`      // 图标
	Path      string `json:"path"`      // 路径
	Component string `json:"component"` // 组件
	Type      int    `json:"type"`      // 类型：0-目录，1-菜单，2-按钮
	Mode      int    `json:"mode"`      // 模式：1-组件，2-内链，3-外链
	Weight    int    `json:"weight"`    // 权重，值越大越靠前
	Level     int    `json:"-"`         // 等级，表示几级菜单
	IsShow    int    `json:"isShow"`    // 是否显示：1-否，2-是
	Status    int    `json:"status"`    // 状态：0-禁用，1-启用
}

// 我的菜单
type myMenu struct {
	ID        int64                  `json:"id"`
	Pid       int64                  `json:"pid"`
	Name      string                 `json:"name"`
	Path      string                 `json:"path"`
	Component string                 `json:"component"`
	Meta      map[string]interface{} `json:"meta"`
	Children  []*myMenu              `json:"children"`
}

// 系统菜单
type sysMenu struct {
	Menu
	Pname    string     `json:"pname"`
	Children []*sysMenu `json:"children"`
}

// NewMenu 初始化菜单
func NewMenu() *Menu {
	return &Menu{}
}

// GetMyMenus 获取我的菜单
func (m *Menu) GetMyMenus(aid int64, username string) (myMenus []*myMenu, err error) {
	var myMenuData []Menu
	var menuModel = helper.GormDefaultDb.Model(NewMenu())

	// 非超级账号，只获取拥有权限的菜单
	if !helper.IsSuperAccount(username) {
		mids, err2 := NewAuthRelation().GetAccountMids(aid)
		if err2 != nil || len(mids) == 0 {
			err = err2
			return
		}

		menuModel.Where("id in ?", mids)
	}

	// 获取菜单数据
	err = menuModel.Where("status = ? and type in ?", MenuStatusOn, []int{MenuTypeDir, MenuTypeMenu}).
		Order("level asc, weight desc, id asc").Find(&myMenuData).Error
	if err != nil || len(myMenuData) == 0 {
		return
	}

	// 生成树结构
	myMenus = generateMyMenusTree(myMenuData, 0)
	return
}

// GetHomePath 获取我的首页路径
func (m *Menu) GetHomePath(myMenus []*myMenu) (path string, err error) {
	for _, myMenu := range myMenus {
		if myMenu.Children == nil && len(myMenu.Children) == 0 {
			return myMenu.Path, nil
		}
		return m.GetHomePath(myMenu.Children)
	}
	return "", err
}

// GetMyPerms 获取我的权限代码
func (m *Menu) GetMyPerms(aid int64, username string) (myPerms []string, err error) {
	myPerms = make([]string, 0)
	var menuModel = helper.GormDefaultDb.Model(NewMenu())

	// 非超级账号，只获取拥有的权限代码
	if !helper.IsSuperAccount(username) {
		mids, err2 := NewAuthRelation().GetAccountMids(aid)
		if err2 != nil || len(mids) == 0 {
			err = err2
			return
		}

		menuModel.Where("id in ?", mids)
	}

	// 获取权限代码
	err = menuModel.Where("type = ? and status = ?", MenuTypeApi, MenuStatusOn).Pluck("path", &myPerms).Error
	return
}

// GetSysMenus 获取系统菜单
func (m *Menu) GetSysMenus(name string, status int, typ []int) (sysMenus []*sysMenu, err error) {
	menuModel := helper.GormDefaultDb.Model(NewMenu())
	if name != "" {
		menuModel.Where("name like ?", "%"+name+"%")
	}
	if len(typ) > 0 {
		menuModel.Where("type in ?", typ)
	}
	if status > 0 {
		menuModel.Where("status = ?", status)
	}

	var sysMenuData []Menu
	err = menuModel.Order("level asc, weight desc, id asc").Find(&sysMenuData).Error
	if err != nil || len(sysMenuData) == 0 {
		return
	}

	// 获取最小父级ID，为了解决搜索名称时，生成菜单树的情况
	minPid := sysMenuData[0].Pid
	sysMenuMap := map[int64]string{0: "顶级目录"}
	for _, sysMenu := range sysMenuData {
		sysMenuMap[sysMenu.ID] = sysMenu.Name
		if sysMenu.Pid < minPid {
			minPid = sysMenu.Pid
		}
	}

	// 生成树结构
	sysMenus = generateSysMenusTree(sysMenuData, sysMenuMap, minPid)
	return
}

// FindMenuInfo 查找菜单信息
func (m *Menu) FindMenuInfo(id int64) (menu *Menu, err error) {
	err = helper.GormDefaultDb.First(&menu, id).Error
	return
}

// FindByNameLevel 基于名称、等级查找菜单信息
func (m *Menu) FindByNameLevel(name string, level int) (menu *Menu, err error) {
	err = helper.GormDefaultDb.Where("name = ? and level = ?", name, level).First(&menu).Error
	return
}

// FindNameByPath 通过路径查询菜单名称
func (m *Menu) FindNameByPath(path string) (name string, err error) {
	err = helper.GormDefaultDb.Model(m).Where("path = ?", path).Pluck("name", &name).Error
	return
}

// 生成我的菜单树
func generateMyMenusTree(menus []Menu, pid int64) (myMenus []*myMenu) {
	for _, menu := range menus {
		if menu.Pid != pid {
			continue
		}

		// 递归处理子菜单
		children := generateMyMenusTree(menus, menu.ID)

		// 处理内嵌网页的菜单
		path, frameSrc := menu.Path, ""
		if menu.Mode == MenuModeInnerLink {
			path, frameSrc = fmt.Sprintf("iframe%d", menu.ID), menu.Path
		}

		// 拼接菜单数据
		myMenus = append(myMenus, &myMenu{
			ID:        menu.ID,
			Pid:       menu.Pid,
			Name:      menu.Name,
			Path:      path,
			Component: menu.Component,
			Meta: map[string]interface{}{
				"title":    menu.Name,
				"icon":     menu.Icon,
				"frameSrc": frameSrc,
				"hideMenu": menu.IsShow == MenuIsShowNo,
			},
			Children: children,
		})
	}
	return
}

// 生成系统菜单树
func generateSysMenusTree(sysMenuData []Menu, sysMenuMap map[int64]string, pid int64) (SysMenus []*sysMenu) {
	for _, menu := range sysMenuData {
		if menu.Pid != pid {
			continue
		}

		// 递归子菜单
		pname := sysMenuMap[menu.Pid]
		children := generateSysMenusTree(sysMenuData, sysMenuMap, menu.ID)
		SysMenus = append(SysMenus, &sysMenu{
			menu,
			pname,
			children,
		})
	}
	return
}
