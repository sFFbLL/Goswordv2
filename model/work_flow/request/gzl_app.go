package request

type EmptyApp struct {
	AppId int `form:"appId"`
}

type AddApp struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}
