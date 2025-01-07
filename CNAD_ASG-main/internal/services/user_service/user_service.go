package user

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

// User struct representing a user in the system
type User struct {
	ID             uint   `json:"id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	MembershipType string `json:"membership_type"`
}

// UserServiceHandler handles requests for user operations
func UserServiceHandler() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			createUser(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	handler.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len("/user/"):])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			getUser(w, uint(id))
			return
		}
		if r.Method == http.MethodPut {
			updateUser(w, r, uint(id))
			return
		}
		if r.Method == http.MethodDelete {
			deleteUser(w, uint(id))
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	})

	return handler
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Save user to database
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getUser(w http.ResponseWriter, id uint) {
	var user User
	// Fetch user from database
	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request, id uint) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update user in database
	user.ID = id
	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func deleteUser(w http.ResponseWriter, id uint) {
	// Delete user from database
	if err := db.Delete(&User{}, id).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func InitUserServiceDB() error {
	var err error
	db, err = database.InitDB() // Initialize database
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	fmt.Println("User service database connection established")
	return nil
}
