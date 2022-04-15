package work_flow

import (
	"errors"
	"gorm.io/gorm"
	"project/global"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
)

type TaskService struct {
}

func (t TaskService) GetDynamic(applicantId, recordId int) (data []WorkFlowReq.Dynamic, err error) {
	db := global.GSD_DB.
		Model(modelWF.GzlTask{}).
		Joins("JOIN sys_users ON sys_users.id = ?", applicantId).
		Joins("JOIN gzl_records ON gzl_records.id = ?", recordId).
		Select("sys_users.username as Applicant", "(gzl_records.created_at - gzl_tasks.created_at) as ConsumeTime",
			"gzl_records.created_at as CreatedAt", "check_state as CheckState", "remarks as Remarks").
		Where("gzl_tasks.record_id = gzl_records.id")
	if err = db.Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
	}

	return
}
