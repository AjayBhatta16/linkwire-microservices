package myfunction

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	utilities.ApplyDefaultHeaders(w, r, "POST")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// validate request
	var req ContactRequest
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if req.ReturnEmail == "" {
		http.Error(w, "Return email is required", http.StatusBadRequest)
		return
	}

	if req.Subject == "" {
		http.Error(w, "Subject is required", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message is required", http.StatusBadRequest)
		return
	}

	// send email to user
	emailRequest := GenerateEmailRequest(req)

	err5 := SendEmail(emailRequest)

	if err5 != nil {
		log.Println("Error sending email:", err5)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// return success response
	w.WriteHeader(http.StatusOK)
}

func GenerateEmailRequest(request ContactRequest) EmailRequest {
	var output EmailRequest

	output.To = "ajay.bhattacharyya.16@gmail.com"
	output.Subject = "LinkWire Contact Form Submission"

	output.Body = "<p><strong>Name:</strong> " + request.Name + "</p>"
	output.Body += "<p><strong>Email:</strong> " + request.ReturnEmail + "</p>"
	output.Body += "<p><strong>Subject:</strong> " + request.Subject + "</p>"
	output.Body += "<br/>"
	output.Body += "<p><strong>Message:</strong> " + request.Message + "</p>"

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