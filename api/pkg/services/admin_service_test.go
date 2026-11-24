package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/config"
	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

const (
	testAdminEmail   = "admin@example.com"
	testSystemUserID = entities.UserID("system-user")
)

type adminUserRepositoryStub struct {
	repositories.UserRepository
	users          map[entities.UserID]*entities.User
	newUserStatus  entities.UserStatus
	statusUpdates  []repositories.UserStatusUpdate
	indexParams    *repositories.UserAdminIndexParams
	phonesCount    int64
	loadOrStoreHit int
}

func newAdminUserRepositoryStub(users ...*entities.User) *adminUserRepositoryStub {
	stub := &adminUserRepositoryStub{users: map[entities.UserID]*entities.User{}, phonesCount: 2}
	for _, user := range users {
		stub.users[user.ID] = user
	}
	return stub
}

func (stub *adminUserRepositoryStub) Load(_ context.Context, userID entities.UserID) (*entities.User, error) {
	user, ok := stub.users[userID]
	if !ok {
		return nil, stacktrace.PropagateWithCodef(gorm.ErrRecordNotFound, repositories.ErrCodeNotFound, "user [%s] not found", userID)
	}
	clone := *user
	return &clone, nil
}

// LoadOrStore never reports a new user so the tests don't need an EventDispatcher
func (stub *adminUserRepositoryStub) LoadOrStore(ctx context.Context, authUser entities.AuthContext, status entities.UserStatus) (*entities.User, bool, error) {
	stub.loadOrStoreHit++
	stub.newUserStatus = status
	if user, err := stub.Load(ctx, authUser.ID); err == nil {
		return user, false, nil
	}
	stub.users[authUser.ID] = &entities.User{ID: authUser.ID, Email: authUser.Email, Status: status}
	user, _ := stub.Load(ctx, authUser.ID)
	return user, false, nil
}

func (stub *adminUserRepositoryStub) UpdateStatus(ctx context.Context, userID entities.UserID, update repositories.UserStatusUpdate) (*entities.User, error) {
	user, ok := stub.users[userID]
	if !ok {
		return nil, stacktrace.PropagateWithCodef(gorm.ErrRecordNotFound, repositories.ErrCodeNotFound, "user [%s] not found", userID)
	}
	stub.statusUpdates = append(stub.statusUpdates, update)
	user.Status = update.Status
	if update.ApprovedAt != nil {
		user.ApprovedAt = update.ApprovedAt
	}
	if update.ApprovedBy != nil {
		user.ApprovedBy = update.ApprovedBy
	}
	return stub.Load(ctx, userID)
}

func (stub *adminUserRepositoryStub) LoadStatus(_ context.Context, userID entities.UserID) (entities.UserStatus, error) {
	return stub.users[userID].Status, nil
}

func (stub *adminUserRepositoryStub) CountPhones(context.Context, entities.UserID) (int64, error) {
	return stub.phonesCount, nil
}

func (stub *adminUserRepositoryStub) AdminIndex(_ context.Context, params repositories.UserAdminIndexParams) ([]repositories.UserWithPhonesCount, error) {
	stub.indexParams = &params
	var result []repositories.UserWithPhonesCount
	for _, user := range stub.users {
		result = append(result, repositories.UserWithPhonesCount{User: user, PhonesCount: stub.phonesCount})
	}
	return result, nil
}

func verifiedAdmin() entities.AuthContext {
	return entities.AuthContext{ID: "admin-user", Email: "Admin@Example.com", EmailVerified: true, FirebaseIDToken: true}
}

func newAdminServiceForTest(logger *recordingLogger, repository repositories.UserRepository, adminEmails string) *AdminService {
	tracer := telemetry.NewOtelLogger("test", logger)
	policy := config.NewAdminPolicy(adminEmails)
	userService := NewUserService(logger, tracer, repository, nil, nil, nil, nil, policy)
	return NewAdminService(logger, tracer, policy, repository, userService, nil, nil, nil, false, func() entities.UserID { return testSystemUserID })
}

func TestAdminService_UpdateStatus_ApproveStampsApprovalFields(t *testing.T) {
	logger := newRecordingLogger()
	repository := newAdminUserRepositoryStub(&entities.User{ID: "pending-user", Email: "user@example.com", Status: entities.UserStatusPending})
	service := newAdminServiceForTest(logger, repository, testAdminEmail)

	before := time.Now().UTC()
	user, err := service.UpdateStatus(context.Background(), "test", verifiedAdmin(), "pending-user", entities.UserStatusActive)

	require.NoError(t, err)
	assert.Equal(t, entities.UserStatusActive, user.Status)
	require.NotNil(t, user.ApprovedAt)
	assert.False(t, user.ApprovedAt.Before(before))
	require.NotNil(t, user.ApprovedBy)
	assert.Equal(t, testAdminEmail, *user.ApprovedBy)
	assert.Equal(t, int64(2), user.PhonesCount)
	assert.False(t, user.IsAdmin)

	logged := strings.Join(logger.infos, "\n")
	assert.Contains(t, logged, "admin [Admin@Example.com] action [set_status:active] target user [pending-user]")
}

