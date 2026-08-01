package myfunction

import (
	"net/http"

	"github.com/AjayBhatta16/linkwire-golang-shared/utilities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	utilities.ApplyDefaultHeaders(w, r, "POST")

	clearCookie := "token=; HttpOnly; Secure; Path=/; Max-Age=0; SameSite=None"

	w.Header().Set("Set-Cookie", clearCookie)
	w.WriteHeader(http.StatusAccepted)
}
