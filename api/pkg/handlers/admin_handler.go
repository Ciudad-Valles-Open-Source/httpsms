package handlers

import (
	"fmt"
	"net/url"

	"github.com/NdoleStudio/httpsms/pkg/entities"
	"github.com/NdoleStudio/httpsms/pkg/repositories"
	"github.com/NdoleStudio/httpsms/pkg/requests"
	"github.com/NdoleStudio/httpsms/pkg/services"
	"github.com/NdoleStudio/httpsms/pkg/telemetry"
	"github.com/NdoleStudio/httpsms/pkg/validators"
	"github.com/NdoleStudio/stacktrace"
	"github.com/davecgh/go-spew/spew"
	"github.com/gofiber/fiber/v3"
)

// AdminHandler handles the administration of user accounts
type AdminHandler struct {
	handler
	logger    telemetry.Logger
	tracer    telemetry.Tracer
	validator *validators.AdminHandlerValidator
	service   *services.AdminService
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(
	logger telemetry.Logger,
	tracer telemetry.Tracer,
	validator *validators.AdminHandlerValidator,
	service *services.AdminService,
) (h *AdminHandler) {
	return &AdminHandler{
		logger:    logger.WithService(fmt.Sprintf("%T", h)),
		tracer:    tracer,
		validator: validator,
		service:   service,
	}
}

// RegisterRoutes registers the routes for the AdminHandler. The middlewares
// must authenticate the request and allow only administrators.
func (h *AdminHandler) RegisterRoutes(router fiber.Router, middlewares ...fiber.Handler) {
	h.register(router, fiber.MethodGet, "/v1/admin/users", middlewares, h.Index)
	h.register(router, fiber.MethodPut, "/v1/admin/users/:userID/status", middlewares, h.UpdateStatus)
	h.register(router, fiber.MethodDelete, "/v1/admin/users/:userID", middlewares, h.Delete)
}

// Index lists the user accounts
// @Summary      List user accounts
// @Description  Lists the user accounts, newest first. Only administrators (a verified email listed in ADMIN_EMAILS, signed in on the web dashboard) can call this endpoint; API keys are rejected.
// @Security	 ApiKeyAuth
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        status		query  string  	false	"filter users by status"							Enums(pending, active, suspended)
// @Param        query		query  string  	false 	"filter users whose email contains query"
// @Param        skip		query  int  	false	"number of users to skip"							minimum(0)
// @Param        limit		query  int  	false	"number of users to return (default 20, capped at 100)"	minimum(1)	maximum(100)
// @Success      200 		{object}	responses.AdminUsersResponse
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 403    	{object}	responses.Forbidden
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /admin/users [get]
func (h *AdminHandler) Index(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	var request requests.AdminUserIndex
	if err := c.Bind().Query(&request); err != nil {
		ctxLogger.Warn(stacktrace.Propagatef(err, "cannot marshall params [%s] into %T", c.OriginalURL(), request))
		return h.responseBadRequest(c, err)
	}

	if errors := h.validator.ValidateIndex(ctx, request.Sanitize()); len(errors) != 0 {
		ctxLogger.Warn(stacktrace.NewErrorf("validation errors [%s], while listing users [%+#v]", spew.Sdump(errors), request))
		return h.responseUnprocessableEntity(c, errors, "validation errors while listing users")
	}

	users, err := h.service.Index(ctx, h.userFromContext(c), request.ToIndexParams())
	if err != nil {
		return h.responseAdminError(c, ctxLogger, err, "cannot list users")
	}

	return h.responseOK(c, fmt.Sprintf("fetched %d %s", len(users), h.pluralize("user", len(users))), users)
}

// UpdateStatus approves or suspends a user account
// @Summary      Approve or suspend a user account
// @Description  Sets the status of a user account to "active" (approve) or "suspended". Approving stamps approved_at and approved_by. Administrators can't change their own account or another administrator's account. Only administrators can call this endpoint; API keys are rejected.
// @Security	 ApiKeyAuth
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param 		 userID 	path		string 							true 	"ID of the user"	default(WB7DRDWrJZRGbYrv2CKGkqbzvqdC)
// @Param        payload   	body 		requests.AdminUserStatusUpdate	true 	"New status of the user"
// @Success      200 		{object}	responses.AdminUserResponse
// @Failure      400		{object}	responses.BadRequest
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 403    	{object}	responses.Forbidden
// @Failure      404		{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /admin/users/{userID}/status [put]
func (h *AdminHandler) UpdateStatus(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := c.Params("userID")
	if errors := h.validator.ValidateUserID(userID); len(errors) != 0 {
		ctxLogger.Warn(stacktrace.NewErrorf("validation errors [%s], while updating the status of user [%s]", spew.Sdump(errors), userID))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating the status of the user")
	}

	var request requests.AdminUserStatusUpdate
	if err := c.Bind().Body(&request); err != nil {
		ctxLogger.Warn(stacktrace.Propagatef(err, "cannot marshall params [%s] into %T", c.OriginalURL(), request))
		return h.responseBadRequest(c, err)
	}

	if errors := h.validator.ValidateStatusUpdate(ctx, request.Sanitize()); len(errors) != 0 {
		ctxLogger.Warn(stacktrace.NewErrorf("validation errors [%s], while updating the status of user [%s]", spew.Sdump(errors), userID))
		return h.responseUnprocessableEntity(c, errors, "validation errors while updating the status of the user")
	}

	user, err := h.service.UpdateStatus(ctx, c.OriginalURL(), h.userFromContext(c), entities.UserID(userID), request.UserStatus())
	if err != nil {
		return h.responseAdminError(c, ctxLogger, err, fmt.Sprintf("cannot update the status of user [%s]", userID))
	}

	return h.responseOK(c, fmt.Sprintf("user status updated to [%s]", user.Status), user)
}

// Delete deletes a user account
// @Summary      Delete a user account
// @Description  Deletes a user account together with all its data. Administrators can't delete their own account or another administrator's account. Only administrators can call this endpoint; API keys are rejected.
// @Security	 ApiKeyAuth
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param 		 userID 	path		string 		true 	"ID of the user"	default(WB7DRDWrJZRGbYrv2CKGkqbzvqdC)
// @Success      204 		{object}	responses.NoContent
// @Failure 	 401    	{object}	responses.Unauthorized
// @Failure 	 403    	{object}	responses.Forbidden
// @Failure      404		{object}	responses.NotFound
// @Failure      422		{object}	responses.UnprocessableEntity
// @Failure      500		{object}	responses.InternalServerError
// @Router       /admin/users/{userID} [delete]
func (h *AdminHandler) Delete(c fiber.Ctx) error {
	ctx, span, ctxLogger := h.tracer.StartFromFiberCtxWithLogger(c, h.logger)
	defer span.End()

	userID := c.Params("userID")
	if errors := h.validator.ValidateUserID(userID); len(errors) != 0 {
		ctxLogger.Warn(stacktrace.NewErrorf("validation errors [%s], while deleting user [%s]", spew.Sdump(errors), userID))
		return h.responseUnprocessableEntity(c, errors, "validation errors while deleting the user")
	}

	if err := h.service.Delete(ctx, c.OriginalURL(), h.userFromContext(c), entities.UserID(userID)); err != nil {
		return h.responseAdminError(c, ctxLogger, err, fmt.Sprintf("cannot delete user [%s]", userID))
	}

	return h.responseNoContent(c, "user deleted successfully")
}

func (h *AdminHandler) responseAdminError(c fiber.Ctx, ctxLogger telemetry.Logger, err error, message string) error {
	switch stacktrace.GetCode(err) {
	case services.ErrCodeNotAdmin:
		ctxLogger.Warn(stacktrace.Propagatef(err, "%s", message))
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "error", "message": "forbidden"})
	case services.ErrCodeProtectedUser:
		ctxLogger.Warn(stacktrace.Propagatef(err, "%s", message))
		return h.responseUnprocessableEntity(c, url.Values{"userID": []string{"Administrators can't change their own account, another administrator's account or the system user."}}, "this user account is protected")
	case repositories.ErrCodeNotFound:
		ctxLogger.Warn(stacktrace.Propagatef(err, "%s", message))
		return h.responseNotFound(c, fmt.Sprintf("cannot find user with ID [%s]", c.Params("userID")))
	default:
		ctxLogger.Error(stacktrace.Propagatef(err, "%s", message))
		return h.responseInternalServerError(c)
	}
}
