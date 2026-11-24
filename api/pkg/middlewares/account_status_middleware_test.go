package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/NdoleStudio/httpsms/pkg/config"
	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

type middlewareNoopLogger struct{}

func (l *middlewareNoopLogger) Error(error)                         {}
func (l *middlewareNoopLogger) WithService(string) telemetry.Logger { return l }
func (l *middlewareNoopLogger) WithString(string, string) telemetry.Logger {
	return l
}
func (l *middlewareNoopLogger) WithSpan(trace.Span) telemetry.Logger { return l }
func (l *middlewareNoopLogger) Trace(string)                         {}
func (l *middlewareNoopLogger) Info(string)                          {}
func (l *middlewareNoopLogger) Warn(error)                           {}
func (l *middlewareNoopLogger) Debug(string)                         {}
func (l *middlewareNoopLogger) Fatal(error)                          {}
func (l *middlewareNoopLogger) Printf(string, ...interface{})        {}

const (
	testAdminEmail   = "admin@example.com"
	testSystemUserID = entities.UserID("system-user")
)

// stubGate implements AccountGate with a real config.AdminPolicy and in-memory statuses
type stubGate struct {
	policy    *config.AdminPolicy
	statuses  map[entities.UserID]entities.UserStatus
	loadErr   error
	loadCalls int
}

func (gate *stubGate) Enabled() bool { return gate.policy.Enabled() }
func (gate *stubGate) IsAdmin(authUser entities.AuthContext) bool {
	return gate.policy.IsAdmin(authUser)
}
func (gate *stubGate) IsSystemUser(userID entities.UserID) bool { return userID == testSystemUserID }
func (gate *stubGate) LoadStatus(_ context.Context, userID entities.UserID) (entities.UserStatus, error) {
	gate.loadCalls++
	if gate.loadErr != nil {
		return "", gate.loadErr
	}
	status, ok := gate.statuses[userID]
	if !ok {
		return "", stacktrace.NewErrorWithCodef(repositories.ErrCodeNotFound, "user [%s] not found", userID)
	}
	return status, nil
}

func newStubGate(adminEmails string) *stubGate {
	return &stubGate{
		policy: config.NewAdminPolicy(adminEmails),
		statuses: map[entities.UserID]entities.UserStatus{
			"active-user":    entities.UserStatusActive,
			"pending-user":   entities.UserStatusPending,
			"suspended-user": entities.UserStatusSuspended,
			"admin-user":     entities.UserStatusPending,
		},
	}
}

func apiKeyAuth(id entities.UserID) *entities.AuthContext {
	return &entities.AuthContext{ID: id, Email: string(id) + "@example.com"}
}

func phoneKeyAuth(id entities.UserID) *entities.AuthContext {
	keyID := uuid.New()
	return &entities.AuthContext{ID: id, Email: string(id) + "@example.com", PhoneAPIKeyID: &keyID, PhoneNumbers: []string{"+18005550199"}}
}

func firebaseAuth(id entities.UserID, email string, verified bool) *entities.AuthContext {
	return &entities.AuthContext{ID: id, Email: email, EmailVerified: verified, FirebaseIDToken: true}
}

func newGateApp(gate AccountGate, authUser *entities.AuthContext) *fiber.App {
	logger := &middlewareNoopLogger{}
	tracer := telemetry.NewOtelLogger("test", logger)
	app := fiber.New()
	// stands in for the global BearerAuth / APIKeyAuth middlewares
	app.Use(func(c fiber.Ctx) error {
		if authUser != nil {
			c.Locals(ContextKeyAuthUserID, *authUser)
		}
		return c.Next()
	})

	ok := func(c fiber.Ctx) error { return c.JSON(fiber.Map{"status": "success"}) }
	gated := AuthenticatedAndApproved(logger, tracer, gate)
	app.Get("/v1/users/me", gated, ok)
	app.Put("/v1/users/me", gated, ok)
	app.Delete("/v1/users/me", gated, ok)
	app.Get("/v1/phones", gated, ok)
	app.Post("/v1/events", gated, ok)
	app.Get("/v1/admin/users", Authenticated(tracer), RequireAdmin(tracer, gate), ok)
	return app
}

type gateResponse struct {
	status int
	body   map[string]any
}

func (response gateResponse) code() any {
	data, ok := response.body["data"].(map[string]any)
	if !ok {
		return nil
	}
	return data["code"]
}

