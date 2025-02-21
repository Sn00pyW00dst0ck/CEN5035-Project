package v1Test

import (
	"app/internal/api"
	v1 "app/internal/api/v1"
	"app/internal/database"
	"context"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

func TestSectorV1(t *testing.T) {
	// Setup temporary database cache location
	tmpDir, clean := database.TestingTempDir(t, "sectordb_cache_test")
	defer clean()

	r := mux.NewRouter().StrictSlash(true)
	testSectorAPI := v1.NewTestingSector(context.Background(), "log_test.txt", tmpDir, t)
	api.AddV1SectorAPIToRouter(r, testSectorAPI)

	s := &http.Server{
		Handler: r,
		Addr:    net.JoinHostPort("127.0.0.1", "3000"),
	}
	log.Fatal(s.ListenAndServe())
	defer s.Close()

	// Setup the test client
	// testClient, err := v1.NewClientWithResponses("localhost:3000")

	// TODO: use auto-generated test-client to query each endpoint and test responses
	cases := []struct{ Name, Other string }{
		{"Get Root", ""},
		{"Get Health", ""},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {

		})
	}
}
