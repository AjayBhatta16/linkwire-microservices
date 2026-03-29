package myfunction

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func ProcessRequest(w http.ResponseWriter, r *http.Request) {
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

	if req.Username == nil {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if req.Password == nil {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	var user User

	usersByUsername := GetItemsByFieldValue[User]("users", "username", req.Username)

	if len(usersByUsername) == 0 {
		usersByEmail := GetItemsByFieldValue[User]("users", "email", req.Username)

		if len(usersByEmail) == 0 {
			http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
			return
		}

		user = usersByEmail[0]
	} else {
		user = usersByUsername[0]
	}

	hashedPassword, _ := HashPassword(req.Password)

	if hashedPassword != user.Password {
		http.Error(w, "Incorrect username or password", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func HashPassword(password string) (string, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10)
}