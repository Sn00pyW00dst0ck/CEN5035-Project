package main

import (
	"app/internal/api"
	v1 "app/internal/api/v1"
	"context"
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/rs/cors"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

func main() {
	port := flag.String("port", "3000", "Port for HTTP server")
	flag.Parse()

	r := mux.NewRouter().StrictSlash(true)
	sectorAPI := v1.NewSector(context.WithoutCancel(openapi3.NewLoader().Context), "log.txt", "cache")
	api.AddV1SectorAPIToRouter(r, sectorAPI)

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
