package accounts

import (
	"errors"
	"ses-go/config"
	"ses-go/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken 토큰 생성
func GenerateToken(userId uint) (string, error) {
	// jwt 토큰 생성
	jti := uuid.Must(uuid.NewV7()).String()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * 365 * 100).Unix(),
		"jti": jti,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"sub": userId,
		"iss": "ses-go",
	})
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GetToken 토큰 저장
func GetToken(userId uint) (string, error) {
	db := config.GetDB()
	token, err := GenerateToken(userId)
	if err != nil {
		return "", err
	}
	// 사용자당 최대 10개까지만 생성
	var count int64
	if err := db.Model(&models.UserToken{}).Where("user_id = ?", userId).Count(&count).Error; err != nil {
		return "", err
	}
	if count >= 10 {
		return "", errors.New("토큰은 사용자당 최대 10개까지만 생성할 수 있습니다")
	}
	// 토큰 생성
	if err := db.Create(&models.UserToken{
		UserId: userId,
		Token:  token,
	}).Error; err != nil {
		return "", err
	}
	return token, nil
}

// ValidateToken 토큰 검증
func ValidateToken(token string) (uint, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		return 0, err
	}
	db := config.GetDB()
	if err := db.Where("user_id = ? AND token = ?", claims["sub"], token).First(&models.UserToken{}).Error; err != nil {
		return 0, err
	}
	return uint(claims["sub"].(float64)), nil
}

// DeleteToken 토큰 삭제
func DeleteToken(userId uint) error {
	db := config.GetDB()
	return db.Where("user_id = ?", userId).Delete(&models.UserToken{}).Error
}
