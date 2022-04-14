package system

import (
	"project/global"
	"project/model/system"
	"strconv"
)

type MenuService struct{}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: GetInfoList
//@description: 获取路由分页
//@return: err error, list interface{}, total int64

func (menuService *MenuService) GetMenuList() (err error, list interface{}, total int64) {
	var menuList []system.SysBaseMenu
	err, treeMap := menuService.getBaseMenuTreeMap()
	menuList = treeMap["0"]
	for i := 0; i < len(menuList); i++ {
		err = menuService.getBaseChildrenList(&menuList[i], treeMap)
	}
	return err, menuList, total
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: getBaseMenuTreeMap
//@description: 获取路由总树map
//@return: err error, treeMap map[string][]model.SysBaseMenu

func (menuService *MenuService) getBaseMenuTreeMap() (err error, treeMap map[string][]system.SysBaseMenu) {
	var allMenu []system.SysBaseMenu
	treeMap = make(map[string][]system.SysBaseMenu)
	err = global.GSD_DB.Order("sort").Preload("Parameters").Find(&allMenu).Error
	for _, v := range allMenu {
		treeMap[v.ParentId] = append(treeMap[v.ParentId], v)
	}
	return err, treeMap
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: getBaseMenuTreeMap
//@description: 获取菜单子菜单
//@param: menu *model.SysBaseMenu, treeMap map[string][]model.SysBaseMenu
//@return: err error

func (menuService *MenuService) getBaseChildrenList(menu *system.SysBaseMenu, treeMap map[string][]system.SysBaseMenu) (err error) {
	menu.Children = treeMap[strconv.Itoa(int(menu.ID))]
	for i := 0; i < len(menu.Children); i++ {
		err = menuService.getBaseChildrenList(&menu.Children[i], treeMap)
	}
	return err
}
