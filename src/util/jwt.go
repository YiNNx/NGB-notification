package util

import (
	"github.com/golang-jwt/jwt"
	"ngb-noti/config"
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
