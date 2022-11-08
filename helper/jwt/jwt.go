package jwt

import (
	"errors"
	"gin-skeleton/model/admin/system"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

var (
	signKey = viper.GetString("JwtToken.SignKey")
	issuer  = viper.GetString("Server.Name")
)

// Claims 声明
type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	jwt.StandardClaims
}

// GenerateJwtToken 生成Token
func GenerateJwtToken(account system.Account, hour int) (string, error) {
	claims := Claims{
		account.ID,
		account.Username,
		account.Nickname,
		jwt.StandardClaims{
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(hour)).Unix(),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(signKey))
}

// ParseJwtToken 解析Token
func ParseJwtToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token is invalid")
}
