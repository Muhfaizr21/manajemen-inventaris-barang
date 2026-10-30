package controllers

import (
	"backend/models"
	"backend/services"
	"backend/views"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthController menangani HTTP request untuk autentikasi (Register & Login)
type AuthController struct {
	authService services.AuthService
}

// NewAuthController membuat instance baru AuthController
func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register menangani endpoint POST /register (dan /api/auth/register)
func (ctrl *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		views.JSONError(c, http.StatusBadRequest, "Validasi input gagal", views.FormatValidationError(err))
		return
	}

	user, err := ctrl.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyRegistered) {
			views.JSONError(c, http.StatusConflict, "Registrasi gagal", err.Error())
			return
		}
		views.JSONError(c, http.StatusInternalServerError, "Gagal mendaftarkan pengguna", err.Error())
		return
	}

	views.JSONSuccess(c, http.StatusCreated, "Registrasi berhasil", user)
}

// Login menangani endpoint POST /login (dan /api/auth/login)
func (ctrl *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		views.JSONError(c, http.StatusBadRequest, "Validasi input gagal", views.FormatValidationError(err))
		return
	}

	authResp, err := ctrl.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			views.JSONError(c, http.StatusUnauthorized, "Autentikasi gagal", err.Error())
			return
		}
		views.JSONError(c, http.StatusInternalServerError, "Gagal memproses login", err.Error())
		return
	}

	views.JSONSuccess(c, http.StatusOK, "Login berhasil", authResp)
}

