package controllers

import (
	"net/http"
	"regexp"
	"strings"
	"time"
	"vstore/backend/database"
	"vstore/backend/models"
	"vstore/backend/utils"

	"github.com/gin-gonic/gin"
)

func SendCode(c *gin.Context) {
	var req models.RequestOTP

	// درست نفرستاده شماره رو
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شماره تلفن اشتباه وارد شده"})
		return
	}

	// شماره رو اشتباه زده
	regex := regexp.MustCompile(`^09[0-9]{9}$`)
	if !regex.MatchString(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شماره تلفن اشتباه وارد شده"})
		return
	}

	if cached, ok, err := utils.GetOTPFromCache(req.Phone); err == nil && ok {
		c.JSON(http.StatusOK, gin.H{
			"message":     "کد ارسال شده",
			"phone":       req.Phone,
			"code":        cached, // اینو یادم باشه کامت کنم
			"from_cached": "true",
		})
		return
	}

	// دو دقیقه تایم اعتبار برای این کد
	expiresAt := time.Now().Add(2 * time.Minute)
	rows, err := database.DB.Query("CALL insert_otp(?,?)", req.Phone, expiresAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در سیستم لطفا دوباره تلاش کنید", "details:": err.Error()})
		return
	}

	defer rows.Close()

	var code string
	var Phone string
	if rows.Next() {
		if err := rows.Scan(&code, &Phone); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در خواندن نتیجه", "details": err.Error()})
			return
		}
	}

	// گزاشتن توی کش
	utils.SetOTPToCache(req.Phone, code, 2*time.Minute)

	c.JSON(http.StatusOK, gin.H{
		"message": "کد برای شما ارسال شد",
		"phone":   Phone,
		"code":    code, // اینو یادم باشه کامنت کنم
	})
}

func CreateUserOrLogin(c *gin.Context) {
	var req models.RequestCreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شماره یا کد به درستی ارسال نشده"})
		return
	}

	// اعتبارسنجی شماره تلفن
	phoneRegex := regexp.MustCompile(`^09[0-9]{9}$`)
	if !phoneRegex.MatchString(req.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "شماره تلفن اشتباه وارد شده"})
		return
	}

	// اعتبارسنجی کد OTP
	codeRegex := regexp.MustCompile(`^[0-9]{5}$`)
	if !codeRegex.MatchString(req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "کد اشتباه ارسال شده"})
		return
	}

	rows, err := database.DB.Query("CALL create_user(?, ?)", req.Phone, req.Code)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid phone number format") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "فرمت شماره تلفن نامعتبر است"})
			return
		}
		if strings.Contains(err.Error(), "Invalid or used OTP code") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "کد OTP نامعتبر یا قبلاً استفاده شده است"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در سیستم لطفاً دوباره تلاش کنید", "details": err.Error()})
		return
	}
	defer rows.Close()

	// خواندن user_id و code
	var userId uint32
	var code string
	if rows.Next() {
		if err := rows.Scan(&userId, &code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در خواندن نتیجه", "details": err.Error()})
			return
		}
	}

	// تولید توکن‌ها
	accessToken, refreshToken, err := utils.GenerateTokens(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ساخت توکن", "details": err.Error()})
		return
	}

	// ذخیره refresh token در دیتابیس
	expiresAt := time.Now().Add(utils.RefreshTokenDuration)
	_, err = database.DB.Exec("CALL insert_refresh_token(?, ?, ?)", userId, refreshToken, expiresAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در ذخیره توکن", "details": err.Error()})
		return
	}

	// ذخیره توکن‌ها در کوکی
	utils.SetJWTCookies(c, accessToken, refreshToken)

	// پاسخ نهایی
	c.JSON(http.StatusOK, gin.H{
		"message": "کاربر با موفقیت ساخته شد",
		"user_id": userId,
		"code":    code,
	})
}

func RefreshToken(c *gin.Context) {
	// گرفتن refresh token از کوکی
	rt, err := c.Cookie("refresh_token")
	if err != nil || rt == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token موجود نیست"})
		return
	}

	// فراخوانی پروسیجر validate_refresh_token
	var dbUserId uint32
	err = database.DB.QueryRow("CALL validate_refresh_token(?)", rt).Scan(&dbUserId)
	if err != nil {
		if strings.Contains(err.Error(), "Refresh Token Expired") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "کد OTP نامعتبر یا قبلاً استفاده شده است"})
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token منقضی یا نامعتبر است"})
		return
	}

	// تولید Access Token جدید
	accessToken, _, err := utils.GenerateTokens(dbUserId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "خطا در تولید Access Token"})
		return
	}

	// ذخیره توکن‌ها در کوکی
	utils.SetJWTCookies(c, accessToken, rt)

	c.JSON(http.StatusOK, gin.H{
		"message":       "Access Token جدید صادر شد",
		"user_id":       dbUserId,
		"access_token":  accessToken,
		"refresh_token": rt,
	})
}
