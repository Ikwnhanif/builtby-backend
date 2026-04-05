package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ikwnhanif/builtby/controllers"
	"github.com/ikwnhanif/builtby/middleware"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// --- PUBLIC ROUTES (Bisa diakses siapa saja) ---
	api.Post("/login", controllers.Login)
	api.Get("/profile", controllers.GetProfile)
	api.Get("/projects", controllers.GetProjects)
	
	// Route baru untuk ambil foto galeri di Home Page
	api.Get("/gallery", controllers.GetGallery) 

	// --- ADMIN ROUTES (Protected by JWT) ---
	admin := api.Group("/", middleware.Protected())
	
	// Profile Management
	admin.Put("/profile", controllers.UpdateProfile)
	
	// Project Management
	admin.Post("/projects", controllers.CreateProject)
	admin.Put("/projects/:id", controllers.UpdateProject)
	admin.Delete("/projects/:id", controllers.DeleteProject)

	// Gallery Management (Baru)
	admin.Post("/gallery", controllers.CreateGallery)    // Upload foto baru
	admin.Delete("/gallery/:id", controllers.DeleteGallery) // Hapus foto
}