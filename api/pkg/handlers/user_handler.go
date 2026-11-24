package handlers

import (
	"fmt"

	"github.com/NdoleStudio/httpsms/pkg/requests"
	"github.com/NdoleStudio/httpsms/pkg/responses"
	"github.com/NdoleStudio/httpsms/pkg/validators"
	"github.com/davecgh/go-spew/spew"

	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/stacktrace"
	"github.com/gofiber/fiber/v3"
)

// UserHandler handles user http requests.
type UserHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	validator *validators.UserHandlerValidator
	service   *services.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.UserHandlerValidator,
	service *services.UserService,
) (h *UserHandler) {
	return &UserHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", h)),
		tracer:    tracer,
		validator: validator,
		service:   service,
	}
}

// RegisterRoutes registers the routes for the MessageHandler
func (h *UserHandler) RegisterRoutes(router fiber.Router, middlewares ...fiber.Handler) {
	h.register(router, fiber.MethodGet, "/v1/users/me", middlewares, h.Show)
	h.register(router, fiber.MethodPut, "/v1/users/me", middlewares, h.Update)
	h.register(router, fiber.MethodDelete, "/v1/users/me", middlewares, h.Delete)
	h.register(router, fiber.MethodDelete, "/v1/users/:userID/api-keys", middlewares, h.DeleteAPIKey)
	h.register(router, fiber.MethodPut, "/v1/users/:userID/notifications", middlewares, h.UpdateNotifications)
}

// Show returns an entities.User
// @Summary      Get current user
// @Description  Get details of the currently authenticated user, including the account "status" (active, pending or suspended) and "is_admin". This endpoint works for pending and suspended accounts too.
// @Security	 ApiKeyAuth
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200 	{object}		responses.CurrentUserResponse
// @Failure      400	{object}		responses.BadRequest
// @Failure 	 401    {object}		responses.Unauthorized
// @Failure      422	{object}		responses.UnprocessableEntity
// @Failure      500	{object}		responses.InternalServerError
// @Router       /users/me [get]
func (h *UserHandler) Show(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	authUser := h.userFromContext(c)
	user, err := h.service.Get(ctx, c.OriginalURL(), authUser)
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot get user with ID [%s]", authUser.ID))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "user fetched successfully", responses.CurrentUser{
		User:    *user,
		IsAdmin: h.service.IsAdmin(authUser),
	})
}

// Update an entities.User
// @Summary      Update a user
// @Description  Updates the details of the currently authenticated user
// @Security	 ApiKeyAuth
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        payload   	body 		requests.UserUpdate  			true 	"Payload of user details to update"
// @Success      200 		{object}	responses.CurrentUserResponse
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /users/me [put]
func (h *UserHandler) Update(c fiber.Ctx) error {
	ctx, span := h.tracer.StartFromFiberCtx(c)
	defer span.End()

	ctxLogger := h.tracer.CtxLogger(h.logger, span)

	var request requests.UserUpdate
	if err := c.Bind().Body(&request); err != nil {
		ctxLogger.Warn(stacktrace.Propagatef(err, "cannot marshall params [%s] into %T", c.OriginalURL(), request))
		return h.responseBadRequest(c, err)
	}

	if errors := h.validator.ValidateUpdate(ctx, request.Sanitize()); len(errors) != 0 {
		ctxLogger.Warn(stacktrace.NewErrorf("validation errors [%s], while updating user [%+#v]", spew.Sdump(errors), request))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating user")
	}

	user, err := h.service.Update(ctx, c.OriginalURL(), h.userFromContext(c), request.ToUpdateParams())
	if err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot update user with params [%+#v]", request))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "user updated successfully", responses.CurrentUser{
		User:    *user,
		IsAdmin: h.service.IsAdmin(h.userFromContext(c)),
	})
}

// Delete an entities.User
// @Summary      Delete a user
// @Description  Deletes the currently authenticated user together with all their data.
// @Security	 ApiKeyAuth
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      201 		{object}	responses.NoContent
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      500		{object}	responses.InternalServerError
// @Router       /users/me [delete]
func (h *UserHandler) Delete(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	if err := h.service.Delete(ctx, c.OriginalURL(), h.userIDFomContext(c)); err != nil {
		ctxLogger.Error(stacktrace.Propagatef(err, "cannot delete user user with ID [%s]", h.userIDFomContext(c)))
		return h.responseInternalServerError(c)
	}

	return h.responseNoContent(c, "user deleted successfully")
}

// UpdateNotifications an entities.User
// @Summary      Update notification settings
// @Description  Update the email notification settings for a user
// @Security	 ApiKeyAuth
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param 		 userID 	path		string 							true 	"ID of the user to update" 				default(32343a19-da5e-4b1b-a767-3298a73703ca)
// @Param        payload   	body 		requests.UserNotificationUpdate	true 	"User notification details to update"
// @Success      200 		{object}	responses.CurrentUserResponse
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /users/{userID}/notifications [put]
func (h *UserHandler) UpdateNotifications(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.UserNotificationUpdate
	if err := c.Bind().Body(&request); err != nil {
		ctxLogger.Warn(stacktrace.Propagatef(err, "cannot marshall params [%s] into %T", c.OriginalURL(), request))
		return h.responseBadRequest(c, err)
	}

	user, err := h.service.UpdateNotificationSettings(ctx, h.userIDFomContext(c), request.ToUserNotificationUpdateParams())
	if err != nil {
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot update notification for [%T] with ID [%s]", user, h.userIDFomContext(c))))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "user notification settings updated successfully", responses.CurrentUser{
		User:    *user,
		IsAdmin: h.service.IsAdmin(h.userFromContext(c)),
	})
}

// DeleteAPIKey rotates the API Key for a user
// @Summary      Rotate the user's API Key
// @Description  Rotate the user's API key in case the current API Key is compromised
// @Security	 ApiKeyAuth
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param 		 userID 	path		string 							true 	"ID of the user to update" 	default(32343a19-da5e-4b1b-a767-3298a73703ca)
// @Success      200 		{object}	responses.CurrentUserResponse
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /users/{userID}/api-keys [delete]
func (h *UserHandler) DeleteAPIKey(c fiber.Ctx) error {
	ctx, span := h.tracer.StartFromFiberCtx(c)
	defer span.End()

	ctxLogger := h.tracer.CtxLogger(h.logger, span)

	if c.Params("userID") != string(h.userIDFomContext(c)) {
		return h.responseUnauthorized(c)
	}

	user, err := h.service.RotateAPIKey(ctx, c.OriginalURL(), h.userIDFomContext(c))
	if err != nil {
		ctxLogger.Error(h.tracer.WrapErrorSpan(span, stacktrace.Propagatef(err, "cannot rotate the api key for [%T] with ID [%s]", user, h.userIDFomContext(c))))
		return h.responseInternalServerError(c)
	}

	return h.responseOK(c, "API Key rotated successfully", responses.CurrentUser{
		User:    *user,
		IsAdmin: h.service.IsAdmin(h.userFromContext(c)),
	})
}
