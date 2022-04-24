package service

import (
	"project/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup system.SysGroup
}

var ServiceGroupApp = new(ServiceGroup)
