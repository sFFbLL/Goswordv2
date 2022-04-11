package work_flow

import (
	"github.com/gin-gonic/gin"
	WorkFlowReq "project/model/work_flow/request"
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
