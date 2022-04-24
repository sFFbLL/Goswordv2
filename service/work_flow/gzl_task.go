package work_flow

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"project/global"
	"project/model/system"
	"project/model/work_flow"
	modelWF "project/model/work_flow"
	WorkFlowRes "project/model/work_flow/response"
	"project/utils"
)

type TaskService struct {
}

// GetDynamic
// @author: [tanshaokang]
// @function: GetDynamic
// @description: 从mysql中获取流程动态数据
// @param: WorkFlowReq.RecordById
// @return: data WorkFlowRes.DynamicList, err error
func (t TaskService) GetDynamic(recordId uint) (data WorkFlowRes.DynamicList, err error) {
	var tasks []modelWF.GzlTask
	db := global.GSD_DB.Preload("Record.App").
		Where("record_id = ?", recordId)
	if err = db.First(&tasks).Error; err != nil {
		return WorkFlowRes.DynamicList{}, errors.New("该数据不存在")
	}
	if len(tasks) > 0 {
		data.Nodes = t.GetMoreNodesName(tasks[0].Record.App.Flow, tasks)
		for _, task := range tasks {
			dynamic := WorkFlowRes.Dynamic{
				CreatedAt:   task.CreatedAt,
				InspectAt:   task.UpdatedAt,
				ConsumeTime: task.UpdatedAt.Unix() - task.CreatedAt.Unix(),
				Applicant:   t.GetUserNickName(task.Inspector),
				CheckState:  task.CheckState,
				Remarks:     task.Remarks,
				AppName:     task.Record.App.Name,
				CurrentNode: t.GetNodeName(task.Record.App.Flow, task.NodeKey),
			}
			data.Dynamics = append(data.Dynamics, dynamic)
		}
	}
	return
}

// GetScheduleList
// @author: [zhaozijie]
// @function: GetScheduleList
// @description: 从mysql中获取待办数据
// @param: WorkFlowReq.Task
// @return: data []WorkFlowReq.ScheduleList, err error
func (t TaskService) GetScheduleList(userId uint) (scheduleData []WorkFlowRes.ScheduleList, err error) {
	var recordIds []uint
	global.GSD_DB.Model(&modelWF.GzlTask{}).Select("record_id").
		Where("node_type = ? AND Inspector = ?", 3, userId).
		Find(&recordIds)
	for i := 0; i < len(recordIds); i++ {
		var task modelWF.GzlTask
		err = global.GSD_DB.
			Where("record_id = ? AND updated_at is NULL", recordIds[i]).
			Preload("Record.App").
			Find(&task).Error
		if err != nil {
			return
		}
		schedule := WorkFlowRes.ScheduleList{
			CreatedAt:   task.CreatedAt,
			RecordId:    task.RecordId,
			Applicant:   t.GetUserNickName(task.Record.CreateBy),
			CurrentNode: t.GetNodeName(task.Record.App.Flow, task.Record.CurrentNode),
			AppName:     task.Record.App.Name,
			CheckState:  task.CheckState,
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
		Where("node_type = ? AND inspector = ?", 3, userId).
		Find(&recordIds)
	for i := 0; i < len(recordIds); i++ {
		var tasks []modelWF.GzlTask
		err = global.GSD_DB.
			Where("record_id = ? AND node_type = ? AND updated_at is NULL", recordIds[i], 3).
			Preload("Record.App").
			Find(&tasks).Error
		if err != nil {
			return
		}
		if len(tasks) > 0 {
			handle := WorkFlowRes.HandleList{
				CreatedAt:    tasks[0].CreatedAt,
				RecordId:     tasks[0].RecordId,
				Applicant:    t.GetUserNickName(tasks[0].Record.CreateBy),
				CurrentState: tasks[0].Record.CurrentState,
				AppName:      tasks[0].Record.App.Name,
				CurrentNode:  t.GetNodeName(tasks[0].Record.App.Flow, tasks[0].Record.CurrentNode),
			}
			for j := 0; j < len(tasks); j++ {
				Inspector := t.GetUserNickName(tasks[j].Inspector)
				if Inspector != "" {
					handle.Inspectors = append(handle.Inspectors, Inspector)
				}
			}
			handleData = append(handleData, handle)
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
		return ProcessFlow(taskInfo.Record, tx)
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
		Where("node_type = ? AND inspector = ?", 4, userId).
		Find(&recordIds)
	for i := 0; i < len(recordIds); i++ {
		var tasks []modelWF.GzlTask
		err = global.GSD_DB.
			Where("record_id = ? AND node_type = ? AND updated_at is NULL", recordIds[i], 3).
			Preload("Record.App").
			Find(&tasks).Error
		if err != nil {
			return
		}
		if len(tasks) > 0 {
			receive := WorkFlowRes.Receive{
				CreatedAt:    tasks[0].CreatedAt,
				RecordId:     tasks[0].RecordId,
				Applicant:    t.GetUserNickName(tasks[0].Record.CreateBy),
				CurrentState: tasks[0].Record.CurrentState,
				AppName:      tasks[0].Record.App.Name,
				CurrentNode:  t.GetNodeName(tasks[0].Record.App.Flow, tasks[0].Record.CurrentNode),
			}
			for j := 0; j < len(tasks); j++ {
				Inspector := t.GetUserNickName(tasks[j].Inspector)
				if Inspector != "" {
					receive.Inspectors = append(receive.Inspectors, Inspector)
				}
			}
			data = append(data, receive)
		}
	}
	return
}

// GetTaskInfo 根据id获取详细信息
func (t TaskService) GetTaskInfo(taskId uint) (task work_flow.GzlTask, err error) {
	err = global.GSD_DB.Preload("Record.App").First(&task, taskId).Error
	return task, err
}

// GetUserNickName 根据用户id获取用户昵称
func (t TaskService) GetUserNickName(userId uint) (nickName string) {
	global.GSD_DB.Model(&system.SysUser{}).
		Where("id = ?", userId).
		Select("nick_name as nickName").
		Find(&nickName)
	return
}

// GetNodeName 根据流程JSON和key获取当前节点名称
func (t TaskService) GetNodeName(flowJson utils.JSON, key string) string {
	var flow Flow
	_ = json.Unmarshal(flowJson, &flow)
	for _, node := range flow.FlowElementList {
		if node.Key == key {
			return node.Name
		}
	}
	return ""
}

// GetMoreNodesName 根据流程JSON获取全部节点名称
func (t TaskService) GetMoreNodesName(flowJson utils.JSON, tasks []modelWF.GzlTask, userId ...uint) (nodes []string) {
	var flow Flow
	_ = json.Unmarshal(flowJson, &flow)
	for _, node := range flow.FlowElementList {
		for _, task := range tasks {
			// 排除结束节点 5, 且 key 值相等
			if node.Type <= 4 && task.NodeKey == node.Key {
				nodeName := node.Properties.Name
				nodes = append(nodes, nodeName)
			}
		}
	}
	return
}
