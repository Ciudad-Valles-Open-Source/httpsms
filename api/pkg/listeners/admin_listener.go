package listeners

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/events"
	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// AdminListener sends the account approval emails
type AdminListener struct {
	logger  telemetry.Logger
	tracer  telemetry.Tracer
	service *services.AdminService
}

// NewAdminListener creates a new instance of AdminListener
func NewAdminListener(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	service *services.AdminService,
) (l *AdminListener, routes map[string]events.EventListener) {
	l = &AdminListener{
		logger:  logger.WithService(fmt.Sprintf("%T", l)),
		tracer:  tracer,
		service: service,
	}

	return l, map[string]events.EventListener{
		events.UserAccountCreated:  l.onUserAccountCreated,
		events.UserAccountApproved: l.onUserAccountApproved,
	}
}

// onUserAccountCreated emails the administrators when a new account is pending approval
func (listener *AdminListener) onUserAccountCreated(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	var payload events.UserAccountCreatedPayload
	if err := event.DataAs(&payload); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot decode [%s] into [%T]", event.Data(), payload))
	}

	listener.service.NotifyAdminsOfPendingUser(ctx, payload.UserID)
	return nil
}

// onUserAccountApproved emails the user when an administrator approves the account
func (listener *AdminListener) onUserAccountApproved(ctx context.Context, event cloudevents.Event) error {
	ctx, span := listener.tracer.Start(ctx)
	defer span.End()

	payload := new(events.UserAccountApprovedPayload)
	if err := event.DataAs(payload); err != nil {
		return listener.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot decode [%s] into [%T]", event.Data(), payload))
	}

	listener.service.SendAccountApprovedEmail(ctx, payload)
	return nil
}
