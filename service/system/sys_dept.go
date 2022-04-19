package system

import (
	"errors"
	"project/global"
	"project/model/common/request"
	"project/model/system"

	"gorm.io/gorm"
)

type DeptService struct {
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: AddDepartment
//@description: 新增部门
//@param: dept system.SysDept
//@return: error

func (departService *DeptService) AddDepartment(dept system.SysDept) (err error) {
	if !errors.Is(global.GSD_DB.Where("dept_name = ?", dept.DeptName).First(&system.SysDept{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	return global.GSD_DB.Create(&dept).Error
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: DeleteDepartment
//@description: 删除部门
//@param: id uint
//@return: error

func (departService *DeptService) DeleteDepartment(id uint) (err error) {
	err = global.GSD_DB.Where("parent_id", id).First(&system.SysDept{}).Error
	if err != nil {
		var dept system.SysDept
		err = global.GSD_DB.Where("id = ?", id).Delete(&dept).Error
		if err != nil {
			return err
		}
	} else {
		return errors.New("存在子部门，不能删除")
	}
	return err
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: UpdateDepartment
//@description: 修改部门
//@param: dept system.SysDept
//@return: error

func (departService *DeptService) UpdateDepartment(dept system.SysDept) (err error) {
	var oldDept system.SysDept
	updateDept := make(map[string]interface{})
	updateDept["dept_name"] = dept.DeptName
	updateDept["parent_id"] = dept.ParentID
	updateDept["dept_sort"] = dept.DeptSort
	db := global.GSD_DB.Where("id = ?", dept.ID).First(&oldDept)
	err = db.Error
	if err != nil {
		return err
	}
	if oldDept.DeptName != dept.DeptName {
		if !errors.Is(global.GSD_DB.Where("dept_name = ?", dept.DeptName).First(&system.SysDept{}).Error, gorm.ErrRecordNotFound) {
			return errors.New("存在重复部门名字，请修改部门名字")
		}
	}
	return db.Updates(&dept).Error
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: GetDeptList
//@description: 查询部门列表
//@param: info request.PageInfo
//@return: error

func (departService *DeptService) GetDeptList(info request.PageInfo, deptId []uint, isAll bool) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GSD_DB.Model(&system.SysDept{})
	var dept []system.SysDept

	if isAll {
		err = db.Limit(limit).Offset(offset).Find(&dept).Order("dept_sort").Error
	} else {
		err = db.Where("id in (?)", deptId).Limit(limit).Offset(offset).Find(&dept).Order("dept_sort").Error
	}
	err = db.Count(&total).Where("parent_id = ?", 0).Error
	return err, dept, total
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: GetDeptListById
//@description: 根据pid查询部门列表
//@param: id uint
//@return: error

func (departService *DeptService) GetDeptListById(id uint, deptId []uint, isAll bool) (err error, list interface{}, total int64) {
	db := global.GSD_DB.Model(&system.SysDept{})
	var dept []system.SysDept
	if isAll {
		err = db.Where("parent_id = ?", id).Find(&dept).Order("dept_sort").Error
	} else {
		err = db.Where("id in (?) AND parent_id = ?", deptId, id).Find(&dept).Order("dept_sort").Error
	}
	err = db.Where("parent_id = ? ", id).Count(&total).Error
	return err, dept, total
}
