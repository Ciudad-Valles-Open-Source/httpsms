package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/config"
	"github.com/NdoleStudio/httpsms/pkg/emails"
	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/events"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
)

const (
	// ErrCodeProtectedUser is returned when an administrator tries to change
	// their own account, another administrator or the system user.
	ErrCodeProtectedUser = stacktrace.ErrorCode(2000)

	// ErrCodeNotAdmin is returned when the caller is not an administrator
	ErrCodeNotAdmin = stacktrace.ErrorCode(2001)

	// AdminUsersDefaultLimit is the default page size when listing users
	AdminUsersDefaultLimit = 20

	// AdminUsersMaxLimit is the maximum page size when listing users
	AdminUsersMaxLimit = 100
)

// AdminService implements the administrator approval gate and the admin user management.
type AdminService struct {
	service
	logger       telemetry.Logger
	tracer       telemetry.Tracer
	policy       *config.AdminPolicy
	repository   repositories.UserRepository
	userService  *UserService
	mailer       emails.Mailer
	emailFactory emails.UserEmailFactory
	dispatcher   *EventDispatcher
	smtpEnabled  bool
	systemUserID func() entities.UserID
}

// NewAdminService creates a new AdminService
func NewAdminService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	policy *config.AdminPolicy,
	repository repositories.UserRepository,
	userService *UserService,
	mailer emails.Mailer,
	emailFactory emails.UserEmailFactory,
	dispatcher *EventDispatcher,
	smtpEnabled bool,
	systemUserID func() entities.UserID,
) (s *AdminService) {
	return &AdminService{
		logger:       logger.WithService(fmt.Sprintf("%T", s)),
		tracer:       tracer,
		policy:       policy,
		repository:   repository,
		userService:  userService,
		mailer:       mailer,
		emailFactory: emailFactory,
		dispatcher:   dispatcher,
		smtpEnabled:  smtpEnabled,
		systemUserID: systemUserID,
	}
}

// Enabled is true when ADMIN_EMAILS is set and new accounts need approval
func (service *AdminService) Enabled() bool {
	return service.policy.Enabled()
}

// IsAdmin is true only for a verified Firebase ID token whose email is in ADMIN_EMAILS
func (service *AdminService) IsAdmin(authUser entities.AuthContext) bool {
	return service.policy.IsAdmin(authUser)
}

// IsSystemUser is true for the internal events queue user (EVENTS_QUEUE_USER_ID)
func (service *AdminService) IsSystemUser(userID entities.UserID) bool {
	if service.systemUserID == nil {
		return false
	}
	systemUserID := service.systemUserID()
	return systemUserID != "" && userID == systemUserID
}

// LoadStatus returns the current entities.UserStatus of a user
func (service *AdminService) LoadStatus(ctx context.Context, userID entities.UserID) (entities.UserStatus, error) {
	return service.repository.LoadStatus(ctx, userID)
}

// AdminUserIndexParams are the parameters for listing users
type AdminUserIndexParams struct {
	Status *entities.UserStatus
	Query  string
	Skip   int
	Limit  int
}

// Index lists users for an administrator, newest first
func (service *AdminService) Index(ctx context.Context, admin entities.AuthContext, params AdminUserIndexParams) ([]*entities.AdminUser, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if !service.IsAdmin(admin) {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.NewErrorWithCodef(ErrCodeNotAdmin, "user [%s] is not an administrator", admin.ID))
	}

	rows, err := service.repository.AdminIndex(ctx, repositories.UserAdminIndexParams{
		Status: params.Status,
		Query:  strings.TrimSpace(params.Query),
		Skip:   max(params.Skip, 0),
		Limit:  ClampAdminUsersLimit(params.Limit),
	})
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot list users with params [%+#v]", params))
	}

	result := make([]*entities.AdminUser, 0, len(rows))
	for _, row := range rows {
		result = append(result, service.toAdminUser(row.User, row.PhonesCount))
	}

	ctxLogger.Info(fmt.Sprintf("admin [%s] action [list_users] listed [%d] users with params [%+v]", admin.Email, len(result), params))
	return result, nil
}

// ClampAdminUsersLimit applies the default (20) and maximum (100) page size
func ClampAdminUsersLimit(limit int) int {
	if limit <= 0 {
		return AdminUsersDefaultLimit
	}
	return min(limit, AdminUsersMaxLimit)
}

// UpdateStatus changes the status of a user. Setting entities.UserStatusActive stamps approved_at/approved_by.
func (service *AdminService) UpdateStatus(ctx context.Context, source string, admin entities.AuthContext, userID entities.UserID, status entities.UserStatus) (*entities.AdminUser, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if status != entities.UserStatusActive && status != entities.UserStatusSuspended {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.NewErrorf("cannot set the status of user [%s] to [%s]", userID, status))
	}

	target, err := service.loadTarget(ctx, admin, userID)
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "admin [%s] cannot change the status of user [%s]", admin.Email, userID))
	}

	update := repositories.UserStatusUpdate{Status: status}
	if status == entities.UserStatusActive {
		now := time.Now().UTC()
		approvedBy := strings.ToLower(strings.TrimSpace(admin.Email))
		update.ApprovedAt = &now
		update.ApprovedBy = &approvedBy
	}

	user, err := service.repository.UpdateStatus(ctx, userID, update)
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot update the status of user [%s] to [%s]", userID, status))
	}

	ctxLogger.Info(fmt.Sprintf("admin [%s] action [set_status:%s] target user [%s] previous status [%s]", admin.Email, status, userID, target.Status))

	if status == entities.UserStatusActive && target.Status != entities.UserStatusActive {
		service.dispatchAccountApproved(ctx, source, user, *update.ApprovedBy)
	}

	phonesCount, err := service.repository.CountPhones(ctx, userID)
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot count the phones of user [%s]", userID))
	}

	return service.toAdminUser(user, phonesCount), nil
}

