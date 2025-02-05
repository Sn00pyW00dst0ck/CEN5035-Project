package base

import (
	"app/internal/database"
	"app/internal/routes/base/health"
	"app/internal/routes/base/root"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	DB     *database.Database
	Logger *zap.Logger
}

func RegisterRoutes(router *mux.Router) {
	baseRouter := router.NewRoute().Subrouter()

	// root handler
	baseRouter.HandleFunc("/", root.Handler).Methods(http.MethodGet)

	// health handler
	baseRouter.HandleFunc("/health", health.Handler).Methods(http.MethodGet, http.MethodHead)
}
