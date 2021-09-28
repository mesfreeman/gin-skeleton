package tool

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

// 密钥
var signKey = viper.GetString("JwtToken.SignKey")

// GenerateJwtToken 生成 Token
func GenerateJwtToken(ID int, hour int) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
	}).SignedString([]byte(signKey))
	if err != nil {
		return "", err
	}

	return token, nil
}

// ParseJwtToken 解析 Token
func ParseJwtToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return 0, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["id"].(float64)), nil
	}
	return 0, errors.New("token expire or invalid")
}
