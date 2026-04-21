package myfunction

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
    "github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	utilities.ApplyDefaultHeaders(w, "POST")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)

	// Request validation
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if req.Password == "" {
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// validate username uniqueness
	usersByUsername, err2 := utilities.GetItemsByFieldValue[models.User, *models.User](constants.USER_CONTAINER_NAME, "username", req.Username)

	if err2 != nil {
		log.Println("Error fetching users by username:", err2)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(usersByUsername) > 0 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// validate email uniqueness
	usersByEmail, err3 := utilities.GetItemsByFieldValue[models.User, *models.User](constants.USER_CONTAINER_NAME, "email", req.Email)

	if err3 != nil {
		log.Println("Error fetching users by email:", err3)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(usersByEmail) > 0 {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// construct user object and hash password
	user := RequestToUser(req)

	// save user to database
	err = utilities.CreateItem(constants.USER_CONTAINER_NAME, user)

	// return created user
	jwt, _ := utilities.GenerateJWT(user.Username)
	cookieHeader := utilities.GetSetCookieHeaderValue(jwt)

	w.Header().Set("Set-Cookie", cookieHeader)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}

func RequestToUser(req Request) models.User {
	var output models.User

	output.Username = req.Username
	output.Email = req.Email
	output.Password = HashPassword(req.Password)
	output.PremiumUser = false
	output.Links = []string{}

	return output
}

func HashPassword(password string) string {
    saltRounds := 10

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), saltRounds)

    return string(hashedPassword)
}