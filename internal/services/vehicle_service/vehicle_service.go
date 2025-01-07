package vehicle

import (
	"car-sharing-system/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// Database connection
var db *gorm.DB

// Vehicle struct representing a vehicle in the system
type Vehicle struct {
	ID           uint    `json:"id"`
	ModelName    string  `json:"model_name"`
	Brand        string  `json:"brand"`
	Year         int     `json:"year"`
	LicensePlate string  `json:"license_plate"`
	Color        string  `json:"color"`
	PricePerHour float64 `json:"price_per_hour"`
	Location     string  `json:"location"`
	Status       string  `json:"status"`
}

// VehicleServiceHandler handles requests for vehicle operations
func VehicleServiceHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/vehicle", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createVehicle(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	handler.HandleFunc("/vehicle/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len("/vehicle/"):])
		if err != nil {
			http.Error(w, "Invalid vehicle ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			getVehicle(w, uint(id))
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	handler.HandleFunc("/vehicle/available", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getAvailableVehicles(w)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return handler
}

func createVehicle(w http.ResponseWriter, r *http.Request) {
	var vehicle Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Save vehicle to database
	if err := db.Create(&vehicle).Error; err != nil {
		http.Error(w, "Failed to create vehicle", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(vehicle)
}

func getVehicle(w http.ResponseWriter, id uint) {
	var vehicle Vehicle
	// Fetch vehicle from database
	if err := db.First(&vehicle, id).Error; err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(vehicle)
}

func getAvailableVehicles(w http.ResponseWriter) {
	var vehicles []Vehicle
	// Fetch available vehicles from database
	if err := db.Where("status = ?", "available").Find(&vehicles).Error; err != nil {
		http.Error(w, "Failed to fetch available vehicles", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(vehicles)
}

func InitVehicleServiceDB() error {
	var err error
	db, err = database.InitDB() // Initialize database
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	fmt.Println("Vehicle service database connection established")
	return nil
}
