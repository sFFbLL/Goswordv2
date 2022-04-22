package response

import "time"

type Dynamic struct {
	CreatedAt          time.Time `json:"createdAt"`          // 创建时间
	InspectAt          time.Time `json:"inspectAt"`          // 审批时间
	CreatedAtFormatStr string    `json:"createdAtFormatStr"` // 格式化后的创建时间 --> yyyy-mm-dd hh:mm
	InspectAtFormatStr string    `json:"inspectAtFormatStr"` // 格式化后的审批时间 --> yyyy-mm-dd hh:mm
	ConsumeTime        int64     `json:"consumeTime"`        // 审批耗时
	Applicant          string    `json:"applicant"`          // 申请人
	CheckState         uint8     `json:"checkState"`         // 审批状态
	Remarks            string    `json:"remarks"`            // 备注
}

type Receive struct {
}
type ScheduleList struct {
	CreatedAt  time.Time `json:"createdAt"`  // 创建时间
	Applicant  string    `json:"applicant"`  // 申请人
	Name       string          `json:"name"` //应用名称
	CheckState uint8     `json:"checkState"` //审批状态
}
type HandleList struct {
	CreatedAt      time.Time `json:"createdAt"`     // 创建时间
	Applicant      string    `json:"applicant"`     // 申请人
	Name           string    `json:"name"`          //应用名称
	CurrentState   uint8     `json:"currentState"`  //当前状态(进行中1默认、已完成2、已结束3)
	InspectorName  []string    `json:"inspector"`     //审批人
	CurrentNode   string     `json:"currentNode"`     //当前节点
}