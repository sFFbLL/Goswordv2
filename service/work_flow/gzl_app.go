package work_flow

import (
	"errors"
	"gorm.io/gorm"
	"project/global"
	modelWF "project/model/work_flow"
	WorkFlowReq "project/model/work_flow/request"
)

type AppService struct {
}

func (a AppService) GetAppEmpty(app WorkFlowReq.App) (data modelWF.JSON, err error) {
	var datas = make([]modelWF.JSON, 1)
	db := global.GSD_DB.Model(&modelWF.GzlApp{}).
		Select("form").
		Where("id = ?", app.AppId)
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
