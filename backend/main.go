package main

import (
	v1 "app/internal/api/v1"
	"app/internal/middleware"
	"context"
	"embed"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/rs/cors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
	oapimiddleware "github.com/oapi-codegen/nethttp-middleware"
)

//go:embed swagger-ui.html
var swaggerUI embed.FS

func main() {
	port := flag.String("port", "3000", "Port for HTTP server")
	flag.Parse()

	// Setup the swagger docs.
	swaggerV1, err := v1.GetSwagger()
	if err != nil {
		panic(err)
	}
	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swaggerV1.Servers = nil
	fixSwaggerPrefix("/v1/api", swaggerV1)

	// Create an instance of the handler which satisfies the generated interface
	sectorAPI := v1.NewSector(context.Background(), "log.txt", "cache", "/orbitdb/12D3KooWFQ9G9VhLeQZX4vtCnUF6bc4ije4ri1oQcpatbWtUbaAn/sectordb")

	// Setup the gorilla mux server with logging.
	r := mux.NewRouter().StrictSlash(true)
	r.Use(middleware.RequestLogger(sectorAPI.Logger))

	// Serve the swagger.json file directly at /docs/swagger.json
	r.HandleFunc("/v1/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(swaggerV1)
	})

	// Serve the Swagger UI using a CDN
	r.HandleFunc("/v1/swagger-ui/", func(w http.ResponseWriter, r *http.Request) {
		data, err := swaggerUI.ReadFile("swagger-ui.html")
		if err != nil {
			http.Error(w, "swagger.html not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(data))
	})

	// Subrouter to validate requests to the /v1/api/
	v1.HandlerWithOptions(sectorAPI, v1.GorillaServerOptions{
		BaseURL:     "/v1/api",
		BaseRouter:  r,
		Middlewares: []v1.MiddlewareFunc{oapimiddleware.OapiRequestValidator(swaggerV1)},
	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
	})

	// Serve HTTP
	s := &http.Server{
		Handler: c.Handler(r),
		Addr:    net.JoinHostPort("127.0.0.1", *port),
	}
	log.Fatal(s.ListenAndServe())
}

// Necessary so that we can use the OapiRequestValidator with a BaseURL
func fixSwaggerPrefix(prefix string, swagger *openapi3.T) {
	var updatedPaths openapi3.Paths = openapi3.Paths{}

	for key, value := range swagger.Paths.Map() {
		updatedPaths.Set(prefix+key, value)
	}

	swagger.Paths = &updatedPaths
}
