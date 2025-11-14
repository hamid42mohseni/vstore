package models

import "time"

type User struct {
	ID           int       `json:"user_id"`
	Phone        string    `json:"phone"`
	NationalCode string    `json:"national_code"`
	CreateAt     time.Time `json:"create_at"`
}

type RequestOTP struct {
	Phone string `json:"phone" binding:"required"`
}

type RequestCreateUser struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}
