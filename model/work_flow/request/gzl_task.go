package request

type Inspect struct {
	TaskId uint  `json:"taskId"`
	State  uint8 `json:"state"`
}
