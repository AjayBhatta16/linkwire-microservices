package myfunction

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	var user User

	usersByUsername, err2 := GetItemsByFieldValue[User, *User]("users", "username", req.Username)

	if err2 != nil {
		log.Println("Error fetching users by username:", err2)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(usersByUsername) == 0 {
		usersByEmail, err3 := GetItemsByFieldValue[User, *User]("users", "email", req.Username)

		if err3 != nil {
			log.Println("Error fetching users by email:", err3)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		user = usersByEmail[0]
	} else {
		user = usersByUsername[0]
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		log.Printf("Password mismatch for user: %s", req.Username)
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	
	json.NewEncoder(w).Encode(user)
}