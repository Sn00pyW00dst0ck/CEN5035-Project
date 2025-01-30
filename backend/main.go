package main

import (
	"app/internal/logger"
	"app/internal/server"
	"app/internal/tools"
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

	logger.Stdout.Info("starting server on port " + port[1:])

	// start listening on port
	if err := server.StartServer(port); err != nil {
		logger.Stderr.Error(err.Error())
		os.Exit(1)
	}
}
