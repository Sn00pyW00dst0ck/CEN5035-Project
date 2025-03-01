package v1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"slices"
	"time"

	orbitdb "berty.tech/go-orbit-db"
	"berty.tech/go-orbit-db/iface"
	"github.com/google/uuid"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/oapi-codegen/runtime/types"
)

var ErrNotFound = errors.New("item for id not found")
var ErrTooMany = errors.New("too many items for id found")

/**
 * Add a new item into the database
 */
func addItem(store orbitdb.DocumentStore, obj interface{}) (interface{}, error) {
	// Check if item with this ID already exists in the database...
	_, err := getItem(store, uuid.Must(uuid.Parse(StructToMap(obj)["id"].(string))))
	if err == nil || err != ErrNotFound {
		return nil, fmt.Errorf("cannot add item to database")
	}

	// Check the type of the given obj, if it is something special then check for its special pre-reqs before adding...
	switch item := obj.(type) {
	case Account:
		// If adding an account there is nothing special to do...
	case Group:
		// Check dependencies when adding a group object
	case Channel:
		// Check dependencies when adding a channel object
	case Message:
		// Check dependencies when adding a message object
	default:
		return nil, fmt.Errorf("cannot add unknown item '%v' type to database", item)
	}

	// Adds the item to the database
	op, err := store.Put(context.Background(), StructToMap(obj))
	if err != nil {
		return nil, err
	}

	// Convert to 'readable' interface format instead of string
	var result map[string]interface{}
	err = json.Unmarshal([]byte(op.GetValue()), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
 * Get an item from the database
 */
func getItem(store orbitdb.DocumentStore, id types.UUID) (interface{}, error) {
	matches, err := store.Get(context.Background(), id.String(), &iface.DocumentStoreGetOptions{})
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, ErrNotFound
	}
	if len(matches) > 1 {
		return nil, ErrTooMany
	}
	return matches[0], nil
}

/**
 * Update an item in the database
 */
func updateItem(store orbitdb.DocumentStore, id types.UUID, obj interface{}) (interface{}, error) {
	// Get the item to update from the DB
	dbItem, err := getItem(store, id)
	if err != nil {
		return nil, err
	}

	// Update values
	updatedItem := dbItem.(map[string]interface{})
	updatesToApply := StructToMap(obj)
	for key, value := range updatesToApply {
		updatedItem[key] = value
	}

	switch item := obj.(type) {
	case AccountUpdate:
		// Check dependencies when updating a account object
	case GroupUpdate:
		// Check dependencies when updating a group object
	case ChannelUpdate:
		// Check dependencies when updating a channel object
	case MessageUpdate:
		// Check dependencies when updating a message object
	default:
		return nil, fmt.Errorf("cannot add unknown item '%v' type to database", item)
	}

	// Updates the item to the database
	op, err := store.Put(context.Background(), updatedItem)
	if err != nil {
		return nil, err
	}

	// Convert to 'readable' interface format instead of string
	var result map[string]interface{}
	err = json.Unmarshal([]byte(op.GetValue()), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

/**
 * Remove an item in the database
 */
func removeItem(store orbitdb.DocumentStore, id types.UUID) error {
	// Get the item to delete from the DB
	dbItem, err := getItem(store, id)
	if err != nil {
		return fmt.Errorf("cannot delete item from database")
	}

	// Based on the type of entry, propogate other deletions or updates
	entry, err := DetectAndUnmarshal(dbItem.(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("cannot determine type of item to delete from database")
	}

	switch entry.(type) {
	case Account:
		// When deleting an account, update groups so that there are no references to the account within the group members list
	case Group:
		// When deleting a group, recursively delete all channels
	case Channel:
		// When deleting a channel, recursively delete all related messages
	case Message:
		// When deleting a message, nothing special is needed
	}

	// Delete the item
	_, err = store.Delete(context.Background(), id.String())
	if err != nil {
		return err
	}
	return nil
}

/**
 * Search for items in the database satisfying some filter!
 *
 * IF a field is null (on object or filter), the filter for it is skipped!
 */
func searchItem(store orbitdb.DocumentStore, filter map[string]interface{}) ([]interface{}, error) {
	containsBehavior := func(entryValue, filterValue interface{}) bool {
		if filterSlice, ok := filterValue.([]interface{}); ok {
			return slices.Contains(filterSlice, entryValue)
		}
		return false
	}

	dateBeforeBehavior := func(entryValue, filterValue interface{}) bool {
		parsedEntryTime, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", entryValue.(string))
		if err != nil {
			return false
		}
		parsedFilterTime, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", filterValue.(string))
		if err != nil {
			return false
		}
		return parsedEntryTime.Before(parsedFilterTime)
	}

	dateAfterBehavior := func(entryValue, filterValue interface{}) bool {
		parsedEntryTime, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", entryValue.(string))
		if err != nil {
			return false
		}
		parsedFilterTime, err := time.Parse("2006-01-02T15:04:05.000000000-07:00", filterValue.(string))
		if err != nil {
			return false
		}
		return parsedEntryTime.After(parsedFilterTime)
	}

	fuzzyMatchBehavior := func(entryValue, filterValue interface{}) bool {
		return fuzzy.Match(filterValue.(string), entryValue.(string))
	}

	exactMatchBehavior := func(entryValue, filterValue interface{}) bool {
		return entryValue == filterValue
	}

	// A filter behavior call will return false if the filter fails, and true if it passes
	filterBehaviors := map[string]func(entryValue, filterValue interface{}) bool{
		"id":       containsBehavior,
		"author":   containsBehavior,
		"channel":  containsBehavior,
		"from":     dateAfterBehavior,
		"until":    dateBeforeBehavior,
		"username": fuzzyMatchBehavior,
		"name":     fuzzyMatchBehavior,
		"body":     fuzzyMatchBehavior,
		"pinned":   exactMatchBehavior,
	}

	// Actually perform the query
	result, err := store.Query(context.Background(), func(doc interface{}) (bool, error) {
		entry, ok := doc.(map[string]interface{})
		if !ok {
			return false, nil
		}

		// Apply filters and discard 'entry' if not a match
		for key, value := range filter {
			entryKey := key
			if key == "from" || key == "until" {
				entryKey = "created_at"
			}

			// If there is a nil, we don't process that filter
			if entry[entryKey] == nil || value == nil {
				continue
			}

			// Use the associated filter behavior to determine if we discard this or not
			if !filterBehaviors[key](entry[entryKey], value) {
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

//

// Whenever we want to convert something from a struct to the database representation use this.
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
	// Marshal map to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Unmarshal JSON into Group struct
	err = json.Unmarshal(jsonData, &obj)
	if err != nil {
		return err
	}
	return nil
}

// Function to determine struct type
func DetectAndUnmarshal(data map[string]interface{}) (interface{}, error) {
	// Convert map back to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// List all possible struct types
	var possibleTypes = []interface{}{&Account{}, &Group{}, &Channel{}, &Message{}}

	// Try unmarshaling into each struct
	for _, candidate := range possibleTypes {
		// Create a new instance of the candidate type
		target := reflect.New(reflect.TypeOf(candidate).Elem()).Interface()

		err := json.Unmarshal(jsonData, target)
		if err == nil {
			// Check if any fields were actually populated
			if !isEmpty(target) {
				return target, nil
			}
		}
	}

	return nil, fmt.Errorf("unknown struct type")
}

// Function to check if a struct is empty (i.e., no fields populated)
func isEmpty(v interface{}) bool {
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsZero() {
			return false
		}
	}
	return true
}
