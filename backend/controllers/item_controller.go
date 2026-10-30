package controllers

import (
	"backend/middlewares"
	"backend/models"
	"backend/repositories"
	"backend/services"
	"backend/views"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ItemController menangani HTTP request untuk manajemen inventaris barang
type ItemController struct {
	itemService services.ItemService
}

// NewItemController membuat instance baru ItemController
func NewItemController(itemService services.ItemService) *ItemController {
	return &ItemController{itemService: itemService}
}

// CreateItem menangani POST /items
func (ctrl *ItemController) CreateItem(c *gin.Context) {
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		views.JSONError(c, http.StatusBadRequest, "Validasi data barang gagal", views.FormatValidationError(err))
		return
	}

	userIDVal, _ := c.Get(middlewares.ContextUserIDKey)
	userID, _ := userIDVal.(uint)

	item, err := ctrl.itemService.CreateItem(c.Request.Context(), userID, req)
	if err != nil {
		if errors.Is(err, services.ErrNegativeStock) || errors.Is(err, services.ErrNegativePrice) {
			views.JSONError(c, http.StatusBadRequest, "Validasi data barang gagal", err.Error())
			return
		}
		views.JSONError(c, http.StatusInternalServerError, "Gagal menambahkan barang", err.Error())
		return
	}

	views.JSONSuccess(c, http.StatusCreated, "Barang berhasil ditambahkan", item)
}

// GetItems menangani GET /items (dengan query param search, page, limit)
func (ctrl *ItemController) GetItems(c *gin.Context) {
	var filter models.ItemQueryFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		views.JSONError(c, http.StatusBadRequest, "Query parameter tidak valid", views.FormatValidationError(err))
		return
	}

	items, meta, err := ctrl.itemService.GetItems(c.Request.Context(), filter)
	if err != nil {
		views.JSONError(c, http.StatusInternalServerError, "Gagal mengambil daftar barang", err.Error())
		return
	}

	views.JSONPaginated(c, http.StatusOK, "Berhasil mengambil daftar barang", items, meta)
}

// GetItemByID menangani GET /items/:id
func (ctrl *ItemController) GetItemByID(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		views.JSONError(c, http.StatusBadRequest, "ID barang tidak valid", err.Error())
		return
	}

	item, err := ctrl.itemService.GetItemByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrItemNotFound) {
			views.JSONError(c, http.StatusNotFound, "Data tidak ditemukan", err.Error())
			return
		}
		views.JSONError(c, http.StatusInternalServerError, "Gagal mengambil detail barang", err.Error())
		return
	}

	views.JSONSuccess(c, http.StatusOK, "Berhasil mengambil detail barang", item)
}

// UpdateItem menangani PUT /items/:id
func (ctrl *ItemController) UpdateItem(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		views.JSONError(c, http.StatusBadRequest, "ID barang tidak valid", err.Error())
		return
	}

	var req models.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		views.JSONError(c, http.StatusBadRequest, "Validasi data barang gagal", views.FormatValidationError(err))
		return
	}

	item, err := ctrl.itemService.UpdateItem(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrItemNotFound):
			views.JSONError(c, http.StatusNotFound, "Gagal memperbarui barang", err.Error())
		case errors.Is(err, services.ErrNegativeStock), errors.Is(err, services.ErrNegativePrice), errors.Is(err, repositories.ErrInsufficientStock):
			views.JSONError(c, http.StatusBadRequest, "Validasi stok/harga gagal", err.Error())
		default:
			views.JSONError(c, http.StatusInternalServerError, "Gagal memperbarui barang", err.Error())
		}
		return
	}

	views.JSONSuccess(c, http.StatusOK, "Data barang berhasil diperbarui", item)
}

// DeleteItem menangani DELETE /items/:id
func (ctrl *ItemController) DeleteItem(c *gin.Context) {
	id, err := parseIDParam(c)
	if err != nil {
		views.JSONError(c, http.StatusBadRequest, "ID barang tidak valid", err.Error())
		return
	}

	if err := ctrl.itemService.DeleteItem(c.Request.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrItemNotFound) {
			views.JSONError(c, http.StatusNotFound, "Gagal menghapus barang", err.Error())
			return
		}
		views.JSONError(c, http.StatusInternalServerError, "Gagal menghapus barang", err.Error())
		return
	}

	views.JSONSuccess(c, http.StatusOK, "Barang berhasil dihapus", nil)
}

func parseIDParam(c *gin.Context) (uint, error) {
	idParam := c.Param("id")
	parsed, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || parsed == 0 {
		return 0, errors.New("parameter ID harus berupa angka bulat positif")
	}
	return uint(parsed), nil
}

