package work_flow

import (
	"github.com/gin-gonic/gin"
	v1 "project/api"
)

type AppRouter struct {
}

func (t *AppRouter) InitAppRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	appRouter := Router.Group("app")
	var appApi = v1.ApiGroupApp.AppApiGroup.AppApi
	{
		appRouter.GET("/empty", appApi.Empty)
		appRouter.POST("/create", appApi.Create)
		appRouter.POST("/", appApi.AddApp)
	}
	return appRouter
}