// Delete deletes a user and all their data using the same flow as a self deletion
func (service *AdminService) Delete(ctx context.Context, source string, admin entities.AuthContext, userID entities.UserID) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if _, err := service.loadTarget(ctx, admin, userID); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "admin [%s] cannot delete user [%s]", admin.Email, userID))
	}

	if err := service.userService.Delete(ctx, source, userID); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "admin [%s] cannot delete user [%s]", admin.Email, userID))
	}

	ctxLogger.Info(fmt.Sprintf("admin [%s] action [delete_user] target user [%s]", admin.Email, userID))
	return nil
}

// loadTarget loads the user an administrator acts on and rejects protected accounts
func (service *AdminService) loadTarget(ctx context.Context, admin entities.AuthContext, userID entities.UserID) (*entities.User, error) {
	if !service.IsAdmin(admin) {
		return nil, stacktrace.NewErrorWithCodef(ErrCodeNotAdmin, "user [%s] is not an administrator", admin.ID)
	}

	if userID == admin.ID {
		return nil, stacktrace.NewErrorWithCode(ErrCodeProtectedUser, "administrators cannot change their own account")
	}

	user, err := service.repository.Load(ctx, userID)
	if err != nil {
		return nil, stacktrace.Propagatef(err, "cannot load user [%s]", userID)
	}

	switch {
	case strings.EqualFold(strings.TrimSpace(user.Email), strings.TrimSpace(admin.Email)):
		return nil, stacktrace.NewErrorWithCode(ErrCodeProtectedUser, "administrators cannot change their own account")
	case service.policy.IsAdminEmail(user.Email):
		return nil, stacktrace.NewErrorWithCodef(ErrCodeProtectedUser, "user [%s] is an administrator listed in ADMIN_EMAILS", userID)
	case service.IsSystemUser(user.ID):
		return nil, stacktrace.NewErrorWithCodef(ErrCodeProtectedUser, "user [%s] is the system user", userID)
	}

	return user, nil
}

func (service *AdminService) toAdminUser(user *entities.User, phonesCount int64) *entities.AdminUser {
	return &entities.AdminUser{
		ID:          user.ID,
		Email:       user.Email,
		Status:      user.Status,
		IsAdmin:     service.policy.IsAdminEmail(user.Email),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		ApprovedAt:  user.ApprovedAt,
		ApprovedBy:  user.ApprovedBy,
		PhonesCount: phonesCount,
	}
}

func (service *AdminService) dispatchAccountApproved(ctx context.Context, source string, user *entities.User, approvedBy string) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	// The event only sends an email, so there is nothing to do without SMTP.
	if !service.smtpEnabled || service.dispatcher == nil {
		return
	}

	event, err := service.createEvent(events.UserAccountApproved, source, &events.UserAccountApprovedPayload{
		UserID:     user.ID,
		UserEmail:  user.Email,
		ApprovedBy: approvedBy,
		Timestamp:  time.Now().UTC(),
	})
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot create event [%s] for user [%s]", events.UserAccountApproved, user.ID))
		return
	}

	if err = service.dispatcher.Dispatch(ctx, event); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", event.Type(), user.ID))
	}
}

// NotifyAdminsOfPendingUser emails every administrator when a new account waits for approval.
// It only sends emails when SMTP is configured and never returns an error, so a failing mail
// server can't affect the registration.
func (service *AdminService) NotifyAdminsOfPendingUser(ctx context.Context, userID entities.UserID) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if !service.smtpEnabled || !service.policy.Enabled() {
		return
	}

	user, err := service.repository.Load(ctx, userID)
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot load user [%s] to notify administrators", userID))
		return
	}

	if user.Status != entities.UserStatusPending {
		return
	}

	for _, adminEmail := range service.policy.Emails() {
		email, err := service.emailFactory.AccountPendingApproval(adminEmail, user)
		if err != nil {
			ctxLogger.Error(stacktrace.Propagatef(err, "cannot create pending account email for user [%s]", userID))
			return
		}

		if err = service.mailer.Send(ctx, email); err != nil {
			ctxLogger.Error(stacktrace.Propagatef(err, "cannot send pending account email for user [%s] to admin [%s]", userID, adminEmail))
			continue
		}

		ctxLogger.Info(fmt.Sprintf("pending account email for user [%s] sent to admin [%s]", userID, adminEmail))
	}
}

// SendAccountApprovedEmail emails a user when their account is approved. Errors are only logged.
func (service *AdminService) SendAccountApprovedEmail(ctx context.Context, payload *events.UserAccountApprovedPayload) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if !service.smtpEnabled {
		return
	}

	email, err := service.emailFactory.AccountApproved(&entities.User{ID: payload.UserID, Email: payload.UserEmail})
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot create account approved email for user [%s]", payload.UserID))
		return
	}

	if err = service.mailer.Send(ctx, email); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot send account approved email to user [%s]", payload.UserID))
		return
	}

	ctxLogger.Info(fmt.Sprintf("account approved email sent to user [%s]", payload.UserID))
}
