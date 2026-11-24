package repositories

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newUserRepositoryForTest(t *testing.T) (UserRepository, *messageThreadTestConnPool, *UserStatusCache) {
	t.Helper()
	pool := &messageThreadTestConnPool{}
	db, err := gorm.Open(
		postgres.New(postgres.Config{Conn: pool, WithoutReturning: true}),
		&gorm.Config{DisableAutomaticPing: true},
	)
	require.NoError(t, err)

	logger := &messageThreadTestLogger{}
	statusCache := NewUserStatusCache(time.Minute)
	return NewGormUserRepository(logger, telemetry.NewOtelLogger("test", logger), nil, statusCache, db), pool, statusCache
}

func TestGormUserRepository_UpdateNeverWritesApprovalColumns(t *testing.T) {
	repository, pool, _ := newUserRepositoryForTest(t)
	approvedBy := "admin@example.com"
	now := time.Now()

	err := repository.Update(context.Background(), &entities.User{
		ID:         "user-id",
		Email:      "user@example.com",
		Timezone:   "UTC",
		Status:     entities.UserStatusActive,
		ApprovedAt: &now,
		ApprovedBy: &approvedBy,
	})

	require.NoError(t, err)
	require.NotEmpty(t, pool.statements)
	update := pool.statements[len(pool.statements)-1].query
	assert.True(t, strings.HasPrefix(update, `UPDATE "users" SET`), update)
	assert.Contains(t, update, `"timezone"=`)
	assert.NotContains(t, update, `"status"=`)
	assert.NotContains(t, update, `"approved_at"=`)
	assert.NotContains(t, update, `"approved_by"=`)
	assert.NotContains(t, update, "ON CONFLICT")
}

func TestGormUserRepository_StoreRequiresAnExplicitStatus(t *testing.T) {
	repository, pool, statusCache := newUserRepositoryForTest(t)

	err := repository.Store(context.Background(), &entities.User{ID: "user-id", Email: "user@example.com"})
	assert.Error(t, err)
	assert.Empty(t, pool.statements)

	err = repository.Store(context.Background(), &entities.User{ID: "user-id", Email: "user@example.com", Status: entities.UserStatusPending})
	require.NoError(t, err)
	insert := pool.statements[len(pool.statements)-1]
	assert.True(t, strings.HasPrefix(insert.query, `INSERT INTO "users"`))
	assert.Contains(t, insert.query, `"status"`)
	assert.Contains(t, insert.args, entities.UserStatusPending)

	status, found := statusCache.get("user-id")
	assert.True(t, found)
	assert.Equal(t, entities.UserStatusPending, status)

	status, err = repository.LoadStatus(context.Background(), "user-id")
	require.NoError(t, err)
	assert.Equal(t, entities.UserStatusPending, status)
}

func TestGormUserRepository_LoadOrStoreRejectsInvalidStatus(t *testing.T) {
	repository, pool, _ := newUserRepositoryForTest(t)

	_, _, err := repository.LoadOrStore(context.Background(), entities.AuthContext{ID: "user-id", Email: "user@example.com"}, "")

	assert.Error(t, err)
	assert.Empty(t, pool.statements)
}

func TestUserStatusCache_ReadsNeverOverwriteWrites(t *testing.T) {
	cache := NewUserStatusCache(time.Minute)

	// an admin suspends the user while a request is still reading the old status from the database
	cache.set("user-id", entities.UserStatusSuspended)
	cache.add("user-id", entities.UserStatusActive)

	status, found := cache.get("user-id")
	assert.True(t, found)
	assert.Equal(t, entities.UserStatusSuspended, status)

	cache.set("user-id", entities.UserStatusActive)
	status, _ = cache.get("user-id")
	assert.Equal(t, entities.UserStatusActive, status)
}

func TestUserStatusCache_Expires(t *testing.T) {
	cache := NewUserStatusCache(20 * time.Millisecond)
	cache.add("user-id", entities.UserStatusActive)

	time.Sleep(40 * time.Millisecond)

	_, found := cache.get("user-id")
	assert.False(t, found)
}

func TestEscapeLikePattern(t *testing.T) {
	assert.Equal(t, `a\%b\_c\\d`, escapeLikePattern(`a%b_c\d`))
}
