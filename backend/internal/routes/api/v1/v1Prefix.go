package v1

import (
	"app/internal/database"
	"app/internal/routes/api/v1/forecast"
	"app/internal/routes/api/v1/temperature"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func RegisterPrefix(router *mux.Router, db *database.Database, logger *zap.Logger) {
	// create the v1 path prefix
	v1Subrouter := router.PathPrefix("/v1/").Subrouter()

	registerRoutes(v1Subrouter, temperature.NewTemperatureHandler(logger), forecast.NewForcastHandler(logger))
}
