package middlewares

import (
	"context"
	"fmt"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/gofiber/fiber/v3"
)

// BearerAuth authenticates a user based on the bearer token
func BearerAuth(logger telemetry.Logger, tracer telemetry.Tracer, authClient *auth.Client) fiber.Handler {
	logger = logger.WithService("middlewares.BearerAuth")
	return func(c fiber.Ctx) error {
		_, span := tracer.StartFromFiberCtx(c, "middlewares.BearerAuth")
		defer span.End()

		authToken := c.Get(authHeaderBearer)
		if !strings.HasPrefix(authToken, bearerScheme) {
			span.AddEvent(fmt.Sprintf("The request header has no [%s] token", bearerScheme))
			return c.Next()
		}

		if len(authToken) > len(bearerScheme)+1 {
			authToken = authToken[len(bearerScheme)+1:]
		}

		ctxLogger := tracer.CtxLogger(logger, span)

		token, err := authClient.VerifyIDToken(context.Background(), authToken)
		if err != nil {
			// Never log the token itself: API keys are also sent as "Bearer <key>".
			ctxLogger.Debug(stacktrace.Propagatef(err, "bearer token is not a valid firebase id token").Error())
			return c.Next()
		}

		span.AddEvent(fmt.Sprintf("[%s] token is valid", bearerScheme))

		authUser := authContextFromToken(token)
		if authUser.IsNoop() {
			ctxLogger.Warn(tracer.WrapErrorSpan(span, stacktrace.NewErrorf("firebase id token for uid [%s] has no email or user_id claim", token.UID)))
			return c.Next()
		}

		c.Locals(ContextKeyAuthUserID, authUser)
		return c.Next()
	}
}

// authContextFromToken builds an entities.AuthContext from a verified Firebase ID token.
// Missing claims or claims with an unexpected type never panic: they become zero values.
func authContextFromToken(token *auth.Token) entities.AuthContext {
	if token == nil {
		return entities.AuthContext{}
	}

	userID := stringClaim(token.Claims, "user_id")
	if userID == "" {
		userID = token.UID
	}

	return entities.AuthContext{
		ID:              entities.UserID(userID),
		Email:           stringClaim(token.Claims, "email"),
		EmailVerified:   boolClaim(token.Claims, "email_verified"),
		FirebaseIDToken: true,
	}
}

func stringClaim(claims map[string]interface{}, name string) string {
	value, ok := claims[name].(string)
	if !ok {
		return ""
	}
	return value
}

func boolClaim(claims map[string]interface{}, name string) bool {
	value, ok := claims[name].(bool)
	return ok && value
}
