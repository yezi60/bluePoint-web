package models

import (
	"github.com/dgrijalva/jwt-go"
)

// CustomClaims payload中存放那些字段信息
type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityID uint
	jwt.StandardClaims
}
