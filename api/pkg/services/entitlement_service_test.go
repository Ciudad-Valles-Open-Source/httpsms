package services

import (
	"context"
	"testing"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type entitlementUserRepository struct {
	repositories.UserRepository
	user      *entities.User
	loadCalls int
}

func (repository *entitlementUserRepository) Load(_ context.Context, _ entities.UserID) (*entities.User, error) {
	repository.loadCalls++
	return repository.user, nil
}

func newEntitlementServiceForTest(enabled bool, repository repositories.UserRepository) *EntitlementService {
	logger := newRecordingLogger()
	return NewEntitlementService(logger, telemetry.NewOtelLogger("test", logger), enabled, repository)
}

func TestEntitlementService_CheckAdditional_DisabledAllowsWithoutLoadingUserOrCount(t *testing.T) {
	repository := &entitlementUserRepository{}
	service := newEntitlementServiceForTest(false, repository)
	countCalls := 0

	result, err := service.CheckAdditional(context.Background(), "user-id", "Contact", 1000, func() (int, error) {
		countCalls++
		return 200, nil
	})

	require.NoError(t, err)
	assert.True(t, result.Allowed)
	assert.Zero(t, repository.loadCalls)
	assert.Zero(t, countCalls)
}

func TestEntitlementService_CheckAdditional_NeverLimitsUsers(t *testing.T) {
	entityNames := []string{
		entities.EntityNameContact,
		entities.EntityNamePhoneAPIKey,
		entities.EntityNameMessageSendSchedule,
	}

	for _, entityName := range entityNames {
		t.Run(entityName, func(t *testing.T) {
			repository := &entitlementUserRepository{
				user: &entities.User{SubscriptionName: entities.SubscriptionNameFree},
			}
			service := newEntitlementServiceForTest(true, repository)

			result, err := service.CheckAdditional(context.Background(), "user-id", entityName, 1_000_000, func() (int, error) {
				return 1_000_000, nil
			})

			require.NoError(t, err)
			assert.True(t, result.Allowed)
			assert.Empty(t, result.Message)
		})
	}
}
