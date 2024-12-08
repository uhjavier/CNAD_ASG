package services

import (
	"car-sharing-system/internal/models"
	"fmt"

	"gorm.io/gorm"
)

type VehicleService struct {
	db *gorm.DB
}

func NewVehicleService(db *gorm.DB) *VehicleService {
	return &VehicleService{db: db}
}

func (s *VehicleService) GetVehicle(id uint) (*models.Vehicle, error) {
	var vehicle models.Vehicle
	if err := s.db.First(&vehicle, id).Error; err != nil {
		return nil, err
	}
	return &vehicle, nil
}

func (s *VehicleService) ListVehicles() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	if err := s.db.Find(&vehicles).Error; err != nil {
		return nil, err
	}
	return vehicles, nil
}

func (s *VehicleService) CreateVehicle(vehicle *models.Vehicle) error {
	// Check if license plate already exists
	var existingVehicle models.Vehicle
	if err := s.db.Where("license_plate = ?", vehicle.LicensePlate).First(&existingVehicle).Error; err == nil {
		return fmt.Errorf("license plate already registered")
	}

	vehicle.Status = "available"
	return s.db.Create(vehicle).Error
}

func (s *VehicleService) UpdateVehicle(vehicle *models.Vehicle) error {
	return s.db.Save(vehicle).Error
}

func (s *VehicleService) DeleteVehicle(id uint) error {
	return s.db.Delete(&models.Vehicle{}, id).Error
}

func (s *VehicleService) GetByID(id uint) (*models.Vehicle, error) {
	return s.GetVehicle(id)
}

func (s *VehicleService) GetAvailableVehicles() ([]models.Vehicle, error) {
	var vehicles []models.Vehicle
	result := s.db.Where("status = ?", "available").Find(&vehicles)
	return vehicles, result.Error
}

func (s *VehicleService) UpdateVehicleStatus(id uint, updates map[string]interface{}) error {
	return s.db.Model(&models.Vehicle{}).Where("id = ?", id).Updates(updates).Error
}

func (s *VehicleService) UpdateVehicleAvailability(id uint, isAvailable bool) error {
	return s.db.Model(&models.Vehicle{}).Where("id = ?", id).Update("is_available", isAvailable).Error
}
