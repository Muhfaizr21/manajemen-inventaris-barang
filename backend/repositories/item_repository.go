package repositories

import (
	"backend/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrItemNotFound      = errors.New("barang tidak ditemukan")
	ErrInsufficientStock = errors.New("stok barang tidak mencukupi (tidak boleh negatif)")
)

// ItemRepository mendefinisikan kontrak akses data Item (Dependency Inversion Principle)
type ItemRepository interface {
	Create(ctx context.Context, item *models.Item) error
	FindAll(ctx context.Context, filter models.ItemQueryFilter) ([]models.Item, int64, error)
	FindByID(ctx context.Context, id uint) (*models.Item, error)
	Update(ctx context.Context, item *models.Item) error
	Delete(ctx context.Context, id uint) error
	AdjustStockAtomically(ctx context.Context, id uint, delta int) (*models.Item, error)
}

type itemRepository struct {
	db *gorm.DB
}

// NewItemRepository membuat instance baru ItemRepository
func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(ctx context.Context, item *models.Item) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *itemRepository) FindAll(ctx context.Context, filter models.ItemQueryFilter) ([]models.Item, int64, error) {
	var items []models.Item
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Item{})

	// Fitur Search berdasarkan nama atau deskripsi (case-insensitive di PostgreSQL dengan ILIKE)
	if search := strings.TrimSpace(filter.Search); search != "" {
		likePattern := fmt.Sprintf("%%%s%%", search)
		query = query.Where("name ILIKE ? OR description ILIKE ?", likePattern, likePattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	if err := query.Order("id DESC").Limit(filter.Limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *itemRepository) FindByID(ctx context.Context, id uint) (*models.Item, error) {
	var item models.Item
	err := r.db.WithContext(ctx).First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *itemRepository) Update(ctx context.Context, item *models.Item) error {
	// Menggunakan transaksi dengan SELECT ... FOR UPDATE (Pessimistic Locking)
	// untuk mencegah race condition saat dua request mengedit barang secara bersamaan
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.Item
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, item.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrItemNotFound
			}
			return err
		}

		if item.Stock < 0 {
			return ErrInsufficientStock
		}

		return tx.Model(&existing).Updates(map[string]interface{}{
			"name":        item.Name,
			"description": item.Description,
			"stock":       item.Stock,
			"price":       item.Price,
		}).Error
	})
}

func (r *itemRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrItemNotFound
	}
	return nil
}

// AdjustStockAtomically mendemonstrasikan penanganan konkurensi menggunakan Pessimistic Locking (SELECT FOR UPDATE)
// dan Atomic SQL Guard agar stok tidak pernah bernilai negatif (-1) saat diakses secara paralel.
func (r *itemRepository) AdjustStockAtomically(ctx context.Context, id uint, delta int) (*models.Item, error) {
	var updatedItem models.Item

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Kunci baris (Row-Level Exclusive Lock: SELECT ... FOR UPDATE)
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&updatedItem, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrItemNotFound
			}
			return err
		}

		// 2. Validasi apakah perubahan stok menyebabkan nilai negatif
		if updatedItem.Stock+delta < 0 {
			return ErrInsufficientStock
		}

		// 3. Atomic Update dengan kondisi guard `stock + delta >= 0`
		res := tx.Model(&updatedItem).
			Where("id = ? AND stock + ? >= 0", id, delta).
			Update("stock", gorm.Expr("stock + ?", delta))
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrInsufficientStock
		}

		updatedItem.Stock += delta
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &updatedItem, nil
}

