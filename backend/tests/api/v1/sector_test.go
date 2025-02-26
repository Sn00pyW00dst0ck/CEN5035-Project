package v1Test

import (
	"app/internal/api"
	v1 "app/internal/api/v1"
	"app/internal/database"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func setupSuite(t *testing.T) (*httptest.Server, *v1.SectorAPI, func(t *testing.T)) {
	tmpDir, clean := database.TestingTempDir(t, "sectordb_cache_test")

	router := mux.NewRouter().StrictSlash(true)
	testSectorAPI := v1.NewTestingSector(context.Background(), "log_test.txt", tmpDir, t)
	api.AddV1SectorAPIToRouter(router, testSectorAPI)

	server := httptest.NewServer(router)

	return server, testSectorAPI, func(t *testing.T) {
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
	server, sectorAPI, teardownSuite := setupSuite(t)
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
			id := uuid.New()
			now := time.Now()
			original := v1.Account{
				Id:         id,
				Username:   "Test Account 2",
				ProfilePic: "",
				CreatedAt:  &now,
			}
			_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
			if err != nil {
				panic(err)
			}

			// Test updating the account!
			body := v1.PutAccountJSONRequestBody{
				Id:       id,
				Username: "Updated Username",
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

			require.Equal(t, original.Id, createdAccount.Id)
			require.Equal(t, body.Username, createdAccount.Username)
			require.Equal(t, original.ProfilePic, createdAccount.ProfilePic)
		})

		t.Run("Get By Id", func(t *testing.T) {
			id := uuid.New()
			now := time.Now()
			original := v1.Account{
				Id:         id,
				Username:   "Test Account 2",
				ProfilePic: "",
				CreatedAt:  &now,
			}
			_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
			if err != nil {
				panic(err)
			}

			// Get the account
			response, err := testClient.GetAccountByIDWithResponse(context.Background(), id)
			if err != nil {
				panic(err)
			}
			require.Equal(t, 200, response.StatusCode())

			var createdAccount v1.Account
			err = json.Unmarshal(response.Body, &createdAccount)
			if err != nil {
				panic(err)
			}

			require.Equal(t, original.Id, createdAccount.Id)
			require.Equal(t, original.Username, createdAccount.Username)
			require.Equal(t, original.ProfilePic, createdAccount.ProfilePic)
		})

		t.Run("Search Accounts", func(t *testing.T) {
			now := time.Now()
			then := now.AddDate(0, 0, -7)
			accounts := make([]interface{}, 5)
			accounts[0] = v1.StructToMap(v1.Account{
				Id:         uuid.New(),
				CreatedAt:  &now,
				Username:   "John Doe",
				ProfilePic: "",
			})
			accounts[1] = v1.StructToMap(v1.Account{
				Id:         uuid.New(),
				CreatedAt:  &now,
				Username:   "Jack Doe",
				ProfilePic: "",
			})
			accounts[2] = v1.StructToMap(v1.Account{
				Id:         uuid.New(),
				CreatedAt:  &then,
				Username:   "Maverick",
				ProfilePic: "",
			})
			accounts[3] = v1.StructToMap(v1.Account{
				Id:         uuid.New(),
				CreatedAt:  &now,
				Username:   "w311un1!k3",
				ProfilePic: "",
			})
			accounts[4] = v1.StructToMap(v1.Account{
				Id:         uuid.New(),
				CreatedAt:  &then,
				Username:   "woefullyconsideringlove",
				ProfilePic: "",
			})
			_, err := sectorAPI.DB.Store.PutAll(context.Background(), accounts)
			if err != nil {
				panic(err)
			}

			// Make a query for a subset (all the 'then' date time ones)
			var timeStart = now.AddDate(0, 0, -10)
			var timeEnd = now.AddDate(0, 0, -5)
			query1 := v1.SearchAccountsJSONRequestBody{
				From:  &timeStart,
				Until: &timeEnd,
			}
			result1, err := testClient.SearchAccountsWithResponse(context.Background(), query1)
			if err != nil {
				panic(err)
			}
			require.Equal(t, 200, result1.StatusCode())

			// Ensure the subset is correct
			var queryResult []v1.Account
			err = json.Unmarshal(result1.Body, &queryResult)
			if err != nil {
				panic(err)
			}
			for i := 0; i < len(queryResult); i++ {
				require.LessOrEqual(t, *(queryResult[i].CreatedAt), timeEnd)
				require.GreaterOrEqual(t, *(queryResult[i].CreatedAt), timeStart)
			}

			// Then make a query for a subset
			var username = "Doe"
			query2 := v1.SearchAccountsJSONRequestBody{
				Username: &username,
			}
			result2, err := testClient.SearchAccountsWithResponse(context.Background(), query2)
			if err != nil {
				panic(err)
			}
			require.Equal(t, 200, result2.StatusCode())

			err = json.Unmarshal(result2.Body, &queryResult)
			if err != nil {
				panic(err)
			}
			for i := 0; i < len(queryResult); i++ {
				// Check that the account username string contains the thing we searched for
				require.True(t, strings.Contains(queryResult[i].Username, username))
			}

		})

		t.Run("Delete Account By Id", func(t *testing.T) {
			id := uuid.New()
			now := time.Now()
			original := v1.Account{
				Id:         id,
				Username:   "Test Account 2",
				ProfilePic: "",
				CreatedAt:  &now,
			}
			_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
			if err != nil {
				panic(err)
			}

			// Delete the account
			response, err := testClient.DeleteAccountByIDWithResponse(context.Background(), id)
			if err != nil {
				panic(err)
			}

			require.Equal(t, 204, response.StatusCode())
		})
	})

}
