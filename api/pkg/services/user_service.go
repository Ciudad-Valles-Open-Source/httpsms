package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"firebase.google.com/go/auth"
	"github.com/NdoleStudio/httpsms/pkg/config"
	"github.com/NdoleStudio/httpsms/pkg/emails"
	"github.com/NdoleStudio/httpsms/pkg/events"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/stacktrace"
	"github.com/google/uuid"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
)

// UserService is handles user requests
type UserService struct {
	service
	logger       telemetry.Logger
	tracer       telemetry.Tracer
	emailFactory emails.UserEmailFactory
	mailer       emails.Mailer
	repository   repositories.UserRepository
	dispatcher   *EventDispatcher
	authClient   *auth.Client
	adminPolicy  *config.AdminPolicy
}

// NewUserService creates a new UserService
func NewUserService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	repository repositories.UserRepository,
	mailer emails.Mailer,
	emailFactory emails.UserEmailFactory,
	dispatcher *EventDispatcher,
	authClient *auth.Client,
	adminPolicy *config.AdminPolicy,
) (s *UserService) {
	return &UserService{
		adminPolicy:  adminPolicy,
		logger:       logger.WithService(fmt.Sprintf("%T", s)),
		tracer:       tracer,
		mailer:       mailer,
		emailFactory: emailFactory,
		repository:   repository,
		dispatcher:   dispatcher,
		authClient:   authClient,
	}
}

// Get fetches or creates an entities.User
func (service *UserService) Get(ctx context.Context, source string, authUser entities.AuthContext) (*entities.User, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, _, err := service.loadOrStore(ctx, source, authUser)
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not get [%T] with from [%+#v]", user, authUser))
	}

	// An administrator whose account was created before their email was verified
	// (or before they were added to ADMIN_EMAILS) is activated on their next login.
	if user.Status != entities.UserStatusActive && service.adminPolicy.IsAdmin(authUser) {
		now := time.Now().UTC()
		approvedBy := strings.ToLower(strings.TrimSpace(authUser.Email))
		activated, err := service.repository.UpdateStatus(ctx, user.ID, repositories.UserStatusUpdate{
			Status:     entities.UserStatusActive,
			ApprovedAt: &now,
			ApprovedBy: &approvedBy,
		})
		if err != nil {
			return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot activate the account of administrator [%s]", user.ID))
		}
		ctxLogger.Info(fmt.Sprintf("activated the account of administrator [%s] with user ID [%s]", authUser.Email, user.ID))
		user = activated
	}

	return user, nil
}

// IsAdmin checks if an authenticated request is made by an administrator
func (service *UserService) IsAdmin(authUser entities.AuthContext) bool {
	return service.adminPolicy.IsAdmin(authUser)
}

// loadOrStore fetches the user or atomically creates it with the status decided by the config.AdminPolicy
func (service *UserService) loadOrStore(ctx context.Context, source string, authUser entities.AuthContext) (*entities.User, bool, error) {
	user, isNew, err := service.repository.LoadOrStore(ctx, authUser, service.adminPolicy.InitialStatus(authUser))
	if err != nil {
		return nil, false, err
	}

	if isNew {
		service.dispatchUserCreatedEvent(ctx, source, user)
	}

	return user, isNew, nil
}

func (service *UserService) dispatchUserCreatedEvent(ctx context.Context, source string, user *entities.User) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	event, err := service.createEvent(events.UserAccountCreated, source, &events.UserAccountCreatedPayload{
		UserID:    user.ID,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot create event [%s] for user [%s]", events.UserAccountCreated, user.ID))
		return
	}

	if err = service.dispatcher.Dispatch(ctx, event); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", event.Type(), user.ID))
		return
	}
}

// GetByID fetches an entities.User
func (service *UserService) GetByID(ctx context.Context, userID entities.UserID) (*entities.User, error) {
	ctx, span, _ := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.repository.Load(ctx, userID)
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not get [%T] with ID [%s]", user, userID))
	}

	return user, nil
}

// UserUpdateParams are parameters for updating an entities.User
type UserUpdateParams struct {
	Timezone      *time.Location
	ActivePhoneID *uuid.UUID
}

// Update an entities.User
func (service *UserService) Update(ctx context.Context, source string, authUser entities.AuthContext, params UserUpdateParams) (*entities.User, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	user, _, err := service.loadOrStore(ctx, source, authUser)
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not get [%T] with from [%+#v]", user, authUser))
	}

	user.Timezone = params.Timezone.String()
	user.ActivePhoneID = params.ActivePhoneID

	if err = service.repository.Update(ctx, user); err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot save user with id [%s]", user.ID))
	}

	ctxLogger.Info(fmt.Sprintf("user saved with id [%s] in the userRepository", user.ID))
	return user, nil
}

// UserNotificationUpdateParams are parameters for updating the notifications of a user
type UserNotificationUpdateParams struct {
	MessageStatusEnabled bool
	WebhookEnabled       bool
	HeartbeatEnabled     bool
	NewsletterEnabled    bool
}

