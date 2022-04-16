package work_flow

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project/global"
	"project/model/common/response"
	WorkFlowReq "project/model/work_flow/request"
	"strconv"
)

type TaskApi struct {
}

// Inspect
// @Tags Task
// @Summary 审批（通过||拒绝）
// @Produce  application/json
// @Param data body WorkFlowReq.Task true "通过||拒绝"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"已审核"}"
// @Router /task/inspect [post]
func (t *TaskApi) Inspect(c *gin.Context) {
	var _ WorkFlowReq.Task
}

// Dynamic
// @Tags Task
// @Summary 流程动态信息
// @Produce  application/json
// @Param data body int true "string"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"ok"}"
// @Router /task/dynamic [get]
func (t *TaskApi) Dynamic(c *gin.Context) {
	var _ []WorkFlowReq.Task
}

// Schedule
// @Tags Task
// @Summary 我的待办
// @Produce  application/json
// @Param data body int true "审批状态, 审批人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询待办任务成功"}"
// @Router /task/schedule [get]
func (t *TaskApi) Schedule(c *gin.Context){
	InspectorId, _ := strconv.Atoi(c.Request.Header.Get("x-user-id"))
	if err, schedule := taskService.GetScheduleList(InspectorId); err != nil {
		global.GSD_LOG.ZapLog.Error( "获取我的待办信息失败", zap.Error(err))
		response.FailWithMessage("获取我的待办信息失败", c)
		return
	} else {
		global.GSD_LOG.ZapLog.Info("获取成功", zap.Any("success", schedule))//打印日志
		response.OkWithDetailed(gin.H{"schedule": schedule}, "获取我的待办信息成功", c)//给前端返回信息
	}
}

// Handle
// @Tags Task
// @Summary 我处理的
// @Produce  application/json
// @Param data body int  true "审批状态, 审批人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我处理的任务成功"}"
// @Router /task/handle [get]
func (t *TaskApi) Handle(c *gin.Context) {

}

// Receive
// @Tags Task
// @Summary 我收到的
// @Produce  application/json
// @Param data body  uint8  true "节点类型"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我收到的任务成功"}"
// @Router /task/receive [get]
func (t *TaskApi) Receive(c *gin.Context) {

}
