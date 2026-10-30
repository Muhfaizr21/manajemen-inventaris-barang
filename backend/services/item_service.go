package services

import (
	"backend/models"
	"backend/repositories"
	"backend/views"
	"context"
	"errors"
	"math"
	"strings"
)

var (
	ErrNegativeStock = errors.New("stok barang tidak boleh bernilai negatif")
	ErrNegativePrice = errors.New("harga barang tidak boleh bernilai negatif")
)

// ItemService mendefinisikan kontrak logika bisnis manajemen inventaris barang
type ItemService interface {
	CreateItem(ctx context.Context, userID uint, req models.CreateItemRequest) (*models.Item, error)
	GetItems(ctx context.Context, filter models.ItemQueryFilter) ([]models.Item, views.PaginationMeta, error)
	GetItemByID(ctx context.Context, id uint) (*models.Item, error)
	UpdateItem(ctx context.Context, id uint, req models.UpdateItemRequest) (*models.Item, error)
	DeleteItem(ctx context.Context, id uint) error
}

type itemService struct {
	itemRepo repositories.ItemRepository
}

// NewItemService menginisialisasi ItemService dengan Dependency Injection
func NewItemService(itemRepo repositories.ItemRepository) ItemService {
	return &itemService{itemRepo: itemRepo}
}

func (s *itemService) CreateItem(ctx context.Context, userID uint, req models.CreateItemRequest) (*models.Item, error) {
	if req.Stock == nil || *req.Stock < 0 {
		return nil, ErrNegativeStock
	}
	if req.Price == nil || *req.Price < 0 {
		return nil, ErrNegativePrice
	}

	item := &models.Item{
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Stock:       *req.Stock,
		Price:       *req.Price,
		CreatedBy:   userID,
	}

	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *itemService) GetItems(ctx context.Context, filter models.ItemQueryFilter) ([]models.Item, views.PaginationMeta, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 10
	} else if filter.Limit > 100 {
		filter.Limit = 100
	}

	items, total, err := s.itemRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, views.PaginationMeta{}, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(filter.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	meta := views.PaginationMeta{
		CurrentPage: filter.Page,
		PerPage:     filter.Limit,
		TotalItems:  total,
		TotalPages:  totalPages,
	}

	return items, meta, nil
}

func (s *itemService) GetItemByID(ctx context.Context, id uint) (*models.Item, error) {
	return s.itemRepo.FindByID(ctx, id)
}

func (s *itemService) UpdateItem(ctx context.Context, id uint, req models.UpdateItemRequest) (*models.Item, error) {
	if req.Stock == nil || *req.Stock < 0 {
		return nil, ErrNegativeStock
	}
	if req.Price == nil || *req.Price < 0 {
		return nil, ErrNegativePrice
	}

	item := &models.Item{
		ID:          id,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		Stock:       *req.Stock,
		Price:       *req.Price,
	}

	if err := s.itemRepo.Update(ctx, item); err != nil {
		return nil, err
	}

	return s.itemRepo.FindByID(ctx, id)
}

func (s *itemService) DeleteItem(ctx context.Context, id uint) error {
	return s.itemRepo.Delete(ctx, id)
}

