package global

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type NewLogger struct {
	ZapLog *zap.Logger
}

type Request struct {
	RequestId string
}

type Field = zapcore.Field

const RequestID = "requestId"

var ErrorGetRequestId = errors.New("requestId设置失败")

func (log *NewLogger) Info(c *gin.Context, msg string, fields ...Field) {
	requestId, _ := GetRequestId(c)
	fields = append(fields, zap.Any("requestId", requestId))
	log.ZapLog.Info(msg, fields...)
}

func (log *NewLogger) Warn(c *gin.Context, msg string, fields ...Field) {
	requestId, _ := GetRequestId(c)
	fields = append(fields, zap.Any("requestId", requestId))
	log.ZapLog.Info(msg, fields...)
}

func (log *NewLogger) Error(c *gin.Context, msg string, fields ...Field) {
	requestId, _ := GetRequestId(c)
	fields = append(fields, zap.Any("requestId", requestId))
	log.ZapLog.Info(msg, fields...)
}

func (log *NewLogger) DPanic(c *gin.Context, msg string, fields ...Field) {
	requestId, _ := GetRequestId(c)
	fields = append(fields, zap.Any("requestId", requestId))
	log.ZapLog.Info(msg, fields...)
}

func (log *NewLogger) Panic(c *gin.Context, msg string, fields ...Field) {
	requestId, _ := GetRequestId(c)
	fields = append(fields, zap.Any("requestId", requestId))
	log.ZapLog.Info(msg, fields...)
}

func (log *NewLogger) Fatal(c *gin.Context, msg string, fields ...Field) {
	requestId, _ := GetRequestId(c)
	fields = append(fields, zap.Any("requestId", requestId))
	log.ZapLog.Info(msg, fields...)
}

func (log *NewLogger) Debug(c *gin.Context, msg string, fields ...Field) {
	requestId, _ := GetRequestId(c)
	fields = append(fields, zap.Any("requestId", requestId))
	log.ZapLog.Info(msg, fields...)
}

func GetRequestId(c *gin.Context) (*Request, error) {
	get, ok := c.Get(RequestID)
	if !ok {
		err := ErrorGetRequestId
		return nil, err
	}
	id := get.(*Request)
	return id, nil
}
