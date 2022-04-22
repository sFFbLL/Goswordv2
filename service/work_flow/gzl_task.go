package work_flow

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"project/global"
	"project/model/system"
	"project/model/work_flow"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	WorkFlowRes "project/model/work_flow/response"
)

type TaskService struct {
}

// GetDynamic
// @author: [tanshaokang]
// @function: GetDynamic
// @description: 从mysql中获取流程动态数据
// @param: WorkFlowReq.RecordById
// @return: data []WorkFlowReq.Dynamic, err error
func (t TaskService) GetDynamic(applicantId, recordId uint) (data []WorkFlowRes.Dynamic, err error) {
	db := global.GSD_DB.
		Model(modelWF.GzlTask{}).
		Joins("JOIN sys_users ON sys_users.id = ?", applicantId).
		Joins("JOIN gzl_records ON gzl_records.id = ?", recordId).
		Select("sys_users.nick_name as Applicant", "gzl_tasks.created_at as InspectAt",
			"gzl_records.created_at as CreatedAt", "check_state as CheckState", "remarks as Remarks").
		Where("gzl_tasks.record_id = gzl_records.id")
	if err = db.Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}
	// 计算耗时
	for i := 0; i < len(data); i++ {
		data[i].ConsumeTime = data[i].InspectAt.Unix() - data[i].CreatedAt.Unix()
	}
	return
}

// GetScheduleList
// @author: [zhaozijie]
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
		err = tx.Updates(&task).Error
		if err != nil {
			return err
		}
		//流程流转
		return ProcessFlow(taskInfo.Record)
	})
}

// GetReceive
// @author: [tanshaokang]
// @function: GetReceive
// @description: 从mysql中获取我收到的信息列表
// @param: WorkFlowReq.RecordById
// @return: data []WorkFlowRes.Receive, err error
func (t TaskService) GetReceive(userId uint) (data []WorkFlowRes.Receive, err error) {
	var recordIds []uint
	global.GSD_DB.Model(&modelWF.GzlTask{}).Select("record_id").
		Where("node_type = ? AND Inspector = ?", 4, userId).
		Find(&recordIds)

	for i := 0; i < len(recordIds); i++ {
		var tasks []modelWF.GzlTask
		global.GSD_DB.
			Where("record_id = ? AND node_type = ?", recordIds[i], 3).
			Preload("Record.App").
			Find(&tasks)
		if len(tasks) > 0 {
			var receive WorkFlowRes.Receive
			// 记录id
			receive.RecordId = tasks[0].RecordId
			// 申请人姓名
			global.GSD_DB.Model(&system.SysUser{}).
				Where("id = ?", tasks[0].Record.CreateBy).
				Select("nick_name as Applicant").
				Find(&receive.Applicant)
			// 应用名称
			receive.AppName = tasks[0].Record.App.Name
			// 当前状态
			receive.CurrentState = tasks[0].Record.CurrentState
			// 获取审批人姓名(可能会有多个, 所以需要遍历)
			for j := 0; j < len(tasks); j++ {
				var Inspector string
				global.GSD_DB.Model(&system.SysUser{}).
					Where("id = ?", tasks[j].Inspector).
					Select("nick_name as Inspector").
					Find(&Inspector)
				if Inspector != "" {
					receive.Inspectors = append(receive.Inspectors, Inspector)
				}
			}
			// 当前节点名称 (解析Flow)
			var flow Flow
			_ = json.Unmarshal(tasks[0].Record.App.Flow, &flow)
			for _, node := range flow.FlowElementList {
				if node.Key == tasks[0].Record.CurrentNode {
					receive.CurrentNode = node.Name
				}
			}
			data = append(data, receive)
		}
	}
	fmt.Println(err)
	return
}

// GetTaskInfo 根据id获取详细信息
func (t TaskService) GetTaskInfo(taskId uint) (task work_flow.GzlTask, err error) {
	err = global.GSD_DB.Preload("Record.App").First(&task, taskId).Error
	return task, err
}
