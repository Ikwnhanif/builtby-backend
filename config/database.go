package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ikwnhanif/builtby/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Load file .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPass, dbHost, dbPort, dbName)

	// GORM Config dengan Logger sederhana
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Gagal menghubungkan ke database:", err)
	}

	// Pengaturan Connection Pool untuk Produksi
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatal("Gagal mendapatkan database instance:", err)
	}

	sqlDB.SetMaxIdleConns(10)                  // Koneksi standby
	sqlDB.SetMaxOpenConns(100)                 // Maksimal koneksi terbuka
	sqlDB.SetConnMaxLifetime(time.Hour)        // Durasi koneksi sebelum di-refresh

	fmt.Println("✅ Database terkoneksi & Connection Pool dikonfigurasi!")
	DB = database
}

func SeedProfile() {
	var count int64
	DB.Model(&models.Profile{}).Count(&count)
	if count == 0 {
		profile := models.Profile{
			FullName: "Nama Lengkap Anda",
			JobTitle: "Fullstack Developer",
			Bio:      "Selamat datang di portfolio saya.",
			PhotoURL: "default-profile.png", // Nama file default
		}
		DB.Create(&profile)
		fmt.Println("🌱 Data profile default berhasil dibuat!")
	}
}