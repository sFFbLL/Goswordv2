package work_flow

import (
	"github.com/gin-gonic/gin"
	v1 "project/api"
)

type RecordRouter struct {
}

func (t *RecordRouter) InitRecordRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	recordRouter := Router.Group("record")
	var recordApi = v1.ApiGroupApp.RecordApiGroup.RecordApi
	{
		recordRouter.POST("submit", recordApi.Submit)
		recordRouter.GET("data", recordApi.Data)
		recordRouter.GET("initiated", recordApi.MyInitiated)
	}
	return recordRouter
}
