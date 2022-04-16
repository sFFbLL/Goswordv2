package work_flow

import (
	"encoding/json"
	"project/global"
	"project/model/work_flow"
)

type Flows struct {
	FlowElementList []Flow `json:"flowElementList"`
}

type Flow struct {
	Incoming   []string `json:"incoming"` //入口
	Outgoing   []string `json:"outgoing"` //出口
	Type       uint8    `json:"type"`     //(线1、开始2、审批3、抄送4、结束5)
	Key        string   `json:"key"`      //Key
	Properties `json:"properties"`
}

type Properties struct {
	Depts         []uint `json:"depts"`         //部门
	Authoritys    []uint `json:"authoritys"`    //角色
	Users         []uint `json:"users"`         //用户
	IsCountersign uint8  `json:"isCountersign"` //是否会签(会签1,或签2)
}

//ProcessFlow 流程流转
func ProcessFlow(recordId uint) {
	var record work_flow.GzlRecord
	var flow Flows
	// 获取流程JSON并转换给结构体
	err := global.GSD_DB.Preload("App").First(&record, recordId).Error
	if err != nil {
		return
	}
	err = json.Unmarshal(record.App.Flow, &flow)
	if err != nil {
		return
	}
	//判读当前节点
	for _, f := range flow.FlowElementList {
		if record.CurrentNode == f.Key {

		}
	}
}
