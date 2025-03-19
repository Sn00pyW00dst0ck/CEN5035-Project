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

func stringPtr(s string) *string {
	return &s
}

func setupTest(t *testing.T, api v1.SectorAPI) ([]interface{}, func(t *testing.T)) {
	now := time.Now()
	then := now.AddDate(0, 0, -7)

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
	}

	group1ID := entries[5].(v1.Group).Id
	group2ID := entries[6].(v1.Group).Id
	group3ID := entries[7].(v1.Group).Id
	group4ID := entries[8].(v1.Group).Id
	group5ID := entries[9].(v1.Group).Id

	account1ID := entries[0].(v1.Account).Id
	account2ID := entries[1].(v1.Account).Id
	account3ID := entries[2].(v1.Account).Id

	entries = append(entries, v1.Channel{
		CreatedAt:   &now,
		Description: stringPtr("Main discussion hub"),
		Group:       group1ID,
		Id:          uuid.New(),
		Name:        "Main",
	})
	entries = append(entries, v1.Channel{
		CreatedAt:   nil,
		Description: stringPtr("Group updates"),
		Group:       group2ID,
		Id:          uuid.New(),
		Name:        "Updates",
	})
	entries = append(entries, v1.Channel{
		CreatedAt:   &then,
		Description: nil,
		Group:       group3ID,
		Id:          uuid.New(),
		Name:        "Chat",
	})
	entries = append(entries, v1.Channel{
		CreatedAt:   &now,
		Description: stringPtr("Tech discussions"),
		Group:       group4ID,
		Id:          uuid.New(),
		Name:        "Tech",
	})
	entries = append(entries, v1.Channel{
		CreatedAt:   &then,
		Description: stringPtr("Strategy planning"),
		Group:       group5ID,
		Id:          uuid.New(),
		Name:        "Strategy",
	})

	mainChannelID := entries[len(entries)-5].(v1.Channel).Id
	chatChannelID := entries[len(entries)-3].(v1.Channel).Id
	techChannelID := entries[len(entries)-2].(v1.Channel).Id

	entries = append(entries, v1.Message{
		Author:    account1ID,
		Body:      "Welcome to the Main channel!",
		Channel:   mainChannelID,
		CreatedAt: &now,
		Id:        uuid.New(),
		Pinned:    true,
	})
	entries = append(entries, v1.Message{
		Author:    account2ID,
		Body:      "Anyone around?",
		Channel:   mainChannelID,
		CreatedAt: nil,
		Id:        uuid.New(),
		Pinned:    false,
	})
	entries = append(entries, v1.Message{
		Author:    account3ID,
		Body:      "What’s up in Chat?",
		Channel:   chatChannelID,
		CreatedAt: &then,
		Id:        uuid.New(),
		Pinned:    false,
	})
	entries = append(entries, v1.Message{
		Author:    account1ID,
		Body:      "New tech ideas here.",
		Channel:   techChannelID,
		CreatedAt: &now,
		Id:        uuid.New(),
		Pinned:    false,
	})

	result := make([]interface{}, len(entries))
	for i, v := range entries {
		result[i] = v1.StructToMap(v)
	}
	_, err := api.DB.Store.PutAll(context.Background(), result)
	require.NoError(t, err)

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
	server, sectorAPI, teardownSuite := setupSuite(t)
	defer teardownSuite(t)

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

			body.Username = "Updated Username!"
			response, err = testClient.PutAccountWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Update Account By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedAccount := entries[0].(v1.Account)

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
			require.Equal(t, expected.Id, updatedAccount.Id)
			require.True(t, expected.CreatedAt.Equal(*updatedAccount.CreatedAt))
			require.Equal(t, expected.Username, updatedAccount.Username)
			require.Equal(t, expected.ProfilePic, updatedAccount.ProfilePic)
		})

		t.Run("Delete Account By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedAccount := entries[0].(v1.Account)

			response, err := testClient.DeleteAccountByIDWithResponse(context.Background(), selectedAccount.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			response, err = testClient.DeleteAccountByIDWithResponse(context.Background(), selectedAccount.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedAccount := entries[0].(v1.Account)

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

				var queryResult []v1.Account
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 3, len(queryResult))
				for i := 0; i < len(queryResult); i++ {
					require.LessOrEqual(t, *(queryResult[i].CreatedAt), timeEnd)
					require.GreaterOrEqual(t, *(queryResult[i].CreatedAt), timeStart)
				}
			})

			t.Run("By username", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

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

			response, err := testClient.DeleteGroupByIDWithResponse(context.Background(), selectedGroup.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			response, err = testClient.DeleteGroupByIDWithResponse(context.Background(), selectedGroup.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get Group By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedGroup := entries[7].(v1.Group)

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
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			groupID := entries[7].(v1.Group).Id
			accountID := entries[2].(v1.Account).Id

			_, err := testClient.AddGroupMemberWithResponse(context.Background(), groupID, accountID)
			require.NoError(t, err)

			result, err := testClient.RemoveGroupMemberWithResponse(context.Background(), groupID, accountID)
			require.NoError(t, err)
			require.Equal(t, 204, result.StatusCode())

			fetchedGroupResp, err := testClient.GetGroupByIDWithResponse(context.Background(), groupID)
			require.NoError(t, err)
			require.Equal(t, 200, fetchedGroupResp.StatusCode())

			var fetchedGroup v1.Group
			err = json.Unmarshal(fetchedGroupResp.Body, &fetchedGroup)
			require.NoError(t, err)
			require.Empty(t, fetchedGroup.Members)
		})
	})

	t.Run("Channel", func(t *testing.T) {
		t.Run("Create Channel", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			body := v1.PutChannelJSONRequestBody{
				Id:          uuid.New(),
				Name:        "New Channel",
				Description: stringPtr("A new test channel"),
				Group:       entries[5].(v1.Group).Id,
			}

			response, err := testClient.PutChannelWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var createdChannel v1.Channel
			err = json.Unmarshal(response.Body, &createdChannel)
			require.NoError(t, err)
			require.Equal(t, body.Id, createdChannel.Id)
			require.NotNil(t, createdChannel.CreatedAt)
			require.Equal(t, body.Name, createdChannel.Name)
			require.Equal(t, *body.Description, *createdChannel.Description)
			require.Equal(t, body.Group, createdChannel.Group)
		})

		t.Run("Update Channel By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedChannel := entries[10].(v1.Channel)

			newName := "Updated Channel Name"
			body := v1.UpdateChannelByIDJSONRequestBody{
				Name: &newName,
			}
			response, err := testClient.UpdateChannelByIDWithResponse(context.Background(), selectedChannel.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var updatedChannel v1.Channel
			err = json.Unmarshal(response.Body, &updatedChannel)
			require.NoError(t, err)
			expected := v1.Channel{
				Id:          selectedChannel.Id,
				CreatedAt:   selectedChannel.CreatedAt,
				Name:        newName,
				Description: selectedChannel.Description,
				Group:       selectedChannel.Group,
			}
			require.Equal(t, expected.Id, updatedChannel.Id)
			require.True(t, expected.CreatedAt.Equal(*updatedChannel.CreatedAt))
			require.Equal(t, expected.Name, updatedChannel.Name)
			require.Equal(t, expected.Description, updatedChannel.Description)
			require.Equal(t, expected.Group, updatedChannel.Group)
		})

		t.Run("Delete Channel By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedChannel := entries[11].(v1.Channel)

			response, err := testClient.DeleteChannelByIDWithResponse(context.Background(), selectedChannel.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			response, err = testClient.DeleteChannelByIDWithResponse(context.Background(), selectedChannel.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get Channel By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedChannel := entries[12].(v1.Channel)

			response, err := testClient.GetChannelByIDWithResponse(context.Background(), selectedChannel.Id)
			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode())

			var fetchedChannel v1.Channel
			err = json.Unmarshal(response.Body, &fetchedChannel)
			require.NoError(t, err)
			require.Equal(t, selectedChannel.Id, fetchedChannel.Id)
			if selectedChannel.CreatedAt != nil {
				require.True(t, selectedChannel.CreatedAt.Equal(*fetchedChannel.CreatedAt))
			}
			require.Equal(t, selectedChannel.Name, fetchedChannel.Name)
			require.Equal(t, selectedChannel.Description, fetchedChannel.Description)
			require.Equal(t, selectedChannel.Group, fetchedChannel.Group)
		})

		t.Run("Search Channels", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var ids = []types.UUID{entries[10].(v1.Channel).Id}
				query := v1.SearchChannelsJSONRequestBody{
					Id: &ids,
				}
				result, err := testClient.SearchChannelsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Channel
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
			})

			t.Run("By creation time", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var timeStart = time.Now().AddDate(0, 0, -10)
				var timeEnd = time.Now().AddDate(0, 0, -5)
				query := v1.SearchChannelsJSONRequestBody{
					From:  &timeStart,
					Until: &timeEnd,
				}
				result, err := testClient.SearchChannelsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Channel
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 2, len(queryResult))
			})

			t.Run("By name", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				searchName := "Main"
				query := v1.SearchChannelsJSONRequestBody{
					Name: &searchName,
				}
				result, err := testClient.SearchChannelsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Channel
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
			})

			t.Run("By group", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var groupIDs = []types.UUID{entries[5].(v1.Group).Id}
				query := v1.SearchChannelsJSONRequestBody{
					Group: &groupIDs,
				}
				result, err := testClient.SearchChannelsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Channel
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
			})
		})
	})

	t.Run("Message", func(t *testing.T) {
		t.Run("Create Message", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			body := v1.PutMessageJSONRequestBody{
				Id:      uuid.New(),
				Author:  entries[0].(v1.Account).Id,
				Body:    "Test message",
				Channel: entries[10].(v1.Channel).Id,
				Pinned:  false,
			}

			response, err := testClient.PutMessageWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var createdMessage v1.Message
			err = json.Unmarshal(response.Body, &createdMessage)
			require.NoError(t, err)
			require.Equal(t, body.Id, createdMessage.Id)
			require.NotNil(t, createdMessage.CreatedAt)
			require.Equal(t, body.Author, createdMessage.Author)
			require.Equal(t, body.Body, createdMessage.Body)
			require.Equal(t, body.Channel, createdMessage.Channel)
			require.Equal(t, body.Pinned, createdMessage.Pinned)
		})

		t.Run("Update Message By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedMessage := entries[15].(v1.Message)

			newBody := "Updated message content"
			body := v1.UpdateMessageByIDJSONRequestBody{
				Body: &newBody,
			}
			response, err := testClient.UpdateMessageByIDWithResponse(context.Background(), selectedMessage.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var updatedMessage v1.Message
			err = json.Unmarshal(response.Body, &updatedMessage)
			require.NoError(t, err)
			expected := v1.Message{
				Id:        selectedMessage.Id,
				Author:    selectedMessage.Author,
				Body:      newBody,
				Channel:   selectedMessage.Channel,
				CreatedAt: selectedMessage.CreatedAt,
				Pinned:    selectedMessage.Pinned,
			}
			require.Equal(t, expected.Id, updatedMessage.Id)
			require.Equal(t, expected.Author, updatedMessage.Author)
			require.Equal(t, expected.Body, updatedMessage.Body)
			require.Equal(t, expected.Channel, updatedMessage.Channel)
			if expected.CreatedAt != nil {
				require.True(t, expected.CreatedAt.Equal(*updatedMessage.CreatedAt))
			}
			require.Equal(t, expected.Pinned, updatedMessage.Pinned)
		})

		t.Run("Delete Message By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedMessage := entries[16].(v1.Message)

			response, err := testClient.DeleteMessageByIDWithResponse(context.Background(), selectedMessage.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			response, err = testClient.DeleteMessageByIDWithResponse(context.Background(), selectedMessage.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get Message By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)

			selectedMessage := entries[17].(v1.Message)

			response, err := testClient.GetMessageByIDWithResponse(context.Background(), selectedMessage.Id)
			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode())

			var fetchedMessage v1.Message
			err = json.Unmarshal(response.Body, &fetchedMessage)
			require.NoError(t, err)
			require.Equal(t, selectedMessage.Id, fetchedMessage.Id)
			require.Equal(t, selectedMessage.Author, fetchedMessage.Author)
			require.Equal(t, selectedMessage.Body, fetchedMessage.Body)
			require.Equal(t, selectedMessage.Channel, fetchedMessage.Channel)
			if selectedMessage.CreatedAt != nil {
				require.True(t, selectedMessage.CreatedAt.Equal(*fetchedMessage.CreatedAt))
			}
			require.Equal(t, selectedMessage.Pinned, fetchedMessage.Pinned)
		})

		t.Run("Search Message", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var ids = []types.UUID{entries[15].(v1.Message).Id}
				query := v1.SearchMessagesJSONRequestBody{
					Id: &ids,
				}
				result, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Message
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
			})

			t.Run("By creation time", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var timeStart = time.Now().AddDate(0, 0, -10)
				var timeEnd = time.Now().AddDate(0, 0, -5)
				query := v1.SearchMessagesJSONRequestBody{
					From:  &timeStart,
					Until: &timeEnd,
				}
				result, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Message
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
			})

			t.Run("By author", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var authorIDs = []types.UUID{entries[0].(v1.Account).Id}
				query := v1.SearchMessagesJSONRequestBody{
					Author: &authorIDs,
				}
				result, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Message
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 2, len(queryResult))
			})

			t.Run("By channel", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				var channelIDs = []types.UUID{entries[10].(v1.Channel).Id}
				query := v1.SearchMessagesJSONRequestBody{
					Channel: &channelIDs,
				}
				result, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Message
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 2, len(queryResult))
			})

			t.Run("By pinned", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				pinned := true
				query := v1.SearchMessagesJSONRequestBody{
					Pinned: &pinned,
				}
				result, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Message
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
			})

			t.Run("By body", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				bodySearch := "Welcome"
				query := v1.SearchMessagesJSONRequestBody{
					Body: &bodySearch,
				}
				result, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var queryResult []v1.Message
				err = json.Unmarshal(result.Body, &queryResult)
				require.NoError(t, err)
				require.Equal(t, 1, len(queryResult))
			})
		})
	})
}
