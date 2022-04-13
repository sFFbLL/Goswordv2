package work_flow

import "project/service"

type ApiGroup struct {
	TaskApi
	FormApi
}

var taskService = service.ServiceGroupApp.WorkFlowServiceGroup.TaskService
