package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoyci/ms-chat/auth-service/service/crypt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hoyci/ms-chat/auth-service/types"
	coreTypes "github.com/hoyci/ms-chat/core/types"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
)

var validate = validator.New()

func passwordValidator(sl validator.StructLevel) {
	data := sl.Current().Interface().(types.CreateUserRequestPayload)
	if data.Password != data.ConfirmPassword {
		sl.ReportError(data.ConfirmPassword, "ConfirmPassword", "ConfirmPassword", "password_mismatch", "")
	}
}

type UserHandler struct {
	userStore       types.UserStore
	passwordHandler crypt.PasswordHandler
}

func NewUserHandler(userStore types.UserStore, passwordHandler crypt.PasswordHandler) *UserHandler {
	validate.RegisterStructValidation(passwordValidator, types.CreateUserRequestPayload{})

	return &UserHandler{userStore: userStore, passwordHandler: passwordHandler}
}

// HandleCreateUser
// @Summary Criar um novo usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param request body types.CreateUserRequestPayload true "Payload contendo os dados do novo usuário"
// @Success 201 {object} types.CreateUserResponse "Usuário criado com sucesso"
// @Failure 400 {object} coreTypes.BadRequestResponse "Body is not a valid json"
// @Failure 400 {object} coreTypes.BadRequestStructResponse "Validation errors for payload"
// @Failure 500 {object} coreTypes.InternalServerErrorResponse "An unexpected error occurred"
// @Failure 503 {object} coreTypes.ContextCanceledResponse "Request canceled"
// @Router /users [post]
func (h *UserHandler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var requestPayload types.CreateUserRequestPayload
	if err := coreUtils.ParseJSON(r, &requestPayload); err != nil {
		coreUtils.WriteError(
			w, http.StatusBadRequest, err, "HandleCreateUser",
			coreTypes.BadRequestResponse{Error: "Body is not a valid json"},
		)
		return
	}

	if err := validate.Struct(requestPayload); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()))
		}

		coreUtils.WriteError(
			w, http.StatusBadRequest, err, "HandleCreateUser", coreTypes.BadRequestStructResponse{Error: errorMessages},
		)
		return
	}

	user, _ := h.userStore.GetByEmail(r.Context(), requestPayload.Email)

	if user != nil {
		coreUtils.WriteError(
			w, http.StatusConflict, fmt.Errorf("email is already in use"), "HandleCreateUser",
			coreTypes.BadRequestResponse{Error: "This email is already in use"},
		)
		return
	}

	hashedPassword, err := h.passwordHandler.HashPassword(r.Context(), requestPayload.Password)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleCreateUser",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	var databasePayload = types.CreateUserDatabasePayload{
		Username:     requestPayload.Username,
		Email:        requestPayload.Email,
		PasswordHash: hashedPassword,
	}

	_, err = h.userStore.Create(r.Context(), databasePayload)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleCreateUser",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}

		if errors.Is(err, sql.ErrConnDone) {
			coreUtils.WriteError(
				w, http.StatusInternalServerError, err, "HandleCreateUser",
				coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
			)
			return
		}

		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleCreateUser",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	_ = coreUtils.WriteJSON(w, http.StatusCreated, types.CreateUserResponse{Message: "User successfully created"})
}

// HandleGetUserByID
// @Summary      Get user by ID
// @Description  Retrieves user details based on the authenticated user's ID extracted from the request context.
// @Tags         Users
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Success      200  {object}  types.UserResponse  "User details successfully retrieved"
// @Failure      400  {object}  coreTypes.BadRequestResponse "Bad request"
// @Failure      401  {object}  types.UnauthorizedResponse "Unauthorized"
// @Failure      404  {object}  coreTypes.NotFoundResponse "User not found"
// @Failure      500  {object}  coreTypes.InternalServerErrorResponse "Internal server error"
// @Router       /users [get]
func (h *UserHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := coreUtils.GetClaimFromContext[string](r, "UserID")
	if !ok || userID == "" {
		coreUtils.WriteError(
			w, http.StatusUnauthorized, fmt.Errorf("failed to retrieve userID from context"),
			"HandleGetUserByID", coreTypes.InternalServerErrorResponse{Error: "Failed to retrieve userID from context"},
		)
		return
	}

	user, err := h.userStore.GetByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleGetUserByID",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}

		if errors.Is(err, sql.ErrConnDone) {
			coreUtils.WriteError(
				w, http.StatusInternalServerError, err, "HandleGetUserByID",
				coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
			)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			coreUtils.WriteError(
				w, http.StatusNotFound, err, "HandleGetUserByID",
				coreTypes.NotFoundResponse{Error: fmt.Sprintf("No user found with ID %s", userID)},
			)
			return
		}
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleGetUserByID",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	_ = coreUtils.WriteJSON(w, http.StatusOK, user)
}

