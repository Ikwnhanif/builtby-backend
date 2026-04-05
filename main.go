package main

import (
	"log"
	"os"
	"time"

	"github.com/ikwnhanif/builtby/config"
	"github.com/ikwnhanif/builtby/models"
	"github.com/ikwnhanif/builtby/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Inisialisasi Environment & Database
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan environment system")
	}

	config.ConnectDatabase()

	// 2. Auto Migration & Seeding Data Awal
	config.DB.AutoMigrate(&models.Profile{}, &models.Project{}, &models.Gallery{})
	config.SeedProfile()

	// 3. Inisialisasi App Fiber
	app := fiber.New(fiber.Config{
		AppName:   "Builtby Portfolio API v1.0",
		BodyLimit: 10 * 1024 * 1024, 
	})

	// 4. Global Middlewares
	app.Use(recover.New())
	app.Use(logger.New())

	// Helmet Config
	app.Use(helmet.New(helmet.Config{
		XSSProtection:             "enabled",
		ContentTypeNosniff:        "enabled",
		XFrameOptions:             "SAMEORIGIN",
		CrossOriginResourcePolicy: "cross-origin",
	}))

	// CORS Config
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Ganti sesuai origin frontend
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// 5. Rate Limiter
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Terlalu banyak permintaan",
			})
		},
	}))

	// 6. FIX STATIC: Gunakan Middleware untuk menyuntikkan Header
	app.Use("/uploads", func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Cross-Origin-Resource-Policy", "cross-origin")
		return c.Next()
	})
	app.Static("/uploads", "./uploads")

	// 7. Setup Routes
	routes.SetupRoutes(app)

	// 8. Menjalankan Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server builtby siap tempur di port %s", port)
	log.Fatal(app.Listen(":" + port))
}