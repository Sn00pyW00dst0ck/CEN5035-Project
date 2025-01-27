package health

import "net/http"

// @description Returns a simple 200 status code back, this is useful for pinging the server for checking health, etc
// @Success 200
// @Router /health [get]
func Handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(200)
}
