package services

import (
	"car-sharing-system/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type BookingService struct {
	db             *gorm.DB
	vehicleService *VehicleService
}

func NewBookingService(db *gorm.DB, vehicleService *VehicleService) *BookingService {
	return &BookingService{
		db:             db,
		vehicleService: vehicleService,
	}
}

func (s *BookingService) CreateBooking(booking *models.Booking) error {
	// Validate vehicle exists and is available
	var vehicle models.Vehicle
	if err := s.db.First(&vehicle, booking.VehicleID).Error; err != nil {
		return fmt.Errorf("vehicle not found")
	}

	if vehicle.Status != "available" {
		return fmt.Errorf("vehicle is not available")
	}

	// Calculate total price
	duration := booking.EndTime.Sub(booking.StartTime).Hours()
	baseRate := 50.0 // SGD per hour
	booking.TotalPrice = baseRate * duration

	// Start transaction
	tx := s.db.Begin()

	// Create booking
	if err := tx.Create(booking).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update vehicle status
	if err := tx.Model(&vehicle).Update("status", "booked").Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (s *BookingService) GetByID(id uint) (*models.Booking, error) {
	var booking models.Booking
	if err := s.db.First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *BookingService) GetUserBookings(userID uint) ([]models.Booking, error) {
	var bookings []models.Booking
	result := s.db.Where("user_id = ?", userID).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}
	return bookings, nil
}

func (s *BookingService) UpdateBookingStatus(id uint, status string) error {
	booking := &models.Booking{}
	if err := s.db.First(booking, id).Error; err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update booking status
		if err := tx.Model(booking).Update("status", status).Error; err != nil {
			return err
		}

		// If booking is completed or cancelled, make vehicle available again
		if status == "completed" || status == "cancelled" {
			return s.vehicleService.UpdateVehicleAvailability(booking.VehicleID, true)
		}
		return nil
	})
}
