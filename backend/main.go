package main

import (
	"app/internal/database"
	"app/internal/logger"
	"app/internal/server"
	"app/internal/tools"
	"context"
	"log"
	"os"
)

// @title Sector Swagger API
// @version 1.0
// @description This is the Sector server API.

// @contact.name API Support
// @contact.url
// @contact.email

// @license.name MIT
// @license.url https://opensource.org/license/mit

// @accept json
// @produce json

// @schemes http https

// @host localhost:3000
// @BasePath /

// Runs the server.
func main() {
	port := tools.EnvPortOr("3000")

	// Setup logger to log everything to a file.
	logger, err := logger.NewLogger("log.txt")
	if err != nil {
		log.Panicln(err)
	}

	// Startup the database.
	db, err := database.NewDatabase(context.Background(), "/orbitdb/bafyreiejrtaennxufa3wvkdvyoj6ywq6nid3lukdqcnx2fc33tckzjzbke/sectordb", "cache", logger)
	if err != nil {
		log.Panicln(err)
	}
	defer db.Disconnect()

	err = db.Connect(func(address string) {
		log.Println("Connected: ", address)
	})
	if err != nil {
		panic(err)
	}

	// Start listening on port.
	if err := server.StartServer(port, db, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
