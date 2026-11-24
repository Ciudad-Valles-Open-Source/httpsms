package middlewares

import (
	"context"
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/gofiber/fiber/v3"
)

const (
	// AccountStatusCodePending is returned when the account waits for an administrator's approval
	AccountStatusCodePending = "account_pending"

	// AccountStatusCodeSuspended is returned when an administrator suspended the account
	AccountStatusCodeSuspended = "account_suspended"

	accountPendingMessage   = "Your account is pending approval by an administrator."
	accountSuspendedMessage = "Your account has been suspended by an administrator."

	userMeRoutePath = "/v1/users/me"
)

// AccountGate decides which authenticated accounts may use the API when the
// administrator approval gate (ADMIN_EMAILS) is enabled.
type AccountGate interface {
	// Enabled is true when the approval gate is turned on
	Enabled() bool

	// IsAdmin is true only for a verified Firebase ID token whose email is in ADMIN_EMAILS
	IsAdmin(authUser entities.AuthContext) bool

	// IsSystemUser is true for the internal events queue user
	IsSystemUser(userID entities.UserID) bool

	// LoadStatus returns the current entities.UserStatus of a user
	LoadStatus(ctx context.Context, userID entities.UserID) (entities.UserStatus, error)
}

// AuthenticatedAndApproved checks that the request is authenticated and, when
// the approval gate is enabled, that the account is active. It is the
// middleware used by every authenticated route, so no route can skip the gate.
//
// Requests are always allowed for administrators, for the system user and for
// GET/DELETE /v1/users/me (so a pending user can see their status or delete
// their account). Everybody else needs an active account.
func AuthenticatedAndApproved(logger telemetry.Logger, tracer telemetry.Tracer, gate AccountGate) fiber.Handler {
	logger = logger.WithService("middlewares.AuthenticatedAndApproved")

	return func(c fiber.Ctx) error {
		ctx, span, ctxLogger := tracer.StartFromFiberCtxWithLogger(c, logger, "middlewares.AuthenticatedAndApproved")
		defer span.End()

		authUser, ok := authContextFromLocals(c)
		if !ok {
			return responseUnauthenticated(c)
		}

		if gate == nil || !gate.Enabled() || gate.IsAdmin(authUser) || gate.IsSystemUser(authUser.ID) || isAccountStatusExemptRoute(c) {
			return c.Next()
		}

		status, err := gate.LoadStatus(ctx, authUser.ID)
		if err != nil && stacktrace.GetCode(err) != repositories.ErrCodeNotFound {
			ctxLogger.Error(tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot load account status for user [%s]", authUser.ID)))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "We ran into an internal error while handling the request.",
			})
		}

		// A user without a row yet (first request is not GET /v1/users/me) is a new registration.
		if err != nil {
			status = entities.UserStatusPending
		}

		if status == entities.UserStatusActive {
			return c.Next()
		}

		span.AddEvent(fmt.Sprintf("account [%s] with status [%s] is not allowed to call [%s %s]", authUser.ID, status, c.Method(), c.Path()))
		return ResponseAccountNotActive(c, status)
	}
}

// RequireAdmin allows only administrators (see AccountGate.IsAdmin). It must run after authentication.
func RequireAdmin(tracer telemetry.Tracer, gate AccountGate) fiber.Handler {
	return func(c fiber.Ctx) error {
		_, span := tracer.StartFromFiberCtx(c, "middlewares.RequireAdmin")
		defer span.End()

		authUser, ok := authContextFromLocals(c)
		if !ok {
			return responseUnauthenticated(c)
		}

		if gate == nil || !gate.IsAdmin(authUser) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "forbidden",
			})
		}

		return c.Next()
	}
}

// ResponseAccountNotActive writes the 403 response for an account which is not active
func ResponseAccountNotActive(c fiber.Ctx, status entities.UserStatus) error {
	code, message := AccountStatusCodePending, accountPendingMessage
	if status == entities.UserStatusSuspended {
		code, message = AccountStatusCodeSuspended, accountSuspendedMessage
	}

	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
		"status":  "error",
		"message": message,
		"data":    fiber.Map{"code": code},
	})
}

func isAccountStatusExemptRoute(c fiber.Ctx) bool {
	method := c.Method()
	if method != fiber.MethodGet && method != fiber.MethodDelete {
		return false
	}

	// Route().Path is the registered path, so this can't be fooled with a different casing or a trailing slash.
	route := c.Route()
	return route != nil && route.Path == userMeRoutePath
}

func authContextFromLocals(c fiber.Ctx) (entities.AuthContext, bool) {
	authUser, ok := c.Locals(ContextKeyAuthUserID).(entities.AuthContext)
	if !ok || authUser.IsNoop() {
		return entities.AuthContext{}, false
	}
	return authUser, true
}

func responseUnauthenticated(c fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status":  "error",
		"message": "You are not authorized to carry out this request.",
		"data":    "Make sure your API key is set in the [x-api-key] header in the request",
	})
}
