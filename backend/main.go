package main

import (
	"backend/config"
	"backend/controllers"
	"backend/models"
	"backend/repositories"
	"backend/routes"
	"backend/services"
	"fmt"
	"log"
)

func main() {
	// 1. Load konfigurasi aplikasi
	cfg := config.LoadConfig()

	// 2. Koneksi ke database PostgreSQL
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("[FATAL] %v", err)
	}

	// 3. Auto-Migrate skema tabel User dan Item (beserta constraint CHECK stock >= 0)
	if err := db.AutoMigrate(&models.User{}, &models.Item{}); err != nil {
		log.Fatalf("[FATAL] Gagal melakukan migrasi database: %v", err)
	}
	log.Println("[DATABASE] Migrasi skema database berhasil!")

	// 4. Dependency Injection (Merakit Repository -> Service -> Controller -> Router)
	userRepo := repositories.NewUserRepository(db)
	itemRepo := repositories.NewItemRepository(db)

	authService := services.NewAuthService(userRepo, cfg)
	itemService := services.NewItemService(itemRepo)

	authController := controllers.NewAuthController(authService)
	itemController := controllers.NewItemController(itemService)

	// 5. Inisialisasi Router
	r := routes.SetupRouter(authController, itemController, authService)

	// 6. Jalankan HTTP Server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("[SERVER] Sistem Manajemen Inventaris berjalan di http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("[FATAL] Gagal menjalankan server: %v", err)
	}
}
