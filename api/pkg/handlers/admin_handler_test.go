package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NdoleStudio/httpsms/pkg/config"
	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/middlewares"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/httpsms/pkg/validators"
	"github.com/NdoleStudio/stacktrace"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const handlerTestAdminEmail = "admin@example.com"

type adminHandlerRepositoryStub struct {
	repositories.UserRepository
	users       map[entities.UserID]*entities.User
	indexParams *repositories.UserAdminIndexParams
	updates     int
}

func (stub *adminHandlerRepositoryStub) Load(_ context.Context, userID entities.UserID) (*entities.User, error) {
	user, ok := stub.users[userID]
	if !ok {
		return nil, stacktrace.PropagateWithCodef(gorm.ErrRecordNotFound, repositories.ErrCodeNotFound, "user [%s] not found", userID)
	}
	clone := *user
	return &clone, nil
}

func (stub *adminHandlerRepositoryStub) LoadOrStore(ctx context.Context, authUser entities.AuthContext, status entities.UserStatus) (*entities.User, bool, error) {
	if _, ok := stub.users[authUser.ID]; !ok {
		stub.users[authUser.ID] = &entities.User{ID: authUser.ID, Email: authUser.Email, Status: status}
	}
	user, err := stub.Load(ctx, authUser.ID)
	return user, false, err
}

func (stub *adminHandlerRepositoryStub) UpdateStatus(ctx context.Context, userID entities.UserID, update repositories.UserStatusUpdate) (*entities.User, error) {
	stub.updates++
	stub.users[userID].Status = update.Status
	stub.users[userID].ApprovedAt = update.ApprovedAt
	stub.users[userID].ApprovedBy = update.ApprovedBy
	return stub.Load(ctx, userID)
}

func (stub *adminHandlerRepositoryStub) LoadStatus(ctx context.Context, userID entities.UserID) (entities.UserStatus, error) {
	user, err := stub.Load(ctx, userID)
	if err != nil {
		return "", err
	}
	return user.Status, nil
}

func (stub *adminHandlerRepositoryStub) CountPhones(context.Context, entities.UserID) (int64, error) {
	return 3, nil
}

func (stub *adminHandlerRepositoryStub) AdminIndex(_ context.Context, params repositories.UserAdminIndexParams) ([]repositories.UserWithPhonesCount, error) {
	stub.indexParams = &params
	return []repositories.UserWithPhonesCount{{User: stub.users["pending-user"], PhonesCount: 1}}, nil
}

func newAdminHandlerTestApp(t *testing.T, authUser entities.AuthContext) (*fiber.App, *adminHandlerRepositoryStub) {
	t.Helper()
	logger := &messageThreadHandlerNoopLogger{}
	tracer := telemetry.NewOtelLogger("test", logger)
	repository := &adminHandlerRepositoryStub{users: map[entities.UserID]*entities.User{
		"admin-user":   {ID: "admin-user", Email: handlerTestAdminEmail, Status: entities.UserStatusActive},
		"pending-user": {ID: "pending-user", Email: "pending@example.com", Status: entities.UserStatusPending},
		"other-admin":  {ID: "other-admin", Email: "Other@Example.com", Status: entities.UserStatusActive},
	}}
	policy := config.NewAdminPolicy(handlerTestAdminEmail + ",other@example.com")
	userService := services.NewUserService(logger, tracer, repository, nil, nil, nil, nil, policy)
	adminService := services.NewAdminService(logger, tracer, policy, repository, userService, nil, nil, nil, false, func() entities.UserID { return "system-user" })

	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals(middlewares.ContextKeyAuthUserID, authUser)
		return c.Next()
	})
	NewAdminHandler(logger, tracer, validators.NewAdminHandlerValidator(logger, tracer), adminService).RegisterRoutes(
		app,
		middlewares.Authenticated(tracer),
		middlewares.RequireAdmin(tracer, adminService),
	)
	NewUserHandler(logger, tracer, validators.NewUserHandlerValidator(logger, tracer, userService), userService).RegisterRoutes(
		app,
		middlewares.AuthenticatedAndApproved(logger, tracer, adminService),
	)
	return app, repository
}

