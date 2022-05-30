package util

import (
	"github.com/golang-jwt/jwt"
	"ngb-noti/config"
	"time"
)

type JwtUserClaims struct {
	Id   int  `json:"id"`
	Role bool `json:"role"`
	jwt.StandardClaims
}

func ParseToken(token string) (*JwtUserClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &JwtUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.Jwt.Secret), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*JwtUserClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func GenerateToken(id int, role bool) string {
	claims := &JwtUserClaims{
		id,
		role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(config.C.Jwt.Secret))
	if err != nil {
		return "error"
	}

	return t
}
