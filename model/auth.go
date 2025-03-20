package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Password  string    `json:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
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

func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