func adminTestRequest(t *testing.T, app *fiber.App, method string, path string, body any) (int, map[string]any) {
	t.Helper()
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		require.NoError(t, err)
		reader = bytes.NewReader(payload)
	}
	request := httptest.NewRequest(method, path, reader)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	require.NoError(t, err)
	defer response.Body.Close()

	result := map[string]any{}
	payload, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	if len(payload) > 0 {
		require.NoError(t, json.Unmarshal(payload, &result))
	}
	return response.StatusCode, result
}

func handlerVerifiedAdmin() entities.AuthContext {
	return entities.AuthContext{ID: "admin-user", Email: handlerTestAdminEmail, EmailVerified: true, FirebaseIDToken: true}
}

func TestAdminHandler_Index_CapsLimitAndValidates(t *testing.T) {
	app, repository := newAdminHandlerTestApp(t, handlerVerifiedAdmin())

	status, body := adminTestRequest(t, app, http.MethodGet, "/v1/admin/users?status=pending&query=pend&limit=500&skip=5", nil)
	require.Equal(t, http.StatusOK, status)
	require.NotNil(t, repository.indexParams)
	assert.Equal(t, 100, repository.indexParams.Limit)
	assert.Equal(t, 5, repository.indexParams.Skip)
	assert.Equal(t, "pend", repository.indexParams.Query)
	require.NotNil(t, repository.indexParams.Status)
	assert.Equal(t, entities.UserStatusPending, *repository.indexParams.Status)
	assert.Equal(t, "success", body["status"])
	users := body["data"].([]any)
	require.Len(t, users, 1)
	user := users[0].(map[string]any)
	for _, key := range []string{"id", "email", "status", "is_admin", "created_at", "updated_at", "approved_at", "approved_by", "phones_count"} {
		assert.Contains(t, user, key)
	}
	assert.NotContains(t, user, "api_key")

	_, _ = adminTestRequest(t, app, http.MethodGet, "/v1/admin/users", nil)
	assert.Equal(t, services.AdminUsersDefaultLimit, repository.indexParams.Limit)
	assert.Nil(t, repository.indexParams.Status)

	for _, query := range []string{"status=deleted", "limit=abc", "skip=-1", "limit=-5"} {
		status, _ = adminTestRequest(t, app, http.MethodGet, "/v1/admin/users?"+query, nil)
		assert.Equal(t, http.StatusUnprocessableEntity, status, query)
	}
}

func TestAdminHandler_UpdateStatus(t *testing.T) {
	app, repository := newAdminHandlerTestApp(t, handlerVerifiedAdmin())

	status, body := adminTestRequest(t, app, http.MethodPut, "/v1/admin/users/pending-user/status", map[string]string{"status": "active"})
	require.Equal(t, http.StatusOK, status)
	data := body["data"].(map[string]any)
	assert.Equal(t, "active", data["status"])
	assert.Equal(t, handlerTestAdminEmail, data["approved_by"])
	assert.NotNil(t, data["approved_at"])
	assert.Equal(t, float64(3), data["phones_count"])

	tests := []struct {
		name   string
		path   string
		body   any
		status int
	}{
		{"self", "/v1/admin/users/admin-user/status", map[string]string{"status": "suspended"}, http.StatusUnprocessableEntity},
		{"other admin", "/v1/admin/users/other-admin/status", map[string]string{"status": "suspended"}, http.StatusUnprocessableEntity},
		{"missing user", "/v1/admin/users/missing-user/status", map[string]string{"status": "active"}, http.StatusNotFound},
		{"invalid status", "/v1/admin/users/pending-user/status", map[string]string{"status": "pending"}, http.StatusUnprocessableEntity},
		{"empty status", "/v1/admin/users/pending-user/status", map[string]string{}, http.StatusUnprocessableEntity},
		{"invalid user id", "/v1/admin/users/bad%20id%27/status", map[string]string{"status": "active"}, http.StatusUnprocessableEntity},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			status, body := adminTestRequest(t, app, http.MethodPut, test.path, test.body)
			assert.Equal(t, test.status, status)
			assert.Equal(t, "error", body["status"])
		})
	}
	assert.Equal(t, 1, repository.updates)
}

