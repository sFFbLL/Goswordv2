package utils

import (
	"project/model/system"
)

type DataScope struct {
	Authorities []system.SysAuthority `json:"authorities"`
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: GetDataScope
//@description: 获取数据权限
//@param:
//@return: []uint 部门id
func (s DataScope) GetDataScope() []uint {
	return nil
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: GetMaxLevel
//@description: 获取最角色最大等级
//@param:
//@return: int 所有角色中的最高等级
func (s DataScope) GetMaxLevel() int {
	return 0
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: CanDo
//@description: 是否有权操作该数据
//@param: users []system.SysUser 操作用户对象
//@return: bool 是否有权操作对象
func (s DataScope) CanDo(users []system.SysUser) bool {
	return false
}
