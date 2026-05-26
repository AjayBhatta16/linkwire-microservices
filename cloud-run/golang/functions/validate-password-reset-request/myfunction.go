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
)

func Handler(w http.ResponseWriter, r *http.Request) {
	utilities.ApplyDefaultHeaders(w, r, "GET")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// validate request path
	log.Printf("Request path: %s", r.URL.Path)
	requestId := GetRequestIDFromPath(r)

	if requestId == "" {
		http.Error(w, "Request ID is required", http.StatusBadRequest)
		return
	}

	// fetch password reset request by ID
	resetRequests, err := utilities.GetItemsByFieldValue[models.PasswordResetRequest, *models.PasswordResetRequest](
		constants.PASSWORD_RESET_REQUEST_CONTAINER_NAME, "requestId", requestId)

	if err != nil {
		log.Println("Handler - Error fetching password reset request by ID:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(resetRequests) == 0 {
		log.Printf("Handler - No password reset request found for ID: %s", requestId)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	resetRequest := resetRequests[0]

	// check if request is still valid (not expired or completed)
	if resetRequest.ExpirationTimestamp < time.Now().Unix() {
		log.Printf("Handler - Password reset request with ID %s is expired", requestId)
		http.Error(w, "Expired", http.StatusBadRequest)
		return
	}

	if resetRequest.ResetCompleted {
		log.Printf("Handler - Password reset request with ID %s is already completed", requestId)
		http.Error(w, "Already Completed", http.StatusBadRequest)
		return
	}

	// return username in response
	response := Response{
		Username: resetRequest.RequestedForUsername,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetRequestIDFromPath(r *http.Request) string {
	parts := strings.Split(r.URL.Path, "/")

	requestId := ""

	for i, part := range parts {
		if part == "reset-request" && i+1 < len(parts) {
			requestId = parts[i+1]
			break
		}
	}

	return requestId
}