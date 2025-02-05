package temperature

import (
	"app/internal/responder"
	"net/http"

	"go.uber.org/zap"
)

type TemperatureHandler struct {
	Logger *zap.Logger
}

func NewTemperatureHandler(logger *zap.Logger) *TemperatureHandler {
	return &TemperatureHandler{Logger: logger}
}

// mock api response
func (h *TemperatureHandler) GetTemperatureHandler(w http.ResponseWriter, _ *http.Request) {
	msg := map[string]string{"temp_c": "18.8", "wind_speed_km": "10", "humidity_percent": "70"}

	if err := responder.JSONPretty(w, msg, http.StatusOK); err != nil {
		h.Logger.Error(err.Error())
	}
}
