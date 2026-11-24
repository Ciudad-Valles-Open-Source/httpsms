package repositories

import (
	"context"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
)

// UserAdminIndexParams are the filters used by administrators to list users
type UserAdminIndexParams struct {
	Status *entities.UserStatus
	Query  string
	Skip   int
	Limit  int
}

// UserWithPhonesCount is an entities.User with the number of phones it owns
type UserWithPhonesCount struct {
	User        *entities.User
	PhonesCount int64
}

// UserStatusUpdate changes the approval status of an entities.User
type UserStatusUpdate struct {
	Status     entities.UserStatus
	ApprovedAt *time.Time
	ApprovedBy *string
}

// UserRepository loads and persists an entities.User
type UserRepository interface {
	// Store a new entities.User. The user.Status must be set explicitly.
	Store(ctx context.Context, user *entities.User) error

	// Update a new entities.User. It never changes the status, approved_at or approved_by columns; use UpdateStatus.
	Update(ctx context.Context, user *entities.User) error

	// UpdateStatus changes the approval status of an entities.User and invalidates the status cache
	UpdateStatus(ctx context.Context, userID entities.UserID, update UserStatusUpdate) (*entities.User, error)

	// LoadStatus returns the entities.UserStatus of a user (cached for a few seconds)
	LoadStatus(ctx context.Context, userID entities.UserID) (entities.UserStatus, error)

	// AdminIndex lists users for administrators, newest first
	AdminIndex(ctx context.Context, params UserAdminIndexParams) ([]UserWithPhonesCount, error)

	// CountPhones returns the number of phones owned by a user
	CountPhones(ctx context.Context, userID entities.UserID) (int64, error)

	// LoadAuthContext fetches an entities.AuthContext by apiKey
	LoadAuthContext(ctx context.Context, apiKey string) (entities.AuthContext, error)

	// Load an entities.User by entities.UserID
	Load(ctx context.Context, userID entities.UserID) (*entities.User, error)

	// RotateAPIKey updates the API Key of a user
	RotateAPIKey(ctx context.Context, userID entities.UserID) (*entities.User, error)

	// LoadOrStore an entities.User by entities.AuthContext. A new user is created atomically with newUserStatus.
	LoadOrStore(ctx context.Context, user entities.AuthContext, newUserStatus entities.UserStatus) (*entities.User, bool, error)

	// LoadByEmail loads a user based on the email
	LoadByEmail(ctx context.Context, email string) (*entities.User, error)

	// Delete an entities.User by entities.UserID
	Delete(ctx context.Context, user *entities.User) error
}
