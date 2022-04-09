package system

type SysAuthorityDept struct {
	SysDeptId               uint   `gorm:"column:sys_dept_id"`
	SysAuthorityAuthorityId string `gorm:"column:sys_authority_authority_id"`
}
