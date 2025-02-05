package forecast

import (
	"app/internal/responder"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ForcastHandler struct {
	Logger *zap.Logger
}

func NewForcastHandler(logger *zap.Logger) *ForcastHandler {
	return &ForcastHandler{Logger: logger}
}

// mock api response
func (h *ForcastHandler) GetThreeDayForcastHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	msg := []map[string]string{}

	switch vars["forecastPeriod"] {
	case "3day":
		msg = append(msg, map[string]string{"temp_c": "18.8", "wind_speed_km": "10", "humidity_percent": "70"})
		msg = append(msg, map[string]string{"temp_c": "25.2", "wind_speed_km": "12", "humidity_percent": "74"})
		msg = append(msg, map[string]string{"temp_c": "27.7", "wind_speed_km": "9", "humidity_percent": "86"})
	default:
		http.Error(w, fmt.Sprintf("unsupported forecast period: %s", vars["forecastPeriod"]), http.StatusBadRequest)
		return
	}

	if err := responder.JSONPretty(w, msg, http.StatusOK); err != nil {
		h.Logger.Error(err.Error())
	}
}
