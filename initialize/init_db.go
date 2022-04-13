package initialize

import (
	"os"
	"project/global"
	"project/model/system"
	"project/source"

	"go.uber.org/zap"
)

func InitDB() {
	err := initDB(
		source.Admin,
		source.BaseMenu,
		source.Authority,
		source.AuthoritiesMenus,
		source.AuthorityMenu,
		source.Dept,
		//source.Api,
		//source.Casbin,
	)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("init table data failed", zap.Any("err", err))
		os.Exit(0)
	}
	global.GSD_LOG.ZapLog.Info("init table data success")
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: initDB
//@description: 初始化超级管理员数据
//@param: InitDBFunctions system.InitDBFunc 初始化方法
//@return: error
func initDB(InitDBFunctions ...system.InitDBFunc) (err error) {
	for _, v := range InitDBFunctions {
		err = v.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
