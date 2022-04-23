package work_flow

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project/global"
	"project/model/common/response"
	"project/model/system"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	"project/utils"
)

type AppApi struct {
}

// Empty
// @Tags App
// @Summary 返回空表单
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param appId query WorkFlowReq.EmptyApp true "应用id"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /app/empty [get]
func (f *AppApi) Empty(c *gin.Context) {
	var app WorkFlowReq.EmptyApp
	_ = c.ShouldBindQuery(&app)
	if err := utils.Verify(app, utils.EmptyAppVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := appService.GetAppEmpty(app.AppId)
	if err != nil {
		global.GSD_LOG.Error("获取空应用表单失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage(err.Error(), c)
	} else {
		global.GSD_LOG.Info("获取空应用表单成功", zap.Any("GetAppEmpty Success", data), utils.GetRequestID(c))
		response.OkWithData(data, c)
	}
}

// Create
// @Tags App
// @Summary 创建表单
// @Produce  application/json
// @Param data body uint true "创建人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建表单成功"}"
// @Router /app/create [post]
func (f *AppApi) Create(c *gin.Context) {

}

// AddApp
// @Tags App
// @Summary 创建应用
// @Produce  application/json
// @Param data body WorkFlowReq.AddApp true "应用名称，应用图标"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建表单成功"}"
// @Router /app/ [post]
func (f *AppApi) AddApp(c *gin.Context) {
	var addApp WorkFlowReq.AddApp
	_ = c.ShouldBindJSON(&addApp)
	if err := utils.Verify(addApp, utils.AddApp); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err := appService.AddApp(modelWF.GzlApp{GSD_MODEL: global.GSD_MODEL{CreateBy: utils.GetUserID(c)}, Name: addApp.Name, Icon: addApp.Icon})
	if err != nil {
		global.GSD_LOG.Error("添加应用失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("操作失败", c)
	} else {
		response.OkWithMessage("添加应用成功", c)
	}
}

// AddForm
// @Tags App
// @Summary 添加应用表单
// @Produce  application/json
// @Param data body WorkFlowReq.AddFlow true "应用id，表单"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建表单成功"}"
// @Router /app/form [post]
func (f *AppApi) AddForm(c *gin.Context) {
	var addForm WorkFlowReq.AddForm
	_ = c.ShouldBindJSON(&addForm)
	if err := utils.Verify(addForm, utils.AddForm); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	form, err := json.Marshal(addForm.Form)
	if err != nil {
		global.GSD_LOG.Error("json解析失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("json格式错误", c)
		return
	}
	err = appService.AddForm(modelWF.GzlApp{GSD_MODEL: global.GSD_MODEL{ID: addForm.AppId}, Form: form})
	if err != nil {
		global.GSD_LOG.Error("添加应用表单失败", zap.Any("err", err))
		response.FailWithMessage("添加应用表单失败", c)
	} else {
		response.OkWithMessage("表单提交成功", c)
	}
}

// AddFlow
// @Tags App
// @Summary 添加应用流程
// @Produce  application/json
// @Param data body WorkFlowReq.AddFlow true "应用id,流程"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建表单成功"}"
// @Router /app/flow [post]
func (f *AppApi) AddFlow(c *gin.Context) {
	var addFlow WorkFlowReq.AddFlow
	_ = c.ShouldBindJSON(&addFlow)
	if err := utils.Verify(addFlow, utils.AddFlow); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	flow, err := json.Marshal(addFlow.Flow)
	if err != nil {
		global.GSD_LOG.Error("json解析失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("json格式错误", c)
		return
	}
	if err = appService.AddForm(modelWF.GzlApp{GSD_MODEL: global.GSD_MODEL{ID: addFlow.AppId}, Flow: flow}); err != nil {
		global.GSD_LOG.Error("添加应用流程失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("添加应用流程失败", c)
	} else {
		response.OkWithMessage("流程提交成功", c)
	}
}

// EnableApp
// @Tags App
// @Summary 启用应用
// @Produce  application/json
// @Param data body WorkFlowReq.StartApp true "应用id,是否启用(1不启用默认，2启用)"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建表单成功"}"
// @Router /app/enable [post]
func (f *AppApi) EnableApp(c *gin.Context) {
	var enableApp WorkFlowReq.EnableApp
	_ = c.ShouldBindJSON(&enableApp)
	if err := utils.Verify(enableApp, utils.StartApp); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := appService.EnableApp(modelWF.GzlApp{GSD_MODEL: global.GSD_MODEL{ID: enableApp.AppId}, IsEnable: enableApp.IsEnable}); err != nil {
		global.GSD_LOG.Error("添加应用流程失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("添加应用流程失败", c)
	} else {
		response.OkWithMessage("流程提交成功", c)
	}
}

// AuthorityApp
// @Tags App
// @Summary 应用权限分配
// @Produce  application/json
// @Param data body WorkFlowReq.AuthorityApp true "应用id,部门s，角色s，用户s"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建表单成功"}"
// @Router /app/authority [post]
func (f *AppApi) AuthorityApp(c *gin.Context) {
	var authorityApp WorkFlowReq.AuthorityApp
	_ = c.ShouldBindJSON(&authorityApp)
	if err := utils.Verify(authorityApp, utils.AuthorityApp); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var depts []system.SysDept
	var authoritys []system.SysAuthority
	var users []system.SysUser
	for _, deptId := range authorityApp.Depts {
		depts = append(depts, system.SysDept{GSD_MODEL: global.GSD_MODEL{ID: deptId}})
	}
	for _, AuthorityId := range authorityApp.Authoritys {
		authoritys = append(authoritys, system.SysAuthority{AuthorityId: AuthorityId})
	}
	for _, userId := range authorityApp.Users {
		users = append(users, system.SysUser{GSD_MODEL: global.GSD_MODEL{ID: userId}})
	}
	if err := appService.AuthorityApp(modelWF.GzlApp{GSD_MODEL: global.GSD_MODEL{ID: authorityApp.AppId}}, depts, authoritys, users); err != nil {
		global.GSD_LOG.Error("添加应用流程失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("添加应用流程失败", c)
	} else {
		response.OkWithMessage("流程提交成功", c)
	}
}
