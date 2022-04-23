package work_flow

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/global"
	"project/model/common/response"
	modelMF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	"project/utils"
)

type RecordApi struct {
}

// Submit
// @Tags Record
// @Summary 提交表单
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param data body WorkFlowReq.RecordSubmit true "string"
// @Success 200 {} json "{"success":true,"data":{},"msg":"null"}"
// @Router /record/submit [post]
func (r *RecordApi) Submit(c *gin.Context) {
	var recordSubmit WorkFlowReq.RecordSubmit
	err := c.ShouldBindJSON(&recordSubmit)
	if err != nil {
		global.GSD_LOG.Error("json解析失败", zap.Any("err", err), utils.GetRequestID(c))
	}
	if err = utils.Verify(recordSubmit, utils.RecordSubmitVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	form, err := json.Marshal(recordSubmit.Form)
	if err != nil {
		return
	}
	// 定义record结构体, 将入参的值赋进去
	var record modelMF.GzlRecord
	record.Form = form
	record.AppId = recordSubmit.AppId
	record.GSD_MODEL.CreateBy = utils.GetUserID(c)
	if err = recordService.Submit(record); err != nil {
		global.GSD_LOG.Error("记录提交失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("提交失败", c)
	} else {
		response.Ok(c)
	}
}

// Data
// @Tags Record
// @Summary 返回之前填写过的表单数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/json
// @Param recordId query WorkFlowReq.RecordById true "记录id"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /record/data [get]
func (r *RecordApi) Data(c *gin.Context) {
	var record WorkFlowReq.RecordById
	_ = c.ShouldBindQuery(&record)
	if err := utils.Verify(record, utils.RecordIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := recordService.GetData(record.RecordId)
	if err != nil && err != gorm.ErrRecordNotFound {
		global.GSD_LOG.Error("表单数据获取失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("该记录不存在", c)
	} else {
		global.GSD_LOG.Info("表单数据获取成功", zap.Any("Form Data: ", data), utils.GetRequestID(c))
		response.OkWithData(data, c)
	}
}

// MyInitiated
// @author: [chenpipi]
// @Tags Record
// @Summary 我发起的
// @Produce  application/json
// @Param data body uint true "创建人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /record/initiated [get]
func (r *RecordApi) MyInitiated(c *gin.Context) {
	myInitiated, err := recordService.MyInitiated(utils.GetUserID(c))
	if err != nil {
		global.GSD_LOG.Error("我发起的列表查询失败", zap.Any("err", err), utils.GetRequestID(c))
		response.FailWithMessage("数据获取失败", c)
	} else {
		global.GSD_LOG.Info("我发起的列表查询成功", zap.Any("Record Form Data([]byte -> string)", myInitiated), utils.GetRequestID(c))
		response.OkWithData(myInitiated, c)
	}
}
