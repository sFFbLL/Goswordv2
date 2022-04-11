package utils

//import (
//	"errors"
//
//	"github.com/gin-gonic/gin"
//)
//
//const RequestID = "requestId"
//
//var ErrorGetRequestId = errors.New("requestId设置失败")
//
//type Request struct {
//	RequestId string
//}
//
//func GetRequestId(c *gin.Context) (*Request, error) {
//	get, ok := c.Get(RequestID)
//	if !ok {
//		err := ErrorGetRequestId
//		return nil, err
//	}
//	id := get.(*Request)
//	return id, nil
//}