// HandleUpdateUserByID
// @Summary      Update user by ID
// @Description  Updates user details based on the authenticated user's ID extracted from the request context.
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param request body  types.UpdateUserPayload  true  "User update payload"
// @Success      200  {object}  types.UserResponse  "User successfully updated"
// @Failure      400  {object}  types.core"Invalid request body"
// @Failure      401  {object}  types.UnauthorizedResponse "Unauthorized"
// @Failure      404  {object}  coreTypes.NotFoundResponse "User not found"
// @Failure      500  {object}  coreTypes.InternalServerErrorResponse "Internal server error"
// @Router       /users [put]
func (h *UserHandler) HandleUpdateUserByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := coreUtils.GetClaimFromContext[string](r, "UserID")
	if !ok || userID == "" {
		coreUtils.WriteError(
			w, http.StatusUnauthorized, fmt.Errorf("failed to retrieve userID from context"),
			"HandleUpdateUserByID",
			coreTypes.InternalServerErrorResponse{Error: "Failed to retrieve userID from context"},
		)
		return
	}

	var payload types.UpdateUserPayload
	if err := coreUtils.ParseJSON(r, &payload); err != nil {
		coreUtils.WriteError(
			w, http.StatusBadRequest, err, "HandleUpdateUserByID",
			coreTypes.BadRequestResponse{Error: "Body is not a valid json"},
		)
		return
	}

	if err := validate.Struct(payload); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessages = append(
				errorMessages, fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", e.Field(), e.Tag()),
			)
		}

		coreUtils.WriteError(
			w, http.StatusBadRequest, err, "HandleUpdateUserByID",
			coreTypes.BadRequestStructResponse{Error: errorMessages},
		)
		return
	}

	user, err := h.userStore.UpdateByID(r.Context(), userID, payload)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleUpdateUserByID",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}

		if errors.Is(err, sql.ErrConnDone) {
			coreUtils.WriteError(
				w, http.StatusInternalServerError, err, "HandleUpdateUserByID",
				coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
			)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			coreUtils.WriteError(
				w, http.StatusNotFound, err, "HandleUpdateUserByID",
				coreTypes.NotFoundResponse{Error: fmt.Sprintf("No user found with ID %s", userID)},
			)
			return
		}
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUpdateUserByID",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	_ = coreUtils.WriteJSON(w, http.StatusOK, user)
}

// HandleDeleteUserByID
// @Summary      Delete user by ID
// @Description  Deletes the user associated with the authenticated user's ID extracted from the request context.
// @Tags         Users
// @Security     BearerAuth
// @Success      204  "User successfully deleted"
// @Failure      401  {object}  types.UnauthorizedResponse "Unauthorized"
// @Failure      404  {object}  coreTypes.NotFoundResponse "User not found"
// @Failure      500  {object}  coreTypes.InternalServerErrorResponse "Internal server error"
// @Router       /users [delete]
func (h *UserHandler) HandleDeleteUserByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := coreUtils.GetClaimFromContext[string](r, "UserID")
	if !ok || userID == "" {
		coreUtils.WriteError(
			w, http.StatusUnauthorized, fmt.Errorf("failed to retrieve userID from context"),
			"HandleDeleteUserByID",
			coreTypes.InternalServerErrorResponse{Error: "Failed to retrieve userID from context"},
		)
		return
	}

	err := h.userStore.DeleteByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleDeleteUserByID",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}

		if errors.Is(err, sql.ErrConnDone) {
			coreUtils.WriteError(
				w, http.StatusInternalServerError, err, "HandleDeleteUserByID",
				coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
			)
			return
		}

		if errors.Is(err, sql.ErrNoRows) {
			coreUtils.WriteError(
				w, http.StatusNotFound, err, "HandleDeleteUserByID",
				coreTypes.NotFoundResponse{Error: fmt.Sprintf("No user found with ID %s", userID)},
			)
			return
		}

		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleDeleteUserByID",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	_ = coreUtils.WriteJSON(w, http.StatusNoContent, nil)
}
