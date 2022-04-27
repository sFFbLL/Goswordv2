package work_flow

//GzlAppAuthority 应用角色中间表
type GzlAppAuthority struct {
	AppId       uint `gorm:"column:gzl_app_id;comment:应用id"`
	AuthorityId uint `gorm:"column:sys_authority_authority_id;comment:角色id"`
}

func (g *GzlAppAuthority) TableName() string {
	return "gzl_app_authority"
}
