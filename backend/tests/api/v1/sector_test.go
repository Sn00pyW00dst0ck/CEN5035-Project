package v1Test

import (
	"Sector/internal/api"
	v1 "Sector/internal/api/v1"
	"Sector/internal/database"
	"context"
	"encoding/base64"
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

	// Create string pointers
	generalName := "General"
	generalDesc := "General discussions"
	randomName := "Random"
	randomDesc := "Random discussions"
	techName := "Tech"
	techDesc := "Tech-related discussions"
	offTopicName := "Off-Topic"
	offTopicDesc := "Off-topic discussions"
	announcementsName := "Announcements"
	announcementsDesc := "Important announcements"

	// Create UUIDs for relationships
	account1ID := uuid.New()
	account2ID := uuid.New()
	account3ID := uuid.New()
	account4ID := uuid.New()
	account5ID := uuid.New()
	group1ID := uuid.New()
	group2ID := uuid.New()
	group3ID := uuid.New()
	advGroup1ID := uuid.New()
	advGroup2ID := uuid.New()
	channel1ID := uuid.New()
	channel2ID := uuid.New()
	channel3ID := uuid.New()
	channel4ID := uuid.New()
	channel5ID := uuid.New()

	entries := []interface{}{
		v1.Account{
			Id:         account1ID,
			CreatedAt:  &now,
			ProfilePic: nil,
			Username:   "John Doe",
		},
		v1.Account{
			Id:         account2ID,
			CreatedAt:  &then,
			ProfilePic: nil,
			Username:   "Jack Doe",
		},
		v1.Account{
			Id:         account3ID,
			CreatedAt:  &now,
			ProfilePic: nil,
			Username:   "Maverick",
		},
		v1.Account{
			Id:         account4ID,
			CreatedAt:  &then,
			ProfilePic: nil,
			Username:   "w311un1!k3",
		},
		v1.Account{
			Id:         account5ID,
			CreatedAt:  &then,
			ProfilePic: nil,
			Username:   "woefullyconsideringlove",
		},
		v1.Group{
			Id:          group1ID,
			CreatedAt:   &then,
			Name:        "Test Group 1",
			Description: "A group for unit testing.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          group2ID,
			CreatedAt:   &now,
			Name:        "Test Group 2",
			Description: "Another unit testing group.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          group3ID,
			CreatedAt:   &now,
			Name:        "Test Group 3",
			Description: "A third group for unit testing.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          advGroup1ID,
			CreatedAt:   &then,
			Name:        "Advanced Test Group 1",
			Description: "For advanced testing.",
			Members:     []types.UUID{},
		},
		v1.Group{
			Id:          advGroup2ID,
			CreatedAt:   &now,
			Name:        "Advanced Test Group 2",
			Description: "For advanced testing.",
			Members:     []types.UUID{},
		},
		v1.Channel{
			Id:          channel1ID,
			CreatedAt:   &now,
			Name:        &generalName,
			Description: &generalDesc,
			Group:       group1ID,
		},
		v1.Channel{
			Id:          channel2ID,
			CreatedAt:   &now,
			Name:        &randomName,
			Description: &randomDesc,
			Group:       group2ID,
		},
		v1.Channel{
			Id:          channel3ID,
			CreatedAt:   &now,
			Name:        &techName,
			Description: &techDesc,
			Group:       group3ID,
		},
		v1.Channel{
			Id:          channel4ID,
			CreatedAt:   &now,
			Name:        &offTopicName,
			Description: &offTopicDesc,
			Group:       advGroup1ID,
		},
		v1.Channel{
			Id:          channel5ID,
			CreatedAt:   &now,
			Name:        &announcementsName,
			Description: &announcementsDesc,
			Group:       advGroup2ID,
		},
		v1.Message{
			Id:        uuid.New(),
			Author:    account1ID,
			Body:      "Hello everyone!",
			Channel:   channel1ID,
			CreatedAt: &now,
			Pinned:    false,
		},
		v1.Message{
			Id:        uuid.New(),
			Author:    account2ID,
			Body:      "Hey John, how's it going?",
			Channel:   channel1ID,
			CreatedAt: &now,
			Pinned:    false,
		},
		v1.Message{
			Id:        uuid.New(),
			Author:    account3ID,
			Body:      "Anyone up for a game?",
			Channel:   channel2ID,
			CreatedAt: &now,
			Pinned:    false,
		},
		v1.Message{
			Id:        uuid.New(),
			Author:    account4ID,
			Body:      "Sure, what are we playing?",
			Channel:   channel2ID,
			CreatedAt: &now,
			Pinned:    false,
		},
		v1.Message{
			Id:        uuid.New(),
			Author:    account1ID,
			Body:      "Important update!",
			Channel:   channel5ID,
			CreatedAt: &now,
			Pinned:    true,
		},
	}

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
	// Setup the server to run for all the tests
	server, sectorAPI, teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	testClient, err := v1.NewClientWithResponses(server.URL,
		v1.WithHTTPClient(server.Client()),
		v1.WithBaseURL(server.URL+"/v1/api"))
	require.NoError(t, err)

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
			account := entries[0].(v1.Account)

			newUsername := "Updated Test Username"
			body := v1.UpdateAccountByIDJSONRequestBody{Username: &newUsername}
			response, err := testClient.UpdateAccountByIDWithResponse(context.Background(), account.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var updated v1.Account
			json.Unmarshal(response.Body, &updated)
			require.Equal(t, newUsername, updated.Username)
		})

		t.Run("Delete Account By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			account := entries[0].(v1.Account)

			response, err := testClient.DeleteAccountByIDWithResponse(context.Background(), account.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			response, err = testClient.DeleteAccountByIDWithResponse(context.Background(), account.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			account := entries[0].(v1.Account)

			response, err := testClient.GetAccountByIDWithResponse(context.Background(), account.Id)
			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode())

			var fetched v1.Account
			json.Unmarshal(response.Body, &fetched)
			require.Equal(t, account.Id, fetched.Id)
		})

		t.Run("Search Accounts", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)
				account := entries[0].(v1.Account)

				query := v1.SearchAccountsJSONRequestBody{Id: &[]types.UUID{account.Id}}
				result, err := testClient.SearchAccountsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var results []v1.Account
				json.Unmarshal(result.Body, &results)
				require.Len(t, results, 1)
				require.Equal(t, account.Id, results[0].Id)
			})

			t.Run("By creation time", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				start := time.Now().Add(-24 * time.Hour)
				end := time.Now()
				query := v1.SearchAccountsJSONRequestBody{From: &start, Until: &end}
				result, err := testClient.SearchAccountsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var results []v1.Account
				json.Unmarshal(result.Body, &results)
				require.Greater(t, len(results), 0)
			})

			t.Run("By username", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				search := "Doe"
				query := v1.SearchAccountsJSONRequestBody{Username: &search}
				result, err := testClient.SearchAccountsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var results []v1.Account
				json.Unmarshal(result.Body, &results)
				require.Len(t, results, 2)
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
				Description: "Test group",
				Members:     []types.UUID{},
			}

			response, err := testClient.PutGroupWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var created v1.Group
			json.Unmarshal(response.Body, &created)
			require.Equal(t, body.Id, created.Id)
		})

		t.Run("Update Group By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			group := entries[5].(v1.Group)

			newName := "Updated Group"
			body := v1.UpdateGroupByIDJSONRequestBody{Name: &newName}
			response, err := testClient.UpdateGroupByIDWithResponse(context.Background(), group.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, response.StatusCode())

			var updated v1.Group
			json.Unmarshal(response.Body, &updated)
			require.Equal(t, newName, updated.Name)
		})

		t.Run("Delete Group By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			group := entries[5].(v1.Group)

			response, err := testClient.DeleteGroupByIDWithResponse(context.Background(), group.Id)
			require.NoError(t, err)
			require.Equal(t, 204, response.StatusCode())

			response, err = testClient.DeleteGroupByIDWithResponse(context.Background(), group.Id)
			require.NoError(t, err)
			require.Equal(t, 500, response.StatusCode())
		})

		t.Run("Get Group By Id", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			group := entries[5].(v1.Group)

			response, err := testClient.GetGroupByIDWithResponse(context.Background(), group.Id)
			require.NoError(t, err)
			require.Equal(t, 200, response.StatusCode())

			var fetched v1.Group
			json.Unmarshal(response.Body, &fetched)
			require.Equal(t, group.Id, fetched.Id)
		})

		t.Run("Search Groups", func(t *testing.T) {
			t.Run("By Id", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)
				group1 := entries[5].(v1.Group)
				group2 := entries[6].(v1.Group)

				query := v1.SearchGroupsJSONRequestBody{Id: &[]types.UUID{group1.Id, group2.Id}}
				result, err := testClient.SearchGroupsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var results []v1.Group
				json.Unmarshal(result.Body, &results)
				require.Len(t, results, 2)
			})

			t.Run("By Name", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				search := "Advanced"
				query := v1.SearchGroupsJSONRequestBody{Name: &search}
				result, err := testClient.SearchGroupsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, result.StatusCode())

				var results []v1.Group
				json.Unmarshal(result.Body, &results)
				require.Greater(t, len(results), 0)
			})
		})

		t.Run("Group Membership", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			group := entries[5].(v1.Group)
			account := entries[0].(v1.Account)

			t.Run("Add Member", func(t *testing.T) {
				resp, err := testClient.AddGroupMemberWithResponse(context.Background(), group.Id, account.Id)
				require.NoError(t, err)
				require.Equal(t, 201, resp.StatusCode())

				// Verify membership
				getResp, err := testClient.GetGroupByIDWithResponse(context.Background(), group.Id)
				require.NoError(t, err)
				var updatedGroup v1.Group
				json.Unmarshal(getResp.Body, &updatedGroup)
				require.Contains(t, updatedGroup.Members, account.Id)
			})

			t.Run("Remove Member", func(t *testing.T) {
				// First add the member
				_, err := testClient.AddGroupMemberWithResponse(context.Background(), group.Id, account.Id)
				require.NoError(t, err)

				// Then remove
				resp, err := testClient.RemoveGroupMemberWithResponse(context.Background(), group.Id, account.Id)
				require.NoError(t, err)
				require.Equal(t, 204, resp.StatusCode())

				// Verify removal
				getResp, err := testClient.GetGroupByIDWithResponse(context.Background(), group.Id)
				require.NoError(t, err)
				var updatedGroup v1.Group
				json.Unmarshal(getResp.Body, &updatedGroup)
				require.NotContains(t, updatedGroup.Members, account.Id)
			})
		})
	})

	t.Run("Channel", func(t *testing.T) {
		t.Run("Create Channel", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			group := entries[5].(v1.Group)

			body := v1.PutChannelJSONRequestBody{
				Id:          uuid.New(),
				Name:        "New Channel",
				Description: "Test channel",
				Group:       group.Id,
			}

			resp, err := testClient.PutChannelWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 201, resp.StatusCode())

			var created v1.Channel
			json.Unmarshal(resp.Body, &created)
			require.Equal(t, body.Id, created.Id)
			require.Equal(t, body.Group, created.Group)
		})

		t.Run("Update Channel", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			channel := entries[10].(v1.Channel)

			newName := "Updated Channel"
			newDesc := "Updated description"
			body := v1.UpdateChannelByIDJSONRequestBody{
				Name:        &newName,
				Description: &newDesc,
			}

			resp, err := testClient.UpdateChannelByIDWithResponse(context.Background(), channel.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, resp.StatusCode())

			var updated v1.Channel
			json.Unmarshal(resp.Body, &updated)
			require.Equal(t, newName, *updated.Name)
			require.Equal(t, newDesc, *updated.Description)
		})

		t.Run("Delete Channel", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			channel := entries[10].(v1.Channel)

			resp, err := testClient.DeleteChannelByIDWithResponse(context.Background(), channel.Id)
			require.NoError(t, err)
			require.Equal(t, 204, resp.StatusCode())

			// Verify deletion
			getResp, err := testClient.GetChannelByIDWithResponse(context.Background(), channel.Id)
			require.Equal(t, 404, getResp.StatusCode())
		})

		t.Run("Get Channel", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			channel := entries[10].(v1.Channel)

			resp, err := testClient.GetChannelByIDWithResponse(context.Background(), channel.Id)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode())

			var fetched v1.Channel
			json.Unmarshal(resp.Body, &fetched)
			require.Equal(t, channel.Id, fetched.Id)
		})

		t.Run("Search Channels", func(t *testing.T) {
			t.Run("By Group", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)
				group := entries[5].(v1.Group)

				query := v1.SearchChannelsJSONRequestBody{Group: &group.Id}
				resp, err := testClient.SearchChannelsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, resp.StatusCode())

				var results []v1.Channel
				json.Unmarshal(resp.Body, &results)
				require.Greater(t, len(results), 0)
				for _, ch := range results {
					require.Equal(t, group.Id, ch.Group)
				}
			})

			t.Run("By Name", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)
				search := "General"

				query := v1.SearchChannelsJSONRequestBody{Name: &search}
				resp, err := testClient.SearchChannelsWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, resp.StatusCode())

				var results []v1.Channel
				json.Unmarshal(resp.Body, &results)
				require.Greater(t, len(results), 0)
			})
		})
	})

	t.Run("Message", func(t *testing.T) {
		t.Run("Create Message", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			channel := entries[10].(v1.Channel)
			author := entries[0].(v1.Account)

			body := v1.PutMessageJSONRequestBody{
				Id:      uuid.New(),
				Author:  author.Id,
				Body:    "Test message",
				Channel: channel.Id,
			}

			resp, err := testClient.PutMessageWithResponse(context.Background(), body)
			require.NoError(t, err)
			require.Equal(t, 201, resp.StatusCode())

			var created v1.Message
			json.Unmarshal(resp.Body, &created)
			require.Equal(t, body.Id, created.Id)
		})

		t.Run("Update Message", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			message := entries[15].(v1.Message)

			newBody := "Updated message content"
			pinned := true
			body := v1.UpdateMessageByIDJSONRequestBody{
				Body:   &newBody,
				Pinned: &pinned,
			}

			resp, err := testClient.UpdateMessageByIDWithResponse(context.Background(), message.Id, body)
			require.NoError(t, err)
			require.Equal(t, 201, resp.StatusCode())

			var updated v1.Message
			json.Unmarshal(resp.Body, &updated)
			require.Equal(t, newBody, updated.Body)
			require.Equal(t, pinned, updated.Pinned)
		})

		t.Run("Delete Message", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			message := entries[15].(v1.Message)

			resp, err := testClient.DeleteMessageByIDWithResponse(context.Background(), message.Id)
			require.NoError(t, err)
			require.Equal(t, 204, resp.StatusCode())

			// Verify deletion
			getResp, err := testClient.GetMessageByIDWithResponse(context.Background(), message.Id)
			require.Equal(t, 404, getResp.StatusCode())
		})

		t.Run("Get Message", func(t *testing.T) {
			entries, teardown := setupTest(t, *sectorAPI)
			defer teardown(t)
			message := entries[15].(v1.Message)

			resp, err := testClient.GetMessageByIDWithResponse(context.Background(), message.Id)
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode())

			var fetched v1.Message
			json.Unmarshal(resp.Body, &fetched)
			require.Equal(t, message.Id, fetched.Id)
		})

		t.Run("Search Messages", func(t *testing.T) {
			t.Run("By Author", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)
				author := entries[0].(v1.Account)

				query := v1.SearchMessagesJSONRequestBody{Author: &author.Id}
				resp, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, resp.StatusCode())

				var results []v1.Message
				json.Unmarshal(resp.Body, &results)
				require.Greater(t, len(results), 0)
			})

			t.Run("By Channel", func(t *testing.T) {
				entries, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)
				channel := entries[10].(v1.Channel)

				query := v1.SearchMessagesJSONRequestBody{Channel: &channel.Id}
				resp, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, resp.StatusCode())

				var results []v1.Message
				json.Unmarshal(resp.Body, &results)
				require.Greater(t, len(results), 0)
			})

			t.Run("By Pinned", func(t *testing.T) {
				_, teardown := setupTest(t, *sectorAPI)
				defer teardown(t)

				pinned := true
				query := v1.SearchMessagesJSONRequestBody{Pinned: &pinned}
				resp, err := testClient.SearchMessagesWithResponse(context.Background(), query)
				require.NoError(t, err)
				require.Equal(t, 200, resp.StatusCode())

				var results []v1.Message
				json.Unmarshal(resp.Body, &results)
				require.Greater(t, len(results), 0)
			})
		})
	})
}

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
