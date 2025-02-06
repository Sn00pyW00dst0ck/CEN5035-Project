// Package models provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package models

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Account User Account Details.
type Account struct {
	Email openapi_types.Email `json:"email"`
	Id    openapi_types.UUID  `json:"id"`
	Name  string              `json:"name"`
}

// AccountFilter User account filtering/searching options.
type AccountFilter struct {
	Name *string `json:"name,omitempty"`
}

// Group A group chat/server of users.
type Group struct {
	Description string             `json:"description"`
	Id          openapi_types.UUID `json:"id"`
	Members     []string           `json:"members"`
	Name        string             `json:"name"`
}

// GroupFilter Group filtering/searching options.
type GroupFilter struct {
	Name *string `json:"name,omitempty"`
}
