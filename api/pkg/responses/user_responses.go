package responses

import (
	"github.com/NdoleStudio/httpsms/pkg/entities"
)

// UserResponse is the payload containing entities.User
type UserResponse struct {
	response
	Data entities.User `json:"data"`
}

// CurrentUser is the authenticated entities.User with the role of the request
type CurrentUser struct {
	entities.User
	// IsAdmin is true only when the request uses a verified Firebase ID token whose email is in ADMIN_EMAILS
	IsAdmin bool `json:"is_admin" example:"false"`
}

// CurrentUserResponse is the payload containing the CurrentUser
type CurrentUserResponse struct {
	response
	Data CurrentUser `json:"data"`
}
