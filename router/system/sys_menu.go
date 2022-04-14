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
		menuRouter.POST("lists", authorityMenuApi.GetMenuList) // 分页获取基础menu列表
	}
}
