package contacts

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/hoyci/ms-chat/contacts-service/mocks"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hoyci/ms-chat/contacts-service/types"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateContact_Success(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	contactId := uuid.New().String()
	userID := uuid.New().String()

	payload := fmt.Sprintf("{\"contact_id\": \"%s\"}", contactId)
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetContactByOwnerID(gomock.Any(), contactId, userID).Return(nil, nil)
	mockStore.EXPECT().CreateContact(gomock.Any(), contactId, userID)

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandleCreateContact_MissingUserID(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	payload := `{"contact_id": "123"}`
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreateContact_InvalidJSON(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	payload := `{"contact_id":`
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "123")
	w := httptest.NewRecorder()

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreateContact_ValidationError(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	payload := `{"contact_id": "456"}`
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", "123")
	w := httptest.NewRecorder()

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleCreateContact_ContactAlreadyExists(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	contactId := uuid.New().String()
	userID := uuid.New().String()

	payload := fmt.Sprintf("{\"contact_id\": \"%s\"}", contactId)
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetContactByOwnerID(gomock.Any(), contactId, userID).Return(
		&types.Contact{ID: "123", OwnerID: "456"}, nil,
	)

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestHandleCreateContact_ContextCanceled(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	contactId := uuid.New().String()
	userID := uuid.New().String()

	payload := fmt.Sprintf("{\"contact_id\": \"%s\"}", contactId)
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetContactByOwnerID(gomock.Any(), contactId, userID).Return(nil, nil)
	mockStore.EXPECT().CreateContact(gomock.Any(), contactId, userID).Return(nil, context.Canceled)

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestHandleCreateContact_DatabaseError(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	contactId := uuid.New().String()
	userID := uuid.New().String()

	payload := fmt.Sprintf("{\"contact_id\": \"%s\"}", contactId)
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetContactByOwnerID(gomock.Any(), contactId, userID).Return(nil, nil)
	mockStore.EXPECT().CreateContact(gomock.Any(), contactId, userID).Return(nil, sql.ErrConnDone)

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleCreateContact_InternalServerError(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	contactId := uuid.New().String()
	userID := uuid.New().String()

	payload := fmt.Sprintf("{\"contact_id\": \"%s\"}", contactId)
	req := httptest.NewRequest(http.MethodPost, "/contacts", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetContactByOwnerID(gomock.Any(), contactId, userID).Return(nil, nil)
	mockStore.EXPECT().CreateContact(gomock.Any(), contactId, userID).Return(
		nil, fmt.Errorf("an unexpected error occurred"),
	)

	handler.HandleCreateContact(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleGetContacts_Success(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	userID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetAllContactsByOwnerID(gomock.Any(), userID).Return(
		[]*types.Contact{
			{
				ID:        "123",
				OwnerID:   userID,
				ContactID: uuid.New().String(),
				Status:    "accepted",
				CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				UpdatedAt: nil,
				DeletedAt: nil,
			},
			{
				ID:        "123",
				OwnerID:   userID,
				ContactID: uuid.New().String(),
				Status:    "accepted",
				CreatedAt: time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC),
				UpdatedAt: nil,
				DeletedAt: nil,
			},
		}, nil,
	)

	handler.HandleGetContacts(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandleGetContacts_MissingUserID(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	w := httptest.NewRecorder()

	handler.HandleGetContacts(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandleGetContacts_ContextCanceled(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	userID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetAllContactsByOwnerID(gomock.Any(), userID).Return(nil, context.Canceled)

	handler.HandleGetContacts(w, req)

	assert.Equal(t, http.StatusServiceUnavailable, w.Code)
}

func TestHandleGetContacts_DatabaseError(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	userID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetAllContactsByOwnerID(gomock.Any(), userID).Return(nil, sql.ErrConnDone)

	handler.HandleGetContacts(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandleGetContacts_InternalServerError(t *testing.T) {
	coreUtils.InitLogger()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mocks.NewMockContactStore(ctrl)
	handler := NewContactHandler(mockStore)

	userID := uuid.New().String()
	req := httptest.NewRequest(http.MethodGet, "/contacts", nil)
	req.Header.Set("X-User-ID", userID)
	w := httptest.NewRecorder()

	mockStore.EXPECT().GetAllContactsByOwnerID(gomock.Any(), userID).Return(
		nil, fmt.Errorf("an unexpected error occurred"),
	)

	handler.HandleGetContacts(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
