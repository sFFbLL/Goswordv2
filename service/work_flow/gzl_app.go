package work_flow

import (
	"errors"
	"gorm.io/gorm"
	"project/global"
	modelWF "project/model/work_flow"
)

type AppService struct {
}

// GetAppEmpty
// @author: [tanshaokang]
// @function: GetAppEmpty
// @description: 从mysql中获取空表单
// @param: WorkFlowReq.EmptyApp
// @return: data utils.JSON, err error
func (a AppService) GetAppEmpty(appId int) (data string, err error) {
	db := global.GSD_DB.Model(&modelWF.GzlApp{}).
		Select("form").
		Where("id = ?", appId)
	if err = db.Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("数据为空")
		}
	}
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
