package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"project/global"
	"project/model/common/request"
	"project/model/system"
	"project/utils"
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type UserService struct {
}

// Register @author: [chenguanglan](https://github.com/sFFbLL)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: err error, userInter model.SysUser
func (userService *UserService) Register(u system.SysUser) (err error, userInter system.SysUser) {
	var user system.SysUser
	if !errors.Is(global.GSD_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return errors.New("用户名已注册"), userInter
	}
	// 否则 附加uuid 密码md5简单加密 注册
	u.Password = utils.MD5V([]byte(u.Password))
	u.UUID = uuid.NewV4()
	err = global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Create(&u).Error
		if TxErr != nil {
			return TxErr
		}
		return nil
	})
	return err, u
}

// Login @author: [chenguanglan](https://github.com/sFFbLL)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser
func (userService *UserService) Login(u *system.SysUser) (err error, userInter *system.SysUser) {
	var user system.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GSD_DB.Where("username = ? AND password = ?", u.Username, u.Password).Preload("Dept").Preload("Authorities.Depts").First(&user).Error
	return err, &user
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserInfoList(info request.PageInfo, deptId []uint, isAll bool) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GSD_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	if isAll {
		err = db.Limit(limit).Offset(offset).Preload("Authorities").Preload("Dept").Find(&userList).Error
	} else {
		err = db.Where("dept_id in (?)", deptId).Limit(limit).Offset(offset).Preload("Authorities").Preload("Dept").Find(&userList).Error
	}
	err = db.Count(&total).Error
	return err, userList, total
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: UpdatePassword
//@description: 用户本人修改密码
//@param: user *system.SysUser, newPassword string
//@return: err error, sysUser *system.SysUser

func (userService *UserService) UpdatePassword(u *system.SysUser, newPassword string) (err error, sysUser *system.SysUser) {
	var user *system.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GSD_DB.Where("username = ? AND password = ?", u.Username, u.Username).First(&user).Update("password", utils.MD5V([]byte(newPassword))).Error
	return err, u
}

// SetUserAuthorities @author: [chenguanglan](https://github.com/sFFbLL)
//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []uint
//@return: err error
func (userService *UserService) SetUserAuthorities(updateUser system.SysUser, authorityIds []uint) (err error) {
	return global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		useAuthority := make([]system.SysAuthority, 0)
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysAuthority{
				AuthorityId: v,
			})
		}
		TxErr := tx.Model(&updateUser).Association("Authorities").Replace(useAuthority)
		if TxErr != nil {
			return TxErr
		}
		//更新缓存
		authorityIdJson, _ := json.Marshal(authorityIds)
		err = global.GSD_REDIS.HSet(context.Background(), updateUser.UUID.String(), "authorityId", authorityIdJson).Err()
		if err != nil {
			return err
		}
		return nil
	})
}

// DeleteUser @author: [chenguanglan](https://github.com/sFFbLL)
//@function: DeleteUser
//@description: 删除用户
//@param: id uint
//@return: err error
func (userService *UserService) DeleteUser(id uint) (err error) {
	return global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		var user system.SysUser
		TxErr := global.GSD_DB.Delete(&user, id).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = global.GSD_DB.Delete(&system.SysUseAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		//TxErr = CasbinServiceApp.DeleteUserAuthority(id)
		//if TxErr != nil {
		//	return tx.Rollback().Error
		//}
		return nil
	})
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: SetUserInfo
//@description: 修改用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetUserInfo(reqUser system.SysUser) (err error, user system.SysUser) {
	tx := global.GSD_DB.Begin()
	err = tx.Updates(&reqUser).Error
	if err != nil {
		tx.Rollback()
	}
	userInfo := reqUser
	//更新deptId
	userInfo.DeptId = reqUser.DeptId
	err = global.GSD_REDIS.HSet(context.Background(), reqUser.UUID.String(), "deptId", reqUser.DeptId).Err()
	if err != nil {
		tx.Rollback()
	}
	tx.Commit()
	return err, reqUser
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser

