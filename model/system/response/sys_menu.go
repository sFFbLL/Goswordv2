package response

import "project/model/system"

type SysBaseMenusResponse struct {
	Menus []system.SysBaseMenu `json:"menus"`
}

type SysMenusResponse struct {
	Menus []system.SysMenu `json:"menus"`
}
