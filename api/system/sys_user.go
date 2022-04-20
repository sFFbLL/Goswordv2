package system

import (
	"project/global"
	"project/middleware"
	"project/model/common/request"
	"project/model/common/response"
	"project/model/system"
	systemReq "project/model/system/request"
	systemRes "project/model/system/response"
	"project/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body systemReq.Login true "用户名, 密码, 验证码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"登陆成功"}"
// @Router /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	_ = c.ShouldBindJSON(&l)

	if err := utils.Verify(l, utils.LoginVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//if !store.Verify(l.CaptchaId, l.Captcha, true) {
	//	response.FailWithMessage("验证码错误", c)
	//	return
	//}
	u := &system.SysUser{Username: l.Username, Password: l.Password}
	if err, user := userService.Login(u); err != nil {
		global.GSD_LOG.Error(c, "登陆失败! 用户名不存在或者密码错误!", zap.Any("err", err))
		response.FailWithMessage("用户名不存在或者密码错误", c)
	} else {
		b.tokenNext(c, *user)
	}
}

// 登录以后签发jwt
func (b *BaseApi) tokenNext(c *gin.Context, user system.SysUser) {
	j := &middleware.JWT{SigningKey: []byte(global.GSD_CONFIG.JWT.SigningKey)} // 唯一签名
	claims := systemReq.CustomClaims{
		UUID:       user.UUID,
		ID:         user.ID,
		Username:   user.Username,
		BufferTime: global.GSD_CONFIG.JWT.BufferTime, // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                              // 签名生效时间
			ExpiresAt: time.Now().Unix() + global.GSD_CONFIG.JWT.ExpiresTime, // 过期时间 7天  配置文件
			Issuer:    "gsdPlus",                                             // 签名的发行者
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GSD_LOG.Error(c, "获取token失败!", zap.Any("err", err))
		response.FailWithMessage("获取token失败", c)
		return
	}
	var authorityIds []uint
	for _, authority := range user.Authorities {
		authorityIds = append(authorityIds, authority.AuthorityId)
	}
	userCache := systemReq.UserCache{
		ID:          user.ID,
		UUID:        user.UUID.String(),
		Authority:   user.Authorities,
		AuthorityId: authorityIds,
		DeptId:      user.DeptId,
	}
	_ = jwtService.SetRedisUserInfo(user.UUID.String(), userCache)
	if !global.GSD_CONFIG.System.UseMultipoint {
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
		return
	}
	if err, jwtStr := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			global.GSD_LOG.Error(c, "设置登录状态失败!", zap.Any("err", err))
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	} else if err != nil {
		global.GSD_LOG.Error(c, "设置登录状态失败!", zap.Any("err", err))
		response.FailWithMessage("设置登录状态失败", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JoinInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("jwt作废失败", c)
			return
		}
		if err := jwtService.SetRedisJWT(token, user.Username); err != nil {
			response.FailWithMessage("设置登录状态失败", c)
			return
		}
		//设置用户缓存
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.StandardClaims.ExpiresAt * 1000,
		}, "登录成功", c)
	}
}

// @Tags SysUser
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body systemReq.Register true "用户名, 昵称, 密码, 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"注册成功"}"
// @Router /user/register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	_ = c.ShouldBindJSON(&r)
	if err := utils.Verify(r, utils.RegisterVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		if err, authority := authorityService.GetAuthorityBasicInfo(system.SysAuthority{
			AuthorityId: v,
		}); err != nil {
			global.GSD_LOG.Error(c, "注册失败, 角色不存在!")
			response.FailWithMessage("注册失败, 角色不存在!", c)
		} else {
			authorities = append(authorities, authority)
		}
	}
	curUser := utils.GetUser(c)
	user := &system.SysUser{GSD_MODEL: global.GSD_MODEL{CreateBy: curUser.ID, UpdateBy: curUser.ID}, Username: r.Username, NickName: r.NickName, Password: r.Password, Authorities: authorities, DeptId: r.DeptId}
	//数据权限校验
	canDo := dataScope.CanDoToTargetUser(curUser, []*system.SysUser{user})
	if !canDo {
		global.GSD_LOG.Error(c, "注册失败, 无权注册该用户!")
		response.FailWithMessage("注册失败, 无权注册该用户!", c)
		return
	}
	err, userReturn := userService.Register(*user, r.AuthorityIds)
	if err != nil {
		global.GSD_LOG.Error(c, "注册失败!", zap.Any("err", err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册失败", c)
	} else {
		response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "注册成功", c)
	}
}

