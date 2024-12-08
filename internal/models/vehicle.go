package models

import (
	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	ModelName    string  `json:"model_name"`
	Brand        string  `json:"brand"`
	Year         int     `json:"year"`
	LicensePlate string  `json:"license_plate" gorm:"unique"`
	Color        string  `json:"color"`
	PricePerHour float64 `json:"price_per_hour"`
	Location     string  `json:"location"`
	Type         string  `json:"type"`   // car, motorcycle, etc.
	Status       string  `json:"status" gorm:"default:available"` // active, maintenance, retired
}
