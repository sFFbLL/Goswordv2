package response

import "time"

type Dynamic struct {
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	InspectAt   time.Time `json:"inspectAt"`   // 审批时间
	ConsumeTime int64     `json:"consumeTime"` // 审批耗时
	Applicant   string    `json:"applicant"`   // 申请人
	CheckState  uint8     `json:"checkState"`  // 审批状态
	Remarks     string    `json:"remarks"`     // 备注
}

type Receive struct {
	Applicant    string   `json:"applicant"`    // 申请人
	Inspectors   []string `json:"inspector"`    // 审批人
	CurrentState uint8    `json:"currentState"` // 审批状态
	AppName      string   `json:"appName"`      // 应用名称
	CurrentNode  string   `json:"currentNode"`  // 当前节点名称
}
