package repositories

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/clause"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbgorm"
	"github.com/dgraph-io/ristretto/v2"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"gorm.io/gorm"
)

// gormUserRepository is responsible for persisting entities.User
type gormUserRepository struct {
	logger      telemetry.Logger
	tracer      telemetry.Tracer
	cache       *ristretto.Cache[string, entities.AuthContext]
	statusCache *UserStatusCache
	db          *gorm.DB
}

// userProtectedColumns are only changed by UpdateStatus so that a generic
// Update can never undo (race with) an administrator's approval or suspension.
var userProtectedColumns = []string{"status", "approved_at", "approved_by"}

// NewGormUserRepository creates the GORM version of the UserRepository
func NewGormUserRepository(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	cache *ristretto.Cache[string, entities.AuthContext],
	statusCache *UserStatusCache,
	db *gorm.DB,
) UserRepository {
	return &gormUserRepository{
		logger:      logger.WithService(fmt.Sprintf("%T", &gormUserRepository{})),
		tracer:      tracer,
		cache:       cache,
		statusCache: statusCache,
		db:          db,
	}
}

// Delete deletes a user
func (repository *gormUserRepository) Delete(ctx context.Context, user *entities.User) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if err := repository.db.WithContext(ctx).Delete(user).Error; err != nil {
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot delete user with ID [%s]", user.ID))
	}

	// The deleted user must not keep access through cached credentials or status.
	repository.statusCache.set(user.ID, userStatusDeleted)
	if user.APIKey != "" && repository.cache != nil {
		repository.cache.Wait()
		repository.cache.Del(user.APIKey)
	}

	return nil
}

// UpdateStatus changes the approval status of a user
func (repository *gormUserRepository) UpdateStatus(ctx context.Context, userID entities.UserID, update UserStatusUpdate) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if !update.Status.IsValid() {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorf("cannot set invalid status [%s] for user [%s]", update.Status, userID))
	}

	user := new(entities.User)
	err := crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Where("id = ?", userID).First(user).Error; err != nil {
			return err
		}

		values := map[string]any{
			"status":     update.Status,
			"updated_at": time.Now().UTC(),
		}
		if update.ApprovedAt != nil {
			values["approved_at"] = update.ApprovedAt.UTC()
		}
		if update.ApprovedBy != nil {
			values["approved_by"] = *update.ApprovedBy
		}

		return tx.WithContext(ctx).Model(user).Where("id = ?", userID).Updates(values).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCodef(err, ErrCodeNotFound, "user with ID [%s] does not exist", userID))
	}
	if err != nil {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot update status of user [%s] to [%s]", userID, update.Status))
	}

	repository.statusCache.set(userID, update.Status)
	return user, nil
}

// LoadStatus returns the status of a user, cached for UserStatusCacheTTL
func (repository *gormUserRepository) LoadStatus(ctx context.Context, userID entities.UserID) (entities.UserStatus, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if status, found := repository.statusCache.get(userID); found {
		if status == userStatusDeleted {
			return "", repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCodef(ErrCodeNotFound, "user with ID [%s] was deleted", userID))
		}
		return status, nil
	}

	user := new(entities.User)
	err := repository.db.WithContext(ctx).Select("id", "status").Where("id = ?", userID).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCodef(err, ErrCodeNotFound, "user with ID [%s] does not exist", userID))
	}
	if err != nil {
		return "", repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot load status of user [%s]", userID))
	}

	status := user.Status
	if !status.IsValid() {
		// fail closed: an unknown status is never treated as active
		status = entities.UserStatusPending
	}

	repository.statusCache.add(userID, status)
	return status, nil
}

