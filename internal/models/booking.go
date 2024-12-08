package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	gorm.Model
	UserID    uint      `json:"user_id"`
	VehicleID uint      `json:"vehicle_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Status    string    `json:"status"`
	TotalPrice float64  `json:"total_price"`
	User       User      `gorm:"foreignKey:UserID"`
	Vehicle    Vehicle   `gorm:"foreignKey:VehicleID"`
}
