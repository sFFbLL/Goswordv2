package work_flow

import (
	"project/global"
)

//GzlTask 任务表
type GzlTask struct {
	global.GSD_MODEL
	RecordId      uint      `json:"recordId" gorm:"not null;comment:记录id"`                                            //记录id
	Record        GzlRecord `json:"record" gorm:"foreignKey:RecordId;references:ID;"`                                 //记录
	CheckState    uint8     `json:"checkState" gorm:"not null;comment:审批状态"`                                          //审批状态(待审批1默认、审批通过2、审批拒绝3、或签已审核4)
	NodeType      uint8     `json:"nodeType" gorm:"not null;default:2;comment:节点类型(连接线1、开始节点2默认、线审批节点3、抄送节点4、结束节点5)"` //节点类型(连接线1、开始节点2默认、线审批节点3、抄送节点4、结束节点5)
	IsCountersign uint8     `json:"isCountersign" gorm:"not null;default:1;comment:是否会签(会签1默认、或签2)"`                  //是否会签(会签1默认、或签2)
	NodeKey       string    `json:"nodeKey" gorm:"not null;comment:节点Key"`                                            //节点Key
	Inspector     uint      `json:"inspector" gorm:"not null;comment:审批人"`                                            //审批人
	Remarks       string    `json:"remarks" gorm:"comment:备注"`                                                        //备注
}
