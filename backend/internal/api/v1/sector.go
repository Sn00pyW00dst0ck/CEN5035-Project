package v1

import (
	"app/internal/database"
	"app/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"testing"
	"time"

	"berty.tech/go-orbit-db/iface"
	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type SectorAPI struct {
	Logger *zap.Logger
	DB     *database.Database
}

// SearchAccounts implements ServerInterface.
func (s *SectorAPI) SearchAccounts(w http.ResponseWriter, r *http.Request) {
	var filter AccountFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return
	}

	accounts, err := s.DB.Store.Query(context.Background(), func(doc interface{}) (bool, error) {
		data, ok := doc.(map[string]interface{})
		if !ok {
			return false, nil
		}

		var account Account
		err := MapToStruct(data, &account)
		if err != nil {
			fmt.Println(err)
			return false, nil
		}

		// Apply filter conditions
		if filter.Ids != nil && !slices.Contains(*filter.Ids, account.Id) {
			return false, nil
		}
		if filter.From != nil && account.CreatedAt.Before(*filter.From) {
			return false, nil
		}
		if filter.Until != nil && account.CreatedAt.After(*filter.Until) {
			return false, nil
		}
		if filter.Username != nil && !fuzzy.Match(*filter.Username, account.Username) {
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not perform database query.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}

// PutAccount implements ServerInterface.
func (s *SectorAPI) PutAccount(w http.ResponseWriter, r *http.Request) {
	var accountDetails Account
	if err := json.NewDecoder(r.Body).Decode(&accountDetails); err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return
	}

	if accountDetails.CreatedAt == nil {
		var now = time.Now()
		accountDetails.CreatedAt = &now
	}

	operation, err := s.DB.Store.Put(context.Background(), StructToMap(accountDetails))
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not update within database.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operation.GetValue())
}

// DeleteAccountByID implements ServerInterface.
func (s *SectorAPI) DeleteAccountByID(w http.ResponseWriter, r *http.Request, id types.UUID) {
	_, err := s.DB.Store.Delete(context.Background(), id.String())
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not delete within database.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetAccountByID implements ServerInterface.
func (s *SectorAPI) GetAccountByID(w http.ResponseWriter, r *http.Request, id types.UUID) {
	account, err := s.DB.Store.Get(context.Background(), id.String(), &iface.DocumentStoreGetOptions{})
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not complete database query.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account[0])
}

// SearchGroups implements ServerInterface.
func (s *SectorAPI) SearchGroups(w http.ResponseWriter, r *http.Request) {
	var filter GroupFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return
	}

	groups, err := s.DB.Store.Query(context.Background(), func(doc interface{}) (bool, error) {
		data, ok := doc.(map[string]interface{})
		if !ok {
			return false, nil
		}

		var group Group
		err := MapToStruct(data, &group)
		if err != nil {
			fmt.Println(err)
			return false, nil
		}

		// Apply filter conditions
		if filter.Id != nil && !slices.Contains(*filter.Id, group.Id) {
			return false, nil
		}
		if filter.From != nil && group.CreatedAt.Before(*filter.From) {
			return false, nil
		}
		if filter.Until != nil && group.CreatedAt.After(*filter.Until) {
			return false, nil
		}
		if filter.Name != nil && !fuzzy.Match(*filter.Name, group.Name) {
			return false, nil
		}

		return true, nil
	})
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not perform database query.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(groups)
}

// PutGroup implements ServerInterface.
func (s *SectorAPI) PutGroup(w http.ResponseWriter, r *http.Request) {
	var groupDetails Group
	if err := json.NewDecoder(r.Body).Decode(&groupDetails); err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return
	}

	if groupDetails.CreatedAt == nil {
		var now = time.Now()
		groupDetails.CreatedAt = &now
	}

	operation, err := s.DB.Store.Put(context.Background(), StructToMap(groupDetails))
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not update within database.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(operation.GetValue())
}

// DeleteGroupByID implements ServerInterface.
func (s *SectorAPI) DeleteGroupByID(w http.ResponseWriter, r *http.Request, id types.UUID) {
	_, err := s.DB.Store.Delete(context.Background(), id.String())
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not delete within database.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetGroupByID implements ServerInterface.
func (s *SectorAPI) GetGroupByID(w http.ResponseWriter, r *http.Request, id types.UUID) {
	group, err := s.DB.Store.Get(context.Background(), id.String(), &iface.DocumentStoreGetOptions{})
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not complete database query.", http.StatusInternalServerError)
		return
	}
	// TODO: might need to do other things to get account return type to work with the openapi validation.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(group[0])
}

// GetHealth implements ServerInterface.
func (s *SectorAPI) GetHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// GetRoot implements ServerInterface.
func (s *SectorAPI) GetRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Hello, World!"}
	json.NewEncoder(w).Encode(response)
}

// Make sure we conform to ServerInterface

var _ ServerInterface = (*SectorAPI)(nil)

// Create a new Sector API instance
func NewSector(ctx context.Context, logfile, dbCache string) *SectorAPI {
	// Setup the logger
	logger, err := logger.NewLogger(logfile)
	if err != nil {
		panic(err)
	}

	// Setup the database
	db, err := database.NewDatabase(ctx, dbCache, logger)
	if err != nil {
		panic(err)
	}
	// defer db.Disconnect() (TODO: FIGURE OUT WHEN TO CALL DISCONNECT)

	err = db.Connect(func(address string) {
		fmt.Println("Connected: ", address)
	})
	if err != nil {
		panic(err)
	}

	return &SectorAPI{
		Logger: logger,
		DB:     db,
	}
}

// Create a new SectorAPI instance for unit testing
func NewTestingSector(ctx context.Context, logfile, dbCache string, t *testing.T) *SectorAPI {
	// Setup the logger
	logger, err := logger.NewLogger(logfile)
	if err != nil {
		panic(err)
	}

	// Setup the database
	db, err := database.NewTestingDatabase(ctx, dbCache, logger, t)
	if err != nil {
		panic(err)
	}

	err = db.Connect(func(address string) {
		fmt.Println("Connected: ", address)
	})
	if err != nil {
		panic(err)
	}

	return &SectorAPI{
		Logger: logger,
		DB:     db,
	}
}

// Helper functions
func StructToMap(obj interface{}) map[string]interface{} {
	// Marshal struct to JSON
	data, _ := json.Marshal(obj)

	// Unmarshal JSON into a map
	var result map[string]interface{}
	json.Unmarshal(data, &result)

	return result
}

// Whenever we want to convert something in the database to a struct use this.
func MapToStruct(data map[string]interface{}, obj interface{}) error {
	// Marshal the map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Unmarshal JSON into the provided struct
	return json.Unmarshal(jsonData, obj)
}
