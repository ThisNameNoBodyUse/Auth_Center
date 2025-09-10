package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"auth-center/config"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	UserID  uint     `json:"user_id"`
	AppID   string   `json:"app_id"`
	Roles   []uint   `json:"roles"`
	JTI     string   `json:"jti"` // JWT ID
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(userID uint, appID string, roles []uint) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		AppID:  appID,
		Roles:  roles,
		JTI:    generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetConfig().JWT.TTL) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-center",
			Subject:   "access-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWT.SecretKey))
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userID uint, appID string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		AppID:  appID,
		JTI:    generateJTI(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.GetConfig().JWT.RefreshTTL) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-center",
			Subject:   "refresh-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().JWT.RefreshSecretKey))
}

// ParseAccessToken 解析访问令牌
func ParseAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.GetConfig().JWT.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.GetConfig().JWT.RefreshSecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证令牌（通用方法）
func ValidateToken(tokenString string, isAccessToken bool) (*JWTClaims, error) {
	if isAccessToken {
		return ParseAccessToken(tokenString)
	}
	return ParseRefreshToken(tokenString)
}

// generateJTI 生成JWT ID
func generateJTI() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
