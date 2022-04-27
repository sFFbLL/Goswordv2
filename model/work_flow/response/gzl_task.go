package response

import "time"

type DynamicList struct {
	Dynamics []Dynamic `json:"dynamics"`
	Nodes    []string  `json:"nodes"` // 全部流程节点
}

type Dynamic struct {
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	InspectAt   time.Time `json:"inspectAt"`   // 审批时间
	ConsumeTime int64     `json:"consumeTime"` // 审批耗时
	Applicant   string    `json:"applicant"`   // 申请人
	CheckState  uint8     `json:"checkState"`  // 审批状态
	Remarks     string    `json:"remarks"`     // 备注
	AppName     string    `json:"appName"`     // 应用名称
	CurrentNode string    `json:"currentNode"` // 当前节点名称
}

type Receive struct {
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
	RecordId     uint      `json:"recordId"`     // 记录id
	Applicant    string    `json:"applicant"`    // 申请人
	Inspectors   []string  `json:"inspector"`    // 审批人
	CurrentState uint8     `json:"currentState"` // 审批状态
	AppName      string    `json:"appName"`      // 应用名称
	CurrentNode  string    `json:"currentNode"`  // 当前节点名称
}
type ScheduleList struct {
	RecordId     uint      `json:"recordId"`     // 记录id
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	Applicant   string    `json:"applicant"`   // 申请人
	Inspectors   []string  `json:"inspector"`    // 审批人
	AppName     string    `json:"name"`        //应用名称
	CurrentNode string    `json:"currentNode"` //当前节点
	CheckState  uint8     `json:"checkState"`  //审批状态
}
type HandleList struct {
	RecordId     uint      `json:"recordId"`     // 记录id
	CreatedAt    time.Time `json:"createdAt"`    // 创建时间
	Applicant    string   `json:"applicant"`    // 申请人
	AppName      string   `json:"name"`         //应用名称
	CurrentState uint8    `json:"currentState"` //当前状态(进行中1默认、已完成2、已结束3)
	Inspectors   []string `json:"inspector"`    //审批人
	CurrentNode  string   `json:"currentNode"`  //当前节点
}
