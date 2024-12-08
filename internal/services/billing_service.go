package services

import (
	"car-sharing-system/internal/models"
	"gorm.io/gorm"
	"time"
)

type BillingService struct {
	db *gorm.DB
	// Rate per hour in USD
	baseHourlyRate float64
	premiumRateMultiplier float64
}

func NewBillingService(db *gorm.DB) *BillingService {
	return &BillingService{
		db: db,
		baseHourlyRate: 10.0,  // $10 per hour for basic members
		premiumRateMultiplier: 0.8,  // 20% discount for premium members
	}
}

func (s *BillingService) CalculateBookingCost(user *models.User, startTime, endTime time.Time) float64 {
	duration := endTime.Sub(startTime).Hours()
	baseRate := 50.0 // Rate per hour in SGD

	// Apply membership discount if any
	var discount float64
	switch user.MembershipType {
	case "PREMIUM":
		discount = 0.2 // 20% discount
	case "GOLD":
		discount = 0.1 // 10% discount
	default:
		discount = 0.0
	}

	totalCost := baseRate * duration
	discountAmount := totalCost * discount
	return totalCost - discountAmount
}
