package v1

import (
	"app/internal/routes/api/v1/forecast"
	"app/internal/routes/api/v1/temperature"
	"net/http"

	"github.com/gorilla/mux"
)

func registerRoutes(router *mux.Router, tempHandler *temperature.TemperatureHandler, forecastHandler *forecast.ForcastHandler) {

	// temperature handler
	router.HandleFunc("/temperature", tempHandler.GetTemperatureHandler).Methods(http.MethodGet)

	// forecast handler
	router.HandleFunc("/forecast/{forecastPeriod}", forecastHandler.GetThreeDayForcastHandler).Methods(http.MethodGet)
}
