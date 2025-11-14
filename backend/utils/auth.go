package utils

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// کلید مخفی JWT (می‌تونی از ENV هم بخونی)
var JwtSecret = []byte("890df93af8e1e008b392c5396b4dfa8e57eeb6eff25538bb7531cee2")

// مدت زمان توکن‌ها
const (
	AccessTokenDuration  = 15 * time.Minute
	RefreshTokenDuration = 7 * 24 * time.Hour
)

// GenerateTokens تولید Access و Refresh Token با user_id
func GenerateTokens(userID uint32) (accessToken string, refreshToken string, err error) {
	// Claims برای Access Token
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(AccessTokenDuration).Unix(),
	}

	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = accessTokenObj.SignedString(JwtSecret)
	if err != nil {
		return
	}

	// Claims برای Refresh Token
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(RefreshTokenDuration).Unix(),
	}

	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = refreshTokenObj.SignedString(JwtSecret)
	if err != nil {
		return
	}

	return
}

func ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return JwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// SetJWTCookies توکن‌ها رو داخل کوکی HTTP قرار می‌دهد
func SetJWTCookies(c *gin.Context, accessToken, refreshToken string) {
	// کوکی Access Token
	c.SetCookie("access_token", accessToken, int(AccessTokenDuration.Seconds()), "/", "", false, true)

	// کوکی Refresh Token
	c.SetCookie("refresh_token", refreshToken, int(RefreshTokenDuration.Seconds()), "/", "", false, true)
}
