package work_flow

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"project/global"
	"project/model/common/response"
	WorkFlowReq "project/model/work_flow/request"
	"project/utils"
	"strconv"
)

type AppApi struct {
}

// Empty
// @author: [tanshaokang](https://github.com/worryfreet)
// @Tags App
// @Summary 返回空表单
// @Produce  application/json
// @Param appId query int true "应用id"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"null"}"
// @Router /app/empty [get]
func (f *AppApi) Empty(c *gin.Context) {
	appId, _ := strconv.Atoi(c.Query("appId"))
	app := WorkFlowReq.EmptyApp{AppId: appId}
	if err := utils.Verify(app, utils.EmptyAppIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := appService.GetAppEmpty(appId)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("获取空应用表单失败", zap.Any("err", err))
		response.FailWithMessage("该应用不存在", c)
	} else {
		global.GSD_LOG.ZapLog.Info("获取空应用表单成功", zap.Any("GetAppEmpty Success", string(data)))
		response.OkWithData(data, c)
	}
}

// Create
// @Tags App
// @Summary 创建表单
// @Produce  application/json
// @Param data body uint true "创建人"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"创建表单成功"}"
// @Router /app/create [post]
func (f *AppApi) Create(c *gin.Context) {

}