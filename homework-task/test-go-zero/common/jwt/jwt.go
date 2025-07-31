package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(secret string, expire int64, uid uint64) (string, error) {
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"exp": now + expire,
		"iat": now,
		"uid": uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
