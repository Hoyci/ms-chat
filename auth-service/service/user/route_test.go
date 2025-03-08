package user_test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/hoyci/ms-chat/auth-service/cmd/api"
	"github.com/hoyci/ms-chat/auth-service/config"
	"github.com/hoyci/ms-chat/auth-service/mocks"
	"github.com/hoyci/ms-chat/auth-service/service/user"
	"github.com/hoyci/ms-chat/auth-service/types"
	"github.com/hoyci/ms-chat/auth-service/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreateUser(t *testing.T) {
	setupTestServer := func() (*mocks.MockUserStore, *httptest.Server, *mux.Router, config.Config) {
		mockUserStore := new(mocks.MockUserStore)
		mockUserHandler := user.NewUserHandler(mockUserStore)
		apiServer := api.NewApiServer(":8080", nil)
		router := apiServer.SetupRouter(nil, mockUserHandler, nil)
		ts := httptest.NewServer(router)
		return mockUserStore, ts, router, apiServer.Config
	}

	t.Run("it should throw an error when body is not a valid JSON", func(t *testing.T) {
		_, ts, router, _ := setupTestServer()
		defer ts.Close()

		invalidBody := bytes.NewReader([]byte("INVALID JSON"))
		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/users", invalidBody)
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
		_, ts, router, _ := setupTestServer()
		defer ts.Close()

		payload := types.CreateUserRequestPayload{}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/users", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":["Field 'Username' is invalid: required", "Field 'Email' is invalid: required", "Field 'Password' is invalid: required", "Field 'ConfirmPassword' is invalid: required"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when body does not contain a valid email", func(t *testing.T) {
		_, ts, router, _ := setupTestServer()
		defer ts.Close()

		payload := types.CreateUserRequestPayload{
			Username:        "JohnDoe",
			Email:           "johndoe",
			Password:        "123mudar",
			ConfirmPassword: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/users", bytes.NewBuffer(marshalled))
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
		_, ts, router, _ := setupTestServer()
		defer ts.Close()

		payload := types.CreateUserRequestPayload{
			Username:        "JohnDoe",
			Email:           "johndoe@email.com",
			Password:        "12345",
			ConfirmPassword: "12345",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/users", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":["Field 'Password' is invalid: min", "Field 'ConfirmPassword' is invalid: min"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should return error when the request context is canceled", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockUserStore.On("Create", mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Err() == context.Canceled
		}), mock.Anything).Return((*types.UserResponse)(nil), context.Canceled)

		mockUserStore.On("GetByEmail", mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Err() == context.Canceled
		}), mock.Anything).Return((*types.GetByEmailResponse)(nil), context.Canceled)

		payload := types.CreateUserRequestPayload{
			Username:        "JohnDoe",
			Email:           "johndoe@email.com",
			Password:        "123mudar",
			ConfirmPassword: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/users", bytes.NewBuffer(marshalled)).WithContext(canceledCtx)
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

	t.Run("it should throw a database connection error", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("GetByEmail", mock.Anything, mock.Anything).Return((*types.GetByEmailResponse)(nil), sql.ErrConnDone)
		mockUserStore.On("Create", mock.Anything, mock.Anything).Return((*types.UserResponse)(nil), sql.ErrConnDone)

		payload := types.CreateUserRequestPayload{
			Username:        "JohnDoe",
			Email:           "johndoe@email.com",
			Password:        "123mudar",
			ConfirmPassword: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/users", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)
		expected := `{"error":"An unexpected error occurred"}`
		assert.JSONEq(t, expected, string(responseBody))
	})

	t.Run("it should successfully create a user", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("GetByEmail", mock.Anything, mock.Anything).Return((*types.GetByEmailResponse)(nil), nil)

		mockUserStore.On("Create", mock.Anything, mock.Anything).Return(
			&types.UserResponse{
				ID:        1,
				Username:  "JohnDoe",
				Email:     "johndoe@email.com",
				CreatedAt: time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
				UpdatedAt: nil,
				DeletedAt: nil,
			},
			nil,
		)

		payload := types.CreateUserRequestPayload{
			Username:        "JohnDoe",
			Email:           "johndoe@email.com",
			Password:        "123mudar",
			ConfirmPassword: "123mudar",
		}
		marshalled, _ := json.Marshal(payload)

		req := httptest.NewRequest(http.MethodPost, ts.URL+"/api/v1/users", bytes.NewBuffer(marshalled))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusCreated, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		var responseMap map[string]interface{}
		err = json.Unmarshal(responseBody, &responseMap)
		if err != nil {
			t.Fatalf("Failed to unmarshal response body: %v", err)
		}

		responseMessage, ok := responseMap["message"].(string)
		if !ok {
			t.Fatalf("Token not found or not a string")
		}
		assert.Equal(t, "User successfully created", responseMessage)
	})
}

