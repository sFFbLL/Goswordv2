package work_flow

import (
	"errors"
	"gorm.io/gorm"
	"project/global"
	"project/model/work_flow"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	WorkFlowRes "project/model/work_flow/response"
	"time"
)

type TaskService struct {
}

// GetDynamic
// @author: [tanshaokang](https://github.com/worryfreet)
// @function: GetDynamic
// @description: 从mysql中获取流程动态数据
// @param: WorkFlowReq.RecordById
// @return: data []WorkFlowReq.Dynamic, err error
func (t TaskService) GetDynamic(applicantId, recordId uint) (data []WorkFlowRes.Dynamic, err error) {
	db := global.GSD_DB.
		Model(modelWF.GzlTask{}).
		Joins("JOIN sys_users ON sys_users.id = ?", applicantId).
		Joins("JOIN gzl_records ON gzl_records.id = ?", recordId).
		Select("sys_users.username as Applicant", "gzl_tasks.created_at as InspectAt",
			"gzl_records.created_at as CreatedAt", "check_state as CheckState", "remarks as Remarks").
		Where("gzl_tasks.record_id = gzl_records.id")
	if err = db.Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	// 计算耗时, 格式化时间
	for i := 0; i < len(data); i++ {
		born := "2006-01-02 15:04:05"
		beginStr := data[i].CreatedAt.Format(born)
		endStr := data[i].InspectAt.Format(born)
		begin, _ := time.ParseInLocation(born, beginStr, time.Local)
		end, _ := time.ParseInLocation(born, endStr, time.Local)
		data[i].ConsumeTime = end.Unix() - begin.Unix()
		data[i].CreatedAtFormatStr = beginStr[:len(beginStr)-3]
		data[i].InspectAtFormatStr = endStr[:len(endStr)-3]
	}
	return
}

// GetScheduleList
// @author: [zhaozijie](https://github.com/worryfreet)
// @function: GetScheduleList
// @description: 从mysql中获取待办数据
// @param: WorkFlowReq.Task
// @return: data []WorkFlowReq.Schedule, err error
func (t *TaskService) GetScheduleList(userId, appid int) (err error, tasks []WorkFlowReq.Function) {
	db := global.GSD_DB.Model(&work_flow.GzlTask{}).
		Joins("JOIN sys_users ON sys_users.id = ?", userId).
		Joins("JOIN gzl_apps ON gzl_apps.id = ?", appid). //连表查询
		Select("sys_users.username as Applicant", "gzl_tasks.created_at as CreatedAt",
			"gzl_apps.name as AppName", "check_state as CheckState").
		Where("gzl_tasks.inspector=Inspector")
	if err = db.Find(&tasks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //如果待办为空，返回空
			return nil, nil
		} else {
			return
		}
	}
	return
}

func (t *TaskService) GetHandleList(userId int, appid int) (err error, tasks []WorkFlowReq.Function) {
	db := global.GSD_DB.Model(&work_flow.GzlTask{}).
		Joins("JOIN sys_users ON sys_users.id = ?", userId).
		Joins("JOIN gzl_apps ON gzl_apps.id = ?", appid). //连表查询
		Select("sys_users.username as Applicant", "gzl_tasks.created_at as CreatedAt",
			"gzl_tasks.inspector as Inspector", "gzl_apps.name as AppName", "check_state as CheckState").
		Where("gzl_tasks.inspector=Inspector")
	if err = db.Find(&tasks).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //如果待办为空，返回空
			return nil, nil
		} else {
			return
		}
	}
	return
}

func (t *TaskService) Inspect(task work_flow.GzlTask) error {
	//获取任务详细信息
	return global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		taskInfo, err := t.GetTaskInfo(task.ID)
		if err != nil {
			return err
		}
		//流程流转
		err = ProcessFlow(taskInfo.Record)
		if err != nil {
			return err
		}
		return tx.Updates(&task).Error
	})
}

// GetReceive
// @author: [tanshaokang](https://github.com/worryfreet)
// @function: GetReceive
// @description: 从mysql中获取我收到的信息列表
// @param: WorkFlowReq.RecordById
// @return: data []WorkFlowReq.Dynamic, err error
func (t TaskService) GetReceive(userId uint) (err error, tasks []modelWF.GzlTask) {
	// 1. 申请人姓名  userId -> sys_users.username
	// 2. 审批人姓名  userId ->
	// 3. 审批状态
	// 4. 应用名称
	// 5. 当前节点
	return
}

// GetTaskInfo 根据id获取信息
func (t TaskService) GetTaskInfo(taskId uint) (task work_flow.GzlTask, err error) {
	err = global.GSD_DB.Preload("RecordById.EmptyApp").First(&task, taskId).Error
	return task, err
}
