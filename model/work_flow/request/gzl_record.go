package request

type Record struct {
	RecordId int `form:"recordId"`
}

type Schedule struct {
	AppId int `form:"appId"`
}