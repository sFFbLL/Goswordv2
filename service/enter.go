package service

import (
	"project/service/system"
	"project/service/work_flow"
)

type ServiceGroup struct {
	SystemServiceGroup   system.SysGroup
	WorkFlowServiceGroup work_flow.GzlGroup
}

var ServiceGroupApp = new(ServiceGroup)
