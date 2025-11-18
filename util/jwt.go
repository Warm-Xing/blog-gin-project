package util

import (
	"blog-gin-project/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(config.LoadConfig().JWTSecret)
var jwtExpireHours = config.LoadConfig().JWTExpireHours

// GenerateJWT 生成JWT令牌
func GenerateJWT(userID uint, username string) (string, error) {
	// 设置令牌过期时间
	expirationTime := time.Now().Add(time.Hour * time.Duration(jwtExpireHours))

	// 创建claims
	claims := &jwt.MapClaims{
		"id":       userID,
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString(jwtSecret)

	return tokenString, err
}

// ParseJWT 解析JWT令牌
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证令牌并返回claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
