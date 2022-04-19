package system

import (
	v1 "project/api"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (s *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	var baseApi = v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		userRouter.POST("register", baseApi.Register)                     // 用户注册账号
		userRouter.POST("getUserList", baseApi.GetUserList)               //用户分页列表
		userRouter.PUT("changePassword", baseApi.UpdatePassword)          //用户本人修改密码
		userRouter.GET("getUserInfo", baseApi.GetUserInfo)                //获取用户信息
		userRouter.PUT("setUserInfo", baseApi.SetUserInfo)                //修改用户信息
		userRouter.DELETE("deleteUser", baseApi.DeleteUser)               //用户删除
		userRouter.POST("setUserAuthorities", baseApi.SetUserAuthorities) //设置用户角色
		userRouter.POST("importExcel", baseApi.ImportExcel)               //导入用户信息
		userRouter.GET("loadExcel", baseApi.LoadExcel)                    //加载excel数据
		userRouter.POST("exportExcel", baseApi.ExportExcel)               //导出用户数据
	}
}
