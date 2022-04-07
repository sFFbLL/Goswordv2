package system

import (
	"project/global"

	uuid "github.com/satori/go.uuid"
)

type SysUser struct {
	global.GSD_MODEL
	UUID        uuid.UUID      `json:"uuid" gorm:"comment:用户UUID"`
	Username    string         `json:"userName" gorm:"uniqueIndex;comment:用户登录名"`                                                 // 用户登录名
	Password    string         `json:"-"  gorm:"comment:用户登录密码"`                                                                  // 用户登录密码
	NickName    string         `json:"nickName" gorm:"not null;default:系统用户;comment:用户昵称"`                                        // 用户昵称
	HeaderImg   string         `json:"headerImg" gorm:"default:http://http://r9qsta3s9.hn-bkt.clouddn.com/head.png;comment:用户头像"` // 用户头像
	DeptId      uint           `json:"deptId" gorm:""`
	Dept        SysDept        `json:"dept" gorm:"foreignKey:DeptId;references:ID;"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
}
