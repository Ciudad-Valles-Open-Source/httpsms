package services

import (
	"context"
	"fmt"
	"time"

	"github.com/NdoleStudio/httpsms/pkg/cache"
	"github.com/NdoleStudio/httpsms/pkg/emails"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/google/uuid"
)

// BillingService is responsible for tracking usages and billing users
type BillingService struct {
	service
	logger                 telemetry.Logger
	tracer                 telemetry.Tracer
	cache                  cache.Cache
	emailFactory           emails.UserEmailFactory
	mailer                 emails.Mailer
	userRepository         repositories.UserRepository
	billingUsageRepository repositories.BillingUsageRepository
}

// NewBillingService creates a new BillingService
func NewBillingService(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	cache cache.Cache,
	mailer emails.Mailer,
	emailFactory emails.UserEmailFactory,
	usageRepository repositories.BillingUsageRepository,
	userRepository repositories.UserRepository,
) (s *BillingService) {
	return &BillingService{
		logger:                 logger.WithService(fmt.Sprintf("%T", s)),
		tracer:                 tracer,
		cache:                  cache,
		emailFactory:           emailFactory,
		mailer:                 mailer,
		userRepository:         userRepository,
		billingUsageRepository: usageRepository,
	}
}

// IsEntitledWithCount always allows sending and receiving messages.
// This fork is 100% free and self-hosted: there are no plan or usage limits.
func (service *BillingService) IsEntitledWithCount(_ context.Context, _ entities.UserID, _ uint) *string {
	return nil
}

// IsEntitled checks if a user can send or receive and SMS message
func (service *BillingService) IsEntitled(ctx context.Context, userID entities.UserID) *string {
	return service.IsEntitledWithCount(ctx, userID, 1)
}

// GetCurrentUsage gets the current billing usage for a user
func (service *BillingService) GetCurrentUsage(ctx context.Context, userID entities.UserID) (*entities.BillingUsage, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	return service.billingUsageRepository.GetCurrent(ctx, userID)
}

// GetUsageHistory gets the billing usage history for a user
func (service *BillingService) GetUsageHistory(ctx context.Context, userID entities.UserID, params repositories.IndexParams) (*[]entities.BillingUsage, error) {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	return service.billingUsageRepository.GetHistory(ctx, userID, params)
}

// RegisterSentMessage records the billing usage for a sent message
func (service *BillingService) RegisterSentMessage(ctx context.Context, messageID uuid.UUID, timestamp time.Time, userID entities.UserID) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	if err := service.billingUsageRepository.RegisterSentMessage(ctx, timestamp, userID); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not register [sent] message with ID [%s] for user with ID [%s]", messageID, userID))
	}

	ctxLogger.Info(fmt.Sprintf("registered [sent] message with ID [%s] for user [%s]", messageID, userID))
	return nil
}

// RegisterReceivedMessage records the billing usage for a received message
func (service *BillingService) RegisterReceivedMessage(ctx context.Context, messageID uuid.UUID, timestamp time.Time, userID entities.UserID) error {
	ctx, span := service.tracer.Start(ctx)
	defer span.End()

	ctxLogger := service.tracer.CtxLogger(service.logger, span)

	if err := service.billingUsageRepository.RegisterReceivedMessage(ctx, timestamp, userID); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not register [received] message with ID [%s] for user with ID [%s]", messageID, userID))
	}

	ctxLogger.Info(fmt.Sprintf("registered [received] message with ID [%s] for user [%s]", messageID, userID))
	return nil
}

// DeleteAllForUser deletes all entities.BillingUsage for an entities.UserID.
func (service *BillingService) DeleteAllForUser(ctx context.Context, userID entities.UserID) error {
	ctx, span, ctxLogger := service.tracer.StartWithLogger(ctx, service.logger)
	defer span.End()

	if err := service.billingUsageRepository.DeleteAllForUser(ctx, userID); err != nil {
		return service.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "could not delete [entities.BillingUsage] for user with ID [%s]", userID))
	}

	ctxLogger.Info(fmt.Sprintf("deleted all [entities.BillingUsage] for user with ID [%s]", userID))
	return nil
}
