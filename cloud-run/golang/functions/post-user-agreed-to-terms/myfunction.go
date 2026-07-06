package myfunction

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
	"github.com/AjayBhatta16/linkwire-golang-shared/utilities"
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

	// validate JWT
	token := utilities.GetTokenFromCookies(w, r)

	notExpired, expErr := utilities.ValidateJWTNotExpired(token)

	if expErr != nil || !notExpired {
		log.Println("Handler - JWT is expired or invalid:", expErr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenUsername, err := utilities.GetJWTUsername(token)

	if err != nil {
		log.Println("Handler - Error validating JWT:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if tokenUsername != username {
		log.Printf("Handler - JWT username %s does not match path username %s", tokenUsername, username)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// fetch user
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

	// update user
	user.AgreedToLatestTerms = true

	// save user
	err4 := utilities.UpdateItem[models.User](
		constants.USER_CONTAINER_NAME, user.FirestoreID, user)

	if err4 != nil {
		http.Error(w, "Failed to update user password", http.StatusInternalServerError)
		return
	}

	// HTTP response
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