package work_flow

import (
	"github.com/gin-gonic/gin"
	WorkFlowReq "project/model/work_flow"
)

type FormApi struct {
}

// FormEmpty
// @Tags form
// @Summary 需要用户填写的空表单
// @Produce  application/json
// @Param data body int true "string"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /form/formempty [get]
func (f *FormApi) FormEmpty(c *gin.Context) {
	var _ WorkFlowReq.GzlApp
}

// FormSave
// @Tags form
// @Summary 保存用户填写的表单
// @Produce  application/json
// @Param data body []WorkFlowReq.GzlFormItem true "string"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /form/formsave [post]
func (f *FormApi) FormSave(c *gin.Context) {
	var _ []WorkFlowReq.GzlFormItem
}

// FormData
// @Tags form
// @Summary 返回之前填写过的表单数据
// @Produce  application/json
// @Param data body WorkFlowReq.GzlAppUser true "string"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /form/formdata [get]
func (f *FormApi) FormData(c *gin.Context) {
	var _ WorkFlowReq.GzlAppUser
}
