package source

import (
	"project/global"
	"project/model/system"
	"time"

	"github.com/gookit/color"

	"gorm.io/gorm"
)

var Authority = new(authority)

type authority struct{}

var authorities = []system.SysAuthority{
	{CreatedAt: time.Now(), UpdatedAt: time.Now(), AuthorityId: 1, AuthorityName: "超级管理员", DataScope: "全部", Level: 1},
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@description: sys_authorities 表数据初始化
func (a *authority) Init() error {
	return global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("authority_id IN ? ", []uint{1}).Find(&[]system.SysAuthority{}).RowsAffected == 1 {
			color.Danger.Println("\n[Mysql] --> sys_authorities 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&authorities).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> sys_authorities 表初始数据成功!")
		return nil
	})
}
