package server

import (
	"app/internal/database"
	"app/internal/router"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func StartServer(address string, db *database.Database, logger *zap.Logger) error {
	muxRouter := router.InitRouter(db, logger)

	srv := &http.Server{
		Handler: &muxRouter,
		Addr:    address,
		// Good practice: enforce timeouts for servers you create!
		// adjust as needed
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}

	logger.Info("Starting server on port: " + address[:])
	log.Println("Server started at http://localhost:3000/")
	return srv.ListenAndServe()
}
