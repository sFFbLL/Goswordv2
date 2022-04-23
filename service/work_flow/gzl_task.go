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
func (t TaskService) GetDynamic(applicantId, recordId int) (data []WorkFlowRes.Dynamic, err error) {
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
// @return: data []WorkFlowReq.ScheduleList, err error
func (t TaskService) GetScheduleList(userId uint) (scheduleData []WorkFlowRes.ScheduleList,err error) {
	var recordIds []uint
	global.GSD_DB.Model(&modelWF.GzlTask{}).Select("record_id").
		Where("node_type = ? AND Inspector = ?", 3, 1).
		Find(&recordIds)

	for i := 0; i < len(recordIds); i++ {
		var tasks modelWF.GzlTask
		global.GSD_DB.
			Where("record_id = ?", recordIds[i]).
			Preload("Record.App").
			Find(&tasks)
		var schedule WorkFlowRes.ScheduleList
		// 申请人姓名
		global.GSD_DB.Model(&system.SysUser{}).
			Where("id = ?", tasks.Record.CreateBy).
			Select("nick_name as Applicant").
			Find(&schedule.Applicant)
		// 应用名称
		schedule.AppName = tasks.Record.App.Name
		//审批状态
		schedule.CheckState = tasks.CheckState
		// 创建时间
		schedule.CreatedAt = tasks.Record.CreatedAt
		// 当前节点名称 (解析Flow)
		var flow Flow
		_ = json.Unmarshal(tasks.Record.App.Flow, &flow)
		for _, node := range flow.FlowElementList {
			if node.Key == tasks.Record.CurrentNode {
				schedule.CurrentNode = node.Name
				break
			}
		}
		scheduleData = append(scheduleData, schedule)
	}
	return
}

// GetHandleList
// @author: [zhaozijie](https://github.com/worryfreet)
// @function: GetHandleList
// @description: 从mysql中获取我处理的数据
// @param: WorkFlowReq.Task
// @return: data []WorkFlowReq.HandleList, err error
func (t TaskService) GetHandleList(userId uint) (handleData []WorkFlowRes.HandleList, err error) {
	var recordIds []uint
	global.GSD_DB.Model(&modelWF.GzlTask{}).Select("record_id").
		Where("node_type = ? AND Inspector = ?", 3, 1).
		Find(&recordIds)

	for i := 0; i < len(recordIds); i++ {
		var tasks []modelWF.GzlTask
		global.GSD_DB.
			Where("record_id = ? AND node_type = ?", recordIds[i], 3).
			Preload("Record.App").
			Find(&tasks)
		if len(tasks) > 0 {
			var handle WorkFlowRes.HandleList
			// 申请人姓名
			global.GSD_DB.Model(&system.SysUser{}).
				Where("id = ?", tasks[0].Record.CreateBy).
				Select("nick_name as Applicant").
				Find(&handle.Applicant)
			// 应用名称
			handle.AppName = tasks[0].Record.App.Name
			// 当前状态
			handle.CurrentState = tasks[0].Record.CurrentState
			// 获取审批人姓名(可能会有多个, 所以需要遍历)
			for j := 0; j < len(tasks); j++ {
				var Inspector string
				global.GSD_DB.Model(&system.SysUser{}).
					Where("id = ?", tasks[j].Inspector).
					Select("nick_name as Inspector").
					Find(&Inspector)
				if Inspector != "" {
					handle.Inspectors = append(handle.Inspectors, Inspector)
				}
			}
			// 当前节点名称 (解析Flow)
			var flow Flow
			_ = json.Unmarshal(tasks[0].Record.App.Flow, &flow)
			for _, node := range flow.FlowElementList {
				if node.Key == tasks[0].Record.CurrentNode {
					handle.CurrentNode = node.Name
				}
			}
			handleData = append(handleData, handle)
		}
	}
	fmt.Println(err)
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
// @author: [tanshaokang](https://github.com/worryfreet)
// @function: GetReceive
// @description: 从mysql中获取我收到的信息列表
// @param: WorkFlowReq.RecordById
// @return: data []WorkFlowReq.Dynamic, err error
func (t TaskService) GetReceive(userId int) (err error, tasks []modelWF.GzlTask) {
	// 1. 申请人姓名  userId -> sys_users.username
	// 2. 审批人姓名  userId ->
	// 3. 审批状态
	// 4. 应用名称
	// 5. 当前节点
	return
}

// GetTaskInfo 根据id获取详细信息
func (t TaskService) GetTaskInfo(taskId uint) (task work_flow.GzlTask, err error) {
	err = global.GSD_DB.Preload("RecordById.EmptyApp").First(&task, taskId).Error
	return task, err
}
