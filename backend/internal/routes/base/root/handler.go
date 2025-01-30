package root

import (
	"app/internal/responder"
	"net/http"
)

// @Router / [get]
func Handler(w http.ResponseWriter, _ *http.Request) {
	// send a simple text response back
	responder.PlainText(w, "Welcome to Gorilla Mux on Railway!", http.StatusOK)
}