func TestAdminService_UpdateStatus_SuspendDoesNotStampApproval(t *testing.T) {
	repository := newAdminUserRepositoryStub(&entities.User{ID: "active-user", Email: "user@example.com", Status: entities.UserStatusActive})
	service := newAdminServiceForTest(newRecordingLogger(), repository, testAdminEmail)

	user, err := service.UpdateStatus(context.Background(), "test", verifiedAdmin(), "active-user", entities.UserStatusSuspended)

	require.NoError(t, err)
	assert.Equal(t, entities.UserStatusSuspended, user.Status)
	require.Len(t, repository.statusUpdates, 1)
	assert.Nil(t, repository.statusUpdates[0].ApprovedAt)
	assert.Nil(t, repository.statusUpdates[0].ApprovedBy)
}

func TestAdminService_ProtectedAccounts(t *testing.T) {
	repository := newAdminUserRepositoryStub(
		&entities.User{ID: "admin-user", Email: testAdminEmail, Status: entities.UserStatusActive},
		&entities.User{ID: "other-admin", Email: "OTHER@example.com", Status: entities.UserStatusActive},
		&entities.User{ID: "same-email", Email: " admin@EXAMPLE.com", Status: entities.UserStatusActive},
		&entities.User{ID: testSystemUserID, Email: "system@httpsms.local", Status: entities.UserStatusActive},
	)
	service := newAdminServiceForTest(newRecordingLogger(), repository, testAdminEmail+",other@example.com")

	for _, userID := range []entities.UserID{"admin-user", "other-admin", "same-email", testSystemUserID} {
		t.Run(string(userID), func(t *testing.T) {
			_, err := service.UpdateStatus(context.Background(), "test", verifiedAdmin(), userID, entities.UserStatusSuspended)
			assert.Equal(t, ErrCodeProtectedUser, stacktrace.GetCode(err))

			err = service.Delete(context.Background(), "test", verifiedAdmin(), userID)
			assert.Equal(t, ErrCodeProtectedUser, stacktrace.GetCode(err))
		})
	}

	assert.Empty(t, repository.statusUpdates)
	assert.Len(t, repository.users, 4)
}

func TestAdminService_RequiresAdmin(t *testing.T) {
	repository := newAdminUserRepositoryStub(&entities.User{ID: "pending-user", Email: "user@example.com", Status: entities.UserStatusPending})
	service := newAdminServiceForTest(newRecordingLogger(), repository, testAdminEmail)

	notAdmins := []entities.AuthContext{
		{ID: "admin-user", Email: testAdminEmail},                                           // user API key
		{ID: "admin-user", Email: testAdminEmail, FirebaseIDToken: true},                    // unverified email
		{ID: "user", Email: "user@example.com", EmailVerified: true, FirebaseIDToken: true}, // not in ADMIN_EMAILS
	}

	for _, authUser := range notAdmins {
		_, err := service.UpdateStatus(context.Background(), "test", authUser, "pending-user", entities.UserStatusActive)
		assert.Equal(t, ErrCodeNotAdmin, stacktrace.GetCode(err))

		_, err = service.Index(context.Background(), authUser, AdminUserIndexParams{})
		assert.Equal(t, ErrCodeNotAdmin, stacktrace.GetCode(err))
	}
	assert.Empty(t, repository.statusUpdates)
}

func TestAdminService_UpdateStatus_Errors(t *testing.T) {
	repository := newAdminUserRepositoryStub(&entities.User{ID: "pending-user", Email: "user@example.com", Status: entities.UserStatusPending})
	service := newAdminServiceForTest(newRecordingLogger(), repository, testAdminEmail)

	_, err := service.UpdateStatus(context.Background(), "test", verifiedAdmin(), "missing-user", entities.UserStatusActive)
	assert.Equal(t, repositories.ErrCodeNotFound, stacktrace.GetCode(err))

	_, err = service.UpdateStatus(context.Background(), "test", verifiedAdmin(), "pending-user", entities.UserStatusPending)
	assert.Error(t, err)
	assert.Empty(t, repository.statusUpdates)
}

