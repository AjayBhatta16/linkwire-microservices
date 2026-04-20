package myfunction

type Request struct {
	RedirectURL string `json:"redirectUrl"`
	Note 	    string `json:"note"`
}

const CODE_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"