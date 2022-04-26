package router

import (
	"project/router/system"
	WorkFlow "project/router/work_flow"
)

type RouterGroup struct {
	System   system.RouterGroup
	WorkFlow WorkFlow.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
