package system

import (
	v1 "project/api"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("register", baseApi.Register) // 用户注册账号
	}
}
