package v1Test

import (
	"app/internal/api"
	v1 "app/internal/api/v1"
	"app/internal/database"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func setupSuite(t *testing.T) (*httptest.Server, func(t *testing.T)) {
	tmpDir, clean := database.TestingTempDir(t, "sectordb_cache_test")

	router := mux.NewRouter().StrictSlash(true)
	testSectorAPI := v1.NewTestingSector(context.Background(), "log_test.txt", tmpDir, t)
	api.AddV1SectorAPIToRouter(router, testSectorAPI)

	server := httptest.NewServer(router)

	return server, func(t *testing.T) {
		testSectorAPI.DB.Disconnect()
		server.Close()
		clean()
	}
}

func setupTest(tb testing.TB) func(tb testing.TB) {

	return func(tb testing.TB) {

	}
}

func TestSectorV1(t *testing.T) {
	// Setup the server to run for all the tests
	server, teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	// Setup the test client
	testClient, err := v1.NewClientWithResponses(server.URL, v1.WithHTTPClient(server.Client()), v1.WithBaseURL(server.URL+"/v1/api"))
	if err != nil {
		panic(err)
	}

	t.Run("Get Root", func(t *testing.T) {
		response, err := testClient.GetRootWithResponse(context.Background())
		if err != nil {
			panic(err)
		}
		require.Equal(t, 200, response.StatusCode())
		require.Equal(t, "{\"message\":\"Hello, World!\"}\n", string(response.Body))
	})

	t.Run("Get Health", func(t *testing.T) {
		response, err := testClient.GetHealthWithResponse(context.Background())
		if err != nil {
			panic(err)
		}
		require.Equal(t, 200, response.StatusCode())
	})

	t.Run("Account", func(t *testing.T) {
		t.Run("Create Account", func(t *testing.T) {
			body := v1.PutAccountJSONRequestBody{
				Id:         uuid.New(),
				Username:   "CreateAccount",
				ProfilePic: "",
			}

			response, err := testClient.PutAccountWithResponse(context.Background(), body)
			if err != nil {
				panic(err)
			}

			require.Equal(t, 200, response.StatusCode())

			var createdAccount v1.Account
			data, err := base64.StdEncoding.DecodeString(string(response.Body)[1 : len(string(response.Body))-2])
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(data, &createdAccount)
			if err != nil {
				panic(err)
			}

			require.Equal(t, body.Id, createdAccount.Id)
			require.Equal(t, body.Username, createdAccount.Username)
			require.Equal(t, body.ProfilePic, createdAccount.ProfilePic)
		})

		t.Run("Update Account", func(t *testing.T) {
			// Add an accounto to the DB
			// Make a query to update this
			// Then check response is correct

		})

		t.Run("Search Account", func(t *testing.T) {
			// Add a bunch of accounts specific for this to the DB
			// Then make a query for a subset
			// Ensure the subset is correct
		})

		t.Run("Delete Account", func(t *testing.T) {
			// Add an accounto to the DB
			// Make a query to delete this
			// Then check response is correct
		})
	})

}
