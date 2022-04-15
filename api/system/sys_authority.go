package system

import (
	"project/global"
	"project/model/common/request"
	"project/model/common/response"
	"project/model/system"
	systemReq "project/model/system/request"
	systemRes "project/model/system/response"
	"project/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityApi struct {
}

// CreateAuthority @Tags Authority
// @Summary 创建角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysAuthority true "权限id, 权限名, 父角色id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /authority/createAuthority [post]
func (a *AuthorityApi) CreateAuthority(c *gin.Context) {
	var r systemReq.SysAuthorityCreateRequest
	_ = c.ShouldBindJSON(&r)
	if err := utils.Verify(r, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//数据权限校验
	curUser := utils.GetUser(c)
	if r.Level < dataScope.GetMaxLevel(curUser.Authority) {
		global.GSD_LOG.Error(c, "无权创建该角色!")
		response.FailWithMessage("当前角色级别过高无权创建", c)
		return
	}
	var depts []system.SysDept
	if r.DataScope == "自定义" {
		for _, v := range r.DeptId {
			depts = append(depts, system.SysDept{
				GSD_MODEL: global.GSD_MODEL{ID: v},
			})
		}
	}
	authority := system.SysAuthority{CreateBy: curUser.ID, UpdateBy: curUser.ID, AuthorityName: r.AuthorityName, DataScope: r.DataScope, Depts: depts, Level: r.Level}
	if err, authBack := authorityService.CreateAuthority(authority); err != nil {
		global.GSD_LOG.Error(c, "创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authBack}, "创建成功", c)
	}
}

// DeleteAuthority @Tags Authority
// @Summary 删除角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysAuthority true "删除角色"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /authority/deleteAuthority [post]
func (a *AuthorityApi) DeleteAuthority(c *gin.Context) {
	var authority system.SysAuthority
	_ = c.ShouldBindJSON(&authority)
	if err := utils.Verify(authority, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//数据权限校验
	curUser := utils.GetUser(c)
	err, deleteAuthority := authorityService.GetAuthorityBasicInfo(system.SysAuthority{AuthorityId: authority.AuthorityId})
	if err != nil {
		global.GSD_LOG.Error(c, "该角色不存在!")
		response.FailWithMessage("该角色不存在", c)
		return
	}
	if deleteAuthority.Level < dataScope.GetMaxLevel(curUser.Authority) {
		global.GSD_LOG.Error(c, "无权创建该角色!")
		response.FailWithMessage("当前角色级别过高无权删除", c)
		return
	}
	authority.UpdateBy = curUser.ID
	if err := authorityService.DeleteAuthority(&authority); err != nil { // 删除角色之前需要判断是否有用户正在使用此角色
		global.GSD_LOG.Error(c, "删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败"+err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateAuthority @Tags Authority
// @Summary 更新角色信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysAuthority true "权限id, 权限名, 父角色id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /authority/updateAuthority [post]
func (a *AuthorityApi) UpdateAuthority(c *gin.Context) {
	var r systemReq.SysAuthorityUpdateRequest
	_ = c.ShouldBindJSON(&r)
	if err := utils.Verify(r, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	curUser := utils.GetUser(c)
	if r.Level < dataScope.GetMaxLevel(curUser.Authority) {
		global.GSD_LOG.Error(c, "更新失败, 修改级别高于当前用户级别!")
		response.FailWithMessage("更新失败, 修改级别高于当前用户级别", c)
		return
	}
	authority := system.SysAuthority{CreateBy: curUser.ID, UpdateBy: curUser.ID, AuthorityName: r.AuthorityName, DataScope: r.DataScope, Level: r.Level, AuthorityId: r.AuthorityId}
	if err, authority := authorityService.UpdateAuthority(authority, r.DeptId); err != nil {
		global.GSD_LOG.Error(c, "更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authority}, "更新成功", c)
	}
}

// GetAuthorityList @Tags Authority
// @Summary 分页获取角色列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /authority/getAuthorityList [post]
func (a *AuthorityApi) GetAuthorityList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, list, total := authorityService.GetAuthorityInfoList(pageInfo); err != nil {
		global.GSD_LOG.Error(c, "获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败"+err.Error(), c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
