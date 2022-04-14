package utils

import (
	"project/global"
	"project/model/system"
	systemReq "project/model/system/request"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// 从Gin的Context中获取从jwt解析出来的用户信息
func GetUser(c *gin.Context) *systemReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		global.GSD_LOG.Error(c, "从Gin的Context中获取从jwt解析出来的用户失败, 请检查路由是否使用jwt中间件!")
		return nil
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse
	}
}

// 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		global.GSD_LOG.Error(c, "从Gin的Context中获取从jwt解析出来的用户ID失败, 请检查路由是否使用jwt中间件!")
		return 0
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.ID
	}
}

// 从Gin的Context中获取从jwt解析出来的用户UUID
func GetUserUuid(c *gin.Context) uuid.UUID {
	if claims, exists := c.Get("claims"); !exists {
		global.GSD_LOG.Error(c, "从Gin的Context中获取从jwt解析出来的用户UUID失败, 请检查路由是否使用jwt中间件!")
		return uuid.UUID{}
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.UUID
	}
}

// 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserAuthorityId(c *gin.Context) []system.SysAuthority {
	if claims, exists := c.Get("claims"); !exists {
		global.GSD_LOG.Error(c, "从Gin的Context中获取从jwt解析出来的用户角色失败, 请检查路由是否使用jwt中间件!")
		return nil
	} else {
		waitUse := claims.(*systemReq.CustomClaims)
		return waitUse.Authority
	}
}
