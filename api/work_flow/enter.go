package work_flow

import "project/service"

type ApiGroup struct {
	TaskApi
	RecordApi
	AppApi
}

var taskService = service.ServiceGroupApp.WorkFlowServiceGroup.TaskService
var recordService = service.ServiceGroupApp.WorkFlowServiceGroup.RecordService
var appService = service.ServiceGroupApp.WorkFlowServiceGroup.AppService
