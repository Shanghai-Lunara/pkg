package signaturetoken

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"strconv"
	"time"
)

type Claims struct {
	Id       int32  `json:"id"`
	Uid      string `json:"uid"`
	Snapshot string `json:"snapshot"`
	jwt.StandardClaims
}

const TokenExpiredTimeInSecond = 3600 * 5
const TokenExpiration = "TOKEN_EXPIRATION"
const SecretSalt = "LUNARA_SECRET"
const Issuer = "Lunara-Issuer"

var secretKey interface{}
var tokenExpiration int64

func Generate(id int32, uid, snapshot string) (string, error) {
	now := time.Now()
	claims := Claims{
		Id:       id,
		Uid:      uid,
		Snapshot: snapshot,
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
