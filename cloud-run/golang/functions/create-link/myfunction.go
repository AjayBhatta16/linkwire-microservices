package myfunction

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
	"github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	utilities.ApplyDefaultHeaders(w, "POST")
	
	// validate request
	var req Request
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RedirectURL == "" {
		http.Error(w, "Redirect URL is required", http.StatusBadRequest)
		return
	}

	// validate JWT
	token := utilities.GetTokenFromCookies(w, r)

	notExpired, err2 := utilities.ValidateJWTNotExpired(token)

	if err2 != nil || !notExpired {
		log.Println("Handler - JWT is expired or invalid:", err2)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tokenUsername, err3 := utilities.GetJWTUsername(token)

	if err3 != nil {
		log.Println("Handler - Error validating JWT:", err3)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// create and save link
	link := RequestToLink(req, tokenUsername)

	err4 := utilities.CreateItem(constants.LINK_CONTAINER_NAME, link)

	if err4 != nil {
		log.Println("Handler - Error creating link:", err4)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func RequestToLink(req Request, createdBy string) models.Link {
	var output models.Link

	output.RedirectURL = req.RedirectURL
	output.Note = req.Note

	output.CreatedBy = createdBy

	output.DisplayID = GenerateID()
	output.TrackingID = GenerateID()

	return output
}

func GenerateID() string {
    newID := make([]byte, 6)

    for i := range newID {
        newID[i] = CODE_CHARS[rand.Intn(len(CODE_CHARS))]
    }

    return string(newID)
}