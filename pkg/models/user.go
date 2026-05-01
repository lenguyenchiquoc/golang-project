package models

import "time"

type User struct {
	ID           string    `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email string    `json:"email" db:"email"`
	Password string `json:"password" binding:"required"`
	RePassword string `json:"repassword" binding:"required"`
}

type RegisterResponse struct {
	UserID string `json:"userid" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email string    `json:"email" db:"email"`
	Message string `json:"message" binding:"required"`
}