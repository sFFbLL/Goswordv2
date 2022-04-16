package work_flow

import (
	"errors"
	"gorm.io/gorm"
	"project/global"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	"project/utils"
)

type RecordService struct {
}

// GetData
// @author: [tanshaokang](https://github.com/worryfreet)
// @function: GetData
// @description: 从mysql中获取record数据
// @param: WorkFlowReq.Record
// @return: data utils.JSON, err error
func (r RecordService) GetData(record WorkFlowReq.Record) (data utils.JSON, err error) {
	var dataArr = make([]utils.JSON, 1)
	db := global.GSD_DB.Model(&modelWF.GzlRecord{}).
		Select("form").
		Where("id = ?", record.RecordId)
	if err = db.Take(&dataArr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据为空")
		} else {
			return
		}
	}
	data = dataArr[0]
	return
}
