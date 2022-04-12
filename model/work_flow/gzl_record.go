package work_flow

import (
	"project/global"
	"project/model/system"
)

//GzlRecord 记录表
type GzlRecord struct {
	global.GSD_MODEL
	AppId        uint           `json:"appId" gorm:"not null;comment:应用id"`            //应用id
	App          GzlApp         `json:"app" gorm:"foreignKey:AppId;references:ID;"`    //应用
	CurrentState uint8          `json:"currentState" gorm:"not null;comment:当前状态"`     //当前状态
	Form         JSON           `json:"form" gorm:"not null;type:json;comment:表单JSON"` //表单JSON
	CurrentNode  string         `json:"currentNode" gorm:"not null;comment:当前节点"`      //当前节点
	DeptId       uint           `json:"depId" gorm:"not null;comment:部门id"`            //部门id
	Dept         system.SysDept `json:"dept" gorm:"foreignKey:DeptId;references:ID;"`  //部门
}
