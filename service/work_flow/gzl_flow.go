package work_flow

import (
	"encoding/json"
	"fmt"
	"project/global"
	"project/model/work_flow"
	"project/service/system"
)

// Flow 流程结构体
type Flow struct {
	FlowElementList []Node `json:"flowElementList"`
}

type Node struct {
	Incoming   []string `json:"incoming"` //入口
	Outgoing   []string `json:"outgoing"` //出口
	Type       uint8    `json:"type"`     //(线1、开始2、审批3、抄送4、结束5)
	Key        string   `json:"key"`      //Key
	Properties `json:"properties"`
}

type Properties struct {
	Depts         []uint   `json:"depts"`         //部门
	Authoritys    []uint   `json:"authoritys"`    //角色
	Users         []uint   `json:"users"`         //用户
	Conditions    []string `json:"conditions"`    //条件
	IsCountersign uint8    `json:"isCountersign"` //是否会签(会签1,或签2)
	Name          string   `json:"name"`          // 节点自定义名称

}

// Form 表单结构体
type Form struct {
	Fields []FormItem `json:"fields"`
}

type FormItem struct {
	Config `json:"__config__"`
}

type Config struct {
	RenderKey    uint   `json:"renderKey"`
	DefaultValue string `json:"defaultValue"`
}

//ProcessFlow 流程流转
func ProcessFlow(record work_flow.GzlRecord) (err error) {
	var flow Flow
	var form Form
	var tasks []work_flow.GzlTask
	// 流程JSON转换结构体
	err = json.Unmarshal(record.App.Flow, &flow)
	if err != nil {
		//	TODO 错误
		fmt.Println("Flow 解析错误")
		return
	}
	// 表单JSON转换结构体
	err = json.Unmarshal(record.App.Form, &form)
	if err != nil {
		//	TODO 错误
		fmt.Println("Form 解析错误")
		return
	}
	// 构造 flow map
	flowMap := make(map[string]Node)
	for _, flowElement := range flow.FlowElementList {
		flowMap[flowElement.Key] = flowElement
	}
	// 构造 form map
	formMap := make(map[uint]FormItem)
	for _, field := range form.Fields {
		formMap[field.RenderKey] = field
	}
	// 首次当前节点
	currentNode := record.CurrentNode
	if currentNode == "" {
		currentNode = flow.FlowElementList[0].Key
	} else {
		// 判断会签或签
		if flowMap[currentNode].IsCountersign == 1 {
			// 判断未签人数
			err = global.GSD_DB.Find(&tasks, "record_id = ? AND check_state = ?", record.ID, 1).Error
			if err != nil {
				return err
			}
			if len(tasks) != 0 {
				return
			}
		}
	}
	// 走节点
	userService := &system.UserService{}
outFlow:
	for true {
		if err != nil {
			break
		}
		for i, key := range flowMap[currentNode].Outgoing {
			line := flowMap[key]
			node := flowMap[line.Outgoing[0]]
			if conditions(line.Conditions) {
				currentNode = node.Key
			} else if i == len(flowMap[currentNode].Outgoing) && !conditions(line.Conditions) {
				// 无路可走
				record.CurrentState = 3
				record.CurrentNode = ""
				err = updateRecord(record)
				if err != nil {
					return
				}
				break outFlow
			}
			if node.Type == 3 {
				//	更新记录数据整理
				record.CurrentNode = node.Key
				//	更新记录表
				err = updateRecord(record)
				if err != nil {
					return
				}
				var userIds []uint
				var ids []uint
				//	通过部门获取用户
				if len(line.Depts) != 0 {
					err, ids = userService.FindUserByDept(line.Depts)
					if err != nil {
						return
					}
				}
				userIds = append(userIds, ids...)
				//	通过角色获取用户
				if len(line.Authoritys) != 0 {
					err, ids = userService.FindUserByAuthority(line.Authoritys)
					if err != nil {
						return
					}
				}
				userIds = append(userIds, ids...)
				//  获取用户
				userIds = append(userIds, node.Users...)
				//	审批任务数据整理
				for _, id := range userIds {
					tasks = append(tasks, work_flow.GzlTask{RecordId: record.ID, NodeType: node.Type, IsCountersign: node.IsCountersign, NodeKey: node.Key, Inspector: id})
				}
				if len(tasks) != 0 {
					err = issueTask(tasks)
					if err != nil {
						return
					}
				}
				//	下发审批任务
				break outFlow
			} else if node.Type == 4 {
				var userIds []uint
				var ids []uint
				//	通过部门获取用户
				if len(line.Depts) != 0 {
					err, ids = userService.FindUserByDept(line.Depts)
					if err != nil {
						return
					}
				}
				userIds = append(userIds, ids...)
				//	通过角色获取用户
				if len(line.Authoritys) != 0 {
					err, ids = userService.FindUserByAuthority(line.Authoritys)
					if err != nil {
						return
					}
				}
				userIds = append(userIds, ids...)
				//  获取用户
				userIds = append(userIds, node.Users...)
				//	抄送任务数据整理
				for _, id := range userIds {
					tasks = append(tasks, work_flow.GzlTask{RecordId: record.ID, NodeType: node.Type, IsCountersign: node.IsCountersign, NodeKey: node.Key, Inspector: id})
				}
				//	下发抄送任务
				if len(tasks) != 0 {
					err = issueTask(tasks)
					if err != nil {
						return
					}
				}
				break
			} else if node.Type == 5 {
				//	更新记录表
				record.CurrentState = 2
				record.CurrentNode = ""
				err = updateRecord(record)
				if err != nil {
					return
				}
				break outFlow
			}
		}
	}
	return
}

// 判断条件
func conditions(str []string) bool {
	// TODO 条件
	return true
}

// 下发任务
func issueTask(tasks []work_flow.GzlTask) error {
	return global.GSD_DB.Create(&tasks).Error
}

// 更新记录
func updateRecord(record work_flow.GzlRecord) error {
	return global.GSD_DB.Updates(&record).Error
}