func callGate(t *testing.T, gate AccountGate, method string, path string, authUser *entities.AuthContext) gateResponse {
	t.Helper()

	response, err := newGateApp(gate, authUser).Test(httptest.NewRequest(method, path, nil))
	require.NoError(t, err)
	defer response.Body.Close()

	payload, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	body := map[string]any{}
	require.NoError(t, json.Unmarshal(payload, &body))
	return gateResponse{status: response.StatusCode, body: body}
}

func TestAuthenticatedAndApproved_BlocksInactiveAccountsForEveryAuthMethod(t *testing.T) {
	gate := newStubGate(testAdminEmail)

	tests := []struct {
		name   string
		auth   *entities.AuthContext
		code   string
		status float64
	}{
		{"pending user api key", apiKeyAuth("pending-user"), AccountStatusCodePending, 0},
		{"pending phone api key", phoneKeyAuth("pending-user"), AccountStatusCodePending, 0},
		{"pending firebase", firebaseAuth("pending-user", "pending@example.com", true), AccountStatusCodePending, 0},
		{"suspended user api key", apiKeyAuth("suspended-user"), AccountStatusCodeSuspended, 0},
		{"suspended phone api key", phoneKeyAuth("suspended-user"), AccountStatusCodeSuspended, 0},
		{"suspended firebase", firebaseAuth("suspended-user", "suspended@example.com", true), AccountStatusCodeSuspended, 0},
		{"user without a row yet", firebaseAuth("new-user", "new@example.com", true), AccountStatusCodePending, 0},
		{"admin email but unverified", firebaseAuth("admin-user", testAdminEmail, false), AccountStatusCodePending, 0},
		{"admin email through api key", apiKeyAuth("admin-user"), AccountStatusCodePending, 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := callGate(t, gate, fiber.MethodGet, "/v1/phones", test.auth)

			assert.Equal(t, fiber.StatusForbidden, response.status)
			assert.Equal(t, "error", response.body["status"])
			assert.NotEmpty(t, response.body["message"])
			assert.Equal(t, test.code, response.code())
		})
	}
}

func TestAuthenticatedAndApproved_AllowsActiveAccounts(t *testing.T) {
	gate := newStubGate(testAdminEmail)

	for _, authUser := range []*entities.AuthContext{
		apiKeyAuth("active-user"),
		phoneKeyAuth("active-user"),
		firebaseAuth("active-user", "active@example.com", false),
	} {
		response := callGate(t, gate, fiber.MethodGet, "/v1/phones", authUser)
		assert.Equal(t, fiber.StatusOK, response.status)
	}
}

func TestAuthenticatedAndApproved_Exceptions(t *testing.T) {
	gate := newStubGate(testAdminEmail)

	t.Run("verified admin with a pending row", func(t *testing.T) {
		response := callGate(t, gate, fiber.MethodGet, "/v1/phones", firebaseAuth("admin-user", "Admin@Example.com", true))
		assert.Equal(t, fiber.StatusOK, response.status)
	})

	t.Run("system user", func(t *testing.T) {
		gate.statuses[testSystemUserID] = entities.UserStatusSuspended
		response := callGate(t, gate, fiber.MethodPost, "/v1/events", apiKeyAuth(testSystemUserID))
		assert.Equal(t, fiber.StatusOK, response.status)
	})

	for _, method := range []string{fiber.MethodGet, fiber.MethodDelete} {
		for _, path := range []string{"/v1/users/me", "/V1/Users/Me", "/v1/users/me/"} {
			t.Run(method+" "+path, func(t *testing.T) {
				for _, authUser := range []*entities.AuthContext{apiKeyAuth("pending-user"), apiKeyAuth("suspended-user"), firebaseAuth("new-user", "new@example.com", false)} {
					response := callGate(t, gate, method, path, authUser)
					assert.Equal(t, fiber.StatusOK, response.status)
				}
			})
		}
	}

	t.Run("PUT /v1/users/me is not an exception", func(t *testing.T) {
		response := callGate(t, gate, fiber.MethodPut, "/v1/users/me", apiKeyAuth("pending-user"))
		assert.Equal(t, fiber.StatusForbidden, response.status)
		assert.Equal(t, AccountStatusCodePending, response.code())
	})
}

