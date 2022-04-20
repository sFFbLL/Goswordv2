package request

import "time"

type Inspect struct {
	TaskId uint  `json:"taskId"`
	State  uint8 `json:"state"`
}

type Function struct {
	CreatedAt  time.Time `json:"createdAt"` // 创建时间
	Applicant  string    `json:"applicant"` // 申请人
	Inspector  uint      `json:"inspector" gorm:"not null;comment:审批人"`
	AppName    string    `json:"appName"`                                                            //应用名称
	CheckState uint8     `json:"checkState" gorm:"not null;comment:审批状态(待审批1默认、审批通过2、审批拒绝3、或签已审核4)"` //审批状态(待审批1默认、审批通过2、审批拒绝3、或签已审核4)
}
