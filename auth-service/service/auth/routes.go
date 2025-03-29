package auth

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/types"
	"github.com/hoyci/ms-chat/auth-service/utils"
	coreTypes "github.com/hoyci/ms-chat/core/types"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"golang.org/x/net/context"
)

var validate = validator.New()

type AuthHandler struct {
	userStore types.UserStore
	authStore types.AuthStore
	UUIDGen   types.UUIDGenerator
}

func NewAuthHandler(
	userStore types.UserStore,
	authStore types.AuthStore,
	UUIDGen types.UUIDGenerator,
) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
		authStore: authStore,
		UUIDGen:   UUIDGen,
	}
}

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
		coreUtils.WriteError(w, http.StatusBadRequest, err, "HandleUserLogin", coreTypes.BadRequestResponse{Error: "Body is not a valid json"})
		return
	}

	if err := validate.Struct(requestPayload); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()))
		}

		coreUtils.WriteError(w, http.StatusBadRequest, err, "HandleUserLogin", coreTypes.BadRequestStructResponse{Error: errorMessages})
		return
	}

	user, err := h.userStore.GetByEmail(r.Context(), requestPayload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			coreUtils.WriteError(w, http.StatusNotFound, err, "HandleUserLogin", coreTypes.NotFoundResponse{Error: "Incorrect credentials. Please try again."})
			return
		}
		if err == context.Canceled {
			coreUtils.WriteError(w, http.StatusServiceUnavailable, err, "HandleUserLogin", coreTypes.ContextCanceledResponse{Error: "Request canceled"})
			return
		}
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
		return
	}

	err = utils.CheckPassword(r.Context(), user.PasswordHash, requestPayload.Password)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "Incorrect credentials. Please try again."})
	}

	accessToken, err := utils.CreateJWT(user.ID, user.Username, user.Email, config.Envs.AccessJWTSecret, int64(config.Envs.AccessJWTExpirationInSeconds), h.UUIDGen, utils.PrivateKeyAccess)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
	}

	_, err = utils.VerifyJWT(accessToken, &utils.PrivateKeyAccess.PublicKey)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
		return
	}

	refreshToken, err := utils.CreateJWT(user.ID, user.Username, user.Email, config.Envs.RefreshJWTSecret, int64(config.Envs.RefreshJWTExpirationInSeconds), h.UUIDGen, utils.PrivateKeyRefresh)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
	}

	refreshTokenClaims, err := utils.VerifyJWT(refreshToken, &utils.PrivateKeyRefresh.PublicKey)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
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
		if err == sql.ErrNoRows {
			coreUtils.WriteError(w, http.StatusNotFound, err, "HandleGetBookByID", coreTypes.NotFoundResponse{Error: fmt.Sprintf("No userID found with ID %d", refreshTokenClaims.UserID)})
			return
		}

		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
		return
	}

	coreUtils.WriteJSON(w, http.StatusOK, types.UserLoginResponse{
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
	})
}

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
		coreUtils.WriteError(w, http.StatusBadRequest, err, "HandleUserLogin", coreTypes.BadRequestResponse{Error: "Body is not a valid json"})
		return
	}

	if err := validate.Struct(requestPayload); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()))
		}

		coreUtils.WriteError(w, http.StatusBadRequest, err, "HandleUserLogin", coreTypes.BadRequestStructResponse{Error: errorMessages})
		return
	}

	claims, err := utils.VerifyJWT(requestPayload.RefreshToken, &utils.PrivateKeyRefresh.PublicKey)
	if err != nil {
		coreUtils.WriteError(w, http.StatusUnauthorized, err, "HandleRefreshToken", coreTypes.UnauthorizedResponse{Error: "Refresh token is invalid or has been expired"})
		return
	}

	storedToken, err := h.authStore.GetRefreshTokenByUserID(r.Context(), claims.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			coreUtils.WriteError(w, http.StatusNotFound, err, "HandleRefreshToken", coreTypes.NotFoundResponse{Error: fmt.Sprintf("No refresh token found with user ID %d", claims.UserID)})
			return
		}
		if err == context.Canceled {
			coreUtils.WriteError(w, http.StatusServiceUnavailable, err, "HandleRefreshToken", coreTypes.ContextCanceledResponse{Error: "Request canceled"})
			return
		}
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleRefreshToken", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
		return
	}

	if storedToken.Jti != claims.RegisteredClaims.ID {
		coreUtils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("stored JTI does not match the claims ID"), "HandleRefreshToken", coreTypes.UnauthorizedResponse{Error: "Refresh token is invalid or has been expired"})
		return
	}

	newAccessToken, err := utils.CreateJWT(claims.UserID, claims.Username, claims.Email, config.Envs.AccessJWTSecret, int64(config.Envs.AccessJWTExpirationInSeconds), h.UUIDGen, utils.PrivateKeyAccess)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "Refresh token is invalid or has been expired"})
	}

	newRefreshToken, err := utils.CreateJWT(claims.UserID, claims.Username, claims.Email, config.Envs.RefreshJWTSecret, int64(config.Envs.RefreshJWTExpirationInSeconds), h.UUIDGen, utils.PrivateKeyRefresh)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "Refresh token is invalid or has been expired"})
	}

	newRefreshTokenClaims, err := utils.VerifyJWT(newRefreshToken, &utils.PrivateKeyRefresh.PublicKey)
	if err != nil {
		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleUserLogin", coreTypes.InternalServerErrorResponse{Error: "Refresh token is invalid or has been expired"})
		return
	}

	err = h.authStore.UpsertRefreshToken(
		r.Context(),
		types.UpdateRefreshTokenPayload{
			UserID:    newRefreshTokenClaims.UserID,
			Jti:       newRefreshTokenClaims.RegisteredClaims.ID,
			ExpiresAt: newRefreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
		})
	if err != nil {
		if err == sql.ErrNoRows {
			coreUtils.WriteError(w, http.StatusNotFound, err, "HandleGetBookByID", coreTypes.NotFoundResponse{Error: fmt.Sprintf("No userID found with ID %d", newRefreshTokenClaims.UserID)})
			return
		}

		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleRefreshToken", coreTypes.InternalServerErrorResponse{Error: "Internal Server Error"})
		return
	}

	coreUtils.WriteJSON(w, http.StatusOK, types.UpdateRefreshTokenResponse{AccessToken: newAccessToken, RefreshToken: newRefreshToken})
}
