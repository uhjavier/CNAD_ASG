package api

import (
	"car-sharing-system/internal/models"
	"car-sharing-system/internal/services"
	"net/http"
	"strconv"
	"time"

	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, us *services.UserService, vs *services.VehicleService, bs *services.BookingService, bls *services.BillingService) {
	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "UP",
			"timestamp": time.Now(),
		})
	})

	api := router.Group("/api")

	// Public routes
	api.POST("/register", handleRegister(us))
	api.POST("/login", handleLogin(us))

	// Protected routes
	protected := api.Group("/")
	protected.Use(authMiddleware())
	{
		// Vehicle routes
		protected.POST("/vehicles", handleCreateVehicle(vs))
		protected.GET("/vehicles", handleListVehicles(vs))
		protected.GET("/vehicles/available", handleGetAvailableVehicles(vs))
		protected.PATCH("/vehicles/:id/status", handleUpdateVehicleStatus(vs))

		// Booking routes
		protected.POST("/bookings", handleCreateBooking(bs, vs))
		protected.GET("/bookings/history", handleGetBookingHistory(bs))

		// Billing routes
		protected.POST("/billing/estimate", handleGetBillingEstimate(bls))
		protected.POST("/billing/calculate/:booking_id", handleCalculateBookingCost(bls))
	}
}

func handleRegister(us *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := us.Register(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func handleLogin(us *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var credentials struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&credentials); err != nil {
			log.Printf("Login - Invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Debug check
		us.DebugUserByEmail(credentials.Email)

		user, err := us.Login(credentials.Email, credentials.Password)
		if err != nil {
			log.Printf("Login failed for email %s: %v", credentials.Email, err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":              user.ID,
			"email":           user.Email,
			"phone":           user.Phone,
			"membership_type": user.MembershipType,
		})
	}
}

func handleListVehicles(vs *services.VehicleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicles, err := vs.ListVehicles()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, vehicles)
	}
}

func handleCreateBooking(bs *services.BookingService, vs *services.VehicleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var booking models.Booking
		if err := c.ShouldBindJSON(&booking); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
			return
		}

		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}
		booking.UserID = user.(*models.User).ID

		// Verify vehicle exists and is available
		vehicle, err := vs.GetVehicle(booking.VehicleID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Vehicle not found: " + err.Error()})
			return
		}

		if vehicle.Status != "available" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Vehicle is not available"})
			return
		}

		// Validate booking times
		if booking.StartTime.IsZero() || booking.EndTime.IsZero() {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start time and end time are required"})
			return
		}

		if booking.EndTime.Before(booking.StartTime) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
			return
		}

		// Set initial booking status
		booking.Status = "active"

		if err := bs.CreateBooking(&booking); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking: " + err.Error()})
			return
		}

		// Update vehicle availability
		if err := vs.UpdateVehicleAvailability(booking.VehicleID, false); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vehicle availability: " + err.Error()})
			return
		}

		c.JSON(http.StatusCreated, booking)
	}
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// For now, just set a dummy user in context
		dummyUser := &models.User{
			Model:          gorm.Model{ID: 1},
			Email:          "user1@example.com",
			MembershipType: "BASIC",
		}
		c.Set("user", dummyUser)

		c.Next()
	}
}

func handleGetAvailableVehicles(vs *services.VehicleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicles, err := vs.GetAvailableVehicles()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, vehicles)
	}
}

func handleGetBookingHistory(bs *services.BookingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// In a real app, you'd get the user ID from the JWT token
		userID := uint(1) // Placeholder

		bookings, err := bs.GetUserBookings(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, bookings)
	}
}

func handleUpdateVehicleStatus(vs *services.VehicleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var updates map[string]interface{}
		if err := c.ShouldBindJSON(&updates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		if err := vs.UpdateVehicleStatus(uint(id), updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Vehicle status updated successfully"})
	}
}

func handleGetBillingEstimate(bls *services.BillingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			VehicleID uint      `json:"vehicle_id" binding:"required"`
			StartTime time.Time `json:"start_time" binding:"required"`
			EndTime   time.Time `json:"end_time" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		// Calculate estimate
		cost := bls.CalculateBookingCost(user.(*models.User), req.StartTime, req.EndTime)

		c.JSON(http.StatusOK, gin.H{
			"estimated_cost": cost,
			"currency":       "SGD",
			"vehicle_id":     req.VehicleID,
			"start_time":     req.StartTime,
			"end_time":       req.EndTime,
		})
	}
}

func handleCalculateBookingCost(bls *services.BillingService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookingID, err := strconv.ParseUint(c.Param("booking_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
			return
		}

		// Get user from context
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
			return
		}

		// For demo, using fixed duration
		startTime := time.Now().Add(-2 * time.Hour) // Assuming booking started 2 hours ago
		endTime := time.Now()                       // Ending now

		cost := bls.CalculateBookingCost(user.(*models.User), startTime, endTime)

		c.JSON(http.StatusOK, gin.H{
			"booking_id":       bookingID,
			"total_cost":       cost,
			"hourly_rate":      50.0,
			"hours":            2.0,
			"membership_type":  user.(*models.User).MembershipType,
			"discount_applied": user.(*models.User).MembershipType == "PREMIUM",
			"currency":         "SGD",
			"start_time":       startTime,
			"end_time":         endTime,
		})
	}
}

func handleCreateVehicle(vs *services.VehicleService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vehicle models.Vehicle
		if err := c.ShouldBindJSON(&vehicle); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := vs.CreateVehicle(&vehicle); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, vehicle)
	}
}
