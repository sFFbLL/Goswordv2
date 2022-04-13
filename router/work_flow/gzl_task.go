package work_flow

import (
	"github.com/gin-gonic/gin"
	v1 "project/api"
)

type TaskRouter struct {
}

func (t *TaskRouter) InitTaskRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	taskRouter := Router.Group("task")
	var taskApi = v1.ApiGroupApp.TaskApiGroup.TaskApi
	{
		taskRouter.POST("inspect", taskApi.Inspect)
		taskRouter.GET("dynamic", taskApi.Dynamic)
	}
	return taskRouter
}
