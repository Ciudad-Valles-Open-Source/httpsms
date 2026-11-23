package listeners

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"

	"github.com/NdoleStudio/httpsms/pkg/events"
	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// UserListener handles cloud events which sends notifications
type UserListener struct {
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.UserService
}

// NewUserListener creates a new instance of UserListener
func NewUserListener(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.UserService,
) (l *UserListener, routes map[string]events.EventListener) {
	l = &UserListener{
		logger:  logger.WithService(fmt.Sprintf("%T", l)),
		tracer:  tracer,
		service: service,
	}

	return l, map[string]events.EventListener{
		events.EventTypePhoneHeartbeatOffline: l.onPhoneHeartbeatDead,
		events.UserAPIKeyRotated:              l.onUserAPIKeyRotated,
		events.UserAccountDeleted:             l.onUserAccountDeleted,
	}
}

// onPhoneHeartbeatDead handles the events.EventTypePhoneHeartbeatOffline event
func (listener *UserListener) onPhoneHeartbeatDead(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.PhoneHeartbeatOfflinePayload
	if err := event.DataAs(&payload); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot decode [%s] into [%T]", event.Data(), payload))
	}

	sendParams := &services.UserSendPhoneDeadEmailParams{
		UserID:                 payload.UserID,
		PhoneID:                payload.PhoneID,
		Owner:                  payload.Owner,
		LastHeartbeatTimestamp: payload.LastHeartbeatTimestamp,
	}

	if err := listener.service.SendPhoneDeadEmail(ctx, sendParams); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot send notification with params [%s] for event with ID [%s]", spew.Sdump(sendParams), event.ID()))
	}

	return nil
}

// onAPIKeyRotated handles the events.UserAPIKeyRotated event
func (listener *UserListener) onUserAPIKeyRotated(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	payload := new(events.UserAPIKeyRotatedPayload)
	if err := event.DataAs(&payload); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot decode [%s] into [%T]", event.Data(), payload))
	}

	if err := listener.service.SendAPIKeyRotatedEmail(ctx, payload); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot send notification with params [%s] for event with ID [%s]", spew.Sdump(payload), event.ID()))
	}

	return nil
}

func (listener *UserListener) onUserAccountDeleted(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.UserAccountDeletedPayload
	if err := event.DataAs(&payload); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot decode [%s] into [%T]", event.Data(), payload))
	}

	if err := listener.service.DeleteAuthUser(ctx, payload.UserID); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot delete [entities.AuthUser] for user [%s] on [%s] event with ID [%s]", payload.UserID, event.Type(), event.ID()))
	}

	return nil
}
