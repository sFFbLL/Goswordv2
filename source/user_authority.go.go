package source

import (
	"project/global"
	"project/model/system"
	"strconv"

	"github.com/gookit/color"
	"gorm.io/gorm"
)

var UserAuthority = new(userAuthority)

type userAuthority struct{}

var userAuthorityModel = []system.SysUseAuthority{
	{1, 1},
}

//@description: user_authority 数据初始化
func (a *userAuthority) Init() error {
	return global.GSD_DB.Model(&system.SysUseAuthority{}).Transaction(func(tx *gorm.DB) error {
		if tx.Where("sys_user_id IN (1)").Find(&[]system.SysUseAuthority{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> sys_user_authority 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&userAuthorityModel).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		if success, err := global.GSD_Casbin.AddRoleForUser(strconv.Itoa(1), strconv.Itoa(1)); !success {
			return err
		}
		color.Info.Println("\n[Mysql] --> sys_user_authority 表初始数据成功!")
		return nil
	})
}
