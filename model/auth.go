package model

import "time"

type Role string

const (
	UserRole  Role = "USER"
	AdminRole Role = "ADMIN"
	DevRole   Role = "DEV"
)

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	Role      Role      `json:"role"`
	Avatar    string    `json:"avatar"`
}

type Creadentials struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Remember bool   `json:"remember"`
}

type Session struct {
	ID        string
	Email     string
	CreatedAt time.Time
	Expiry    time.Time
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

type ForgottPassword struct {
	Email string `json:"email"`
}

type ResetPassword struct {
	Password        string     `json:"password" binding:"required"`
	ConfirmPassword string     `json:"confirmPassword" binding:"required"`
	ExpiresAt       *time.Time `json:"expires_at"`
}
