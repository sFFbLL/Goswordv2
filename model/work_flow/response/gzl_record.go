package response

type MyInitiated struct {
	AppName      string   `json:"appName"`
	CurrentState uint8    `json:"currentState"` //当前状态(进行中1默认、已完成2、已结束3)
	CurrentNode  string   `json:"currentNode"`  //当前节点
	Inspector    []string `json:"inspector"`
}
