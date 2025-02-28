package v1Test

import (
	"app/internal/api"
	v1 "app/internal/api/v1"
	"app/internal/database"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/oapi-codegen/runtime/types"
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
			require.NoError(t, err)
			defer sectorAPI.DB.Store.Delete(context.Background(), body.Id.String())
			require.Equal(t, 201, response.StatusCode())

			var createdAccount v1.Account
			err = json.Unmarshal(response.Body, &createdAccount)
			require.NoError(t, err)
			require.Equal(t, body.Id, createdAccount.Id)
			require.NotNil(t, createdAccount.CreatedAt)
			require.Equal(t, body.Username, createdAccount.Username)
			require.Equal(t, body.ProfilePic, createdAccount.ProfilePic)

			// Test updating the account via this route (should fail)
			body.Username = "Updated Username!"
			response, err = testClient.PutAccountWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Update Account By Id", func(t *testing.T) {
			original := v1.Account{
				Id:         uuid.New(),
				Username:   "Test Account 2",
				ProfilePic: "",
			}
			_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
			require.NoError(t, err)
			defer sectorAPI.DB.Store.Delete(context.Background(), original.Id.String())

			// Test updating the account username
			newUsername := "Updated Test Username"
			body := v1.UpdateAccountByIDJSONRequestBody{
				Username: &(newUsername),
			}
			response, err := testClient.UpdateAccountByIDWithResponse(context.Background(), original.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var updatedAccount v1.Account
			err = json.Unmarshal(response.Body, &updatedAccount)
			require.NoError(t, err)
			require.Equal(t, original.Id, updatedAccount.Id)
			require.Equal(t, newUsername, updatedAccount.Username)
			require.Equal(t, original.ProfilePic, updatedAccount.ProfilePic)
		})

		t.Run("Delete Account By Id", func(t *testing.T) {
			now := time.Now()
			original := v1.Account{
				Id:         uuid.New(),
				Username:   "Test Account 2",
				ProfilePic: "",
				CreatedAt:  &now,
			}
			_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
			require.NoError(t, err)

			// Delete the account
			response, err := testClient.DeleteAccountByIDWithResponse(context.Background(), original.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			// Second time should give internal server error
			response, err = testClient.DeleteAccountByIDWithResponse(context.Background(), original.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get By Id", func(t *testing.T) {
			now := time.Now()
			original := v1.Account{
				Id:         uuid.New(),
				Username:   "Test Account 2",
				ProfilePic: "",
				CreatedAt:  &now,
			}
			_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
			require.NoError(t, err)

			// Get the account
			response, err := testClient.GetAccountByIDWithResponse(context.Background(), original.Id)
			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode())

			var createdAccount v1.Account
			err = json.Unmarshal(response.Body, &createdAccount)
			require.NoError(t, err)

			require.Equal(t, original.Id, createdAccount.Id)
			require.Equal(t, original.Username, createdAccount.Username)
			require.Equal(t, original.ProfilePic, createdAccount.ProfilePic)

			// Remove all store content
			sectorAPI.DB.Store.Delete(context.Background(), original.Id.String())
		})

		t.Run("Search Accounts", func(t *testing.T) {
			queryID := uuid.New()
			now := time.Now()
			then := now.AddDate(0, 0, -7)
			accounts := make([]interface{}, 5)
			accounts[0] = v1.StructToMap(v1.Account{
				Id:         queryID,
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
			require.NoError(t, err)

			// Make a query for a subset (all the 'then' date time ones)
			var timeStart = now.AddDate(0, 0, -10)
			var timeEnd = now.AddDate(0, 0, -5)
			query1 := v1.SearchAccountsJSONRequestBody{
				From:  &timeStart,
				Until: &timeEnd,
			}
			result1, err := testClient.SearchAccountsWithResponse(context.Background(), query1)
			require.NoError(t, err)
			require.Equal(t, 200, result1.StatusCode())

			// Ensure the subset is correct
			var queryResult []v1.Account
			err = json.Unmarshal(result1.Body, &queryResult)
			require.NoError(t, err)
			require.Equal(t, 2, len(queryResult))
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
			require.NoError(t, err)
			require.Equal(t, 200, result2.StatusCode())

			err = json.Unmarshal(result2.Body, &queryResult)
			require.NoError(t, err)
			require.Equal(t, 2, len(queryResult))
			for i := 0; i < len(queryResult); i++ {
				// Check that the account username string contains the thing we searched for
				require.True(t, strings.Contains(queryResult[i].Username, username))
			}

			// Make a third query to test id searches
			var ids = []types.UUID{queryID}
			query3 := v1.SearchAccountsJSONRequestBody{
				Id: &ids,
			}
			result3, err := testClient.SearchAccountsWithResponse(context.Background(), query3)
			require.NoError(t, err)
			require.Equal(t, 200, result3.StatusCode())

			err = json.Unmarshal(result3.Body, &queryResult)
			require.NoError(t, err)
			require.Equal(t, 1, len(queryResult))

			// Clean up
			for _, value := range accounts {
				if val, ok := value.(v1.Account); ok {
					sectorAPI.DB.Store.Delete(context.Background(), val.Id.String())
				}
			}
		})
	})

	/*

		t.Run("Group", func(t *testing.T) {
			t.Run("Create Group", func(t *testing.T) {
				body := v1.PutGroupJSONRequestBody{
					Id:          uuid.New(),
					Name:        "Test Group",
					Description: "Used for unit tests!",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}

				response, err := testClient.PutGroupWithResponse(context.Background(), body)
				if err != nil {
					panic(err)
				}

				require.Equal(t, 201, response.StatusCode())

				var createdGroup v1.Group
				data, err := base64.StdEncoding.DecodeString(string(response.Body)[1 : len(string(response.Body))-2])
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(data, &createdGroup)
				if err != nil {
					panic(err)
				}

				require.Equal(t, body.Id, createdGroup.Id)
				require.Equal(t, body.Name, createdGroup.Name)
				require.Equal(t, body.Description, createdGroup.Description)
				require.Equal(t, body.Members, createdGroup.Members)
				require.Equal(t, body.Channels, createdGroup.Channels)

				// Test updating the group by wrong endpoing (should fail)!
				body.Name = "Updated Group Name!"
				response, err = testClient.PutGroupWithResponse(context.Background(), body)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 500, response.StatusCode())

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), body.Id.String())
			})

			t.Run("Get Group By Id", func(t *testing.T) {
				id := uuid.New()
				now := time.Now()
				original := v1.Group{
					Id:          id,
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
				if err != nil {
					panic(err)
				}

				// Get the account
				response, err := testClient.GetGroupByIDWithResponse(context.Background(), original.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 200, response.StatusCode())

				var fetchedGroup v1.Group
				err = json.Unmarshal(response.Body, &fetchedGroup)
				if err != nil {
					panic(err)
				}

				require.Equal(t, original.Id, fetchedGroup.Id)
				require.Equal(t, original.Name, fetchedGroup.Name)
				require.Equal(t, original.Description, fetchedGroup.Description)
				require.Equal(t, original.Members, fetchedGroup.Members)
				require.Equal(t, original.Channels, fetchedGroup.Channels)

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), id.String())
			})

			t.Run("Update Group By Id", func(t *testing.T) {
				id := uuid.New()
				now := time.Now()
				original := v1.Group{
					Id:          id,
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
				if err != nil {
					panic(err)
				}

				// Update the group
				body := v1.UpdateGroupByIDJSONRequestBody{
					Id:          id,
					CreatedAt:   &now,
					Name:        "Updated Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				response, err := testClient.UpdateGroupByIDWithResponse(context.Background(), body.Id, body)
				if err != nil {
					sectorAPI.DB.Store.Delete(context.Background(), id.String())
					panic(err)
				}
				require.Equal(t, 201, response.StatusCode())

				var updatedGroup v1.Group
				data, err := base64.StdEncoding.DecodeString(string(response.Body)[1 : len(string(response.Body))-2])
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(data, &updatedGroup)
				if err != nil {
					panic(err)
				}

				require.Equal(t, body.Id, updatedGroup.Id)
				require.Equal(t, body.Name, updatedGroup.Name)
				require.Equal(t, original.Description, updatedGroup.Description)
				require.Equal(t, original.Members, updatedGroup.Members)
				require.Equal(t, original.Channels, updatedGroup.Channels)

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), id.String())
			})

			t.Run("Delete Group By Id", func(t *testing.T) {
				id := uuid.New()
				now := time.Now()
				original := v1.Group{
					Id:          id,
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(original))
				if err != nil {
					panic(err)
				}

				// Delete the account
				response, err := testClient.DeleteGroupByIDWithResponse(context.Background(), id)
				if err != nil {
					sectorAPI.DB.Store.Delete(context.Background(), id.String())
					panic(err)
				}
				require.Equal(t, 204, response.StatusCode())

				// Second time should give internal server error
				response, err = testClient.DeleteGroupByIDWithResponse(context.Background(), id)
				if err != nil {
					sectorAPI.DB.Store.Delete(context.Background(), id.String())
					panic(err)
				}
				require.Equal(t, 500, response.StatusCode())
			})

			// TODO: search test

			t.Run("Add Member", func(t *testing.T) {
				now := time.Now()
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}

				originalAccout := v1.Account{
					Id:         uuid.New(),
					Username:   "Test Account 2",
					ProfilePic: "",
					CreatedAt:  &now,
				}

				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalAccout))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}

				response, err := testClient.AddGroupMemberWithResponse(context.Background(), originalGroup.Id, originalAccout.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 201, response.StatusCode())
				// TODO: Check that member list updated
			})

			t.Run("Remove Member", func(t *testing.T) {
				now := time.Now()
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}

				originalAccout := v1.Account{
					Id:         uuid.New(),
					Username:   "Test Account 2",
					ProfilePic: "",
					CreatedAt:  &now,
				}

				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalAccout))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}

				response, err := testClient.RemoveGroupMemberWithResponse(context.Background(), originalGroup.Id, originalAccout.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 204, response.StatusCode())
				// TODO: Check that member list updated
			})
		})

		t.Run("Channel", func(t *testing.T) {
			t.Run("Create Channel", func(t *testing.T) {
				now := time.Now()
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}

				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}

				desc := "Channel for unit testing"
				body := v1.PutChannelJSONRequestBody{
					Id:             uuid.New(),
					CreatedAt:      &now,
					Name:           "Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				response, err := testClient.PutChannelWithResponse(context.Background(), originalGroup.Id, body)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 201, response.StatusCode())

				var createdChannel v1.Channel
				data, err := base64.StdEncoding.DecodeString(string(response.Body)[1 : len(string(response.Body))-2])
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(data, &createdChannel)
				if err != nil {
					panic(err)
				}

				require.Equal(t, body.Id, createdChannel.Id)
				require.Equal(t, body.Name, createdChannel.Name)
				require.Equal(t, body.Description, createdChannel.Description)
				require.Equal(t, body.Messages, createdChannel.Messages)
				require.Equal(t, body.PinnedMessages, createdChannel.PinnedMessages)

				// Get group by ID and check list has updated
				response2, err := testClient.GetGroupByIDWithResponse(context.Background(), originalGroup.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 200, response2.StatusCode())

				var fetchedGroup v1.Group
				err = json.Unmarshal(response2.Body, &fetchedGroup)
				if err != nil {
					panic(err)
				}

				require.Equal(t, originalGroup.Id, fetchedGroup.Id)
				require.Equal(t, originalGroup.Name, fetchedGroup.Name)
				require.Equal(t, originalGroup.Description, fetchedGroup.Description)
				require.Equal(t, originalGroup.Members, fetchedGroup.Members)
				require.NotEqual(t, originalGroup.Channels, fetchedGroup.Channels)

				// Cleanup
				sectorAPI.DB.Store.Delete(context.Background(), body.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
			})

			t.Run("Get Channel By Id", func(t *testing.T) {
				now := time.Now()
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				desc := "Channel for unit testing"
				originalChannel := v1.Channel{
					Id:             uuid.New(),
					CreatedAt:      &now,
					Name:           "Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalChannel))
				if err != nil {
					panic(err)
				}

				// Get the channel
				response, err := testClient.GetChannelByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 200, response.StatusCode())

				var fetchedChannel v1.Channel
				err = json.Unmarshal(response.Body, &fetchedChannel)
				if err != nil {
					panic(err)
				}

				require.Equal(t, originalChannel.Id, fetchedChannel.Id)
				require.Equal(t, originalChannel.Name, fetchedChannel.Name)
				require.Equal(t, originalChannel.Description, fetchedChannel.Description)
				require.Equal(t, originalChannel.Messages, fetchedChannel.Messages)
				require.Equal(t, originalChannel.PinnedMessages, fetchedChannel.PinnedMessages)

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalChannel.Id.String())
			})

			t.Run("Update Channel By Id", func(t *testing.T) {
				now := time.Now()
				desc := "Channel for unit testing"
				originalChannel := v1.Channel{
					Id:             uuid.New(),
					CreatedAt:      &now,
					Name:           "Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{originalChannel.Id},
				}

				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalChannel))
				if err != nil {
					panic(err)
				}

				// Update the Channel
				body := v1.UpdateChannelByIDJSONRequestBody{
					Id:             originalChannel.Id,
					CreatedAt:      &now,
					Name:           "Updated Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				response, err := testClient.UpdateChannelByIDWithResponse(context.Background(), body.Id, originalChannel.Id, body)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 201, response.StatusCode())

				var updatedChannel v1.Channel
				data, err := base64.StdEncoding.DecodeString(string(response.Body)[1 : len(string(response.Body))-2])
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(data, &updatedChannel)
				if err != nil {
					panic(err)
				}

				require.Equal(t, body.Id, updatedChannel.Id)
				require.Equal(t, body.Name, updatedChannel.Name)
				require.Equal(t, originalChannel.Description, updatedChannel.Description)

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalChannel.Id.String())
			})

			t.Run("Delete Channel By Id", func(t *testing.T) {
				now := time.Now()
				desc := "Channel for unit testing"
				originalChannel := v1.Channel{
					Id:             uuid.New(),
					CreatedAt:      &now,
					Name:           "Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{originalChannel.Id},
				}

				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalChannel))
				if err != nil {
					panic(err)
				}

				// Attempt delete once
				response, err := testClient.DeleteChannelByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 204, response.StatusCode())

				// Check group is updated
				response2, err := testClient.GetGroupByIDWithResponse(context.Background(), originalGroup.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 200, response2.StatusCode())

				var fetchedGroup v1.Group
				err = json.Unmarshal(response2.Body, &fetchedGroup)
				if err != nil {
					panic(err)
				}

				require.Equal(t, originalGroup.Id, fetchedGroup.Id)
				require.Equal(t, originalGroup.Name, fetchedGroup.Name)
				require.Equal(t, originalGroup.Description, fetchedGroup.Description)
				require.Equal(t, originalGroup.Members, fetchedGroup.Members)
				require.NotEqual(t, originalGroup.Channels, fetchedGroup.Channels)

				// Attempt delete twice (should fail)
				response, err = testClient.DeleteChannelByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 500, response.StatusCode())

				// Remove all leftover store content
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
			})

			// TODO: test for search

		})

		t.Run("Message", func(t *testing.T) {
			t.Run("Create Message", func(t *testing.T) {
				now := time.Now()
				originalAccount := v1.Account{
					Id:         uuid.New(),
					CreatedAt:  &now,
					Username:   "Test User",
					ProfilePic: "",
				}
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				desc := "Channel for unit testing"
				originalChannel := v1.Channel{
					Id:             uuid.New(),
					CreatedAt:      &now,
					Name:           "Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalAccount))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalChannel))
				if err != nil {
					panic(err)
				}

				body := v1.PutMessageJSONRequestBody{
					Id:        uuid.New(),
					CreatedAt: &now,
					Author:    originalAccount.Id,
					Body:      "Message Content",
				}
				response, err := testClient.PutMessageWithResponse(context.Background(), originalGroup.Id, originalChannel.Id, body)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 201, response.StatusCode())

				var createdMessage v1.Message
				data, err := base64.StdEncoding.DecodeString(string(response.Body)[1 : len(string(response.Body))-2])
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(data, &createdMessage)
				if err != nil {
					panic(err)
				}

				// Check content
				require.Equal(t, body.Id, createdMessage.Id)
				require.Equal(t, body.Author, createdMessage.Author)
				require.Equal(t, body.Body, createdMessage.Body)

				// Check channel's message list updated
				response2, err := testClient.GetChannelByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 200, response2.StatusCode())

				var fetchedChannel v1.Channel
				err = json.Unmarshal(response2.Body, &fetchedChannel)
				if err != nil {
					panic(err)
				}

				require.Equal(t, originalChannel.Id, fetchedChannel.Id)
				require.Equal(t, originalChannel.Name, fetchedChannel.Name)
				require.Equal(t, originalChannel.Description, fetchedChannel.Description)
				require.Equal(t, 1, len(fetchedChannel.Messages))
				require.Equal(t, 0, len(fetchedChannel.PinnedMessages))

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), originalAccount.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalChannel.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), body.Id.String())
			})

			t.Run("Get Message By Id", func(t *testing.T) {
				now := time.Now()
				originalAccount := v1.Account{
					Id:         uuid.New(),
					CreatedAt:  &now,
					Username:   "Test User",
					ProfilePic: "",
				}
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				desc := "Channel for unit testing"
				originalChannel := v1.Channel{
					Id:             uuid.New(),
					CreatedAt:      &now,
					Name:           "Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				originalMessage := v1.Message{
					Id:        uuid.New(),
					CreatedAt: &now,
					Author:    originalAccount.Id,
					Body:      "Test Message",
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalAccount))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalChannel))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalMessage))
				if err != nil {
					panic(err)
				}

				// Get the message
				response, err := testClient.GetMessageByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id, originalMessage.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 200, response.StatusCode())

				var fetchedMessage v1.Message
				err = json.Unmarshal(response.Body, &fetchedMessage)
				if err != nil {
					panic(err)
				}

				require.Equal(t, originalMessage.Id, fetchedMessage.Id)
				require.Equal(t, originalMessage.Author, fetchedMessage.Author)
				require.Equal(t, originalMessage.Body, fetchedMessage.Body)

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), originalAccount.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalChannel.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalMessage.Id.String())
			})

			t.Run("Update Message By Id", func(t *testing.T) {
				now := time.Now()
				originalAccount := v1.Account{
					Id:         uuid.New(),
					CreatedAt:  &now,
					Username:   "Test User",
					ProfilePic: "",
				}
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				desc := "Channel for unit testing"
				originalChannel := v1.Channel{
					Id:             uuid.New(),
					CreatedAt:      &now,
					Name:           "Test Channel",
					Description:    &desc,
					Messages:       []types.UUID{},
					PinnedMessages: []types.UUID{},
				}
				originalMessage := v1.Message{
					Id:        uuid.New(),
					CreatedAt: &now,
					Author:    originalAccount.Id,
					Body:      "Test Message",
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalAccount))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalChannel))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalMessage))
				if err != nil {
					panic(err)
				}

				// Update the Message
				body := v1.UpdateMessageByIDJSONRequestBody{
					Id:        originalChannel.Id,
					Author:    originalMessage.Author,
					CreatedAt: &now,
					Body:      "Updated Test Message",
				}
				response, err := testClient.UpdateMessageByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id, originalMessage.Id, body)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 201, response.StatusCode())

				var updatedMessage v1.Message
				data, err := base64.StdEncoding.DecodeString(string(response.Body)[1 : len(string(response.Body))-2])
				if err != nil {
					panic(err)
				}
				err = json.Unmarshal(data, &updatedMessage)
				if err != nil {
					panic(err)
				}

				require.Equal(t, body.Id, updatedMessage.Id)
				require.Equal(t, body.Author, updatedMessage.Author)
				require.Equal(t, body.Body, updatedMessage.Body)

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), originalAccount.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalChannel.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalMessage.Id.String())
			})

			t.Run("Delete Message By Id", func(t *testing.T) {
				now := time.Now()
				originalAccount := v1.Account{
					Id:         uuid.New(),
					CreatedAt:  &now,
					Username:   "Test User",
					ProfilePic: "",
				}
				originalGroup := v1.Group{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Group",
					Description: "Group for unit testing",
					Members:     []types.UUID{},
					Channels:    []types.UUID{},
				}
				desc := "Channel for unit testing"
				originalChannel := v1.Channel{
					Id:          uuid.New(),
					CreatedAt:   &now,
					Name:        "Test Channel",
					Description: &desc,
				}
				originalMessage := v1.Message{
					Id:        uuid.New(),
					Channel:   originalChannel.Id,
					CreatedAt: &now,
					Author:    originalAccount.Id,
					Body:      "Test Message",
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalAccount))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalGroup))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalChannel))
				if err != nil {
					panic(err)
				}
				_, err = sectorAPI.DB.Store.Put(context.Background(), v1.StructToMap(originalMessage))
				if err != nil {
					panic(err)
				}

				// Attempt delete once
				response, err := testClient.DeleteMessageByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id, originalMessage.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 204, response.StatusCode())

				// Check group is updated
				response2, err := testClient.GetChannelByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 200, response2.StatusCode())

				var fetchedChannel v1.Channel
				err = json.Unmarshal(response2.Body, &fetchedChannel)
				if err != nil {
					panic(err)
				}
				require.Equal(t, originalChannel.Id, fetchedChannel.Id)
				require.Equal(t, originalChannel.Name, fetchedChannel.Name)
				require.Equal(t, originalChannel.Description, fetchedChannel.Description)

				// Attempt delete twice (should fail)
				response, err = testClient.DeleteMessageByIDWithResponse(context.Background(), originalGroup.Id, originalChannel.Id, originalMessage.Id)
				if err != nil {
					panic(err)
				}
				require.Equal(t, 500, response.StatusCode())

				// Remove all store content
				sectorAPI.DB.Store.Delete(context.Background(), originalAccount.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalGroup.Id.String())
				sectorAPI.DB.Store.Delete(context.Background(), originalChannel.Id.String())
			})

			// TODO: test for search

		})

	*/
}