// AdminIndex lists users for administrators, newest first
func (repository *gormUserRepository) AdminIndex(ctx context.Context, params UserAdminIndexParams) ([]UserWithPhonesCount, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	query := repository.db.WithContext(ctx).Model(&entities.User{})
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	if params.Query != "" {
		query = query.Where("email ILIKE ?", "%"+escapeLikePattern(params.Query)+"%")
	}

	users := make([]*entities.User, 0, params.Limit)
	err := query.Order("created_at DESC").Order("id DESC").Limit(params.Limit).Offset(params.Skip).Find(&users).Error
	if err != nil {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot list users with params [%+#v]", params))
	}

	counts := map[entities.UserID]int64{}
	if len(users) > 0 {
		ids := make([]entities.UserID, 0, len(users))
		for _, user := range users {
			ids = append(ids, user.ID)
		}

		var rows []struct {
			UserID entities.UserID
			Count  int64
		}
		err = repository.db.WithContext(ctx).
			Model(&entities.Phone{}).
			Select("user_id, COUNT(*) AS count").
			Where("user_id IN ?", ids).
			Group("user_id").
			Scan(&rows).Error
		if err != nil {
			return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot count phones of [%d] users", len(ids)))
		}
		for _, row := range rows {
			counts[row.UserID] = row.Count
		}
	}

	result := make([]UserWithPhonesCount, 0, len(users))
	for _, user := range users {
		result = append(result, UserWithPhonesCount{User: user, PhonesCount: counts[user.ID]})
	}
	return result, nil
}

// CountPhones returns the number of phones owned by a user
func (repository *gormUserRepository) CountPhones(ctx context.Context, userID entities.UserID) (int64, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	var count int64
	if err := repository.db.WithContext(ctx).Model(&entities.Phone{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot count phones of user [%s]", userID))
	}
	return count, nil
}

// escapeLikePattern escapes the LIKE wildcards so user input is matched literally
func escapeLikePattern(value string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(value)
}

func (repository *gormUserRepository) RotateAPIKey(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	apiKey, err := repository.generateAPIKey(64)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate apiKey for user [%s]", userID)
	}

	user := new(entities.User)
	var oldAPIKey string
	err = crdbgorm.ExecuteTx(
		ctx, repository.db, nil,
		func(tx *gorm.DB) error {
			if err := tx.WithContext(ctx).Where("id = ?", userID).First(user).Error; err != nil {
				return err
			}
			oldAPIKey = user.APIKey
			return tx.WithContext(ctx).Model(user).
				Clauses(clause.Returning{}).
				Where("id = ?", userID).
				Update("api_key", "uk_"+apiKey).Error
		},
	)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCodef(err, ErrCodeNotFound, "user with ID [%s] does not exist", userID))
	}

	if err == nil && oldAPIKey != "" {
		// Flush pending ristretto Set operations before Del to avoid a
		// buffered Set re-adding the entry after removal.
		repository.cache.Wait()
		repository.cache.Del(oldAPIKey)
	}

	return user, nil
}

func (repository *gormUserRepository) LoadByEmail(ctx context.Context, email string) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	user := new(entities.User)
	err := repository.db.WithContext(ctx).
		Where("email = ?", email).
		First(user).
		Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCodef(err, ErrCodeNotFound, "user with email [%s] does not exist", email))
	}

	if err != nil {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot load user with email ID [%s]", email))
	}

	return user, nil
}

func (repository *gormUserRepository) Store(ctx context.Context, user *entities.User) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	// Never rely on the column default ("active") for a new user.
	if !user.Status.IsValid() {
		return repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorf("cannot save user with ID [%s] and invalid status [%s]", user.ID, user.Status))
	}

	if err := repository.db.WithContext(ctx).Create(user).Error; err != nil {
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot save user with ID [%s]", user.ID))
	}

	repository.statusCache.set(user.ID, user.Status)
	return nil
}

func (repository *gormUserRepository) Update(ctx context.Context, user *entities.User) error {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	// Select("*") updates every column (like Save) without Save's upsert fallback,
	// and Omit keeps the approval columns out of reach of generic updates.
	err := repository.db.WithContext(ctx).Select("*").Omit(userProtectedColumns...).Save(user).Error
	if err != nil {
		return repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot update user with ID [%s]", user.ID))
	}

	return nil
}

