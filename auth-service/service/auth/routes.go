package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/hoyci/ms-chat/auth-service/keys"
	"github.com/hoyci/ms-chat/auth-service/service/crypt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/types"
	coreTypes "github.com/hoyci/ms-chat/core/types"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"golang.org/x/net/context"
)

var validate = validator.New()

type AuthHandler struct {
	userStore       types.UserStore
	authStore       types.AuthStore
	UUIDGen         coreTypes.UUIDGenerator
	passwordHandler crypt.PasswordHandler
}

func NewAuthHandler(
	userStore types.UserStore,
	authStore types.AuthStore,
	UUIDGen coreTypes.UUIDGenerator,
	passwordHandler crypt.PasswordHandler,
) *AuthHandler {
	return &AuthHandler{
		userStore:       userStore,
		authStore:       authStore,
		UUIDGen:         UUIDGen,
		passwordHandler: passwordHandler,
	}
}

// HandleUserLogin
// @Summary Realizar login do usuário
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body coreTypes.UserLoginPayload true "Dados para login do usuário"
// @Success 200 {object} coreTypes.UserLoginResponse "Tokens de acesso e refresh"
// @Failure 400 {object} coreTypes.BadRequestResponse "Body is not a valid json"
// @Failure 400 {object} coreTypes.BadRequestStructResponse "Validation errors for payload"
// @Failure 404 {object} coreTypes.NotFoundResponse "No user found with the given email"
// @Failure 500 {object} coreTypes.InternalServerErrorResponse "An unexpected error occurred"
// @Router /auth/login [post]
func (h *AuthHandler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var requestPayload types.UserLoginPayload
	if err := coreUtils.ParseJSON(r, &requestPayload); err != nil {
		coreUtils.WriteError(
			w, http.StatusBadRequest, err, "HandleUserLogin",
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
			w, http.StatusBadRequest, err, "HandleUserLogin", coreTypes.BadRequestStructResponse{Error: errorMessages},
		)
		return
	}

	user, err := h.userStore.GetByEmail(r.Context(), requestPayload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			coreUtils.WriteError(
				w, http.StatusNotFound, err, "HandleUserLogin",
				coreTypes.NotFoundResponse{Error: "Incorrect credentials. Please try again."},
			)
			return
		}
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleUserLogin",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	err = h.passwordHandler.CheckPassword(r.Context(), user.PasswordHash, requestPayload.Password)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "Incorrect credentials. Please try again."},
		)
	}

	accessToken, err := coreUtils.CreateJWT(
		user.ID, user.Username, user.Email, config.Envs.AccessJWTSecret,
		int64(config.Envs.AccessJWTExpirationInSeconds), h.UUIDGen, keys.PrivateKeyAccess,
	)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
	}

	_, err = coreUtils.VerifyJWT(accessToken, keys.PublicKeyAccess)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	refreshToken, err := coreUtils.CreateJWT(
		user.ID, user.Username, user.Email, config.Envs.RefreshJWTSecret,
		int64(config.Envs.RefreshJWTExpirationInSeconds), h.UUIDGen, keys.PrivateKeyRefresh,
	)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
	}

	refreshTokenClaims, err := coreUtils.VerifyJWT(refreshToken, keys.PublicKeyRefresh)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	err = h.authStore.UpsertRefreshToken(
		r.Context(),
		types.UpdateRefreshTokenPayload{
			UserID:    refreshTokenClaims.UserID,
			Jti:       refreshTokenClaims.RegisteredClaims.ID,
			ExpiresAt: refreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			coreUtils.WriteError(
				w, http.StatusNotFound, err, "HandleGetBookByID",
				coreTypes.NotFoundResponse{Error: fmt.Sprintf("No userID found with ID %s", refreshTokenClaims.UserID)},
			)
			return
		}

		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	_ = coreUtils.WriteJSON(
		w, http.StatusOK, types.UserLoginResponse{
			User: types.UserResponse{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
				UpdatedAt: user.UpdatedAt,
				DeletedAt: user.DeletedAt,
			},
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	)
}

// HandleRefreshToken
// @Summary Atualizar tokens (Refresh Token)
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body coreTypes.RefreshTokenPayload true "Payload contendo o refresh token"
// @Success 200 {object} coreTypes.UpdateRefreshTokenResponse "Novos tokens de acesso e refresh"
// @Failure 400 {object} coreTypes.BadRequestResponse "Body is not a valid json"
// @Failure 400 {object} coreTypes.BadRequestStructResponse "Validation errors for payload"
// @Failure 401 {object} coreTypes.UnauthorizedResponse "Refresh token is invalid or has been expired"
// @Failure 404 {object} coreTypes.NotFoundResponse "No refresh token found with the given user ID"
// @Failure 500 {object} coreTypes.InternalServerErrorResponse "An unexpected error occurred"
// @Router /auth/refresh [post]
func (h *AuthHandler) HandleRefreshToken(w http.ResponseWriter, r *http.Request) {
	var requestPayload types.RefreshTokenPayload
	if err := coreUtils.ParseJSON(r, &requestPayload); err != nil {
		coreUtils.WriteError(
			w, http.StatusBadRequest, err, "HandleUserLogin",
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
			w, http.StatusBadRequest, err, "HandleUserLogin", coreTypes.BadRequestStructResponse{Error: errorMessages},
		)
		return
	}

	claims, err := coreUtils.VerifyJWT(requestPayload.RefreshToken, keys.PublicKeyRefresh)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusUnauthorized, err, "HandleRefreshToken",
			coreTypes.UnauthorizedResponse{Error: "Refresh token is invalid or has been expired"},
		)
		return
	}

	storedToken, err := h.authStore.GetRefreshTokenByUserID(r.Context(), claims.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			coreUtils.WriteError(
				w, http.StatusNotFound, err, "HandleRefreshToken",
				coreTypes.NotFoundResponse{Error: fmt.Sprintf("No refresh token found with user ID %s", claims.UserID)},
			)
			return
		}
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleRefreshToken",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleRefreshToken",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	if storedToken.Jti != claims.RegisteredClaims.ID {
		coreUtils.WriteError(
			w, http.StatusUnauthorized, fmt.Errorf("stored JTI does not match the claims ID"), "HandleRefreshToken",
			coreTypes.UnauthorizedResponse{Error: "Refresh token is invalid or has been expired"},
		)
		return
	}

	newAccessToken, err := coreUtils.CreateJWT(
		claims.UserID, claims.Username, claims.Email, config.Envs.AccessJWTSecret,
		int64(config.Envs.AccessJWTExpirationInSeconds), h.UUIDGen, keys.PrivateKeyAccess,
	)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "Refresh token is invalid or has been expired"},
		)
	}

	newRefreshToken, err := coreUtils.CreateJWT(
		claims.UserID, claims.Username, claims.Email, config.Envs.RefreshJWTSecret,
		int64(config.Envs.RefreshJWTExpirationInSeconds), h.UUIDGen, keys.PrivateKeyRefresh,
	)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "Refresh token is invalid or has been expired"},
		)
	}

	newRefreshTokenClaims, err := coreUtils.VerifyJWT(newRefreshToken, keys.PublicKeyRefresh)
	if err != nil {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleUserLogin",
			coreTypes.InternalServerErrorResponse{Error: "Refresh token is invalid or has been expired"},
		)
		return
	}

	err = h.authStore.UpsertRefreshToken(
		r.Context(),
		types.UpdateRefreshTokenPayload{
			UserID:    newRefreshTokenClaims.UserID,
			Jti:       newRefreshTokenClaims.RegisteredClaims.ID,
			ExpiresAt: newRefreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
		},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			coreUtils.WriteError(
				w, http.StatusNotFound, err, "HandleGetBookByID", coreTypes.NotFoundResponse{
					Error: fmt.Sprintf(
						"No userID found with ID %s", newRefreshTokenClaims.UserID,
					),
				},
			)
			return
		}

		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleRefreshToken",
			coreTypes.InternalServerErrorResponse{Error: "Internal Server Error"},
		)
		return
	}

	_ = coreUtils.WriteJSON(
		w, http.StatusOK, types.UpdateRefreshTokenResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken},
	)
}
