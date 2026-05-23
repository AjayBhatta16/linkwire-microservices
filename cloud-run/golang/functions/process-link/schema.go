package myfunction

type Request struct {
	LinkID string `json:"linkID"`
}

type PageMeta struct {
	Title       string
	Description string
	OGImage     string
	OGTitle     string
}