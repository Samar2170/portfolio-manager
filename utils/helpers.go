package utils

import (
	"encoding/json"
	"math/rand"

	"github.com/Samar2170/portfolio-manager/account"
	"github.com/golang-jwt/jwt"
)

func UnwrapToken(token *jwt.Token) (account.JwtCustomClaims, error) {
	var claims account.JwtCustomClaims
	tmp, _ := json.Marshal(token.Claims)
	_ = json.Unmarshal(tmp, &claims)
	return claims, nil
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func CreateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
