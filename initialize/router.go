package initialize

import (
	"net/http"
	_ "project/docs"
	"project/global"
	"project/middleware"
	"project/router"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// 初始化总路由

func Routers() *gin.Engine {
	var Router = gin.New()
	Router.Use(middleware.GinLogger(), middleware.GinRecovery(true))
	Router.StaticFS(global.GSD_CONFIG.Local.Path, http.Dir(global.GSD_CONFIG.Local.Path)) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了
	global.GSD_LOG.ZapLog.Info("use middleware logger")
	// 跨域
	Router.Use(middleware.Cors()) // 如需跨域可以打开
	global.GSD_LOG.ZapLog.Info("use middleware cors")
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GSD_LOG.ZapLog.Info("register swagger handler")
	// 方便统一添加路由组前缀 多服务器上线使用

	//获取路由组实例
	systemRouter := router.RouterGroupApp.System
	PublicGroup := Router.Group("")
	{
		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
	}
	PrivateGroup := Router.Group("")
	PrivateGroup.Use(middleware.JWTAuth())
	{
		//systemRouter.InitApiRouter(PrivateGroup)                    // 注册功能api路由
		systemRouter.InitJwtRouter(PrivateGroup)  // jwt相关路由
		systemRouter.InitUserRouter(PrivateGroup) // 注册用户路由
		//systemRouter.InitMenuRouter(PrivateGroup)                   // 注册menu路由
		//systemRouter.InitEmailRouter(PrivateGroup)                  // 邮件相关路由
		//systemRouter.InitSystemRouter(PrivateGroup)                 // system相关路由
		//systemRouter.InitCasbinRouter(PrivateGroup)                 // 权限相关路由
		//systemRouter.InitAutoCodeRouter(PrivateGroup)               // 创建自动化代码
		//systemRouter.InitAuthorityRouter(PrivateGroup)              // 注册角色路由
		//systemRouter.InitSysDictionaryRouter(PrivateGroup)          // 字典管理
		//systemRouter.InitSysOperationRecordRouter(PrivateGroup)     // 操作记录
		//systemRouter.InitSysDictionaryDetailRouter(PrivateGroup)    // 字典详情管理
	}
	global.GSD_LOG.ZapLog.Info("router register success")
	return Router
}
