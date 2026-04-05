package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/ikwnhanif/builtby/config"
	"github.com/ikwnhanif/builtby/models"
)

// Get semua foto untuk ditampilkan di Gallery Page
func GetGallery(c *fiber.Ctx) error {
	var gallery []models.Gallery
	config.DB.Order("created_at desc").Find(&gallery)
	return c.JSON(gallery)
}

// Upload foto baru ke Gallery
func CreateGallery(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "File image diperlukan"})
	}

	uniqueId := uuid.New().String()
	fileName := fmt.Sprintf("gallery-%s-%s", uniqueId, file.Filename)
	c.SaveFile(file, fmt.Sprintf("./uploads/%s", fileName))

	gallery := models.Gallery{
		Title:    c.FormValue("title"),
		Category: c.FormValue("category"),
		ImageURL: fileName,
	}

	config.DB.Create(&gallery)
	return c.JSON(gallery)
}

// Hapus foto dari Gallery
func DeleteGallery(c *fiber.Ctx) error {
	id := c.Params("id")
	var gallery models.Gallery
	if err := config.DB.First(&gallery, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Foto tidak ditemukan"})
	}

	config.DB.Delete(&gallery)
	return c.JSON(fiber.Map{"message": "Foto berhasil dihapus"})
}