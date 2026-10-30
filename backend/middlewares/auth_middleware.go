package middlewares

import (
	"backend/services"
	"backend/views"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const ContextUserIDKey = "user_id"

// AuthMiddleware memvalidasi token JWT Bearer sebelum request masuk ke Controller
func AuthMiddleware(authService services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			views.JSONError(c, http.StatusUnauthorized, "Akses ditolak: Header Authorization tidak ditemukan", "Format wajib: Bearer <token>")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			views.JSONError(c, http.StatusUnauthorized, "Format header Authorization tidak valid", "Gunakan format: Bearer <token>")
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(strings.TrimSpace(parts[1]))
		if err != nil {
			views.JSONError(c, http.StatusUnauthorized, "Token tidak valid atau sudah kedaluwarsa", err.Error())
			c.Abort()
			return
		}

		// Simpan identitas user ke dalam context request
		c.Set(ContextUserIDKey, claims.UserID)
		c.Next()
	}
}

