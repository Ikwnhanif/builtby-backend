package controllers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ikwnhanif/builtby/config"
	"github.com/ikwnhanif/builtby/models"
)

// --- PROFILE HANDLERS ---

func GetProfile(c *fiber.Ctx) error {
    var profile models.Profile
    config.DB.First(&profile)
    
    // AMBIL URL DARI ENV
    baseURL := os.Getenv("APP_URL") 
    if baseURL == "" {
        baseURL = "https://api-builtby.outsys.space"
    }

    // Suntikkan URL lengkap hanya saat output JSON
    if profile.PhotoURL != "" {
        profile.PhotoURL = fmt.Sprintf("%s/uploads/%s", baseURL, profile.PhotoURL)
    }

    return c.JSON(profile)
}

func UpdateProfile(c *fiber.Ctx) error {
    var profile models.Profile
    config.DB.First(&profile) // Selalu ambil ID 1

    profile.FullName = c.FormValue("full_name")
    profile.JobTitle = c.FormValue("job_title")
    profile.Bio = c.FormValue("bio")

    // Cek upload foto
    file, err := c.FormFile("image")
    if err == nil {
        uniqueId := uuid.New().String()
        fileName := fmt.Sprintf("%s-%s", uniqueId, file.Filename)
        c.SaveFile(file, fmt.Sprintf("./uploads/%s", fileName))
        profile.PhotoURL = fileName
    }

    config.DB.Save(&profile)
    return c.JSON(profile)
}

// --- PROJECT HANDLERS (CRUD) ---

// --- PROJECT HANDLERS ---
func GetProjects(c *fiber.Ctx) error {
    var projects []models.Project
    config.DB.Find(&projects)
    
    baseURL := os.Getenv("APP_URL")
    if baseURL == "" {
        baseURL = "https://api-builtby.outsys.space"
    }

    // Loop untuk menambahkan domain ke setiap thumbnail project
    for i := range projects {
        if projects[i].ImageURL != "" {
            projects[i].ImageURL = fmt.Sprintf("%s/uploads/%s", baseURL, projects[i].ImageURL)
        }
    }

    return c.JSON(projects)
}

func CreateProject(c *fiber.Ctx) error {
	// 1. Ambil data teks dari Form
	title := c.FormValue("title")
	desc := c.FormValue("description")
	link := c.FormValue("link")

	// 2. Handle Upload Gambar
	file, err := c.FormFile("image")
	var fileName string

	if err == nil {
		// Buat nama file unik pakai UUID agar tidak bentrok
		uniqueId := uuid.New().String()
		fileName = fmt.Sprintf("%s-%s", uniqueId, file.Filename)
		
		// Simpan file ke folder ./uploads
		c.SaveFile(file, fmt.Sprintf("./uploads/%s", fileName))
	}

	// 3. Simpan ke Database
	newProject := models.Project{
		Title:       title,
		Description: desc,
		ImageURL:    fileName, // Simpan nama filenya saja
		Link:        link,
	}

	if err := config.DB.Create(&newProject).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal simpan ke database"})
	}

	return c.JSON(newProject)
}

func DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project

	// Cari datanya dulu untuk ambil nama file
	if err := config.DB.First(&project, id).Error; err == nil {
		// Hapus file fisik
		os.Remove(fmt.Sprintf("./uploads/%s", project.ImageURL))
	}

	config.DB.Delete(&models.Project{}, id)
	return c.JSON(fiber.Map{"message": "Project dan file gambar berhasil dihapus"})
}

func UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var project models.Project

	// Cari project berdasarkan ID
	if err := config.DB.First(&project, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Project tidak ditemukan"})
	}

	// Update field teks
	project.Title = c.FormValue("title")
	project.Description = c.FormValue("description")
	project.Link = c.FormValue("link")

	// Cek apakah ada file gambar baru yang diunggah
	file, err := c.FormFile("image")
	if err == nil {
		// Hapus file lama (opsional, tapi bagus untuk hemat storage)
		// os.Remove(fmt.Sprintf("./uploads/%s", project.ImageURL))

		uniqueId := uuid.New().String()
		fileName := fmt.Sprintf("%s-%s", uniqueId, file.Filename)
		c.SaveFile(file, fmt.Sprintf("./uploads/%s", fileName))
		project.ImageURL = fileName
	}

	config.DB.Save(&project)
	return c.JSON(project)
}