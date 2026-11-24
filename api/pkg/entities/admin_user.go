package entities

import "time"

// AdminUser is the view of an account shown to administrators
type AdminUser struct {
	ID          UserID     `json:"id" example:"WB7DRDWrJZRGbYrv2CKGkqbzvqdC"`
	Email       string     `json:"email" example:"name@email.com"`
	Status      UserStatus `json:"status" example:"pending"`
	IsAdmin     bool       `json:"is_admin" example:"false"`
	CreatedAt   time.Time  `json:"created_at" example:"2022-06-05T14:26:02.302718+03:00"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2022-06-05T14:26:10.303278+03:00"`
	ApprovedAt  *time.Time `json:"approved_at" example:"2022-06-05T14:26:02.302718+03:00" validate:"optional"`
	ApprovedBy  *string    `json:"approved_by" example:"admin@email.com" validate:"optional"`
	PhonesCount int64      `json:"phones_count" example:"1"`
}
