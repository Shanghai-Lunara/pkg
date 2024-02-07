package jwttoken

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

type Claims struct {
	UserId   int64  `json:"userId"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

const TokenExpiredTimeInSecond = 3600
const TokenExpiration = "TOKEN_EXPIRATION"
const SecretSalt = "LUNARA_SECRET"
const Issuer = "Lunara-Issue"

var secretKey interface{}
var tokenExpiration int64

func Generate(username string, id int64, isAdmin bool) (string, error) {
	now := time.Now()
	claims := Claims{
		UserId:   id,
		Username: username,
		IsAdmin:  isAdmin,
		StandardClaims: jwt.StandardClaims{
			Audience:  "",
			ExpiresAt: now.Add(time.Second * time.Duration(GetTokenExpirationFromEnv())).Unix(),
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

func GetTokenExpirationFromEnv() int64 {
	if tokenExpiration != 0 {
		return tokenExpiration
	}
	tokenExpiration = TokenExpiredTimeInSecond
	key := os.Getenv(TokenExpiration)
	if key == "" {
		return tokenExpiration
	}
	t, err := strconv.Atoi(key)
	if err != nil {
		return tokenExpiration
	}
	tokenExpiration = int64(t)
	return tokenExpiration
}
