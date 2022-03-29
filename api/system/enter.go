package system

import "project/service"

type ApiGroup struct {
	BaseApi
	JwtApi
}

var jwtService = service.ServiceGroupApp.SystemServiceGroup.JwtService
var userService = service.ServiceGroupApp.SystemServiceGroup.UserService
