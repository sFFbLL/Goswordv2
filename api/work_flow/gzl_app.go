package work_flow

import (
	"github.com/gin-gonic/gin"
	WorkFlowReq "project/model/work_flow"
)

type AppApi struct {
}

// Empty
// @Tags app
// @Summary 返回空表单
// @Produce  application/json
// @Param data body int true "string"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /app/empty [get]
func (f *AppApi) Empty(c *gin.Context) {
	var _ WorkFlowReq.GzlApp
}
