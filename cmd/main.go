package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"car-sharing-system/internal/database"
	"car-sharing-system/api"
	"car-sharing-system/internal/services"
)

func main() {
	// Initialize router
	router := gin.Default()

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize services
	userService := services.NewUserService(db)
	vehicleService := services.NewVehicleService(db)
	bookingService := services.NewBookingService(db, vehicleService)
	billingService := services.NewBillingService(db)

	// Setup routes
	api.SetupRoutes(router, userService, vehicleService, bookingService, billingService)

	// Start server
	router.Run(":8080")
}
