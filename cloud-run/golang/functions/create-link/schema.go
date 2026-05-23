package myfunction

type Request struct {
	RedirectURL string `json:"redirectUrl"`
	Note 	    string `json:"note"`
}

type ProcessLinkRequest struct {
	LinkID string `json:"linkID"`
}

const CODE_CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"