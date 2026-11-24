package events

import (
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
)

// UserAccountApproved is raised when an administrator approves a user's account.
const UserAccountApproved = "user.account.approved"

// UserAccountApprovedPayload stores the data for the UserAccountApproved event
type UserAccountApprovedPayload struct {
	UserID     entities.UserID `json:"user_id"`
	UserEmail  string          `json:"user_email"`
	ApprovedBy string          `json:"approved_by"`
	Timestamp  time.Time       `json:"timestamp"`
}
