package request

import (
	"project/model/system"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// Custom claims structure
type CustomClaims struct {
	UUID       uuid.UUID
	ID         uint
	DeptId     uint
	Username   string
	NickName   string
	Authority  []system.SysAuthority
	BufferTime int64
	jwt.StandardClaims
}
