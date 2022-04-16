package router

import (
	"project/router/system"
	"project/router/work_flow"
)

type RouterGroup struct {

	System system.RouterGroup
	WorkFlow work_flow.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