func TestAuthenticatedAndApproved_DisabledGateKeepsCurrentBehaviour(t *testing.T) {
	gate := newStubGate("")

	response := callGate(t, gate, fiber.MethodGet, "/v1/phones", apiKeyAuth("suspended-user"))

	assert.Equal(t, fiber.StatusOK, response.status)
	assert.Zero(t, gate.loadCalls)
}

func TestAuthenticatedAndApproved_Unauthenticated(t *testing.T) {
	response := callGate(t, newStubGate(testAdminEmail), fiber.MethodGet, "/v1/phones", nil)
	assert.Equal(t, fiber.StatusUnauthorized, response.status)

	response = callGate(t, newStubGate(testAdminEmail), fiber.MethodGet, "/v1/users/me", nil)
	assert.Equal(t, fiber.StatusUnauthorized, response.status)
}

func TestAuthenticatedAndApproved_FailsClosedOnLookupError(t *testing.T) {
	gate := newStubGate(testAdminEmail)
	gate.loadErr = errors.New("database is down")

	response := callGate(t, gate, fiber.MethodGet, "/v1/phones", apiKeyAuth("active-user"))

	assert.Equal(t, fiber.StatusInternalServerError, response.status)
}

func TestRequireAdmin(t *testing.T) {
	gate := newStubGate(testAdminEmail)

	tests := []struct {
		name   string
		auth   *entities.AuthContext
		status int
	}{
		{"verified firebase admin", firebaseAuth("admin-user", testAdminEmail, true), fiber.StatusOK},
		{"unverified firebase admin", firebaseAuth("admin-user", testAdminEmail, false), fiber.StatusForbidden},
		{"admin's user api key", apiKeyAuth("admin-user"), fiber.StatusForbidden},
		{"admin's phone api key", &entities.AuthContext{ID: "admin-user", Email: testAdminEmail, PhoneAPIKeyID: &uuid.Nil}, fiber.StatusForbidden},
		{"active non admin", firebaseAuth("active-user", "active@example.com", true), fiber.StatusForbidden},
		{"pending non admin", firebaseAuth("pending-user", "pending@example.com", true), fiber.StatusForbidden},
		{"unauthenticated", nil, fiber.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := callGate(t, gate, fiber.MethodGet, "/v1/admin/users", test.auth)
			assert.Equal(t, test.status, response.status)
			if test.status == fiber.StatusForbidden {
				assert.Equal(t, map[string]any{"status": "error", "message": "forbidden"}, response.body)
			}
		})
	}

	response := callGate(t, newStubGate(""), fiber.MethodGet, "/v1/admin/users", firebaseAuth("admin-user", testAdminEmail, true))
	assert.Equal(t, fiber.StatusForbidden, response.status)
}

func TestAuthContextFromToken_NeverPanics(t *testing.T) {
	assert.Equal(t, entities.AuthContext{}, authContextFromToken(nil))

	tests := []struct {
		name     string
		token    *auth.Token
		expected entities.AuthContext
	}{
		{
			"verified email",
			&auth.Token{UID: "uid", Claims: map[string]interface{}{"user_id": "uid", "email": "a@example.com", "email_verified": true}},
			entities.AuthContext{ID: "uid", Email: "a@example.com", EmailVerified: true, FirebaseIDToken: true},
		},
		{
			"unverified email",
			&auth.Token{UID: "uid", Claims: map[string]interface{}{"user_id": "uid", "email": "a@example.com", "email_verified": false}},
			entities.AuthContext{ID: "uid", Email: "a@example.com", FirebaseIDToken: true},
		},
		{
			"email_verified is not a bool",
			&auth.Token{UID: "uid", Claims: map[string]interface{}{"user_id": "uid", "email": "a@example.com", "email_verified": "true"}},
			entities.AuthContext{ID: "uid", Email: "a@example.com", FirebaseIDToken: true},
		},
		{
			"missing claims",
			&auth.Token{UID: "uid", Claims: map[string]interface{}{}},
			entities.AuthContext{ID: "uid", FirebaseIDToken: true},
		},
		{
			"claims with wrong types",
			&auth.Token{UID: "uid", Claims: map[string]interface{}{"user_id": 42, "email": []string{"a@example.com"}}},
			entities.AuthContext{ID: "uid", FirebaseIDToken: true},
		},
		{
			"nil claims",
			&auth.Token{UID: "uid"},
			entities.AuthContext{ID: "uid", FirebaseIDToken: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotPanics(t, func() {
				assert.Equal(t, test.expected, authContextFromToken(test.token))
			})
		})
	}
}