// @Tags SysUser
// @Summary 删除用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "用户ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /user/deleteUser [delete]
func (b *BaseApi) DeleteUser(c *gin.Context) {
	var reqId request.GetById
	_ = c.ShouldBindJSON(&reqId)
	if err := utils.Verify(reqId, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jwtId := utils.GetUserID(c)
	if jwtId == reqId.ID {
		response.FailWithMessage("拒绝自杀", c)
		return
	}
	//获取需要删除用户的信息
	curUser := utils.GetUser(c)
	err, deleteUser := userService.FindUserById(reqId.ID)
	if err != nil {
		global.GSD_LOG.Error(c, "删除失败, 该用户不存在!")
		response.FailWithMessage("删除失败, 该用户不存在!", c)
		return
	}
	//数据权限校验
	canDo := dataScope.CanDoToTargetUser(curUser, []*system.SysUser{deleteUser})
	if !canDo {
		global.GSD_LOG.Error(c, "删除失败, 无权删除该用户!")
		response.FailWithMessage("删除失败, 无权删除该用户!", c)
		return
	}
	//删除用户
	if err := userService.DeleteUser(reqId.ID); err != nil {
		global.GSD_LOG.Error(c, "删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		//删除用户缓存
		_ = jwtService.DelRedisUserInfo(deleteUser.UUID.String())
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags SysUser
// @Summary 设置用户角色
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body systemReq.SetUserAuthorities true "用户UUID, 角色ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setUserAuthorities [post]
func (b *BaseApi) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities
	_ = c.ShouldBindJSON(&sua)
	if err := utils.Verify(sua, utils.SetUserAuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	curUser := utils.GetUser(c)
	err, updateUser := userService.FindUserById(sua.ID)
	if err != nil {
		global.GSD_LOG.Error(c, "修改失败!", zap.Any("err", err))
		response.FailWithMessage("操作用户不存在", c)
		return
	}
	//校验数据权限
	canDo := dataScope.CanDoToTargetUser(curUser, []*system.SysUser{updateUser})
	if !canDo {
		global.GSD_LOG.Error(c, "修改失败, 无权修改该用户!")
		response.FailWithMessage("操作失败, 无权操作该用户!", c)
		return
	}
	var updateAuthorities []system.SysAuthority
	for _, authorityId := range sua.AuthorityIds {
		if err, authority := authorityService.GetAuthorityBasicInfo(system.SysAuthority{AuthorityId: authorityId}); err != nil {
			global.GSD_LOG.Error(c, "设置角色不存在!")
			response.FailWithMessage("设置角色不存在!", c)
			return
		} else {
			updateAuthorities = append(updateAuthorities, authority)
		}
	}
	//校验目标level是否垂直越权
	if dataScope.GetMaxLevel(updateAuthorities) < dataScope.GetMaxLevel(curUser.Authority) {
		global.GSD_LOG.Error(c, "设置角色级别高于当前用户级别!")
		response.FailWithMessage("设置角色级别高于当前用户级别!", c)
		return
	}
	if err := userService.SetUserAuthorities(*updateUser, sua.AuthorityIds); err != nil {
		global.GSD_LOG.Error(c, "修改失败!", zap.Any("err", err))
		response.FailWithMessage("修改失败", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// @Tags SysUser
// @Summary 分页获取用户列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserList [post]
func (b *BaseApi) GetUserList(c *gin.Context) {
	var pageInfo request.PageInfo
	_ = c.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	curUser := utils.GetUser(c)
	deptId, isAll := dataScope.GetDataScope(curUser)
	if err, list, total := userService.GetUserInfoList(pageInfo, deptId, isAll); err != nil {
		global.GSD_LOG.Error(c, "获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// @Tags SysUser
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body systemReq.ChangePasswordStruct true "用户名, 原密码, 新密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/changePassword [put]
func (b *BaseApi) UpdatePassword(c *gin.Context) {
	var user systemReq.ChangePasswordStruct
	_ = c.ShouldBindJSON(&user)
	if err := utils.Verify(user, utils.ChangePasswordVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	u := &system.SysUser{
		Username: user.Username,
		Password: user.Password,
	}
	if err, _ := userService.UpdatePassword(u, user.NewPassword); err != nil {
		global.GSD_LOG.Error(c, "修改失败", zap.Error(err))
		response.FailWithMessage("修改失败， 原密码与当前账户不符", c)
	} else {
		response.OkWithMessage("修改成功", c)
	}
}

// @Tags SysUser
// @Summary 获取用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /user/getUserInfo [get]
func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	if err, userInfo := userService.GetUserInfo(uuid); err != nil {
		global.GSD_LOG.Error(c, "获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败", c)
		return
	} else {
		response.OkWithDetailed(gin.H{"userInfo": userInfo}, "获取用户信息成功", c)
	}
}

// @Tags SysUser
// @Summary 设置用户信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body system.SysUser true "ID, 用户名, 昵称, 头像链接"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"设置成功"}"
// @Router /user/setUserInfo [put]
func (b *BaseApi) SetUserInfo(c *gin.Context) {
	var user system.SysUser
	_ = c.ShouldBindJSON(&user)
	if err := utils.Verify(user, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user.Username = ""
	user.Password = ""
	curUser := utils.GetUser(c)
	user.CreateBy = curUser.ID
	if err, sysUser := userService.SetUserInfo(user); err != nil {
		global.GSD_LOG.Error(c, "设置失败", zap.Error(err))
		response.FailWithMessage("设置失败", c)
	} else {
		response.OkWithDetailed(gin.H{"userinfo": sysUser}, "设置成功", c)
	}
}

// @Tags SysUser
// @Summary 导入用户Excel文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "导入Excel文件"
// @Success 200 {object} response.Response{msg=string} "导入Excel文件"
// @Router /user/importExcel [post]
func (b *BaseApi) ImportExcel(c *gin.Context) {
	_, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GSD_LOG.Error(c, "接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	_ = c.SaveUploadedFile(header, global.GSD_CONFIG.Excel.Dir+"ExcelImport.xlsx")
	response.OkWithMessage("导入成功", c)
}

// @Tags SysUser
// @Summary 加载Excel数据
// @Security ApiKeyAuth
// @Produce  application/json
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "加载Excel数据,返回包括列表,总数,页码,每页数量"
// @Router /user/loadExcel [get]
func (b *BaseApi) LoadExcel(c *gin.Context) {
	list, err := userService.ParseExcelToDataList()
	if err != nil {
		global.GSD_LOG.Error(c, "加载数据失败", zap.Error(err))
		response.FailWithMessage("加载数据失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    int64(len(list)),
		Page:     1,
		PageSize: 999,
	}, "加载数据成功", c)
}

// @Tags SysUser
// @Summary 导出Excel
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/octet-stream
// @Param data body example.ExcelInfo true "导出Excel文件信息"
// @Success 200
// @Router /user/exportExcel [post]
func (b *BaseApi) ExportExcel(c *gin.Context) {
	var excelInfo system.ExcelInfo
	_ = c.ShouldBindJSON(&excelInfo)
	filePath := global.GSD_CONFIG.Excel.Dir + excelInfo.FileName
	if err := userService.ParseDataListToExcel(excelInfo.InfoList, filePath); err != nil {
		global.GSD_LOG.Error(c, "导出excel失败", zap.Error(err))
		response.FailWithMessage("导出excel失败", c)
		return
	}
	c.Writer.Header().Add("success", "true")
	c.File(filePath)

}

// @Tags SysUser
// @Summary 下载模板
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/octet-stream
// @Param data body example.ExcelInfo true "下载模板信息"
// @Success 200
// @Router /user/downloadTemplate [post]
func (b *BaseApi) DownloadTemplate(c *gin.Context) {
	name := c.Query("fileName")
	filePath := global.GSD_CONFIG.Excel.Dir + name
	if err := userService.Template(filePath); err != nil {
		global.GSD_LOG.Error(c, "模板下载失败", zap.Error(err))
		response.FailWithMessage("模板下载失败", c)
		return
	}
	c.Writer.Header().Add("success", "true")
	c.File(filePath)
}
