package v1Test

import (
	"Sector/internal/api"
	v1 "Sector/internal/api/v1"
	"Sector/internal/database"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"berty.tech/go-orbit-db/iface"
	"berty.tech/go-orbit-db/stores/documentstore"
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

func setupTest(t *testing.T, api v1.SectorAPI) ([]interface{}, func(t *testing.T)) {
	now := time.Now()
	then := now.AddDate(0, 0, -7)

	// Define the entries of the database for every (TODO: have this be a better set of data)
	entries := []interface{}{
		v1.Account{
			Id:         uuid.New(),
			CreatedAt:  &now,
			ProfilePic: "",
			Username:   "John Doe",
		},
		v1.Account{
			Id:         uuid.New(),
			CreatedAt:  &then,
			ProfilePic: "",
			Username:   "Jack Doe",
		},
		v1.Account{
			Id:         uuid.New(),
			CreatedAt:  &now,
			ProfilePic: "",
			Username:   "Maverick",
		},
		v1.Account{
			Id:         uuid.New(),
			CreatedAt:  &then,
			ProfilePic: "",
			Username:   "w311un1!k3",
		},
		v1.Account{
			Id:         uuid.New(),
			CreatedAt:  &then,
			ProfilePic: "",
			Username:   "woefullyconsideringlove",
		},
		v1.Group{
			Id:          uuid.New(),
			CreatedAt:   &then,
			Name:        "Test Group 1",
			Description: "A group for unit testing.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          uuid.New(),
			CreatedAt:   &now,
			Name:        "Test Group 2",
			Description: "Another unit testing group.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          uuid.New(),
			CreatedAt:   &now,
			Name:        "Test Group 3",
			Description: "A third group for unit testing.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          uuid.New(),
			CreatedAt:   &then,
			Name:        "Advanced Test Group 1",
			Description: "For advanced testing.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          uuid.New(),
			CreatedAt:   &now,
			Name:        "Advanced Test Group 2",
			Description: "For advanced testing.",
			Members:     []types.UUID{},
		},
		// v1.Channel{},
		// v1.Channel{},
		// v1.Channel{},
		// v1.Channel{},
		// v1.Channel{},
		// v1.Message{},
		// v1.Message{},
		// v1.Message{},
		// v1.Message{},
		// v1.Message{},
	}

	// Add the content to the database (have to convert to maps for DB entry)
	result := make([]interface{}, len(entries))
	for i, v := range entries {
		result[i] = v1.StructToMap(v)
	}
	_, err := api.DB.Store.PutAll(context.Background(), result)
	require.NoError(t, err)

	// Clean up all the resources
	return entries, func(t *testing.T) {
		err := api.DB.Store.Drop()
		require.NoError(t, err)
		store, err := api.DB.OrbitDB.Docs(context.Background(), api.DB.Store.Address().String(), &iface.CreateDBOptions{
			StoreSpecificOpts: documentstore.DefaultStoreOptsForMap("id"),
		})
		require.NoError(t, err)
		api.DB.Store = store
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
			_, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			body := v1.PutAccountJSONRequestBody{
				Id:         uuid.New(),
				Username:   "CreateAccount",
				ProfilePic: "",
			}

			response, err := testClient.PutAccountWithResponse(context.Background(), body)
			require.NoError(t, err)
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
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedAccount := entries[0].(v1.Account)

			// Test updating the account username
			newUsername := "Updated Test Username"
			body := v1.UpdateAccountByIDJSONRequestBody{
				Username: &(newUsername),
			}
			response, err := testClient.UpdateAccountByIDWithResponse(context.Background(), selectedAccount.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var updatedAccount v1.Account
			err = json.Unmarshal(response.Body, &updatedAccount)
			require.NoError(t, err)
			expected := v1.Account{
				Id:         selectedAccount.Id,
				CreatedAt:  selectedAccount.CreatedAt,
				Username:   newUsername,
				ProfilePic: selectedAccount.ProfilePic,
			}
			// Must compare all fields individually because timestamps are tricky and suck
			require.Equal(t, expected.Id, updatedAccount.Id)
			require.True(t, expected.CreatedAt.Equal(*updatedAccount.CreatedAt))
			require.Equal(t, expected.Username, updatedAccount.Username)
			require.Equal(t, expected.ProfilePic, updatedAccount.ProfilePic)
		})

		t.Run("Delete Account By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedAccount := entries[0].(v1.Account)

			// Delete the account
			response, err := testClient.DeleteAccountByIDWithResponse(context.Background(), selectedAccount.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			// Second time should give internal server error
			response, err = testClient.DeleteAccountByIDWithResponse(context.Background(), selectedAccount.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedAccount := entries[0].(v1.Account)

			// Get the account
			response, err := testClient.GetAccountByIDWithResponse(context.Background(), selectedAccount.Id)
			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode())

			var fetchedAccount v1.Account
			err = json.Unmarshal(response.Body, &fetchedAccount)
			require.NoError(t, err)
			require.Equal(t, selectedAccount.Id, fetchedAccount.Id)
			require.True(t, selectedAccount.CreatedAt.Equal(*fetchedAccount.CreatedAt))
			require.Equal(t, selectedAccount.Username, fetchedAccount.Username)
			require.Equal(t, selectedAccount.ProfilePic, fetchedAccount.ProfilePic)
		})

		t.Run("Search Accounts", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var ids = []types.UUID{entries[0].(v1.Account).Id}
				query := v1.SearchAccountsJSONRequestBody{
					Id: &ids,
				}
				result, err := testClient.SearchAccountsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Account
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
				// TODO: could check that fields are equal too
			})

			t.Run("By creation time", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var timeStart = time.Now().AddDate(0, 0, -10)
				var timeEnd = time.Now().AddDate(0, 0, -5)
				query := v1.SearchAccountsJSONRequestBody{
					From:  &timeStart,
					Until: &timeEnd,
				}
				result, err := testClient.SearchAccountsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				// Ensure the subset is correct
				var queryResult []v1.Account
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				// We should have 3 results, each with proper creation timestamp
				require.Equal(t, 3, len(queryResult))
				for i := 0; i < len(queryResult); i++ {
					require.LessOrEqual(t, *(queryResult[i].CreatedAt), timeEnd)
					require.GreaterOrEqual(t, *(queryResult[i].CreatedAt), timeStart)
				}
			})

			t.Run("By username", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				// Query for subset based on username
				var username = "Doe"
				query := v1.SearchAccountsJSONRequestBody{
					Username: &username,
				}
				result, err := testClient.SearchAccountsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Account
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 2, len(queryResult))
				for i := 0; i < len(queryResult); i++ {
					// Check that the account username string contains the thing we searched for
					require.True(t, strings.Contains(queryResult[i].Username, username))
				}
			})
		})
	})

	t.Run("Group", func(t *testing.T) {
		t.Run("Create Group", func(t *testing.T) {
			_, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			body := v1.PutGroupJSONRequestBody{
				Id:          uuid.New(),
				Name:        "New Group",
				Description: "",
				Members:     []types.UUID{},
			}

			response, err := testClient.PutGroupWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var createdGroup v1.Group
			err = json.Unmarshal(response.Body, &createdGroup)
			require.NoError(t, err)
			require.Equal(t, body.Id, createdGroup.Id)
			require.NotNil(t, createdGroup.CreatedAt)
			require.Equal(t, body.Description, createdGroup.Description)
			require.Equal(t, body.Members, createdGroup.Members)
		})

		t.Run("Update Group By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedGroup := entries[8].(v1.Group)

			// Test updating the account username
			newName := "Updated Group Name"
			body := v1.UpdateGroupByIDJSONRequestBody{
				Name: &(newName),
			}
			response, err := testClient.UpdateGroupByIDWithResponse(context.Background(), selectedGroup.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var updatedGroup v1.Group
			err = json.Unmarshal(response.Body, &updatedGroup)
			require.NoError(t, err)
			expected := v1.Group{
				Id:          selectedGroup.Id,
				CreatedAt:   selectedGroup.CreatedAt,
				Name:        newName,
				Description: selectedGroup.Description,
				Members:     selectedGroup.Members,
			}
			// Must compare all fields individually because timestamps are tricky and suck
			require.Equal(t, expected.Id, updatedGroup.Id)
			require.True(t, expected.CreatedAt.Equal(*updatedGroup.CreatedAt))
			require.Equal(t, expected.Name, updatedGroup.Name)
			require.Equal(t, expected.Description, updatedGroup.Description)
			require.Equal(t, expected.Members, updatedGroup.Members)
		})

		t.Run("Delete Group By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedGroup := entries[5].(v1.Group)

			// Delete the account
			response, err := testClient.DeleteGroupByIDWithResponse(context.Background(), selectedGroup.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			// Second time should give internal server error
			response, err = testClient.DeleteGroupByIDWithResponse(context.Background(), selectedGroup.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())

			// TODO: add more advanced tests for variations when we have to also delete channels, etc...
		})

		t.Run("Get Group By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedGroup := entries[7].(v1.Group)

			// Get the account
			response, err := testClient.GetGroupByIDWithResponse(context.Background(), selectedGroup.Id)
			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode())

			var fetchedGroup v1.Group
			err = json.Unmarshal(response.Body, &fetchedGroup)
			require.NoError(t, err)
			require.Equal(t, selectedGroup.Id, fetchedGroup.Id)
			require.True(t, selectedGroup.CreatedAt.Equal(*fetchedGroup.CreatedAt))
			require.Equal(t, selectedGroup.Name, fetchedGroup.Name)
			require.Equal(t, selectedGroup.Description, fetchedGroup.Description)
			require.Equal(t, selectedGroup.Members, fetchedGroup.Members)
		})

		t.Run("Search Groups", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var ids = []types.UUID{entries[6].(v1.Group).Id, entries[9].(v1.Group).Id}
				query := v1.SearchGroupsJSONRequestBody{
					Id: &ids,
				}
				result, err := testClient.SearchGroupsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Group
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 2, len(queryResult))
				// TODO: could check that fields are equal too
			})

			t.Run("By creation time", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var timeStart = time.Now().AddDate(0, 0, -10)
				var timeEnd = time.Now().AddDate(0, 0, -5)
				query := v1.SearchGroupsJSONRequestBody{
					From:  &timeStart,
					Until: &timeEnd,
				}
				result, err := testClient.SearchGroupsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Group
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 2, len(queryResult))
			})

			t.Run("By name", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				searchName := "Advanced"
				query := v1.SearchGroupsJSONRequestBody{
					Name: &searchName,
				}
				result, err := testClient.SearchGroupsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Group
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 2, len(queryResult))
			})

			t.Run("By members", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var ids = []types.UUID{entries[2].(v1.Account).Id}
				query := v1.SearchGroupsJSONRequestBody{
					Members: &ids,
				}
				result, err := testClient.SearchGroupsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Group
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 0, len(queryResult))
				// TODO: could check that fields are equal too
				// TODO: actually have an account be a member of a group
			})
		})

		t.Run("Add Member", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			result, err := testClient.AddGroupMemberWithResponse(context.Background(), entries[7].(v1.Group).Id, entries[2].(v1.Account).Id)
			require.NoError(t, err)
			require.Equal(t, 201, result.StatusCode())
		})

		t.Run("Remove Member", func(t *testing.T) {

		})
	})

	t.Run("Channel", func(t *testing.T) {
		t.Run("Create Channel", func(t *testing.T) {

		})

		t.Run("Update Channel By Id", func(t *testing.T) {

		})

		t.Run("Delete Channel By Id", func(t *testing.T) {

		})

		t.Run("Get Channel By Id", func(t *testing.T) {

		})

		t.Run("Search Channels", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {

			})

			t.Run("By creation time", func(t *testing.T) {

			})

			t.Run("By name", func(t *testing.T) {

			})

			t.Run("By group", func(t *testing.T) {

			})
		})
	})

	t.Run("Message", func(t *testing.T) {
		t.Run("Create Message", func(t *testing.T) {

		})

		t.Run("Update Message By Id", func(t *testing.T) {

		})

		t.Run("Delete Message By Id", func(t *testing.T) {

		})

		t.Run("Get Message By Id", func(t *testing.T) {

		})

		t.Run("Search Message", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {

			})

			t.Run("By creation time", func(t *testing.T) {

			})

			t.Run("By author", func(t *testing.T) {

			})

			t.Run("By channel", func(t *testing.T) {

			})

			t.Run("By pinned", func(t *testing.T) {

			})

			t.Run("By body", func(t *testing.T) {

			})
		})
	})

	/*

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
