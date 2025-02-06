package main

import (
	"app/internal/api"
	"app/internal/middleware"
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	nethttp_middleware "github.com/oapi-codegen/nethttp-middleware"
)

func main() {
	port := flag.String("port", "3000", "Port for HTTP server")
	flag.Parse()

	// Setup the swagger docs.
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create an instance of the handler which satisfies the generated interface
	sectorAPI := api.NewSector(context.Background(), "log.txt", "cache", "/orbitdb/bafyreiejrtaennxufa3wvkdvyoj6ywq6nid3lukdqcnx2fc33tckzjzbke/sectordb")

	// Setup the gorilla mux server.
	// Use validation middleware to check all requests against the OpenAPI schema.
	// Then define the sectorAPI as the one to handle that schema.
	r := mux.NewRouter()
	r.Use(nethttp_middleware.OapiRequestValidator(swagger))
	r.Use(middleware.RequestLogger(sectorAPI.Logger))
	api.HandlerFromMux(sectorAPI, r)

	// Serve HTTP
	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("127.0.0.1", *port),
	}
	log.Fatal(s.ListenAndServe())
}
