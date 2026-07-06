package myfunction

type ContactRequest struct {
	Name         string `json:"name"`
	ReturnEmail  string `json:"returnEmail"`
	Subject	     string `json:"subject"`
	Message      string `json:"message"`
}

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}