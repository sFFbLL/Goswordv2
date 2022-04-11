package work_flow

//GzlAppDept 应用部门中间表
type GzlAppDept struct {
	AppId  uint `gorm:"column:gzl_app_id;comment:应用id"`
	DeptId uint `gorm:"column:sys_dept_id;comment:部门id"`
}

func (g *GzlAppDept) TableName() string {
	return "gzl_app_dept"
}
