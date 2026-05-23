package myfunction

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"

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

	// trigger pub/sub for link processing
	err5 := SubmitLinkForProcessing(link.DisplayID)

	if err5 != nil {
		// don't return a 500 response, since the link has already been created
		// eventually, we will set up a dead letter queue and/or email alerts for failures here
		log.Println("Handler - Error submitting link for processing:", err5)
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

func SubmitLinkForProcessing(linkID string) error {
	publisher, err := utilities.NewPublisher(context.Background())

    if err != nil {
        log.Fatalf("failed to create publisher: %v", err)
		return err
    }

    defer publisher.Close()

	var payload ProcessLinkRequest

	payload.LinkID = linkID

    err = utilities.Publish(publisher, "process-link-topic", payload)

	return err
}