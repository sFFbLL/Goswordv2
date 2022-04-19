package work_flow

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/global"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	"project/utils"
	"time"
)

type RecordService struct {
}

// GetData
// @author: [tanshaokang](https://github.com/worryfreet)
// @function: GetData
// @description: 从mysql中获取record数据
// @param: recordId int
// @return: data utils.JSON, err error
func (r RecordService) GetData(recordId int) (data utils.JSON, err error) {
	var dataArr = make([]utils.JSON, 1)
	db := global.GSD_DB.Model(&modelWF.GzlRecord{}).
		Select("form").
		Where("id = ?", recordId)
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

// Submit
// @author: [tanshaokang](https://github.com/worryfreet)
// @function: Submit
// @description: 提交记录
// @param: WorkFlowReq.RecordSubmit
// @return: err error
func (r RecordService) Submit(recordSubmit WorkFlowReq.RecordSubmit) (err error) {
	// 先存进 gzl_record 数据库中
	recordSubmit.CreatedAt = time.Now()
	// 保证多表操作时是原子操作, 开启事务
	global.GSD_DB.Begin()
	global.GSD_DB.Model(&modelWF.GzlRecord{}).Create(&recordSubmit)
	var recordId int
	global.GSD_DB.Model(&modelWF.GzlRecord{}).Select("id").Last(&recordId)
	// 解析里面的form, 然后逐个添加进 gzl_form_items 里
	var formItems []WorkFlowReq.FormItem
	err = json.Unmarshal(recordSubmit.Form, &formItems)
	if err != nil {
		global.GSD_LOG.ZapLog.Error("表单解析错误", zap.Any("err", err))
		// 遇见解析错误, 回滚
		global.GSD_DB.Rollback()
		return
	}
	for i := 0; i < len(formItems); i++ {
		formItems[i].CreateBy = recordSubmit.CreateBy
		formItems[i].CreatedAt = recordSubmit.CreatedAt
		formItems[i].DeptId = recordSubmit.DeptId
		formItems[i].RecordId = uint(recordId)
		formItems[i].Form = recordSubmit.Form
	}
	// 批量插入
	global.GSD_DB.Model(&modelWF.GzlFormItem{}).CreateInBatches(formItems, len(formItems))
	// 成功, 提交
	global.GSD_DB.Commit()
	return
}