// UpdateNotificationSettings for an entities.User
func (service *UserService) UpdateNotificationSettings(ctx context.Context, userID entities.UserID, params *UserNotificationUpdateParams) (*entities.User, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.repository.Load(ctx, userID)
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not load [%T] with ID [%s]", user, userID))
	}

	user.NotificationWebhookEnabled = params.WebhookEnabled
	user.NotificationHeartbeatEnabled = params.HeartbeatEnabled
	user.NotificationMessageStatusEnabled = params.MessageStatusEnabled
	user.NotificationNewsletterEnabled = params.NewsletterEnabled

	if err = service.repository.Update(ctx, user); err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot save user with id [%s] in [%T]", user.ID, service.repository))
	}

	ctxLogger.Info(fmt.Sprintf("updated notification settings for [%T] with ID [%s] in the [%T]", user, user.ID, service.repository))
	return user, nil
}

// RotateAPIKey for an entities.User
func (service *UserService) RotateAPIKey(ctx context.Context, source string, userID entities.UserID) (*entities.User, error) {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.repository.RotateAPIKey(ctx, userID)
	if err != nil {
		return nil, service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not rotate API key for [%T] with ID [%s]", user, userID))
	}

	ctxLogger.Info(fmt.Sprintf("rotated the api key for [%T] with ID [%s] in the [%T]", user, user.ID, service.repository))

	event, err := service.createEvent(events.UserAPIKeyRotated, source, &events.UserAPIKeyRotatedPayload{
		UserID:    user.ID,
		Email:     user.Email,
		Timestamp: time.Now().UTC(),
		Timezone:  user.Timezone,
	})
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot create event [%s] for user [%s]", events.UserAPIKeyRotated, user.ID))
		return user, nil
	}

	if err = service.dispatcher.Dispatch(ctx, event); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", event.Type(), user.ID))
		return user, nil
	}

	return user, nil
}

// Delete an entities.User
func (service *UserService) Delete(ctx context.Context, source string, userID entities.UserID) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	user, err := service.repository.Load(ctx, userID)
	if err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot load user with ID [%s] from the [%T]", userID, service.repository))
	}

	if err = service.repository.Delete(ctx, user); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not delete user with ID [%s] from the [%T]", userID, service.repository))
	}

	ctxLogger.Info(fmt.Sprintf("sucessfully deleted user with ID [%s] in the [%T]", userID, service.repository))

	event, err := service.createEvent(events.UserAccountDeleted, source, &events.UserAccountDeletedPayload{
		UserID:    userID,
		UserEmail: user.Email,
		Timestamp: time.Now().UTC(),
	})
	if err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot create event [%s] for user [%s]", events.UserAccountDeleted, userID))
	}

	if err = service.dispatcher.Dispatch(ctx, event); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot dispatch [%s] event for user [%s]", event.Type(), userID))
	}

	return nil
}

// SendAPIKeyRotatedEmail sends an email to an entities.User when the API key is rotated
func (service *UserService) SendAPIKeyRotatedEmail(ctx context.Context, payload *events.UserAPIKeyRotatedPayload) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	email, err := service.emailFactory.APIKeyRotated(payload.Email, payload.Timestamp, payload.Timezone)
	if err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot create api key rotated email for user [%s]", payload.UserID))
	}

	if err = service.mailer.Send(ctx, email); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "canot create api key rotated email to user [%s]", payload.UserID))
	}

	ctxLogger.Info(fmt.Sprintf("api key rotated email sent successfully to [%s] with user ID  [%s]", payload.Email, payload.UserID))
	return nil
}

// UserSendPhoneDeadEmailParams are parameters for notifying a user when a phone is dead
type UserSendPhoneDeadEmailParams struct {
	UserID                 entities.UserID
	PhoneID                uuid.UUID
	Owner                  string
	LastHeartbeatTimestamp time.Time
}

// SendPhoneDeadEmail sends an email to an entities.User when a phone is dead
func (service *UserService) SendPhoneDeadEmail(ctx context.Context, params *UserSendPhoneDeadEmailParams) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	user, err := service.repository.Load(ctx, params.UserID)
	if err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not get [%T] with ID [%s]", user, params.UserID))
	}

	if !user.NotificationHeartbeatEnabled {
		ctxLogger.Info(fmt.Sprintf("[%s] email notifications disabled for user [%s] with owner [%s]", events.EventTypePhoneHeartbeatOffline, params.UserID, params.Owner))
		return nil
	}

	email, err := service.emailFactory.PhoneDead(user, params.LastHeartbeatTimestamp, params.Owner)
	if err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot create phone dead email for user [%s]", params.UserID))
	}

	if err = service.mailer.Send(ctx, email); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "canot send phone dead notification to user [%s]", params.UserID))
	}

	ctxLogger.Info(fmt.Sprintf("phone dead notification sent successfully to [%s] about [%s]", user.Email, params.Owner))
	return nil
}

// DeleteAuthUser deletes an entities.AuthContext from firebase
func (service *UserService) DeleteAuthUser(ctx context.Context, userID entities.UserID) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if err := service.authClient.DeleteUser(ctx, userID.String()); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not delete [entities.AuthContext] from firebase with ID [%s]", userID))
	}

	ctxLogger.Info(fmt.Sprintf("deleted [entities.AuthContext] from firebase for user with ID [%s]", userID))
	return nil
}
