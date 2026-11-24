package validators

import (
	"context"
	"fmt"
	"net/url"
	"regexp"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/requests"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/thedevsaddam/govalidator"
)

// userIDPattern matches Firebase UIDs and the UUID of the system user
var userIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

// AdminHandlerValidator validates models used in handlers.AdminHandler
type AdminHandlerValidator struct {
	validator
	logger telemetry.Logger
	tracer telemetry.Tracer
}

// NewAdminHandlerValidator creates a new handlers.AdminHandler validator
func NewAdminHandlerValidator(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
) (v *AdminHandlerValidator) {
	return &AdminHandlerValidator{
		logger: logger.WithService(fmt.Sprintf("%T", v)),
		tracer: tracer,
	}
}

// ValidateIndex validates requests.AdminUserIndex
func (validator *AdminHandlerValidator) ValidateIndex(_ context.Context, request requests.AdminUserIndex) url.Values {
	v := govalidator.New(govalidator.Options{
		Data: &request,
		Rules: govalidator.MapData{
			"status": []string{
				fmt.Sprintf("in:%s,%s,%s", entities.UserStatusPending, entities.UserStatusActive, entities.UserStatusSuspended),
			},
			"query": []string{
				"max:100",
			},
			"skip": []string{
				"required",
				"regex:^[0-9]{1,9}$",
			},
			"limit": []string{
				"required",
				"regex:^[0-9]{1,9}$",
			},
		},
	})
	return v.ValidateStruct()
}

// ValidateStatusUpdate validates requests.AdminUserStatusUpdate
func (validator *AdminHandlerValidator) ValidateStatusUpdate(_ context.Context, request requests.AdminUserStatusUpdate) url.Values {
	v := govalidator.New(govalidator.Options{
		Data: &request,
		Rules: govalidator.MapData{
			"status": []string{
				"required",
				fmt.Sprintf("in:%s,%s", entities.UserStatusActive, entities.UserStatusSuspended),
			},
		},
	})
	return v.ValidateStruct()
}

// ValidateUserID validates the ID of a user in the URL
func (validator *AdminHandlerValidator) ValidateUserID(userID string) url.Values {
	result := url.Values{}
	if !userIDPattern.MatchString(userID) {
		result.Add("userID", "The userID field must contain 1 to 128 letters, digits, '-' or '_'")
	}
	return result
}
