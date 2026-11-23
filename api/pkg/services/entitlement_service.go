package services

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
)

// EntitlementCheckResult holds the outcome of an entitlement check.
type EntitlementCheckResult struct {
	Allowed bool
	Message string
}

// EntitlementService checks whether a user can create more of a given entity
// based on their subscription plan.
type EntitlementService struct {
	service
	logger         telemetry.Logger
	tracer         telemetry.Tracer
	enabled        bool
	userRepository repositories.UserRepository
}

// NewEntitlementService creates a new EntitlementService.
// The enabled flag should come from the ENTITLEMENT_ENABLED environment variable.
func NewEntitlementService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	enabled bool,
	userRepository repositories.UserRepository,
) *EntitlementService {
	return &EntitlementService{
		logger:         logger.WithService(fmt.Sprintf("%T", &EntitlementService{})),
		tracer:         tracer,
		enabled:        enabled,
		userRepository: userRepository,
	}
}

// Check verifies if the user can create another instance of the given entity.
func (service *EntitlementService) Check(
	ctx context.Context,
	userID entities.UserID,
	entityName string,
	countFunc func() (int, error),
) (*EntitlementCheckResult, error) {
	return service.CheckAdditional(ctx, userID, entityName, 1, countFunc)
}

// CheckAdditional always allows creating more entities.
// This fork is 100% free and self-hosted: there are no plan-based limits.
func (service *EntitlementService) CheckAdditional(
	ctx context.Context,
	_ entities.UserID,
	_ string,
	_ int,
	_ func() (int, error),
) (*EntitlementCheckResult, error) {
	_, span := service.tracer.Start(ctx)
	defer span.End()

	return &EntitlementCheckResult{Allowed: true}, nil
}
