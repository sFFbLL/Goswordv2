package response

type MyInitiated struct {
	CurrentState  uint8    `json:"currentState"` //当前状态(进行中1默认、已完成2、已结束3)
	InspectorName []string `json:"inspector"`
	CurrentNode   string   `json:"currentNode"` //当前节点
}
