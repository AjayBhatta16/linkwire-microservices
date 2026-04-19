package myfunction

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
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

	usersByUsername, _ := GetItemsByFieldValue[User]("users", "username", req.Username)

	if len(usersByUsername) == 0 {
		usersByEmail, _ := GetItemsByFieldValue[User]("users", "email", req.Username)

		if len(usersByEmail) == 0 {
			http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
			return
		}

		user = usersByEmail[0]
	} else {
		user = usersByUsername[0]
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}