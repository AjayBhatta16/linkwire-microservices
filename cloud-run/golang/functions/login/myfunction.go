package myfunction

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
    "github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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

	var user models.User

	usersByUsername, err2 := utilities.GetItemsByFieldValue[models.User, *models.User]("users", "username", req.Username)

	if err2 != nil {
		log.Println("Error fetching users by username:", err2)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(usersByUsername) == 0 {
		usersByEmail, err3 := utilities.GetItemsByFieldValue[models.User, *models.User]("users", "email", req.Username)

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

	jwt, _ := utilities.GenerateJWT(user.Username)
	cookieHeader := utilities.GetSetCookieHeaderValue(jwt)

	w.Header().Set("Set-Cookie", cookieHeader)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}