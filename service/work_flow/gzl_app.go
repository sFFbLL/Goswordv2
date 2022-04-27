package work_flow

import (
	"errors"
	"gorm.io/gorm"
	"project/global"
	"project/model/system"
	modelWF "project/model/work_flow"
)

type AppService struct {
}

// GetAppEmpty
// @author: [tanshaokang]
// @function: GetAppEmpty
// @description: 从mysql中获取空表单
// @param: WorkFlowReq.EmptyApp
// @return: data string, err error
func (a AppService) GetAppEmpty(appId uint) (data string, err error) {
	db := global.GSD_DB.Model(&modelWF.GzlApp{}).
		Select("form").
		Where("id = ?", appId)
	if err = db.First(&data).Error; err != nil {
		return "", errors.New("数据不存在")
	}
	return
}

// AddApp
// @author: [chenpipi]
// @function: AddApp
// @description: 添加应用
// @return: err error
func (a AppService) AddApp(app modelWF.GzlApp) (err error) {
	if err = global.GSD_DB.Create(&app).Error; err != nil {
		return
	}
	return
}

// AddForm
// @author: [chenpipi]
// @function: AddForm
// @description: 添加应用
// @return: err error
func (a AppService) AddForm(app modelWF.GzlApp) (err error) {
	if err = global.GSD_DB.Where("is_enable = ?", 1).Updates(&app).Error; err != nil {
		return
	}
	return
}

// AddFlow
// @author: [chenpipi]
// @function: AddFlow
// @description: 添加应用
// @return: err error
func (a AppService) AddFlow(app modelWF.GzlApp) (err error) {
	if err = global.GSD_DB.Where("is_enable = ?", 1).Updates(&app).Error; err != nil {
		return
	}
	return
}

// EnableApp
// @author: [chenpipi]
// @function: EnableApp
// @description: 添加应用
// @return: err error
func (a AppService) EnableApp(app modelWF.GzlApp) (err error) {
	if err = global.GSD_DB.Updates(&app).Error; err != nil {
		return
	}
	return
}

// AuthorityApp
// @author: [chenpipi]
// @function: AuthorityApp
// @description: 权限分配
// @return: err error
func (a AppService) AuthorityApp(app modelWF.GzlApp, depts []system.SysDept, authoritys []system.SysAuthority, users []system.SysUser) (err error) {
	return global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&app)
		if err = tx.Association("Depts").Replace(&depts); err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		if err = tx.Association("Authoritys").Replace(&authoritys); err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		if err = tx.Association("Users").Replace(&users); err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		return nil
	})
}
