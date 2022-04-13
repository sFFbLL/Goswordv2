package main

import (
	"project/core"
	"project/global"
	"project/initialize"
)

//go:generate go env -w GO111MODULE=on
//go:generate go env -w GOPROXY=https://goproxy.cn,direct
//go:generate go mod tidy
//go:generate go mod download

// @title                       Work Flow API
// @version                     0.0.1
// @description                 This is a sample Server pets
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        x-token
// @BasePath                    /
func main() {
	global.GSD_VP = core.Viper()      // 初始化Viper
	global.GSD_LOG = core.Zap()       // 初始化zap日志库
	global.GSD_DB = initialize.Gorm() // gorm连接数据库
	if global.GSD_DB != nil {
		initialize.MysqlTables(global.GSD_DB) // 初始化表
		initialize.InitDB()                   // 初始化表数据
		// 程序结束前关闭数据库链接
		db, _ := global.GSD_DB.DB()
		defer db.Close()
	}
	core.RunWindowsServer()
}
