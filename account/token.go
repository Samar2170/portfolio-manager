package account

import (
	"encoding/json"

	"github.com/golang-jwt/jwt"
)

type JwtCustomClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func UnwrapToken(token *jwt.Token) (JwtCustomClaims, error) {
	var claims JwtCustomClaims
	tmp, _ := json.Marshal(token.Claims)
	_ = json.Unmarshal(tmp, &claims)
	return claims, nil
}
