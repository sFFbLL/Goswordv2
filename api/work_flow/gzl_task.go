package work_flow

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project/global"
	"project/model/common/response"
	WorkFlow "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	"project/utils"
	"strconv"
)

type TaskApi struct {
}

// Inspect
// @Tags Task
// @Summary 审批（通过||拒绝）
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
	gzlTask := &WorkFlow.GzlTask{GSD_MODEL: global.GSD_MODEL{ID: inspect.TaskId, UpdateBy: utils.GetUserID(c)}, CheckState: inspect.State}
	if err := taskService.Inspect(*gzlTask); err != nil {
		global.GSD_LOG.Error(c, "审批错误", zap.Any("err", err))
		response.FailWithMessage("审批错误", c)
		return
	} else {
		response.OkWithMessage("审批成功", c)
	}

}

// Dynamic
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags Task
// @Summary 流程动态信息
// @Produce  application/json
// @Param recordId query int true "记录id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"ok"}"
// @Router /task/dynamic [get]
func (t *TaskApi) Dynamic(c *gin.Context) {
	_, _ = strconv.Atoi(c.Request.Header.Get("x-user-id"))
	recordId, _ := strconv.Atoi(c.Query("recordId"))
	record := WorkFlowReq.RecordById{RecordId: recordId}
	if err := utils.Verify(record, utils.RecordIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tasks, err := taskService.GetDynamic(1, recordId)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("获取流程动态错误", zap.Any("err", err))
		response.FailWithMessage("数据不存在", c)
	} else {
		global.GSD_LOG.ZapLog.Info("流程动态信息成功返回", zap.Any("success", tasks))
		response.OkWithData(tasks, c)
	}
}

// Schedule
// @Tags Task
// @Summary 我的待办
// @Produce  application/json
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询待办任务成功"}"
// @Router /task/schedule [get]
func (t *TaskApi) Schedule(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("x-user-id"))
	var app WorkFlowReq.EmptyApp
	_ = c.ShouldBindJSON(&app)
	if err, schedule := taskService.GetScheduleList(userId, app.AppId); err != nil {
		global.GSD_LOG.ZapLog.Error("获取我的待办信息失败", zap.Error(err))
		response.FailWithMessage("获取我的待办信息失败", c)
		return
	} else {
		global.GSD_LOG.ZapLog.Info("获取成功", zap.Any("success", schedule))      //打印日志
		response.OkWithDetailed(gin.H{"schedule": schedule}, "获取我的待办信息成功", c) //给前端返回信息
	}
}

// Handle
// @Tags Task
// @Summary 我处理的
// @Produce  application/json
// @Param data body int true "审批人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我处理的任务成功"}"
// @Router /task/handle [get]
func (t *TaskApi) Handle(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("x-user-id"))
	var app WorkFlowReq.EmptyApp
	_ = c.ShouldBindJSON(&app)
	if err, handle := taskService.GetHandleList(userId, app.AppId); err != nil {
		global.GSD_LOG.ZapLog.Error("获取我处理的信息失败", zap.Error(err))
		response.FailWithMessage("获取我处理的信息失败", c)
		return
	} else {
		global.GSD_LOG.ZapLog.Info("获取成功", zap.Any("success", handle))    //打印日志
		response.OkWithDetailed(gin.H{"handle": handle}, "获取我处理的信息成功", c) //给前端返回信息
	}
}

// Receive
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags Task
// @Summary 我收到的
// @Produce  application/json
// @Param data query int  true "节点类型"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我收到的任务成功"}"
// @Router /task/receive [get]
func (t *TaskApi) Receive(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("x-user-id"))
	tasks, err := taskService.GetReceive(userId)
	// TODO ReceiveVerify && params select

	if err != nil {
		global.GSD_LOG.ZapLog.Error("获取我收到的信息列表错误", zap.Any("err", err))
		response.FailWithMessage("数据不存在", c)
	} else {
		global.GSD_LOG.ZapLog.Info("我收到的信息成功返回", zap.Any("success", tasks))
		response.OkWithData(tasks, c)
	}
}