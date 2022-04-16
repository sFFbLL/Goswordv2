package work_flow

import (
	"errors"
	"gorm.io/gorm"

	"project/global"
	"project/model/work_flow"
)

type TaskService struct {
}

func (taskService *TaskService) GetScheduleList(InspectorId int) (err error, tasks []work_flow.GzlTask) {
	db := global.GSD_DB.Model(&work_flow.GzlTask{}) //查表GzlTask
	if err = db.Find(&tasks, "inspector = ?", InspectorId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //如果待办为空，返回空
			return nil, nil
		} else {
			return
		}
	}
	return
}

//func (taskService *TaskService) GetHandleList(InspectorId int) (err error, tasks []work_flow.GzlTask) {
//	db := global.GSD_DB.Model(&work_flow.GzlTask{}) //查表GzlTask
//	if err = db.Find(&tasks, "inspector = ?", InspectorId).Error; err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) { //如果待办为空，返回空
//			return nil, nil
//		} else {
//			return
//		}
//	}
//	return
//}

func (taskService *TaskService) Inspect(task work_flow.GzlTask) (err error) {
	return global.GSD_DB.Updates(&task).Error
}
