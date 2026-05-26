package myfunction

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/AjayBhatta16/linkwire-golang-shared/constants"
	"github.com/AjayBhatta16/linkwire-golang-shared/models"
	"github.com/AjayBhatta16/linkwire-golang-shared/utilities"

	"github.com/google/uuid"
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

	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// check if user exists
	usersByEmail, err3 := utilities.GetItemsByFieldValue[models.User, *models.User]("users", "email", req.Email)

	if err3 != nil {
		log.Println("Error fetching users by email:", err3)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(usersByEmail) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	user := usersByEmail[0]

	// create password reset request
	passwordResetRequest := CreatePasswordResetRequest(user)

	// save password reset request to database
	err4 := utilities.CreateItem[models.PasswordResetRequest](constants.PASSWORD_RESET_REQUEST_CONTAINER_NAME, passwordResetRequest)

	if err4 != nil {
		log.Println("Error saving password reset request:", err4)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// send email to user
	emailRequest := GenerateEmailRequest(user, passwordResetRequest)

	err5 := SendEmail(emailRequest)

	if err5 != nil {
		log.Println("Error sending email:", err5)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreatePasswordResetRequest(user models.User) models.PasswordResetRequest {
	var output models.PasswordResetRequest

	output.RequestedForUsername = user.Username
	output.RequestedForEmail = user.Email

	output.RequestedTimestamp = time.Now().Unix()
	output.ExpirationTimestamp = time.Now().Add(4 * time.Hour).Unix()

	output.RequestId = uuid.NewString()

	output.ResetCompleted = false

	return output
}

func GenerateEmailRequest(user models.User, request models.PasswordResetRequest) EmailRequest {
	var output EmailRequest

	output.To = user.Email
	output.Subject = "Reset Your LinkWire Password"

	output.Body = "<p>Click the link below to reset your password.</p>"
	output.Body += "<p>This link will expire in 4 hours.</p>"
	output.Body += "<br/>"
	output.Body += "<a href='https://app.linkwire.cc/password-resets/" + request.RequestId + "'>Reset Password</a>"

	return output
}

func SendEmail(payload EmailRequest) error {
	publisher, err := utilities.NewPublisher(context.Background())

    if err != nil {
        log.Fatalf("failed to create publisher: %v", err)
		return err
    }

    defer publisher.Close()

    err = utilities.Publish(publisher, "send-email-topic", payload)

	return err
}