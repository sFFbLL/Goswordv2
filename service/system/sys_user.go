package system

import (
	"errors"
	"project/global"
	"project/model/common/request"
	"project/model/system"
	"project/utils"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: Register
//@description: 用户注册
//@param: u model.SysUser
//@return: err error, userInter model.SysUser

type UserService struct {
}

func (userService *UserService) Register(u system.SysUser, roles []uint) (err error, userInter system.SysUser) {
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
		TxErr = CasbinServiceApp.UpdateUserAuthority(u.ID, roles)
		if TxErr != nil {
			return TxErr
		}
		return nil
	})
	return err, u
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: Login
//@description: 用户登录
//@param: u *model.SysUser
//@return: err error, userInter *model.SysUser

func (userService *UserService) Login(u *system.SysUser) (err error, userInter *system.SysUser) {
	var user system.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GSD_DB.Where("username = ? AND password = ?", u.Username, u.Password).Preload("Dept").Preload("Authorities").First(&user).Error
	return err, &user
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: GetUserInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: err error, list interface{}, total int64

func (userService *UserService) GetUserInfoList(info request.PageInfo) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GSD_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Preload("Authorities").Find(&userList).Error
	return err, userList, total
}

//@function: UpdatePassword
//@description: 用户修改密码
//@param: user *system.SysUser, newPassword string
//@return: err error, sysUser *system.SysUser

func (userService *UserService) UpdatePassword(u *system.SysUser, newPassword string) (err error, sysUser *system.SysUser) {
	var user *system.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GSD_DB.Where("username = ? AND password = ?", u.Username, u.Username).First(&user).Update("password", utils.MD5V([]byte(newPassword))).Error
	return err, u
}

//@function: SetUserAuthority
//@description: 设置一个用户的权限
//@param: uuid uuid.UUID, authorityId string
//@return: err error

func (userService *UserService) SetUserAuthority(id uint, uuid uuid.UUID, authorityId string) (err error) {
	return err
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: SetUserAuthorities
//@description: 设置一个用户的权限
//@param: id uint, authorityIds []uint
//@return: err error

func (userService *UserService) SetUserAuthorities(id uint, authorityIds []uint) (err error) {
	return global.GSD_DB.Transaction(func(tx *gorm.DB) error {
		TxErr := tx.Delete(&[]system.SysUseAuthority{}, "sys_user_id = ?", id).Error
		if TxErr != nil {
			return TxErr
		}
		var useAuthority []system.SysUseAuthority
		for _, v := range authorityIds {
			useAuthority = append(useAuthority, system.SysUseAuthority{
				SysUserId:               id,
				SysAuthorityAuthorityId: v,
			})
		}
		TxErr = tx.Create(&useAuthority).Error
		if TxErr != nil {
			return TxErr
		}
		TxErr = CasbinServiceApp.UpdateUserAuthority(id, authorityIds)
		if TxErr != nil {
			return TxErr
		}
		return nil
	})
}

//@author: [chenguanglan](https://github.com/sFFbLL)
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
		TxErr = CasbinServiceApp.DeleteUserAuthority(id)
		if TxErr != nil {
			return tx.Rollback().Error
		}
		return nil
	})
}

//@function: SetUserInfo
//@description: 修改用户信息
//@param: reqUser model.SysUser
//@return: err error, user model.SysUser

func (userService *UserService) SetUserInfo(reqUser system.SysUser) (err error, user system.SysUser) {
	err = global.GSD_DB.Updates(&reqUser).Error
	return err, reqUser
}

//@function: GetUserInfo
//@description: 获取用户信息
//@param: uuid uuid.UUID
//@return: err error, user system.SysUser

func (userService *UserService) GetUserInfo(uuid uuid.UUID) (err error, user system.SysUser) {
	var reqUser system.SysUser
	err = global.GSD_DB.Preload("Authorities").First(&reqUser, "uuid = ?", uuid).Error
	return err, reqUser
}

//@function: FindUserById
//@description: 通过id获取用户信息
//@param: id int
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserById(id uint) (err error, user *system.SysUser) {
	var u system.SysUser
	err = global.GSD_DB.Where("`id` = ?", id).First(&u).Error
	return err, &u
}

//@function: FindUserByUuid
//@description: 通过uuid获取用户信息
//@param: uuid string
//@return: err error, user *model.SysUser

func (userService *UserService) FindUserByUuid(uuid string) (err error, user *system.SysUser) {
	var u system.SysUser
	if err = global.GSD_DB.Where("`uuid` = ?", uuid).First(&u).Error; err != nil {
		return errors.New("用户不存在"), &u
	}
	return nil, &u
}
