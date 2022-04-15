package request

import "time"

type Task struct {
	State uint8
}

type Dynamic struct {
	CreatedAt   time.Time `json:"createdAt"`   // 创建时间
	ConsumeTime int       `json:"consumeTime"` // 审批耗时
	Applicant   string    `json:"applicant"`   // 申请人
	CheckState  uint8     `json:"checkState"`  // 审批状态
	Remarks     string    `json:"remarks"`     // 备注
}
