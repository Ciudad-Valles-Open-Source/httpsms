package entities

import "github.com/google/uuid"

// AuthContext is the user gotten from an auth request
type AuthContext struct {
	ID            UserID     `json:"id"`
	PhoneAPIKeyID *uuid.UUID `json:"phone_api_key_id"`
	PhoneNumbers  []string   `json:"phone_numbers"`
	Email         string     `json:"email"`

	// EmailVerified is the "email_verified" claim of a Firebase ID token. It is
	// always false for requests authenticated with an API key (user or phone).
	EmailVerified bool `json:"email_verified"`

	// FirebaseIDToken is true only when the request was authenticated with a
	// verified Firebase ID token (Bearer) and not with any API key. It is never
	// serialized so it can't be set from untrusted input.
	FirebaseIDToken bool `json:"-"`
}

// IsNoop checks if a user is empty
func (user AuthContext) IsNoop() bool {
	return user.ID == "" || user.Email == ""
}
