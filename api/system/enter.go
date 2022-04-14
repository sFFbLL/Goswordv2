package system

import (
	"project/service"
	"project/utils"
)

type ApiGroup struct {
	BaseApi
	JwtApi
}

var dataScope = utils.DataScope{}

var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
var userService = service.ServiceGroupApp.SystemServiceGroup.UserService
var authorityService = service.ServiceGroupApp.SystemServiceGroup.AuthorityService
