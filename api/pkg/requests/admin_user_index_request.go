package requests

import (
	"strconv"
	"strings"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/services"
)

// AdminUserIndex is the payload for listing users as an administrator
type AdminUserIndex struct {
	request
	Status string `json:"status" query:"status"`
	Query  string `json:"query" query:"query"`
	Skip   string `json:"skip" query:"skip"`
	Limit  string `json:"limit" query:"limit"`
}

// Sanitize sets defaults to AdminUserIndex. A limit above the maximum is capped to services.AdminUsersMaxLimit.
func (input *AdminUserIndex) Sanitize() AdminUserIndex {
	input.Status = strings.ToLower(strings.TrimSpace(input.Status))
	input.Query = strings.TrimSpace(input.Query)

	input.Skip = strings.TrimSpace(input.Skip)
	if input.Skip == "" {
		input.Skip = "0"
	}

	input.Limit = strings.TrimSpace(input.Limit)
	if input.Limit == "" {
		input.Limit = strconv.Itoa(services.AdminUsersDefaultLimit)
	}
	if limit, err := strconv.Atoi(input.Limit); err == nil && limit > services.AdminUsersMaxLimit {
		input.Limit = strconv.Itoa(services.AdminUsersMaxLimit)
	}

	return *input
}

// ToIndexParams converts AdminUserIndex to services.AdminUserIndexParams
func (input *AdminUserIndex) ToIndexParams() services.AdminUserIndexParams {
	params := services.AdminUserIndexParams{
		Query: input.Query,
		Skip:  input.getInt(input.Skip),
		Limit: services.ClampAdminUsersLimit(input.getInt(input.Limit)),
	}
	if input.Status != "" {
		status := entities.UserStatus(input.Status)
		params.Status = &status
	}
	return params
}

// AdminUserStatusUpdate is the payload for changing the status of a user
type AdminUserStatusUpdate struct {
	request
	Status string `json:"status" example:"active"`
}

// Sanitize trims the AdminUserStatusUpdate
func (input *AdminUserStatusUpdate) Sanitize() AdminUserStatusUpdate {
	input.Status = strings.ToLower(strings.TrimSpace(input.Status))
	return *input
}

// UserStatus returns the requested entities.UserStatus
func (input *AdminUserStatusUpdate) UserStatus() entities.UserStatus {
	return entities.UserStatus(input.Status)
}
