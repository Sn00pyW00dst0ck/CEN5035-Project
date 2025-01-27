package router

import (
	"app/internal/middleware"
	"app/internal/routes/api"
	"app/internal/routes/base"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "app/docs"
)

// initialize a new mux router here
var MuxRouter = mux.NewRouter().StrictSlash(true)

func init() {
	// global middleware is registered here
	MuxRouter.Use(handlers.RecoveryHandler())
	MuxRouter.Use(middleware.TrustProxy(middleware.PrivateRanges()))
	MuxRouter.Use(middleware.Logger())

	MuxRouter.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// register base routes here, eg '/' and '/health'
	base.RegisterRoutes(MuxRouter)

	// register route prefixes here, eg '/api/...'
	api.RegisterPrefix(MuxRouter)
}
