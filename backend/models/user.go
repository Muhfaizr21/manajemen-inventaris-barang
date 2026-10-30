package models

import "time"

// User merepresentasikan entitas pengguna di tabel `users`
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"` // Disembunyikan dari serialisasi JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RegisterRequest merupakan DTO untuk validasi input registrasi user
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email,max=100"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// LoginRequest merupakan DTO untuk validasi input login user
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse merupakan DTO response saat registrasi atau login berhasil
type AuthResponse struct {
	Token string `json:"token,omitempty"`
	User  User   `json:"user"`
}