func (userService *UserService) GetUserInfo(uuid string) (err error, user system.SysUser) {
	var reqUser system.SysUser
	err = global.GSD_DB.Preload("Authorities").Preload("Dept").First(&reqUser, "uuid = ?", uuid).Error
	return err, reqUser
}

//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserById(id uint) (err error, user *system.SysUser) {
	var u system.SysUser
	err = global.GSD_DB.Preload("Authorities").Preload("Dept").Where("`id` = ?", id).First(&u).Error
	return err, &u
}

// FindUserByUuid @function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.SysUser
func (userService *UserService) FindUserByUuid(uuid string) (err error, user *system.SysUser) {
	var u system.SysUser
	if err = global.GSD_DB.Preload("Authorities").Preload("Dept").Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return errors.New("用户不存在"), &u
	}
	return nil, &u
}

// FindUserByDept @function: FindUserByDept
//@description: 通过dept获取用户信息
//@param: deptId uint
//@return: err error, userId []uint
func (userService *UserService) FindUserByDept(deptId []uint) (err error, userId []uint) {
	err = global.GSD_DB.Select("id").Where("`deptId` in ?", deptId).Find(&userId).Error
	return
}

// FindUserByAuthority @function: FindUserByAuthority
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user []uint
func (userService *UserService) FindUserByAuthority(authorityId []uint) (err error, userId []uint) {
	global.GSD_DB.Model(system.SysUseAuthority{}).Distinct("sys_user_id").Where("`sys_authority_authority_id` in (?)", authorityId).Find(&userId)
	return
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: ParseExcelToDataList
//@description: 加载excel
//@return: user []system.SysUser, err error
func (userService *UserService) ParseExcelToDataList() ([]system.SysUser, error) {
	skipHeader := true
	fixedHeader := []string{"ID", "用户名", "昵称", "电话号", "邮箱", "部门名称"}
	file, err := excelize.OpenFile(global.GSD_CONFIG.Excel.Dir + "ExcelImport.xlsx")
	if err != nil {
		return nil, err
	}
	users := make([]system.SysUser, 0)
	rows, err := file.Rows("Sheet1")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		if skipHeader {
			if utils.CompareStrSlice(row, fixedHeader) {
				skipHeader = false
				continue
			} else {
				return nil, errors.New("excel格式错误")
			}
		}
		if len(row) != len(fixedHeader) {
			continue
		}
		id, _ := strconv.Atoi(row[0])
		user := system.SysUser{
			GSD_MODEL: global.GSD_MODEL{
				ID: uint(id),
			},
			Username: row[1],
			NickName: row[2],
			Phone:    row[3],
			Email:    row[4],
			Dept:     system.SysDept{DeptName: row[5]},
		}
		users = append(users, user)
	}
	return users, nil
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: ParseDataListToExcel
//@description: 导出excel
//@param: info []system.SysUser, path string
//@return: err error
func (userService *UserService) ParseDataListToExcel(info []system.SysUser, path string) error {
	excel := excelize.NewFile()
	excel.SetSheetRow("Sheet1", "A1", &[]string{"ID", "用户名", "昵称", "电话号", "邮箱", "部门名称"})
	for i, user := range info {
		axis := fmt.Sprintf("A%d", i+2)
		excel.SetSheetRow("Sheet1", axis, &[]interface{}{
			user.ID,
			user.Username,
			user.NickName,
			user.Phone,
			user.Email,
			user.Dept.DeptName,
		})
	}
	err := excel.SaveAs(path)
	return err
}

//@author: [houruotong](https://github.com/Monkey-Pear)
//@function: Template
//@description: 下载模板
//@param: path string
//@return: err error
func (userService *UserService) Template(path string) error {
	excel := excelize.NewFile()
	excel.SetSheetRow("Sheet1", "A1", &[]string{"ID", "用户名", "昵称", "电话号", "邮箱", "部门名称"})
	axis := fmt.Sprintf("A%d", 2)
	excel.SetSheetRow("Sheet1", axis, &[]interface{}{
		"1",
		"test",
		"测试用户",
		"156***",
		"xxx@xx.com",
		"顶级部门",
	})

	err := excel.SaveAs(path)
	return err
}
