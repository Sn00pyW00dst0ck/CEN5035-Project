package router

import (
	"app/internal/database"
	"app/internal/middleware"
	"app/internal/routes/api"
	"app/internal/routes/base"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	_ "app/docs"
)

// Creates a new mux router with the given databse and logger.
func InitRouter(db *database.Database, logger *zap.Logger) mux.Router {
	var MuxRouter = mux.NewRouter().StrictSlash(true)

	// Global middleware is registered here.
	MuxRouter.Use(handlers.RecoveryHandler())
	MuxRouter.Use(middleware.TrustProxy(middleware.PrivateRanges(), logger))
	MuxRouter.Use(middleware.Logger(logger))

	MuxRouter.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// register base routes here, eg '/' and '/health'
	base.RegisterRoutes(MuxRouter)

	// register route prefixes here, eg '/api/...'
	api.RegisterPrefix(MuxRouter, db, logger)

	return *MuxRouter
}
