package request

// RecordById 通过记录id找到record记录
type RecordById struct {
	RecordId uint `form:"recordId"`
}

// RecordSubmit 提交表单(记录)时的入参结构体
type RecordSubmit struct {
	AppId uint        `json:"appId"`
	Form  interface{} `json:"form"`
}

type Form struct {
	Fields []FormItem `json:"fields"`
}

type FormItem struct {
	Config `json:"__config__"`
}

type Config struct {
	Name          string `json:"label"`
	DefaultValue  string `json:"defaultValue"`
	ComponentType string `json:"tagIcon"`
	IsRequired    bool   `json:"required"`
}
