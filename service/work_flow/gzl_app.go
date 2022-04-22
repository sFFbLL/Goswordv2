package work_flow

import (
	"errors"
	"gorm.io/gorm"
	"project/global"
	modelWF "project/model/work_flow"
	"project/utils"
)

type AppService struct {
}

// GetAppEmpty
// @author: [tanshaokang]
// @function: GetAppEmpty
// @description: 从mysql中获取空表单
// @param: WorkFlowReq.EmptyApp
// @return: data utils.JSON, err error
func (a AppService) GetAppEmpty(appId int) (data utils.JSON, err error) {
	var datas = make([]utils.JSON, 1)
	db := global.GSD_DB.Model(&modelWF.GzlApp{}).
		Select("form").
		Where("id = ?", appId)
	if err = db.First(&datas).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("数据为空")
		} else {
			return
		}
	}
	data = datas[0]
	return
}

// AddApp
// @author: [chenpipi]
// @function: AddApp
// @description: 添加应用
// @param: WorkFlowReq.AddApp
// @return: err error
func (a AppService) AddApp(app modelWF.GzlApp) (err error) {
	if err = global.GSD_DB.Create(&app).Error; err != nil {
		return
	}
	return
}