func TestHandleGetUser(t *testing.T) {
	setupTestServer := func() (*mocks.MockUserStore, *httptest.Server, *mux.Router, config.Config) {
		mockUserStore := new(mocks.MockUserStore)
		mockUserHandler := user.NewUserHandler(mockUserStore)
		apiServer := api.NewApiServer(":8080", nil)
		router := apiServer.SetupRouter(nil, mockUserHandler, nil)
		ts := httptest.NewServer(router)
		return mockUserStore, ts, router, apiServer.Config
	}

	token := utils.GenerateTestToken(1, "JohnDoe", "johndoe@example.com", utils.PrivateKeyAccess)

	t.Run("it should return error when context is canceled", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockUserStore.On("GetByID", mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Err() == context.Canceled
		}), 1).Return(&types.UserResponse{}, context.Canceled)

		req := httptest.NewRequest(http.MethodGet, ts.URL+"/api/v1/users", nil).WithContext(canceledCtx)
		req.Header.Set("Authorization", "Bearer "+token)
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

	t.Run("it should throw a database connection error", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("GetByID", mock.Anything, mock.Anything).Return(&types.UserResponse{}, sql.ErrConnDone)

		req := httptest.NewRequest(http.MethodGet, ts.URL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":"An unexpected error occurred"}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when database not found user by ID", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("GetByID", mock.Anything, mock.Anything).Return(&types.UserResponse{}, sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodGet, ts.URL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error": "No user found with ID 1"}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should successfully get a user by ID", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("GetByID", mock.Anything, mock.Anything).Return(&types.UserResponse{
			ID:        1,
			Username:  "johndoe",
			Email:     "johndoe@email.com",
			CreatedAt: time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
			DeletedAt: nil,
			UpdatedAt: nil,
		}, nil)

		req := httptest.NewRequest(http.MethodGet, ts.URL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{
			"id": 1,
			"username": "johndoe",
			"email": "johndoe@email.com",
			"createdAt": "0001-01-01T00:00:00Z",
			"deletedAt": null,
			"updatedAt": null
		}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})
}

func TestHandleUpdateUser(t *testing.T) {
	setupTestServer := func() (*mocks.MockUserStore, *httptest.Server, *mux.Router, config.Config) {
		mockUserStore := new(mocks.MockUserStore)
		mockUserHandler := user.NewUserHandler(mockUserStore)
		apiServer := api.NewApiServer(":8080", nil)
		router := apiServer.SetupRouter(nil, mockUserHandler, nil)
		ts := httptest.NewServer(router)
		return mockUserStore, ts, router, apiServer.Config
	}

	token := utils.GenerateTestToken(1, "JohnDoe", "johndoe@example.com", utils.PrivateKeyAccess)

	t.Run("it should throw an error when no fields are provided for update", func(t *testing.T) {
		_, ts, router, _ := setupTestServer()
		defer ts.Close()

		emptyPayload := `{}`
		req := httptest.NewRequest(http.MethodPut, ts.URL+"/api/v1/users", bytes.NewBufferString(emptyPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error": ["Field validation for 'Username' failed on the 'required' tag", "Field validation for 'Email' failed on the 'required' tag"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should throw an error when body is invalid", func(t *testing.T) {
		_, ts, router, _ := setupTestServer()
		defer ts.Close()

		invalidPayload := `{"username": "", "email": ""}`
		req := httptest.NewRequest(http.MethodPut, ts.URL+"/api/v1/users", bytes.NewBufferString(invalidPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error":["Field validation for 'Username' failed on the 'required' tag", "Field validation for 'Email' failed on the 'required' tag"]}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should return error when the request context is canceled", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockUserStore.On("UpdateByID", mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Err() == context.Canceled
		}), mock.Anything, mock.Anything).Return(&types.UserResponse{}, context.Canceled)

		validPayload := `{
			"username": "johndoe - updated",
			"email": "johndoeupdated@email.com"
		}`

		req := httptest.NewRequest(http.MethodPut, ts.URL+"/api/v1/users", bytes.NewBufferString(validPayload)).WithContext(canceledCtx)
		req.Header.Set("Authorization", "Bearer "+token)
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

	t.Run("it should throw a database connection error", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("UpdateByID", mock.Anything, mock.Anything).Return(&types.UserResponse{}, sql.ErrConnDone)

		validPayload := `{
			"username": "johndoe - updated",
			"email": "johndoeupdated@email.com"
		}`

		req := httptest.NewRequest(http.MethodPut, ts.URL+"/api/v1/users", bytes.NewBufferString(validPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)
		expected := `{"error":"An unexpected error occurred"}`
		assert.JSONEq(t, expected, string(responseBody))
	})

	t.Run("it should throw an error when call endpoint with a non-existent user ID", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("UpdateByID", mock.Anything, mock.Anything).Return(&types.UserResponse{}, sql.ErrNoRows)

		validPayload := `{
			"username": "johndoe - updated",
			"email": "johndoeupdated@email.com"
		}`
		req := httptest.NewRequest(http.MethodPut, ts.URL+"/api/v1/users", bytes.NewBufferString(validPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error": "No user found with ID 1"}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should return successfully status and body when the user is updated", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockedDate := time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC)

		mockUserStore.On("UpdateByID", mock.Anything, mock.Anything).Return(&types.UserResponse{
			ID:        1,
			Username:  "johndoe - updated",
			Email:     "johndoeupdated@email.com",
			CreatedAt: time.Date(0001, 01, 01, 0, 0, 0, 0, time.UTC),
			DeletedAt: nil,
			UpdatedAt: &mockedDate,
		}, nil)

		validPayload := `{
			"username": "johndoe - updated",
			"email": "johndoeupdated@email.com"
		}`
		req := httptest.NewRequest(http.MethodPut, ts.URL+"/api/v1/users", bytes.NewBufferString(validPayload))
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{
			"id": 1,
			"username":  "johndoe - updated",
			"email": "johndoeupdated@email.com",
			"createdAt": "0001-01-01T00:00:00Z",
			"updatedAt": "0001-01-01T00:00:00Z",
			"deletedAt": null
		}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})
}

func TestHandleDeleteUser(t *testing.T) {
	setupTestServer := func() (*mocks.MockUserStore, *httptest.Server, *mux.Router, config.Config) {
		mockUserStore := new(mocks.MockUserStore)
		mockUserHandler := user.NewUserHandler(mockUserStore)
		apiServer := api.NewApiServer(":8080", nil)
		router := apiServer.SetupRouter(nil, mockUserHandler, nil)
		ts := httptest.NewServer(router)
		return mockUserStore, ts, router, apiServer.Config
	}

	token := utils.GenerateTestToken(1, "JohnDoe", "johndoe@example.com", utils.PrivateKeyAccess)

	t.Run("it should return error when the request context is canceled", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		canceledCtx, cancel := context.WithCancel(context.Background())
		cancel()

		mockUserStore.On("DeleteByID", mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Err() == context.Canceled
		}), 1).Return(context.Canceled)

		req := httptest.NewRequest(http.MethodDelete, ts.URL+"/api/v1/users", nil).WithContext(canceledCtx)
		req.Header.Set("Authorization", "Bearer "+token)
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

	t.Run("it should throw a database connection error", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("DeleteByID", mock.Anything, mock.Anything).Return(sql.ErrConnDone)

		req := httptest.NewRequest(http.MethodDelete, ts.URL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)

		responseBody, _ := io.ReadAll(res.Body)
		expected := `{"error":"An unexpected error occurred"}`
		assert.JSONEq(t, expected, string(responseBody))
	})

	t.Run("it should throw an error when call endpoint with a non-existent user ID", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("DeleteByID", mock.Anything, mock.Anything).Return(sql.ErrNoRows)

		req := httptest.NewRequest(http.MethodDelete, ts.URL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNotFound, res.StatusCode)

		responseBody, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		expectedResponse := `{"error": "No user found with ID 1"}`
		assert.JSONEq(t, expectedResponse, string(responseBody))
	})

	t.Run("it should return succssefully status and body when call endpoint with valid body", func(t *testing.T) {
		mockUserStore, ts, router, _ := setupTestServer()
		defer ts.Close()

		mockUserStore.On("DeleteByID", mock.Anything, mock.Anything).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, ts.URL+"/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		res := w.Result()
		defer res.Body.Close()

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
	})
}
