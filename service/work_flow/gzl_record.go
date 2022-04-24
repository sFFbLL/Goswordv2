package work_flow

import (
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"project/global"
	"project/model/system"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
	WorkFlowRes "project/model/work_flow/response"
	"reflect"
)

type RecordService struct {
}

// GetData
// @author: [tanshaokang]
// @function: GetData
// @description: 从mysql中获取record数据
// @param: recordId int
// @return: data string, err error
func (r RecordService) GetData(recordId uint) (data string, err error) {
	db := global.GSD_DB.Model(&modelWF.GzlRecord{}).
		Select("form").
		Where("id = ?", recordId)
	if err = db.First(&data).Error; err != nil {
		return "", errors.New("该数据不存在")
	}
	return
}

// Submit
// @author: [tanshaokang]
// @function: Submit
// @description: 提交记录
// @param: WorkFlowReq.RecordSubmit
// @return: err error
func (r RecordService) Submit(record modelWF.GzlRecord) (err error) {
	// 判断appId在数据库中是否存在
	if err = global.GSD_DB.Where("id = ?", record.AppId).First(&modelWF.GzlApp{}).Error; err != nil {
		return errors.New("该数据不存在")
	}
	// 开启提交表单事务
	err = global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		// 插入部门id, 并创建一条新记录
		tx.Model(&system.SysUser{}).
			Select("dept_id").
			Where("id = ?", record.CreateBy).
			Find(&record.DeptId)
		// 在gzl_records中添加记录
		tx.Create(&record)
		var formItems []modelWF.GzlFormItem
		var form WorkFlowReq.Form
		_ = json.Unmarshal(record.Form, &form)
		for _, item := range form.Fields {
			var isRequired uint8 = 1
			if !item.IsRequired {
				isRequired = 2
			}
			gfi := modelWF.GzlFormItem{
				GSD_MODEL:     global.GSD_MODEL{CreateBy: record.CreateBy},
				RecordId:      record.ID,
				Name:          item.Name,
				DataType:      reflect.TypeOf(item.DefaultValue).String(),
				ComponentType: item.ComponentType,
				Content:       item.DefaultValue,
				IsRequired:    isRequired,
				Form:          record.Form, // 目前先不存
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
func (r RecordService) MyInitiated(userid uint) (myInitiated []WorkFlowRes.MyInitiated, err error) {
	var records []modelWF.GzlRecord
	if err = global.GSD_DB.Preload("App").Find(&records, userid).Error; err != nil {
		return
	}
	for _, record := range records {
		var tasks []modelWF.GzlTask
		var names []string
		var flow Flow
		err = json.Unmarshal(record.App.Flow, &flow)
		if err != nil {
			//	TODO 错误
			return
		}
		// 构造 flow map
		flowMap := make(map[string]Node)
		for _, flowElement := range flow.FlowElementList {
			flowMap[flowElement.Key] = flowElement
		}
		err = global.GSD_DB.Preload("User").Find(&tasks, "node_key = ?", record.CurrentNode).Error
		if err != nil {
			return
		}
		for _, task := range tasks {
			names = append(names, task.User.NickName)
		}
		myInitiated = append(myInitiated, WorkFlowRes.MyInitiated{
			AppName:      record.App.Name,
			CurrentState: record.CurrentState,
			CurrentNode:  flowMap[record.CurrentNode].Name,
			Inspector:    names,
		})
	}
	return myInitiated, nil
}
