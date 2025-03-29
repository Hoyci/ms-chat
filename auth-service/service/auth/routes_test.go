package auth_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/hoyci/ms-chat/auth-service/cmd/api"
	"github.com/hoyci/ms-chat/auth-service/mocks"
	"github.com/hoyci/ms-chat/auth-service/service/auth"
	"github.com/hoyci/ms-chat/auth-service/types"
	"github.com/hoyci/ms-chat/auth-service/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleUserLogin(t *testing.T) {
	setupTestServer := func() (*mocks.MockUserStore, *mocks.MockAuthStore, *mocks.MockUUIDGenerator, *httptest.Server, *mux.Router) {
		mockUUID := new(mocks.MockUUIDGenerator)
		mockAuthStore := new(mocks.MockAuthStore)
		mockUserStore := new(mocks.MockUserStore)
		mockAuthHandler := auth.NewAuthHandler(mockUserStore, mockAuthStore, mockUUID)
		apiServer := api.NewApiServer(":8080", nil)
		router := apiServer.SetupRouter(nil, nil, mockAuthHandler)
		ts := httptest.NewServer(router)
		return mockUserStore, mockAuthStore, mockUUID, ts, router
	}

	t.Run("it should throw an error when body is not a valid JSON", func(t *testing.T) {
		_, _, _, ts, router := setupTestServer()
		defer ts.Close()

		invalidBody := bytes.NewReader([]byte("INVALID JSON"))
		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", invalidBody)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":"Body is not a valid json"}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when body is a valid JSON but missing key", func(t *testing.T) {
		_, _, _, ts, router := setupTestServer()
		defer ts.Close()

		payload := types.UserLoginPayload{}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":["Field 'Email' is invalid: required", "Field 'Password' is invalid: required"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when body does not contain a valid email", func(t *testing.T) {
		_, _, _, ts, router := setupTestServer()
		defer ts.Close()

		payload := types.UserLoginPayload{
			Email:    "johndoe",
			Password: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":["Field 'Email' is invalid: email"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when password or confirmPassword is smaller than 8 chars", func(t *testing.T) {
		_, _, _, ts, router := setupTestServer()
		defer ts.Close()

		payload := types.UserLoginPayload{
			Email:    "johndoe@email.com",
			Password: "12345",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":["Field 'Password' is invalid: min"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when call endpoint with a non-existent user ID", func(t *testing.T) {
		mockUserStore, _, _, ts, router := setupTestServer()
		defer ts.Close()

		mockUserStore.On("GetByEmail", mock.Anything, mock.Anything).Return((*types.GetByEmailResponse)(nil), sql.ErrNoRows)

		payload := types.UserLoginPayload{
			Email:    "johndoe@email.com",
			Password: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)
		expected := `{"error": "Incorrect credentials. Please try again."}`
		assert.JSONEq(t, expected, string(responseBody))
	})

	t.Run("it should return error when the request context is canceled during the process of get user by email", func(t *testing.T) {
		mockUserStore, _, _, ts, router := setupTestServer()
		defer ts.Close()

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockUserStore.On("GetByEmail", mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Err() == context.Canceled
		}), mock.Anything).Return((*types.GetByEmailResponse)(nil), context.Canceled)

		payload := types.UserLoginPayload{
			Email:    "johndoe@email.com",
			Password: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled)).WithContext(canceledCtx)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		expected := `{"error":"Request canceled"}`
		assert.JSONEq(t, expected, string(responseBody))
	})

	t.Run("it should throw a database find error", func(t *testing.T) {
		mockUserStore, _, _, ts, router := setupTestServer()
		defer ts.Close()

		mockUserStore.On("GetByEmail", mock.Anything, mock.Anything).Return((*types.GetByEmailResponse)(nil), fmt.Errorf("no row found with email: 'johndoe@email.com'"))

		payload := types.UserLoginPayload{
			Email:    "johndoe@email.com",
			Password: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)
		expected := `{"error":"An unexpected error occurred"}`
		assert.JSONEq(t, expected, string(responseBody))
	})

	t.Run("it should successfully authenticate a user", func(t *testing.T) {
		mockUserStore, mockAuthStore, mockUUID, ts, router := setupTestServer()
		defer ts.Close()

		mockUUID.On("New").Return("mocked-uuid")

		hashedPassword, _ := utils.HashPassword(context.Background(), "123mudar")

		mockUserStore.On("GetByEmail", mock.Anything, mock.Anything).Return(
			&types.GetByEmailResponse{
				ID:           1,
				Username:     "JohnDoe",
				Email:        "johndoe@email.com",
				PasswordHash: hashedPassword,
				CreatedAt:    time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    nil,
				DeletedAt:    nil,
			},
			nil,
		)

		mockAuthStore.On("UpsertRefreshToken", mock.Anything, mock.Anything).Return(nil)

		payload := types.UserLoginPayload{
			Email:    "johndoe@email.com",
			Password: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var responseMap types.UserLoginResponse
		err = json.Unmarshal(responseBody, &responseMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		accessTokenClaims, err := utils.VerifyJWT(responseMap.AccessToken, &utils.PrivateKeyAccess.PublicKey)
		assert.NoError(t, err, "Failed to verify JWT token")

		assert.Equal(t, "johndoe@email.com", accessTokenClaims.Email, "Email claim mismatch")
		assert.Equal(t, "JohnDoe", accessTokenClaims.Username, "Username claim mismatch")
		assert.Equal(t, 1, accessTokenClaims.UserID, "UserID claim mismatch")

		refresh_token_claims, err := utils.VerifyJWT(responseMap.RefreshToken, &utils.PrivateKeyRefresh.PublicKey)
		assert.NoError(t, err, "Failed to verify JWT token")

		assert.Equal(t, 1, refresh_token_claims.UserID, "UserID claim mismatch")
	})
}

func TestHandleRefreshToken(t *testing.T) {
	setupTestServer := func() (*mocks.MockUserStore, *mocks.MockAuthStore, *mocks.MockUUIDGenerator, *httptest.Server, *mux.Router) {
		mockUUID := new(mocks.MockUUIDGenerator)
		mockAuthStore := new(mocks.MockAuthStore)
		mockUserStore := new(mocks.MockUserStore)
		mockAuthHandler := auth.NewAuthHandler(mockUserStore, mockAuthStore, mockUUID)
		apiServer := api.NewApiServer(":8080", nil)
		router := apiServer.SetupRouter(nil, nil, mockAuthHandler)
		ts := httptest.NewServer(router)
		return mockUserStore, mockAuthStore, mockUUID, ts, router
	}

	t.Run("it should throw an error when body is not a valid JSON", func(t *testing.T) {
		_, _, _, ts, router := setupTestServer()
		defer ts.Close()

		invalidBody := bytes.NewReader([]byte("INVALID JSON"))
		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/refresh", invalidBody)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":"Body is not a valid json"}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when body is a valid JSON but missing key", func(t *testing.T) {
		_, _, _, ts, router := setupTestServer()
		defer ts.Close()

		payload := types.UserLoginPayload{}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/refresh", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":["Field 'RefreshToken' is invalid: required"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when body does not contain a valid token", func(t *testing.T) {
		_, _, _, ts, router := setupTestServer()
		defer ts.Close()

		payload := types.RefreshTokenPayload{
			RefreshToken: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/refresh", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":"Refresh token is invalid or has been expired"}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should return error when the request context is canceled during the process of get refresh token by user id", func(t *testing.T) {
		token := utils.GenerateTestToken(1, "JohnDoe", "johndoe@example.com", utils.PrivateKeyRefresh)
		_, mockAuthStore, _, ts, router := setupTestServer()
		defer ts.Close()

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockAuthStore.On("GetRefreshTokenByUserID", mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Err() == context.Canceled
		}), 1).Return((*types.RefreshToken)(nil), context.Canceled)

		payload := types.RefreshTokenPayload{
			RefreshToken: token,
		}
		payloadMarshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/refresh", bytes.NewBuffer(payloadMarshalled)).WithContext(canceledCtx)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusServiceUnavailable, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		assert.NoError(t, err)

		expected := `{"error":"Request canceled"}`
		assert.JSONEq(t, expected, string(responseBody))
	})

	t.Run("it should successfully refresh user token", func(t *testing.T) {
		mockUserStore, mockAuthStore, mockUUID, ts, router := setupTestServer()
		defer ts.Close()

		mockUUID.On("New").Return("mocked-uuid")

		hashedPassword, _ := utils.HashPassword(context.Background(), "123mudar")

		mockUserStore.On("GetByEmail", mock.Anything, mock.Anything).Return(
			&types.GetByEmailResponse{
				ID:           1,
				Username:     "JohnDoe",
				Email:        "johndoe@email.com",
				PasswordHash: hashedPassword,
				CreatedAt:    time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    nil,
				DeletedAt:    nil,
			},
			nil,
		)

		mockAuthStore.On("GetRefreshTokenByUserID", mock.Anything, mock.Anything).Return(
			&types.RefreshToken{
				ID:        1,
				UserID:    1,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				Jti:       "mocked-uuid",
			},
			nil,
		)

		mockAuthStore.On("UpsertRefreshToken", mock.Anything, mock.Anything).Return(nil)

		userLoginPayload := types.UserLoginPayload{
			Email:    "johndoe@email.com",
			Password: "123mudar",
		}
		marshalled, _ := json.Marshal(userLoginPayload)

		userLoginReq := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		userLoginW := httptest.NewRecorder()

		router.ServeHTTP(userLoginW, userLoginReq)

		resUserLogin := userLoginW.Result()
		defer resUserLogin.Body.Close()

		assert.Equal(t, http.StatusOK, resUserLogin.StatusCode)

		responseUserLoginBody, err := io.ReadAll(resUserLogin.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var responseUserLoginMap types.UserLoginResponse
		err = json.Unmarshal(responseUserLoginBody, &responseUserLoginMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.NotEmpty(t, responseUserLoginMap.AccessToken, "Access token should not be empty")
		assert.NotEmpty(t, responseUserLoginMap.RefreshToken, "Refresh token should not be empty")

		userRefreshTokenPayload := types.RefreshTokenPayload{
			RefreshToken: responseUserLoginMap.RefreshToken,
		}
		userRefreshTokenMarshalled, _ := json.Marshal(userRefreshTokenPayload)

		reqRefreshToken := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/refresh", bytes.NewBuffer(userRefreshTokenMarshalled))
		wRefreshToken := httptest.NewRecorder()

		router.ServeHTTP(wRefreshToken, reqRefreshToken)

		resRefreshToken := wRefreshToken.Result()
		defer resRefreshToken.Body.Close()

		assert.Equal(t, http.StatusOK, resRefreshToken.StatusCode)

		responseRefreshTokenBody, err := io.ReadAll(resRefreshToken.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var responseRefreshTokenMap map[string]interface{}
		err = json.Unmarshal(responseRefreshTokenBody, &responseRefreshTokenMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		access_token, ok := responseRefreshTokenMap["access_token"].(string)
		if !ok {
			t.Fatalf("access_token not found or not a string")
		}

		refresh_token, ok := responseRefreshTokenMap["refresh_token"].(string)
		if !ok {
			t.Fatalf("refresh_token not found or not a string")
		}

		access_token_claims, err := utils.VerifyJWT(access_token, &utils.PrivateKeyAccess.PublicKey)
		assert.NoError(t, err, "Failed to verify JWT token")

		assert.Equal(t, "johndoe@email.com", access_token_claims.Email, "Email claim mismatch")
		assert.Equal(t, "JohnDoe", access_token_claims.Username, "Username claim mismatch")
		assert.Equal(t, 1, access_token_claims.UserID, "UserID claim mismatch")

		refresh_token_claims, err := utils.VerifyJWT(refresh_token, &utils.PrivateKeyRefresh.PublicKey)
		assert.NoError(t, err, "Failed to verify JWT token")
		assert.Equal(t, 1, refresh_token_claims.UserID, "UserID claim mismatch")
	})

	t.Run("it should not be able refresh user with expired token", func(t *testing.T) {
		mockUserStore, mockAuthStore, mockUUID, ts, router := setupTestServer()
		defer ts.Close()

		mockUUID.On("New").Return("mocked-uuid")

		hashedPassword, _ := utils.HashPassword(context.Background(), "123mudar")

		mockUserStore.On("GetByEmail", mock.Anything, mock.Anything).Return(
			&types.GetByEmailResponse{
				ID:           1,
				Username:     "JohnDoe",
				Email:        "johndoe@email.com",
				PasswordHash: hashedPassword,
				CreatedAt:    time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
				UpdatedAt:    nil,
				DeletedAt:    nil,
			},
			nil,
		)

		mockAuthStore.On("GetRefreshTokenByUserID", mock.Anything, mock.Anything).Return(
			&types.RefreshToken{
				ID:        1,
				UserID:    1,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				Jti:       "mocked-uuid",
			},
			nil,
		).Once()

		mockAuthStore.On("UpsertRefreshToken", mock.Anything, mock.Anything).Return(nil)

		userLoginPayload := types.UserLoginPayload{
			Email:    "johndoe@email.com",
			Password: "123mudar",
		}
		marshalled, _ := json.Marshal(userLoginPayload)

		userLoginReq := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth", bytes.NewBuffer(marshalled))
		userLoginW := httptest.NewRecorder()

		router.ServeHTTP(userLoginW, userLoginReq)

		resUserLogin := userLoginW.Result()
		defer resUserLogin.Body.Close()

		assert.Equal(t, http.StatusOK, resUserLogin.StatusCode)

		responseUserLoginBody, err := io.ReadAll(resUserLogin.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var responseUserLoginMap types.UserLoginResponse
		err = json.Unmarshal(responseUserLoginBody, &responseUserLoginMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		assert.NotEmpty(t, responseUserLoginMap.AccessToken, "Access token should not be empty")
		assert.NotEmpty(t, responseUserLoginMap.RefreshToken, "Refresh token should not be empty")

		userRefreshTokenPayload := types.RefreshTokenPayload{
			RefreshToken: responseUserLoginMap.RefreshToken,
		}
		userRefreshTokenMarshalled, _ := json.Marshal(userRefreshTokenPayload)

		reqRefreshToken1 := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/refresh", bytes.NewBuffer(userRefreshTokenMarshalled))
		wRefreshToken1 := httptest.NewRecorder()

		router.ServeHTTP(wRefreshToken1, reqRefreshToken1)

		resRefreshToken1 := wRefreshToken1.Result()
		defer resRefreshToken1.Body.Close()

		assert.Equal(t, http.StatusOK, resRefreshToken1.StatusCode)

		mockAuthStore.On("GetRefreshTokenByUserID", mock.Anything, mock.Anything).Return(
			&types.RefreshToken{
				ID:        1,
				UserID:    1,
				CreatedAt: time.Now(),
				ExpiresAt: time.Now().Add(24 * time.Hour),
				Jti:       "mocked-uuid2",
			},
			nil,
		).Once()

		reqRefreshToken2 := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/auth/refresh", bytes.NewBuffer(userRefreshTokenMarshalled))
		wRefreshToken2 := httptest.NewRecorder()

		router.ServeHTTP(wRefreshToken2, reqRefreshToken2)

		resRefreshToken2 := wRefreshToken2.Result()
		defer resRefreshToken2.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resRefreshToken2.StatusCode)
	})
}
