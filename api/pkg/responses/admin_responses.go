package responses

import "github.com/NdoleStudio/httpsms/pkg/entities"

// AdminUsersResponse is the payload containing a list of entities.AdminUser
type AdminUsersResponse struct {
	response
	Data []entities.AdminUser `json:"data"`
}

// AdminUserResponse is the payload containing an entities.AdminUser
type AdminUserResponse struct {
	response
	Data entities.AdminUser `json:"data"`
}

// Forbidden is the response with status code 403
type Forbidden struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"forbidden"`
}

// AccountNotActiveData explains why an account can't use the API
type AccountNotActiveData struct {
	Code string `json:"code" example:"account_pending"`
}

// AccountNotActive is the 403 response for a pending or suspended account
type AccountNotActive struct {
	Status  string               `json:"status" example:"error"`
	Message string               `json:"message" example:"Your account is pending approval by an administrator."`
	Data    AccountNotActiveData `json:"data"`
}
