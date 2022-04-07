package system

type SysAuthority struct {
	AuthorityId   uint          `json:"authorityId" gorm:"not null;primary_key;comment:角色ID;size:90"` // 角色ID
	AuthorityName string        `json:"authorityName" gorm:"comment:角色名"`                             // 角色名
	SysBaseMenus  []SysBaseMenu `json:"menus" gorm:"many2many:sys_authority_menus;"`
	DefaultRouter string        `json:"defaultRouter" gorm:"comment:默认菜单;default:dashboard"` // 默认菜单(默认dashboard)
}
