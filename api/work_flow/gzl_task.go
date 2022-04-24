package work_flow

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project/global"
	"project/model/common/response"
	WorkFlow "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	"project/utils"
)

type TaskApi struct {
}

// Inspect
// @Tags Task
// @Summary 审批（通过||拒绝）
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body WorkFlowReq.Inspect true "任务id，状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"审批成功"}"
// @Router /task/inspect [put]
func (t *TaskApi) Inspect(c *gin.Context) {
	var inspect WorkFlowReq.Inspect
	_ = c.ShouldBindJSON(&inspect)
	if err := utils.Verify(inspect, utils.InspectVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := taskService.Inspect(WorkFlow.GzlTask{GSD_MODEL: global.GSD_MODEL{ID: inspect.TaskId, UpdateBy: utils.GetUserID(c)}, CheckState: inspect.State}); err != nil {
		global.GSD_LOG.Error("审批错误", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("审批错误", c)
		return
	} else {
		response.OkWithMessage("审批成功", c)
	}
}

// Dynamic
// @Tags Task
// @Summary 流程动态信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param recordId query int true "记录id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"ok"}"
// @Router /task/dynamic [get]
func (t *TaskApi) Dynamic(c *gin.Context) {
	var record WorkFlowReq.RecordById
	_ = c.ShouldBindQuery(&record)
	if err := utils.Verify(record, utils.RecordIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tasks, err := taskService.GetDynamic(record.RecordId)
	if err != nil {
		global.GSD_LOG.Error("获取流程动态错误", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage(err.Error(), c)
	} else {
		global.GSD_LOG.Info("流程动态信息成功返回", zap.Any("success", tasks), utils.GetRequestID(c))
		response.OkWithData(tasks, c)
	}
}

// Schedule
// @Tags Task
// @Summary 我的待办
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询待办任务成功"}"
// @Router /task/schedule [get]
func (t *TaskApi) Schedule(c *gin.Context) {
	userId := utils.GetUserID(c)
	var app WorkFlowReq.EmptyApp
	_ = c.ShouldBindJSON(&app)
	if err, schedule := taskService.GetScheduleList(userId, app.AppId); err != nil {
		global.GSD_LOG.Error("获取我的待办信息失败", zap.Error(err), utils.GetRequestID(c))
		response.FailWithMessage("获取我的待办信息失败", c)
		return
	} else {
		global.GSD_LOG.Info("获取成功", zap.Any("success", schedule), utils.GetRequestID(c)) //打印日志
		response.OkWithDetailed(gin.H{"schedule": schedule}, "获取我的待办信息成功", c)            //给前端返回信息
	}
}

// Handle
// @Tags Task
// @Summary 我处理的
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data query WorkFlowReq.EmptyApp true "审批人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我处理的任务成功"}"
// @Router /task/handle [get]
func (t *TaskApi) Handle(c *gin.Context) {
	userId := utils.GetUserID(c)
	var app WorkFlowReq.EmptyApp
	_ = c.ShouldBindJSON(&app)
	if err, handle := taskService.GetHandleList(userId, app.AppId); err != nil {
		global.GSD_LOG.Error("获取我处理的信息失败", zap.Error(err), utils.GetRequestID(c))
		response.FailWithMessage("获取我处理的信息失败", c)
		return
	} else {
		global.GSD_LOG.Info("获取成功", zap.Any("success", handle), utils.GetRequestID(c)) //打印日志
		response.OkWithDetailed(gin.H{"handle": handle}, "获取我处理的信息成功", c)              //给前端返回信息
	}
}

// Receive
// @Tags Task
// @Summary 我收到的
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我收到的任务成功"}"
// @Router /task/receive [get]
func (t *TaskApi) Receive(c *gin.Context) {
	userId := utils.GetUserID(c)
	tasks, err := taskService.GetReceive(userId)
	if err != nil {
		global.GSD_LOG.Error("获取我收到的信息列表错误", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("数据不存在", c)
	} else {
		global.GSD_LOG.Info("我收到的信息成功返回", zap.Any("success", tasks), utils.GetRequestID(c))
		response.OkWithData(tasks, c)
	}
}
