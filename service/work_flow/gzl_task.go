package work_flow

import (
	"encoding/json"
	"fmt"
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
// @return: data []WorkFlowReq.Dynamic, err error
func (t TaskService) GetDynamic(recordId uint) (data []WorkFlowRes.Dynamic, err error) {
	var tasks []modelWF.GzlTask
	global.GSD_DB.Preload("Record.App").
		Where("recordId = ?", recordId).
		Find(&tasks)
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
			Nodes:       t.GetMoreNodesName(task.Record.App.Flow, tasks),
		}
		data = append(data, dynamic)
	}
	return
}

// GetScheduleList
// @author: [zhaozijie]
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
			Where("record_id = ? AND node_type = ? AND update_at is NULL", recordIds[i], 3).
			Preload("Record.App").
			Find(&tasks)
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
func (t TaskService) GetMoreNodesName(flowJson utils.JSON, tasks []modelWF.GzlTask) (nodes []string) {
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
