package account

import "github.com/golang-jwt/jwt"

type JwtCustomClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}
