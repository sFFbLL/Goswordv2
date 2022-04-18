package system

import (
	"github.com/gin-gonic/gin"
	v1 "project/api"
)

type MenuRouter struct {
}

func (s *MenuRouter) InitMenuRouter(Router *gin.RouterGroup) {
	menuRouter := Router.Group("menu")
	authorityMenuApi := v1.ApiGroupApp.SystemApiGroup.AuthorityMenuApi
	{
		menuRouter.POST("getMenuList", authorityMenuApi.GetMenuList)   // 分页获取基础menu列表
		menuRouter.POST("addBaseMenu", authorityMenuApi.AddMenu)       //新增菜单
		menuRouter.POST("deleteBaseMenu", authorityMenuApi.DeleteMenu) //删除菜单
		menuRouter.POST("updateBaseMenu", authorityMenuApi.UpdateMenu) //修改菜单
		menuRouter.POST("getMenu", authorityMenuApi.GetUserMenuTree)   //当前用户菜单
	}
}
