package main

import (
	"log"
	"net/http"

	billing "car-sharing-system/internal/services/billing_service"
	booking "car-sharing-system/internal/services/booking_service"
	user "car-sharing-system/internal/services/user_service"
	vehicle "car-sharing-system/internal/services/vehicle_service"
)

func main() {
	// Initialize service databases
	if err := user.InitUserServiceDB(); err != nil {
		log.Fatalf("User service database initialization failed: %v", err)
	}
	if err := billing.InitBillingServiceDB(); err != nil {
		log.Fatalf("Billing service database initialization failed: %v", err)
	}
	if err := booking.InitBookingServiceDB(); err != nil {
		log.Fatalf("Booking service database initialization failed: %v", err)
	}
	if err := vehicle.InitVehicleServiceDB(); err != nil {
		log.Fatalf("Vehicle service database initialization failed: %v", err)
	}

	// Start services
	go func() {
		log.Println("Starting User Service on port 8081...")
		err := http.ListenAndServe(":8081", user.UserServiceHandler())
		if err != nil {
			log.Fatalf("User Service failed: %v", err)
		}
	}()

	go func() {
		log.Println("Starting Billing Service on port 8082...")
		err := http.ListenAndServe(":8082", billing.BillingServiceHandler())
		if err != nil {
			log.Fatalf("Billing Service failed: %v", err)
		}
	}()

	go func() {
		log.Println("Starting Booking Service on port 8083...")
		err := http.ListenAndServe(":8083", booking.BookingServiceHandler())
		if err != nil {
			log.Fatalf("Booking Service failed: %v", err)
		}
	}()

	go func() {
		log.Println("Starting Vehicle Service on port 8084...")
		err := http.ListenAndServe(":8084", vehicle.VehicleServiceHandler())
		if err != nil {
			log.Fatalf("Vehicle Service failed: %v", err)
		}
	}()

	// Block main goroutine to keep services running
	select {}
}