func TestAdminHandler_Delete_Protections(t *testing.T) {
	app, repository := newAdminHandlerTestApp(t, handlerVerifiedAdmin())

	for path, expected := range map[string]int{
		"/v1/admin/users/admin-user":   http.StatusUnprocessableEntity,
		"/v1/admin/users/other-admin":  http.StatusUnprocessableEntity,
		"/v1/admin/users/missing-user": http.StatusNotFound,
	} {
		status, _ := adminTestRequest(t, app, http.MethodDelete, path, nil)
		assert.Equal(t, expected, status, path)
	}
	assert.Len(t, repository.users, 3)
}

func TestAdminHandler_RejectsNonAdmins(t *testing.T) {
	for name, authUser := range map[string]entities.AuthContext{
		"admin user api key":   {ID: "admin-user", Email: handlerTestAdminEmail},
		"unverified admin":     {ID: "admin-user", Email: handlerTestAdminEmail, FirebaseIDToken: true},
		"verified non admin":   {ID: "pending-user", Email: "pending@example.com", EmailVerified: true, FirebaseIDToken: true},
		"pending user api key": {ID: "pending-user", Email: "pending@example.com"},
	} {
		t.Run(name, func(t *testing.T) {
			app, repository := newAdminHandlerTestApp(t, authUser)
			for _, request := range []struct{ method, path string }{
				{http.MethodGet, "/v1/admin/users"},
				{http.MethodPut, "/v1/admin/users/pending-user/status"},
				{http.MethodDelete, "/v1/admin/users/pending-user"},
			} {
				status, body := adminTestRequest(t, app, request.method, request.path, map[string]string{"status": "active"})
				assert.Equal(t, http.StatusForbidden, status)
				assert.Equal(t, map[string]any{"status": "error", "message": "forbidden"}, body)
			}
			assert.Zero(t, repository.updates)
			assert.Nil(t, repository.indexParams)
		})
	}
}

func TestUserHandler_Show_ReturnsIsAdminAndStatus(t *testing.T) {
	tests := []struct {
		name    string
		auth    entities.AuthContext
		isAdmin bool
		status  string
	}{
		{"verified admin", handlerVerifiedAdmin(), true, "active"},
		{"admin's user api key", entities.AuthContext{ID: "admin-user", Email: handlerTestAdminEmail}, false, "active"},
		{"pending user", entities.AuthContext{ID: "pending-user", Email: "pending@example.com"}, false, "pending"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			app, _ := newAdminHandlerTestApp(t, test.auth)

			status, body := adminTestRequest(t, app, http.MethodGet, "/v1/users/me", nil)

			require.Equal(t, http.StatusOK, status)
			data := body["data"].(map[string]any)
			assert.Equal(t, test.isAdmin, data["is_admin"])
			assert.Equal(t, test.status, data["status"])
			assert.Contains(t, data, "api_key")
		})
	}

	app, _ := newAdminHandlerTestApp(t, entities.AuthContext{ID: "pending-user", Email: "pending@example.com"})
	status, body := adminTestRequest(t, app, http.MethodPut, "/v1/users/me", map[string]string{"timezone": "UTC"})
	assert.Equal(t, http.StatusForbidden, status)
	assert.Equal(t, map[string]any{"code": "account_pending"}, body["data"])
}
