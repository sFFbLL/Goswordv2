package system

import (
	"project/global"

	uuid "github.com/satori/go.uuid"
)

type SysUser struct {
	global.GSD_MODEL
	UUID        uuid.UUID      `json:"uuid" gorm:"not null;comment:用户UUID"`
	Username    string         `json:"userName" gorm:"not null;uniqueIndex;comment:用户登录名"`                                          // 用户登录名
	Password    string         `json:"-"  gorm:"not null;comment:用户登录密码"`                                                           // 用户登录密码
	NickName    string         `json:"nickName" gorm:"not null;default:系统用户;comment:用户昵称"`                                          // 用户昵称
	HeaderImg   string         `json:"headerImg" gorm:"not null;default:http://r9qsta3s9.hn-bkt.clouddn.com/head.png;comment:用户头像"` // 用户头像
	Email       string         `json:"email" gorm:"not null;comment:用户邮箱"`
	Phone       string         `json:"phone" gorm:"comment:用户手机号"`
	DeptId      uint           `json:"deptId" gorm:"comment:用户所属部门"`
	Dept        SysDept        `json:"dept" gorm:"foreignKey:DeptId;references:ID;"`
	Authorities []SysAuthority `json:"authorities" gorm:"many2many:sys_user_authority;"`
}
