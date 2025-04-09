package contacts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/hoyci/ms-chat/contacts-service/types"
	coreTypes "github.com/hoyci/ms-chat/core/types"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"log"
	"net/http"
)

var validate = validator.New()

type ContactHandler struct {
	contactStore types.ContactStore
}

func NewContactHandler(contactStore types.ContactStore) *ContactHandler {
	return &ContactHandler{contactStore: contactStore}
}

func (h *ContactHandler) HandleCreateContact(w http.ResponseWriter, r *http.Request) {
	userID, ok := coreUtils.GetClaimFromContext[string](r, "UserID")
	if !ok {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve userID from context"),
			"HandleCreateContact", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	var requestPayload types.CreateContactPayload
	if err := coreUtils.ParseJSON(r, &requestPayload); err != nil {
		coreUtils.WriteError(
			w, http.StatusBadRequest, err, "HandleCreateContact",
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
			w, http.StatusBadRequest, err, "HandleCreateContact",
			coreTypes.BadRequestStructResponse{Error: errorMessages},
		)
		return
	}

	contact, _ := h.contactStore.GetContactByOwnerID(r.Context(), requestPayload.ContactID, userID)

	if contact != nil {
		coreUtils.WriteError(
			w, http.StatusConflict, fmt.Errorf("you have already registered this contact"), "HandleCreateContact",
			coreTypes.InternalServerErrorResponse{Error: "You have already registered this contact"},
		)
		return
	}

	contact, err := h.contactStore.CreateContact(r.Context(), requestPayload.ContactID, userID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleCreateContact",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}

		if errors.Is(err, sql.ErrConnDone) {
			log.Printf("Database connection error: %v", err)
			coreUtils.WriteError(
				w, http.StatusInternalServerError, err, "HandleCreateContact",
				coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
			)
			return
		}

		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleCreateContact",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	_ = coreUtils.WriteJSON(
		w, http.StatusCreated, types.CreateContactResponse{Contact: contact},
	)
}

func (h *ContactHandler) HandleGetContacts(w http.ResponseWriter, r *http.Request) {
	userID, ok := coreUtils.GetClaimFromContext[string](r, "UserID")
	if !ok {
		coreUtils.WriteError(
			w, http.StatusInternalServerError, fmt.Errorf("failed to retrieve userID from context"),
			"HandleGetContacts", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	contacts, err := h.contactStore.GetAllContactsByOwnerID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(
				w, http.StatusServiceUnavailable, err, "HandleGetContacts",
				coreTypes.ContextCanceledResponse{Error: "Request canceled"},
			)
			return
		}

		if errors.Is(err, sql.ErrConnDone) {
			log.Printf("Database connection error: %v", err)
			coreUtils.WriteError(
				w, http.StatusInternalServerError, err, "HandleGetContacts",
				coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
			)
			return
		}

		coreUtils.WriteError(
			w, http.StatusInternalServerError, err, "HandleGetContacts",
			coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"},
		)
		return
	}

	_ = coreUtils.WriteJSON(
		w, http.StatusOK, types.GetContactResponse{Contact: contacts},
	)
}
