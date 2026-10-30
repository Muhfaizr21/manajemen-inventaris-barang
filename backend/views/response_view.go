package views

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIResponse merupakan struktur standar JSON response (View Layer)
type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Errors  interface{}     `json:"errors,omitempty"`
}

// PaginationMeta menyimpan metadata paginasi untuk endpoint list
type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

// JSONSuccess mengirimkan response JSON sukses standar
func JSONSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// JSONPaginated mengirimkan response JSON sukses beserta metadata paginasi
func JSONPaginated(c *gin.Context, statusCode int, message string, data interface{}, meta PaginationMeta) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

// JSONError mengirimkan response JSON error standar
func JSONError(c *gin.Context, statusCode int, message string, errDetails interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

// FormatValidationError mengubah error validator menjadi map pesan yang ramah pengguna
func FormatValidationError(err error) interface{} {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]string, len(ve))
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				out[field] = fmt.Sprintf("Field '%s' wajib diisi", field)
			case "email":
				out[field] = "Format email tidak valid"
			case "min":
				out[field] = fmt.Sprintf("Field '%s' minimal %s karakter/nilai", field, fe.Param())
			case "max":
				out[field] = fmt.Sprintf("Field '%s' maksimal %s karakter/nilai", field, fe.Param())
			case "gte":
				out[field] = fmt.Sprintf("Field '%s' tidak boleh kurang dari %s (tidak boleh negatif)", field, fe.Param())
			default:
				out[field] = fmt.Sprintf("Field '%s' tidak memenuhi aturan '%s'", field, fe.Tag())
			}
		}
		return out
	}
	return err.Error()
}
