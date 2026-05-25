package myfunction

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/resend/resend-go/v3"

    "github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	subscriber, err := utilities.NewSubscriber(context.Background())

	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating subscriber: %v", err), http.StatusInternalServerError)
		return
	}

	defer subscriber.Close()

	req, _, err := utilities.Receive[Request](subscriber, r)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error receiving message: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Received message - To: %s, Subject: %s", req.To, req.Subject)

	err2 := SendEmail(req.To, req.Subject, req.Body)

	if err2 != nil {
		http.Error(w, fmt.Sprintf("Error sending email: %v", err2), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func SendEmail(to string, subject string, htmlContent string) error {
	resendKey := os.Getenv("RESEND_API_KEY")

	if resendKey == "" {
		return fmt.Errorf("RESEND_API_KEY environment variable is not set")
	}

	client := resend.NewClient(resendKey)

    params := &resend.SendEmailRequest{
        From:    "LinkWire <notifications@linkwire.cc>",
        To:      []string{to},
        Html:    htmlContent,
        Subject: subject,
    }

    sent, err := client.Emails.Send(params)

    if err != nil {
        fmt.Println(err.Error())
        return err
    }

    fmt.Println("Email sent with ID:", sent.Id)

	return nil
}