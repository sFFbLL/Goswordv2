package system

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project/global"
	"project/model/common/request"
	"project/model/common/response"
	"project/model/system"
	"project/utils"
)

type AuthorityMenuApi struct {
}

// @Tags Menu
// @Summary 分页获取基础menu列表
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取基础menu列表,返回包括列表,总数,页码,每页数量"
// @Router /menu/getMenuList [post]

func (a *AuthorityMenuApi) GetMenuList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, list, total := menuService.GetMenuList(); err != nil {
		global.GSD_LOG.Error(c, "获取失败", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags Menu
// @Summary 菜单管理-新增菜单
// @Produce application/json
// @Param data body system.SysBaseMenu true "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Success 200 {object} response.Response{msg=string} "新增菜单"
// @Router /menu/addition [post]

func (a *AuthorityMenuApi) AddMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	_ = c.ShouldBindJSON(&menu)
	if err := utils.Verify(menu, utils.MenuVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(menu.Meta, utils.MenuMetaVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := menuService.AddMenu(menu); err != nil {
		global.GSD_LOG.Error(c, "添加菜单失败", zap.Error(err))
		response.FailWithMessage("添加菜单失败", c)
	} else {
		response.OkWithMessage("添加菜单成功", c)
	}
}

// @Tags Menu
// @Summary 菜单管理-删除菜单
// @Produce application/json
// @Param data body request.GetById true "菜单id"
// @Success 200 {object} response.Response{msg=string} "删除菜单"
// @Router /menu/deleteMenu [post]

func (a *AuthorityMenuApi) DeleteMenu(c *gin.Context) {
	var menu request.GetById
	_ = c.ShouldBindJSON(&menu)
	if err := utils.Verify(menu, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := menuService.DeleteMenu(menu.ID); err != nil {
		global.GSD_LOG.Error(c, "删除菜单失败", zap.Error(err))
		response.FailWithMessage("菜单删除失败", c)
		return
	} else {
		response.OkWithMessage("删除菜单成功", c)
	}
}

// @Tags Menu
// @Summary 菜单管理-更新菜单
// @Produce application/json
// @Param data body system.SysBaseMenu true "路由path, 父菜单ID, 路由name, 对应前端文件路径, 排序标记"
// @Success 200 {object} response.Response{msg=string} "更新菜单"
// @Router /menu/updateMenu [post]

func (a *AuthorityMenuApi) UpdateMenu(c *gin.Context) {
	var menu system.SysBaseMenu
	_ = c.ShouldBindJSON(&menu)
	if err := utils.Verify(menu, utils.MenuVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(menu.Meta, utils.MenuMetaVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := menuService.UpdateMenu(menu); err != nil {
		global.GSD_LOG.Error(c, "更新菜单失败", zap.Error(err))
		response.FailWithMessage("菜单更新失败", c)
	} else {
		response.OkWithMessage("菜单更新成功", c)
	}

}
