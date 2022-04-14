package system

import "project/global"

type SysDept struct {
	global.GSD_MODEL
	PartenID uint   `json:"partenID" gorm:"not null;default:0"`
	DeptName string `json:"deptName" gorm:"not null;default:''"`
	DeptSort uint   `json:"deptSort" gorm:"not null;default:1"`
}
