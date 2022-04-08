package system

import "project/global"

type SysAuthority struct {
	global.GSD_MODEL
	AuthorityId   uint          `json:"authorityId" gorm:"not null;primary_key;comment:角色ID;size:90"` // 角色ID
	Level         uint          `json:"level" gorm:"comment:角色等级0最大"`
	AuthorityName string        `json:"authorityName" gorm:"comment:角色名"` // 角色名
	DataScope     string        `json:"dataScope"`
	SysBaseMenus  []SysBaseMenu `json:"menus" gorm:"many2many:sys_authority_menus;"`
}