func (repository *gormUserRepository) LoadAuthContext(ctx context.Context, apiKey string) (entities.AuthContext, error) {
	ctx, span, ctxLogger := repository.tracer.StartWithLogger(ctx, repository.logger)
	defer span.End()

	if authUser, found := repository.cache.Get(apiKey); found {
		if authUser.IsNoop() {
			return authUser, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorf("user with api key [%s] does not exist", entities.SecretHint(apiKey)))
		}
		return authUser, nil
	}

	user := new(entities.User)
	err := repository.db.WithContext(ctx).Where("api_key = ?", apiKey).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		repository.cache.SetWithTTL(apiKey, entities.AuthContext{}, 1, 2*time.Hour)

		return entities.AuthContext{}, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCodef(err, ErrCodeNotFound, "user with api key [%s] does not exist", entities.SecretHint(apiKey)))
	}

	if err != nil {
		return entities.AuthContext{}, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot load user with api key [%s]", entities.SecretHint(apiKey)))
	}

	authUser := entities.AuthContext{
		ID:    user.ID,
		Email: user.Email,
	}

	if result := repository.cache.SetWithTTL(apiKey, authUser, 1, 2*time.Hour); !result {
		ctxLogger.Error(repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorf("cannot cache [%T] with ID [%s] and result [%t]", authUser, user.ID, result)))
	}

	return authUser, nil
}

func (repository *gormUserRepository) Load(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	user := new(entities.User)
	err := repository.db.WithContext(ctx).First(user, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.PropagateWithCodef(err, ErrCodeNotFound, "user with ID [%s] does not exist", userID))
	}

	if err != nil {
		return nil, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot load user with ID [%s]", userID))
	}

	return user, nil
}

func (repository *gormUserRepository) LoadOrStore(ctx context.Context, authUser entities.AuthContext, newUserStatus entities.UserStatus) (*entities.User, bool, error) {
	ctx, span := repository.tracer.Start(ctx)
	defer span.End()

	if !newUserStatus.IsValid() {
		return nil, false, repository.tracer.WrapErrorSpan(span, stacktrace.NewErrorf("cannot create user [%s] with invalid status [%s]", authUser.ID, newUserStatus))
	}

	user, err := repository.Load(ctx, authUser.ID)
	if err == nil {
		return user, false, nil
	}

	apiKey, err := repository.generateAPIKey(64)
	if err != nil {
		return nil, false, stacktrace.Propagatef(err, "cannot generate apiKey for user [%s]", authUser.ID)
	}

	user = &entities.User{
		ID:               authUser.ID,
		Email:            authUser.Email,
		APIKey:           "uk_" + apiKey,
		SubscriptionName: entities.SubscriptionNameFree,
		Status:           newUserStatus,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	isNew := false
	err = crdbgorm.ExecuteTx(ctx, repository.db, nil, func(tx *gorm.DB) error {
		result := tx.WithContext(ctx).Where(entities.User{ID: user.ID}).FirstOrCreate(user)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected > 0 {
			isNew = true
		}
		return result.Error
	})
	if err != nil {
		return user, isNew, repository.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot create user from auth user [%+#v]", authUser))
	}

	if isNew {
		repository.statusCache.set(user.ID, user.Status)
	}

	return user, isNew, nil
}

// generateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func (repository *gormUserRepository) generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	// Note that err == nil only if we read len(b) bytes.
	if _, err := rand.Read(b); err != nil {
		return nil, stacktrace.Propagatef(err, "cannot generate [%d] random bytes", n)
	}

	return b, nil
}

// generateAPIKey returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func (repository *gormUserRepository) generateAPIKey(n int) (string, error) {
	b, err := repository.generateRandomBytes(n)
	return base64.URLEncoding.EncodeToString(b)[0:n], stacktrace.Propagatef(err, "cannot generate random bytes")
}
