package jwttoken

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type Claims struct {
	UserId   int    `json:"userId"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

const TokenExpiredTimeInSecond = 3600
const SecretSalt = "Lunara-Secret"
const Issuer = "Lunara-Issue"

var secretKey interface{}

func Generate(username string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:   0,
		Username: username,
		IsAdmin:  true,
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: now.Add(time.Second * TokenExpiredTimeInSecond).Unix(),
			Id:        "",
			IssuedAt:  now.Unix(),
			Issuer:    Issuer,
			NotBefore: now.Unix(),
			Subject:   "",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(GetSecretKeyFromEnv())
}

func Parse(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return GetSecretKeyFromEnv(), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("failed vaild tokenClaims:%#v", tokenClaims)
}

func GetSecretKeyFromEnv() interface{} {
	if secretKey != nil {
		return secretKey
	}
	secretKey = []byte(SecretSalt)
	key := os.Getenv(SecretSalt)
	if key == "" {
		return secretKey
	}
	secretKey = []byte(key)
	return secretKey
}
