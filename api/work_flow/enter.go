package work_flow

import "project/service"

type ApiGroup struct {
	TaskApi
}

var taskService = service.ServiceGroupApp.WorkFlowServiceGroup.TaskService
