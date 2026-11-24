package middlewares

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type apiKeyAuthUserRepository struct {
	repositories.UserRepository
	authUser entities.AuthContext
	calls    int
}

func (repository *apiKeyAuthUserRepository) LoadAuthContext(context.Context, string) (entities.AuthContext, error) {
	repository.calls++
	return repository.authUser, nil
}

func runAPIKeyAuth(t *testing.T, firebaseUser *entities.AuthContext, repository *apiKeyAuthUserRepository) entities.AuthContext {
	t.Helper()
	logger := &middlewareNoopLogger{}
	tracer := telemetry.NewOtelLogger("test", logger)
	app := fiber.New()
	app.Use(func(c fiber.Ctx) error { // stands in for BearerAuth
		if firebaseUser != nil {
			c.Locals(ContextKeyAuthUserID, *firebaseUser)
		}
		return c.Next()
	})
	app.Use(APIKeyAuth(logger, tracer, repository))
	app.Get("/whoami", func(c fiber.Ctx) error {
		authUser, _ := c.Locals(ContextKeyAuthUserID).(entities.AuthContext)
		return c.JSON(fiber.Map{"id": authUser.ID, "firebase": authUser.FirebaseIDToken})
	})

	request := httptest.NewRequest("GET", "/whoami", nil)
	request.Header.Set("x-api-key", "uk_api-key-of-the-same-user")
	response, err := app.Test(request)
	require.NoError(t, err)

	var body struct {
		ID       entities.UserID `json:"id"`
		Firebase bool            `json:"firebase"`
	}
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))
	return entities.AuthContext{ID: body.ID, FirebaseIDToken: body.Firebase}
}

func TestAPIKeyAuth_KeepsFirebaseIdentityWhenBothArePresent(t *testing.T) {
	repository := &apiKeyAuthUserRepository{authUser: entities.AuthContext{ID: "api-key-user", Email: "admin@example.com"}}
	firebase := firebaseAuth("firebase-user", "admin@example.com", true)

	authUser := runAPIKeyAuth(t, firebase, repository)

	assert.Equal(t, entities.UserID("firebase-user"), authUser.ID)
	assert.True(t, authUser.FirebaseIDToken, "the dashboard must keep its Firebase identity (and admin rights)")
	assert.Zero(t, repository.calls, "the API key must not be looked up")
}

func TestAPIKeyAuth_AuthenticatesWithAPIKeyAlone(t *testing.T) {
	repository := &apiKeyAuthUserRepository{authUser: entities.AuthContext{ID: "api-key-user", Email: "user@example.com"}}

	authUser := runAPIKeyAuth(t, nil, repository)

	assert.Equal(t, entities.UserID("api-key-user"), authUser.ID)
	assert.False(t, authUser.FirebaseIDToken)
	assert.Equal(t, 1, repository.calls)
}
