package work_flow

//GzlAppUser 应用用户中间表
type GzlAppUser struct {
	AppId  uint `gorm:"column:gzl_app_id;comment:应用id"`
	UserId uint `gorm:"column:sys_user_id;comment:用户id"`
}

func (g *GzlAppUser) TableName() string {
	return "gzl_app_user"
}
