package request

type EmptyApp struct {
	AppId int `form:"AppId"`
}

type AddApp struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}
