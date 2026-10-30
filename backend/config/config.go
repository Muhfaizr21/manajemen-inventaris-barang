package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AppConfig menyimpan seluruh konfigurasi environment aplikasi
type AppConfig struct {
	Port           string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPass         string
	DBName         string
	JWTSecret      string
	JWTExpiryHours int
}

// LoadConfig memuat konfigurasi dari environment variables dengan fallback default
func LoadConfig() *AppConfig {
	expiryHours, err := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	if err != nil {
		expiryHours = 24
	}

	cfg := &AppConfig{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "muhfaiizr"),
		DBPass:         getEnv("DB_PASS", "admin"),
		DBName:         getEnv("DB_NAME", "tes_backend"),
		JWTSecret:      getEnv("JWT_SECRET", "super-secret-inventory-jwt-key-2026"),
		JWTExpiryHours: expiryHours,
	}

	log.Printf("[CONFIG] Configuration loaded (Port: %s, DB: %s)", cfg.Port, cfg.DBName)
	return cfg
}

// ConnectDatabase menginisialisasi koneksi GORM ke PostgreSQL dan mengembalikan instance *gorm.DB
func ConnectDatabase(cfg *AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database PostgreSQL: %w", err)
	}

	log.Println("[DATABASE] Berhasil terhubung ke database PostgreSQL!")
	return db, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
