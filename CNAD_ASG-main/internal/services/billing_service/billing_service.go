package billing

import (
	"car-sharing-system/internal/database"
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// Database connection
var db *gorm.DB

// BillingRequest represents a request to estimate booking cost
type BillingRequest struct {
	VehicleID uint    `json:"vehicle_id"`
	Hours     float64 `json:"hours"`
}

// BillingResponse represents the cost estimation response
type BillingResponse struct {
	EstimatedCost float64 `json:"estimated_cost"`
}

// BillingServiceHandler handles requests for billing operations
func BillingServiceHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/billing/estimate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			estimateCost(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return handler
}

func estimateCost(w http.ResponseWriter, r *http.Request) {
	var req BillingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Fetch vehicle details from database
	var vehicle struct {
		PricePerHour float64
	}
	if err := db.Table("vehicles").Select("price_per_hour").Where("id = ?", req.VehicleID).Scan(&vehicle).Error; err != nil {
		http.Error(w, "Vehicle not found", http.StatusNotFound)
		return
	}

	// Calculate estimated cost
	estimatedCost := vehicle.PricePerHour * req.Hours
	resp := BillingResponse{EstimatedCost: estimatedCost}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func InitBillingServiceDB() error {
	var err error
	db, err = database.InitDB() // Initialize database
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	fmt.Println("Billing service database connection established")
	return nil
}
