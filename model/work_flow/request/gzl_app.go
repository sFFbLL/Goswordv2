package request

type EmptyApp struct {
	AppId uint `form:"appId"`
}

type AddApp struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type AddForm struct {
	AppId uint        `json:"appId"`
	Form  interface{} `json:"form"`
}

type AddFlow struct {
	AppId uint        `json:"appId"`
	Flow  interface{} `json:"flow"`
}

type EnableApp struct {
	AppId    uint  `json:"appId"`
	IsEnable uint8 `json:"isEnable"`
}

type AuthorityApp struct {
	AppId      uint   `json:"appId"`
	Depts      []uint `json:"depts"`
	Authoritys []uint `json:"authoritys"`
	Users      []uint `json:"users"`
}
