package work_flow

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"project/global"
	"project/model/common/response"
	WorkFlowReq "project/model/work_flow/request"
)

type RecordApi struct {
}

// Submit
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags Record
// @Summary 提交表单
// @Produce  application/json
// @Param data body uint true "string"  //TODO 修改入参结构体
// @Success 200 {} json "{"success":true,"data":{},"msg":"null"}"
// @Router /record/submit [post]
func (r *RecordApi) Submit(c *gin.Context) {
	//var gzlRecord modelWF.GzlRecord
	//_ = c.ShouldBindJSON(&gzlRecord)

}

// Data
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags Record
// @Summary 返回之前填写过的表单数据
// @Produce  application/json
// @Param id query WorkFlowReq.Record true "记录id"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /record/data [get]
func (r *RecordApi) Data(c *gin.Context) {
	var record WorkFlowReq.Record
	_ = c.ShouldBindJSON(&record)
	data, err := recordService.GetData(record)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("表单数据获取失败", zap.Any("err", err))
		response.FailWithMessage("该记录不存在", c)
	} else {
		global.GSD_LOG.ZapLog.Info("表单数据获取成功", zap.Any("Record Form Data([]byte -> string)", string(data)))
		response.OkWithData(data, c)
	}
}

// Launch
// @Tags record
// @Summary 我发起的
// @Produce  application/json
// @Param data body uint true "创建人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"查询我发起的任务成功"}"
// @Router /record/schedule [get]
func (r *RecordApi) Launch(c *gin.Context) {

}
