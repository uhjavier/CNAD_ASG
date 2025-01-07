package booking

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

// Booking struct representing a booking in the system
type Booking struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	VehicleID uint   `json:"vehicle_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
}

// BookingServiceHandler handles requests for booking operations
func BookingServiceHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/booking", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createBooking(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	handler.HandleFunc("/booking/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len("/booking/"):])
		if err != nil {
			http.Error(w, "Invalid booking ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			getBooking(w, uint(id))
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return handler
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	var booking Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Save booking to database
	if err := db.Create(&booking).Error; err != nil {
		http.Error(w, "Failed to create booking", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func getBooking(w http.ResponseWriter, id uint) {
	var booking Booking
	// Fetch booking from database
	if err := db.First(&booking, id).Error; err != nil {
		http.Error(w, "Booking not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(booking)
}

func InitBookingServiceDB() error {
	var err error
	db, err = database.InitDB() // Initialize database
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	fmt.Println("Booking service database connection established")
	return nil
}
