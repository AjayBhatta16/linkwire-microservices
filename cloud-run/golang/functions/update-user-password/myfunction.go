package myfunction

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
	"github.com/AjayBhatta16/linkwire-golang-shared/utilities"

	"golang.org/x/crypto/bcrypt"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	utilities.ApplyDefaultHeaders(w, r, "POST")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// validate request path
	log.Printf("Request path: %s", r.URL.Path)
	username := GetUsernameFromPath(r)

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// validate request body
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ResetRequestID == "" {
		http.Error(w, "Reset request ID is required", http.StatusBadRequest)
		return
	}

	if req.NewPassword == "" {
		http.Error(w, "New password is required", http.StatusBadRequest)
		return
	}

	// fetch password reset request from database and validate
	passwordResetRequests, err2 := utilities.GetItemsByFieldValue[models.PasswordResetRequest, *models.PasswordResetRequest](
		constants.PASSWORD_RESET_REQUEST_CONTAINER_NAME, "resetRequestId", req.ResetRequestID)

	if err2 != nil {
		http.Error(w, "Failed to fetch password reset request", http.StatusInternalServerError)
		return
	}

	if len(passwordResetRequests) == 0 {
		http.Error(w, "Invalid reset request ID", http.StatusBadRequest)
		return
	}

	passwordResetRequest := passwordResetRequests[0]

	if passwordResetRequest.RequestedForUsername != username {
		http.Error(w, "Reset request does not match username", http.StatusBadRequest)
		return
	}

	if passwordResetRequest.ExpirationTimestamp < time.Now().Unix() {
		http.Error(w, "Reset request has expired", http.StatusBadRequest)
		return
	}

	// update user password
	users, err3 := utilities.GetItemsByFieldValue[models.User, *models.User](
		constants.USER_CONTAINER_NAME, "username", username)

	if err3 != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user := users[0]

	user.Password = HashPassword(req.NewPassword)

	err4 := utilities.UpdateItem[models.User](
		constants.USER_CONTAINER_NAME, user.FirestoreID, user)

	if err4 != nil {
		http.Error(w, "Failed to update user password", http.StatusInternalServerError)
		return
	}

	// update password reset request to mark it as used
	passwordResetRequest.ResetCompleted = true

	err5 := utilities.UpdateItem[models.PasswordResetRequest](
		constants.PASSWORD_RESET_REQUEST_CONTAINER_NAME, 
		passwordResetRequest.FirestoreID, 
		passwordResetRequest)

	if err5 != nil {
		http.Error(w, "Failed to update password reset request", http.StatusInternalServerError)
		return
	}

	// return user with new JWT
	jwt, _ := utilities.GenerateJWT(user.Username)
	cookieHeader := utilities.GetSetCookieHeaderValue(jwt)

	w.Header().Set("Set-Cookie", cookieHeader)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}

func GetUsernameFromPath(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")

	username := ""

	for i, part := range parts {
		if part == "users" && i+1 < len(parts) {
			username = parts[i+1]
			break
		}
	}

	return username
}

func HashPassword(password string) string {
    saltRounds := 10

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), saltRounds)

    return string(hashedPassword)
}