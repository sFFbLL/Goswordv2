package request

// RecordById 通过记录id找到record记录
type RecordById struct {
	RecordId uint `form:"recordId"`
}

// RecordSubmit 提交表单(记录)时的入参结构体
type RecordSubmit struct {
	AppId uint `json:"appId"`
	Form  Form `json:"form"`
}

// Form 表单结构体
type Form struct {
	Fields      []FormItem  `json:"fields"`
	CreateBy    uint        //创建人
	AppId       uint        `json:"appId"`
	Form        interface{} `json:"form"`
	CurrentNode string      `json:"currentNode" gorm:"comment:当前节点"`
	DeptId      uint        `json:"deptId" gorm:"comment:部门id"`
}

type FormItem struct {
	Config `json:"__config__"`
}

type Config struct {
	Name          string `json:"label"`
	RenderKey     string `json:"renderKey"`
	DefaultValue  string `json:"defaultValue"`
	ComponentType string `json:"tagIcon"`
	IsRequired    uint8  `json:"required"`
}
