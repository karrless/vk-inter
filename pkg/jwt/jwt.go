package jwt

import (
	"github.com/golang-jwt/jwt/v5"
)

// NewAccessToken returns a new access token
func NewAccessToken(secret string) string {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = jwt.MapClaims{}
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}
