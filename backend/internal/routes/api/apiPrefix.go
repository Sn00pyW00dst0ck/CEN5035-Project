package api

import (
	"app/internal/database"
	v1 "app/internal/routes/api/v1"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func RegisterPrefix(router *mux.Router, db *database.Database, logger *zap.Logger) {
	// create the api prefix
	apiPrefix := router.PathPrefix("/api/").Subrouter()

	// register the v1 prefix
	v1.RegisterPrefix(apiPrefix, db, logger)
}
