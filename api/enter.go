package v1

import (
	"project/api/system"
	"project/api/work_flow"
)

type ApiGroup struct {
	SystemApiGroup system.ApiGroup
	TaskApiGroup   work_flow.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
