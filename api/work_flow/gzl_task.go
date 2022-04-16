package work_flow

import (
	"fmt"
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
// @Success 200 {string} string "{"success":true,"data":{},"msg":"已审核"}"
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
	} else {
		response.OkWithMessage("审批成功", c)
	}
	//流程流转
	//go func() {
	//	work_flow.ProcessFlow(1)
	//}()
}

// Dynamic
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags Task
// @Summary 流程动态信息
// @Produce  application/json
// @Param data query WorkFlowReq.Record true "记录id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"ok"}"
// @Router /task/dynamic [get]
func (t *TaskApi) Dynamic(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.Header.Get("x-user-id"))
	var record WorkFlowReq.Record
	_ = c.ShouldBindJSON(&record)
	fmt.Println(userId)
	tasks, err := taskService.GetDynamic(1, record.RecordId)
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
// @Param data body int true "审批状态, 审批人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询待办任务成功"}"
// @Router /task/schedule [get]
func (t *TaskApi) Schedule(c *gin.Context) {
	var _ WorkFlow.GzlTask
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
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags Task
// @Summary 我收到的
// @Produce  application/json
// @Param data body  uint8  true "节点类型"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我收到的任务成功"}"
// @Router /task/receive [get]
func (t *TaskApi) Receive(c *gin.Context) {

}
