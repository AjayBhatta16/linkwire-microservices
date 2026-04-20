package myfunction

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
    "github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// validate request
	id := GetLinkIDFromPath(r)

	if id == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
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

	// fetch link by id
	data, err2 := utilities.GetItemsByFieldValue[models.Link, *models.Link](constants.LINK_CONTAINER_NAME, "trackingID", id)

	if err2 != nil {
		log.Println("Handler - Error fetching link by id:", err2)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if len(data) == 0 {
		http.Error(w, "Link not found", http.StatusNotFound)
		return
	}

	// validate link
	link := data[0]

	if link.CreatedBy != tokenUsername {
		log.Printf("Handler - JWT username %s does not match link's createdBy %s", tokenUsername, link.CreatedBy)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// populate click data
	clicks, err3 := utilities.GetItemsByFieldValue[models.Click, *models.Click](constants.CLICK_CONTAINER_NAME, "linkID", link.DisplayID)

	if err3 != nil {
		log.Println("Handler - Error fetching clicks for link:", err3)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	link.Clicks = clicks

	// return link as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

func GetLinkIDFromPath(r *http.Request) string {
	url := r.URL.Path
	
	tokens := strings.Split(url, "/")
	slices.Reverse(tokens)
	lastToken := tokens[0]

	if lastToken == "links" || lastToken == "links/" {
		return ""
	}

	return lastToken
}