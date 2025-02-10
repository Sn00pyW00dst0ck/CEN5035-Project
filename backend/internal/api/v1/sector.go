package v1

import (
	"app/internal/database"
	"app/internal/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"berty.tech/go-orbit-db/iface"
	"github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"
)

type SectorAPI struct {
	Logger *zap.Logger
	DB     *database.Database
}

// PutAccount implements ServerInterface.
func (s *SectorAPI) PutAccount(w http.ResponseWriter, r *http.Request) {
	var account_details Account
	if err := json.NewDecoder(r.Body).Decode(&account_details); err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not parse request body.", http.StatusBadRequest)
		return
	}

	// TODO: this fails every time because the 'Account' generated has 'Id' and not 'id' so we have to convert that. UGH!
	created_account, err := s.DB.Store.Put(context.Background(), StructToMap(account_details))
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not update within databse.", http.StatusInternalServerError)
		return
	}
	// TODO: might need to do other things to get account return type to work with the openapi validation.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(created_account)
}

// DeleteAccountByID implements ServerInterface.
func (s *SectorAPI) DeleteAccountByID(w http.ResponseWriter, r *http.Request, id types.UUID) {
	_, err := s.DB.Store.Delete(context.Background(), id.String())
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not delete within databse.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetAccountByID implements ServerInterface.
func (s *SectorAPI) GetAccountByID(w http.ResponseWriter, r *http.Request, id types.UUID) {
	account, err := s.DB.Store.Get(context.Background(), id.String(), &iface.DocumentStoreGetOptions{})
	if err != nil {
		s.Logger.Debug(err.Error())
		http.Error(w, "Could not complete databse query.", http.StatusInternalServerError)
		return
	}
	// TODO: might need to do other things to get account return type to work with the openapi validation.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
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
func NewSector(ctx context.Context, logfile, dbCache, dbConnectionString string) *SectorAPI {
	// Setup the logger
	logger, err := logger.NewLogger(logfile)
	if err != nil {
		panic(err)
	}

	// Setup the database
	db, err := database.NewDatabase(ctx, dbConnectionString, dbCache, logger)
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

// Whenever we want to convert a struct to a thing to put into the database use this.
func StructToMap(obj interface{}) map[string]interface{} {
	// Marshal struct to JSON
	data, _ := json.Marshal(obj)

	// Unmarshal JSON into a map
	var result map[string]interface{}
	json.Unmarshal(data, &result)

	return result
}
