package utils

import (
	"encoding/json"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/golang-jwt/jwt"
)

func UnwrapToken(token *jwt.Token) (account.JwtCustomClaims, error) {
	var claims account.JwtCustomClaims
	tmp, _ := json.Marshal(token.Claims)
	_ = json.Unmarshal(tmp, &claims)
	return claims, nil
}
