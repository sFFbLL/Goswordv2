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
<<<<<<< Updated upstream
		userRouter.POST("register", baseApi.Register)       // 用户注册账号
		userRouter.POST("lists", baseApi.GetUserList)       //用户分页列表
		userRouter.POST("password", baseApi.UpdatePassword) //用户修改密码
		userRouter.GET("infos", baseApi.GetUserInfo)        //获取用户信息
		userRouter.PUT("infos", baseApi.SetUserInfo)        //修改用户信息
		
=======
		userRouter.POST("register", baseApi.Register)                     // 用户注册账号
		userRouter.POST("lists", baseApi.GetUserList)                     //用户分页列表
		userRouter.DELETE("deleteUser", baseApi.DeleteUser)               //用户删除
		userRouter.POST("setUserAuthorities", baseApi.SetUserAuthorities) //设置用户角色
>>>>>>> Stashed changes
	}
}
