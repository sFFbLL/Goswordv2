package system

import (
	"context"
	"encoding/json"
	"errors"
	"project/global"
	"project/model/system"
	"project/model/system/request"
	"time"

	"go.uber.org/zap"

	"gorm.io/gorm"
)

type JwtService struct {
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: JsonInBlacklist
//@description: 拉黑jwt
//@param: jwtList model.JwtBlacklist
//@return: err error

func (jwtService *JwtService) JsonInBlacklist(jwtList system.JwtBlacklist) (err error) {
	err = global.GSD_DB.Create(&jwtList).Error
	return
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: 判断JWT是否在黑名单内部
//@param: jwt string
//@return: bool

func (jwtService *JwtService) IsBlacklist(jwt string) bool {
	err := global.GSD_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	return !isNotFound
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: GetRedisJWT
//@description: 从redis取jwt
//@param: userName string
//@return: err error, redisJWT string

func (jwtService *JwtService) GetRedisJWT(userName string) (err error, redisJWT string) {
	redisJWT, err = global.GSD_REDIS.Get(context.Background(), userName).Result()
	return err, redisJWT
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: SetRedisJWT
//@description: jwt存入redis并设置过期时间
//@param: jwt string, userName string
//@return: err error

func (jwtService *JwtService) SetRedisJWT(jwt string, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.GSD_CONFIG.JWT.ExpiresTime) * time.Second
	err = global.GSD_REDIS.Set(context.Background(), userName, jwt, timer).Err()
	return err
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: GetRedisClaims
//@description: 从redis取claims
//@param: token string
//@return: err error, request

func (jwtService *JwtService) GetRedisClaims(token string) (error, *request.CustomClaims) {
	claimsStr, err := global.GSD_REDIS.Get(context.Background(), token).Result()
	if err != nil {
		global.GSD_LOG.ZapLog.Error("获取载荷失败", zap.Any("error", err))
		return err, &request.CustomClaims{}
	}
	claims := request.CustomClaims{}
	if err := json.Unmarshal([]byte(claimsStr), &claims); err != nil {
		global.GSD_LOG.ZapLog.Error("载荷反序列化失败", zap.Any("error", err))
		return err, &request.CustomClaims{}
	}
	return err, &claims
}

//@author: [chenguanglan](https://github.com/sFFbLL)
//@function: SetRedisClaims
//@description: claims存入redis并设置过期时间
//@param: claims request.CustomClaims string, token string
//@return: err error

func (jwtService *JwtService) SetRedisClaims(calims request.CustomClaims, token string) (err error) {
	// 此处过期时间等于jwt过期时间
	timer := time.Duration(global.GSD_CONFIG.JWT.ExpiresTime) * time.Second
	err = global.GSD_REDIS.Set(context.Background(), token, calims, timer).Err()
	return err
}
