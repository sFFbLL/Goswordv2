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
	// 保证多表操作时是原子操作, 开启事务
	global.GSD_DB.Begin()
	global.GSD_DB.Model(&modelWF.GzlRecord{}).Create(&recordSubmit)
	var recordId int
	global.GSD_DB.Model(&modelWF.GzlRecord{}).Select("id").Last(&recordId)
	// 解析里面的form, 然后逐个添加进 gzl_form_items 里
	var formItems []WorkFlowReq.FormItem
	form, _ := json.Marshal(recordSubmit.Form)
	global.GSD_DB.Where("id = ?", 1).Updates(&modelWF.GzlApp{Form: form})
	err = json.Unmarshal(form, &formItems)
	if err != nil {
		global.GSD_LOG.Error("表单解析错误", zap.Any("err", err))
		// 遇见解析错误, 回滚
		global.GSD_DB.Rollback()
		return
	}
	for i := 0; i < len(formItems); i++ {
		formItems[i].CreateBy = recordSubmit.CreateBy
		formItems[i].DeptId = recordSubmit.DeptId
		formItems[i].RecordId = uint(recordId)
		//formItems[i].Form = recordSubmit.Form
	}
	// 批量插入
	global.GSD_DB.Model(&modelWF.GzlFormItem{}).CreateInBatches(formItems, len(formItems))
	// 成功, 提交
	global.GSD_DB.Commit()
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
