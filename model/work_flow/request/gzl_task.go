package request

import (
	"time"
)

type Inspect struct {
	TaskId uint  `json:"taskId"`
	State  uint8 `json:"state"`
}

type Dynamic struct {
	CreatedAt          time.Time `json:"createdAt"`          // 创建时间
	InspectAt          time.Time `json:"inspectAt"`          // 审批时间
	CreatedAtFormatStr string    `json:"createdAtFormatStr"` //格式化后的创建时间 --> yyyy-mm-dd hh:mm
	InspectAtFormatStr string    `json:"inspectAtFormatStr"` //格式化后的审批时间 --> yyyy-mm-dd hh:mm
	ConsumeTime        int64     `json:"consumeTime"`        // 审批耗时
	Applicant          string    `json:"applicant"`          // 申请人
	CheckState         uint8     `json:"checkState"`         // 审批状态
	Remarks            string    `json:"remarks"`            // 备注

}


