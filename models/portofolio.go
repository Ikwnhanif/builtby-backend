package models

import "time"

type Profile struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FullName  string    `json:"full_name"`
	JobTitle  string    `json:"job_title"`
	Bio       string    `json:"bio"`
	PhotoURL  string    `json:"photo_url"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Project struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
}

type Gallery struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title"`
	Category  string    `json:"category"` // Contoh: Street, Landscape, Portrait
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}