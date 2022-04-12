package work_flow

import (
	"project/global"
	"project/model/system"
)

//GzlFormItem 表单组件表
type GzlFormItem struct {
	global.GSD_MODEL
	RecordId      uint           `json:"recordId" gorm:"not null;comment:记录id"`                         //记录id
	Record        GzlRecord      `json:"record" gorm:"foreignKey:RecordId;references:ID;"`              //记录
	Name          string         `json:"name" gorm:"not null;comment:名称"`                               //名称
	DataType      string         `json:"type" gorm:"not null;comment:数据类型"`                             //数据类型
	ComponentType string         `json:"componentType" gorm:"not null;comment:组件类型"`                    //组件类型
	Content       string         `json:"content" gorm:"not null;comment:内容"`                            //内容
	IsRequired    uint8          `json:"isRequired" gorm:"not null;default:1;comment:是否必填(必填1默认、不必填2)"` //是否必填(必填1默认、不必填2)
	Form          JSON           `json:"form" gorm:"not null;type:json;comment:表单"`                     //表单
	DeptId        uint           `json:"deptId"  gorm:"not null;comment:部门id"`                          //部门id
	Dept          system.SysDept `json:"dept" gorm:"foreignKey:DeptId;references:ID;"`                  //部门
}
