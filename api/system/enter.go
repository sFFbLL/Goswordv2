package system

import (
	"project/service"
	"project/utils"
)

type ApiGroup struct {
	BaseApi
	JwtApi
	AuthorityMenuApi
	AuthorityApi
}

var dataScope = utils.DataScope{}

var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
var userService = service.ServiceGroupApp.SystemServiceGroup.UserService
var authorityService = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
var menuService = service.ServiceGroupApp.SystemServiceGroup.MenuService
