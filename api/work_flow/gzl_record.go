package work_flow

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project/utils"
	"strconv"

	"project/global"
	"project/model/common/response"
	WorkFlowReq "project/model/work_flow/request"
)

type RecordApi struct {
}

// Submit
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags RecordById
// @Summary 提交表单
// @Produce  application/json
// @Param data body WorkFlowReq.RecordSubmit true "string"
// @Success 200 {} json "{"success":true,"data":{},"msg":"null"}"
// @Router /record/submit [post]
func (r *RecordApi) Submit(c *gin.Context) {
	_, _ = strconv.Atoi(c.Request.Header.Get("x-user-id"))
	var recordSubmit WorkFlowReq.RecordSubmit
	err := c.ShouldBindJSON(&recordSubmit)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("json解析失败", zap.Any("err", err))
	}
	fmt.Println("recordSubmit: ", recordSubmit)
	recordSubmit.CreateBy = uint(1)
	//if err = utils.Verify(recordSubmit, utils.RecordSubmitVerify); err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	err = recordService.Submit(recordSubmit)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("记录提交失败", zap.Any("err", err))
		response.FailWithMessage("提交失败", c)
	} else {
		response.OkWithMessage("提交成功", c)
	}
}

// Data
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags RecordById
// @Summary 返回之前填写过的表单数据
// @Produce  application/json
// @Param recordId query int true "记录id"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /record/data [get]
func (r *RecordApi) Data(c *gin.Context) {
	recordId, _ := strconv.Atoi(c.Query("recordId"))
	record := WorkFlowReq.RecordById{RecordId: recordId}
	if err := utils.Verify(record, utils.RecordIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := recordService.GetData(recordId)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("表单数据获取失败", zap.Any("err", err))
		response.FailWithMessage("该记录不存在", c)
	} else {
		global.GSD_LOG.ZapLog.Info("表单数据获取成功", zap.Any("RecordById Form Data([]byte -> string)", string(data)))
		response.OkWithData(data, c)
	}
}

// Launch
// @Tags RecordById
// @Summary 我发起的
// @Produce  application/json
// @Param data body uint true "创建人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我发起的任务成功"}"
// @Router /record/schedule [get]
func (r *RecordApi) Launch(c *gin.Context) {

}
