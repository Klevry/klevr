package common

import (
	"fmt"

	//"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	CustomClaims map[string]string `json:"cusotm,omitempty"`
	jwt.StandardClaims
}

type JWTHelper struct {
	Secret         []byte
	Claims         map[string]string
	ExpirationTime int64
}

func NewJWTHelper(secret []byte) *JWTHelper {
	return &JWTHelper{Secret: secret, Claims: make(map[string]string)}
}

func (j *JWTHelper) AddClaims(key, value string) *JWTHelper {
	j.Claims[key] = value
	return j
}

func (j *JWTHelper) SetExpirationTime(t int64) *JWTHelper {
	j.ExpirationTime = t
	return j
}

func (j *JWTHelper) GenToken() (string, error) {
	claims := &Claims{
		CustomClaims: j.Claims,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: j.ExpirationTime,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.Secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func (j *JWTHelper) GetClaims(tokenString string) (map[string]string, error) {
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.Secret, nil
	})

	if !tkn.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	if err != nil {
		return nil, err
	}

	return claims.CustomClaims, nil
}
