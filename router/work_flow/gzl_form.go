package work_flow

import (
	"github.com/gin-gonic/gin"
	v1 "project/api"
)

type FormRouter struct {
}

func (t *TaskRouter) InitFormRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	formRouter := Router.Group("form")
	var formApi = v1.ApiGroupApp.FormApiGroup.FormApi
	{
		formRouter.POST("formempty", formApi.FormEmpty)
		formRouter.GET("formsave", formApi.FormSave)
		formRouter.GET("formdata", formApi.FormData)
	}
	return formRouter
}
