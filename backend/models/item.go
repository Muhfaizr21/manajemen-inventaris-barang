package models

import (
	"time"

	"gorm.io/gorm"
)

// Item merepresentasikan entitas barang inventaris di tabel `items`
// Menggunakan constraint CHECK pada level database untuk menjamin stock >= 0 dan price >= 0
type Item struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:150;not null;index" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Stock       int            `gorm:"not null;default:0;check:stock >= 0" json:"stock"`
	Price       float64        `gorm:"type:numeric(15,2);not null;default:0;check:price >= 0" json:"price"`
	CreatedBy   uint           `gorm:"not null;index" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateItemRequest merupakan DTO untuk validasi pembuatan barang baru
// Menggunakan pointer (*int, *float64) agar nilai 0 tetap valid namun nilai negatif (-1) atau kosong ditolak oleh validator
type CreateItemRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=150"`
	Description string   `json:"description" binding:"max=1000"`
	Stock       *int     `json:"stock" binding:"required,gte=0"`
	Price       *float64 `json:"price" binding:"required,gte=0"`
}

// UpdateItemRequest merupakan DTO untuk validasi pembaruan data barang (PUT /items/:id)
type UpdateItemRequest struct {
	Name        string   `json:"name" binding:"required,min=2,max=150"`
	Description string   `json:"description" binding:"max=1000"`
	Stock       *int     `json:"stock" binding:"required,gte=0"`
	Price       *float64 `json:"price" binding:"required,gte=0"`
}

// ItemQueryFilter menampung query parameter untuk pencarian dan paginasi pada GET /items
type ItemQueryFilter struct {
	Search string `form:"search"`
	Page   int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit  int    `form:"limit,default=10" binding:"omitempty,min=1,max=100"`
}