func TestAdminService_Index_ClampsLimit(t *testing.T) {
	tests := []struct {
		limit    int
		skip     int
		expected int
		skipped  int
	}{
		{0, 0, AdminUsersDefaultLimit, 0},
		{-5, -3, AdminUsersDefaultLimit, 0},
		{50, 10, 50, 10},
		{100, 0, 100, 0},
		{101, 0, AdminUsersMaxLimit, 0},
		{100000, 0, AdminUsersMaxLimit, 0},
	}

	for _, test := range tests {
		repository := newAdminUserRepositoryStub(&entities.User{ID: "admin-user", Email: testAdminEmail, Status: entities.UserStatusActive})
		service := newAdminServiceForTest(newRecordingLogger(), repository, testAdminEmail)

		users, err := service.Index(context.Background(), verifiedAdmin(), AdminUserIndexParams{Limit: test.limit, Skip: test.skip, Query: "  example "})

		require.NoError(t, err)
		require.NotNil(t, repository.indexParams)
		assert.Equal(t, test.expected, repository.indexParams.Limit)
		assert.Equal(t, test.skipped, repository.indexParams.Skip)
		assert.Equal(t, "example", repository.indexParams.Query)
		require.Len(t, users, 1)
		assert.True(t, users[0].IsAdmin)
	}
}

func TestAdminService_IsSystemUser(t *testing.T) {
	service := newAdminServiceForTest(newRecordingLogger(), newAdminUserRepositoryStub(), testAdminEmail)
	assert.True(t, service.IsSystemUser(testSystemUserID))
	assert.False(t, service.IsSystemUser("other"))

	tracer := telemetry.NewOtelLogger("test", newRecordingLogger())
	empty := NewAdminService(newRecordingLogger(), tracer, config.NewAdminPolicy(testAdminEmail), nil, nil, nil, nil, nil, false, func() entities.UserID { return "" })
	assert.False(t, empty.IsSystemUser(""))
}

func TestUserService_Get_NewUserStatus(t *testing.T) {
	tests := []struct {
		name        string
		adminEmails string
		authUser    entities.AuthContext
		expected    entities.UserStatus
	}{
		{"gate disabled", "", entities.AuthContext{ID: "u", Email: "user@example.com", FirebaseIDToken: true}, entities.UserStatusActive},
		{"gate enabled, regular user", testAdminEmail, entities.AuthContext{ID: "u", Email: "user@example.com", EmailVerified: true, FirebaseIDToken: true}, entities.UserStatusPending},
		{"gate enabled, verified admin", testAdminEmail, entities.AuthContext{ID: "u", Email: "ADMIN@example.com", EmailVerified: true, FirebaseIDToken: true}, entities.UserStatusActive},
		{"gate enabled, unverified admin", testAdminEmail, entities.AuthContext{ID: "u", Email: testAdminEmail, FirebaseIDToken: true}, entities.UserStatusPending},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logger := newRecordingLogger()
			repository := newAdminUserRepositoryStub()
			service := NewUserService(logger, telemetry.NewOtelLogger("test", logger), repository, nil, nil, nil, nil, config.NewAdminPolicy(test.adminEmails))

			user, err := service.Get(context.Background(), "test", test.authUser)

			require.NoError(t, err)
			assert.Equal(t, test.expected, repository.newUserStatus)
			assert.Equal(t, test.expected, user.Status)
		})
	}
}

func TestUserService_Get_ActivatesPendingAdministrator(t *testing.T) {
	logger := newRecordingLogger()
	repository := newAdminUserRepositoryStub(&entities.User{ID: "admin-user", Email: testAdminEmail, Status: entities.UserStatusPending})
	service := NewUserService(logger, telemetry.NewOtelLogger("test", logger), repository, nil, nil, nil, nil, config.NewAdminPolicy(testAdminEmail))

	user, err := service.Get(context.Background(), "test", verifiedAdmin())

	require.NoError(t, err)
	assert.Equal(t, entities.UserStatusActive, user.Status)
	require.NotNil(t, user.ApprovedBy)
	assert.Equal(t, testAdminEmail, *user.ApprovedBy)

	// an API key of the same administrator never activates the account
	repository.users["admin-user"].Status = entities.UserStatusPending
	user, err = service.Get(context.Background(), "test", entities.AuthContext{ID: "admin-user", Email: testAdminEmail})
	require.NoError(t, err)
	assert.Equal(t, entities.UserStatusPending, user.Status)
	assert.False(t, service.IsAdmin(entities.AuthContext{ID: "admin-user", Email: testAdminEmail}))
}
