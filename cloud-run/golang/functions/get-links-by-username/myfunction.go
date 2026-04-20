package myfunction

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
    "github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// validate request
	username := utilities.GetVariableFromPath(r, "get-links-by-username")

	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// validate JWT
	token := utilities.GetTokenFromCookies(w, r)

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

	// fetch links for username
	links, err2 := utilities.GetItemsByFieldValue[models.Link, *models.Link](constants.LINK_CONTAINER_NAME, "createdBy", username)

	if err2 != nil {
		log.Println("Handler - Error fetching links for username:", err2)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// return links as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}