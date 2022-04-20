package request

import (
	"project/utils"
	"time"
)

type RecordById struct {
	RecordId int `form:"recordId"`
}

type RecordSubmit struct {
	CreateBy    uint        //创建人
	AppId       uint        `json:"appId"`
	Form        interface{} `json:"form"`
	CurrentNode string      `json:"currentNode" gorm:"comment:当前节点"`
	DeptId      uint        `json:"deptId" gorm:"comment:部门id"`
}

type FormItem struct {
	CreateBy      uint       //创建人
	CreatedAt     time.Time  // 创建时间
	RecordId      uint       `json:"recordId"`      //记录id
	Name          string     `json:"name"`          //名称
	DataType      string     `json:"type"`          //数据类型
	ComponentType string     `json:"componentType"` //组件类型
	Content       string     `json:"content"`       //内容
	IsRequired    uint8      `json:"isRequired" `   //是否必填(必填1默认、不必填2)
	Form          utils.JSON `json:"form"`          //表单
	DeptId        uint       `json:"deptId"`        //部门id
}
