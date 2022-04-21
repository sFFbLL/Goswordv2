package work_flow

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"project/global"
	"project/model/system"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	WorkFlowRes "project/model/work_flow/response"
	"project/utils"
	"reflect"
)

type RecordService struct {
}

// GetData
// @author: [tanshaokang](https://github.com/worryfreet)
// @function: GetData
// @description: 从mysql中获取record数据
// @param: recordId int
// @return: data utils.JSON, err error
func (r RecordService) GetData(recordId uint) (data utils.JSON, err error) {
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
func (r RecordService) Submit(record modelWF.GzlRecord) (err error) {
	// 开启提交表单事务
	err = global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		// 插入部门id, 并创建一条新记录
		tx.Model(&system.SysUser{}).Select("dept_id").Where("id = ?", record.CreateBy).Find(&record.DeptId)
		tx.Create(&record)
		var recordId uint
		tx.Model(&modelWF.GzlRecord{}).Select("id").Last(&recordId)

		// 解析里面的form, 然后逐个添加进 gzl_form_items 里
		var formItems []modelWF.GzlFormItem
		var form WorkFlowReq.Form
		err = json.Unmarshal(record.Form, &form)
		if err != nil {
			global.GSD_LOG.ZapLog.Error("表单解析错误", zap.Any("err", err))
			// 遇见解析错误, 回滚
			tx.Rollback()
			return err
		}
		for _, item := range form.Fields {
			gfi := modelWF.GzlFormItem{
				GSD_MODEL:     global.GSD_MODEL{CreateBy: record.CreateBy},
				RecordId:      recordId,
				Name:          item.Name,
				DataType:      reflect.TypeOf(item.DefaultValue).String(),
				ComponentType: item.ComponentType,
				Content:       item.DefaultValue,
				IsRequired:    item.IsRequired,
				Form:          record.Form,
				DeptId:        record.DeptId,
			}
			formItems = append(formItems, gfi)
		}
		// 批量插入
		tx.Model(&modelWF.GzlFormItem{}).CreateInBatches(formItems, len(formItems))
		return err
	})
	return
}

// MyInitiated
// @author: [chenpipi]
// @function: MyInitiated
// @description: 从mysql中获取record数据
// @param: record WorkFlowRes.MyInitiated
// @return: data utils.JSON, err error
func (r RecordService) MyInitiated(id uint) (myInitiated WorkFlowRes.MyInitiated, err error) {
	var ids []uint
	if err = global.GSD_DB.Select("current_state,current_node").Model(&modelWF.GzlRecord{}).Find(&myInitiated, "create_by = ?", id).Error; err != nil {
		return
	}
	if err = global.GSD_DB.Select("inspector").Model(&modelWF.GzlTask{}).Find(&ids, "node_key = ?", myInitiated.CurrentNode).Error; err != nil {
		return
	}
	if err = global.GSD_DB.Select("nick_name").Model(&system.SysUser{}).Find(&myInitiated.InspectorName, "id in ?", ids).Error; err != nil {
		return
	}
	return
}
