package work_flow

import (
	"project/global"
	"project/model/work_flow"
)

type TaskService struct {
}

func (taskService *TaskService) Inspect(task work_flow.GzlTask) (err error) {
	return global.GSD_DB.Updates(&task).Error
}
