package system

import (
	"project/global"

	uuid "github.com/satori/go.uuid"
)

type SysUser struct {
	global.GSD_MODEL
	UUID        uuid.UUID      `json:"uuid" gorm:"comment:用户UUID"`                                                    // 用户UUID
	Username    string         `json:"userName" gorm:"comment:用户登录名"`                                                 // 用户登录名
	Password    string         `json:"-"  gorm:"comment:用户登录密码"`                                                      // 用户登录密码
	NickName    string         `json:"nickName" gorm:"default:系统用户;comment:用户昵称"`                                     // 用户昵称
	HeaderImg   string         `json:"headerImg" gorm:"default:http://qmplusimg.henrongyi.top/head.png;comment:用户头像"` // 用户头像
	ActiveColor string         `json:"activeColor" gorm:"default:#1890ff;comment:用户角色ID"`                             // 活跃颜色
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
}
