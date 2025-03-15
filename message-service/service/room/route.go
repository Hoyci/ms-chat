package room

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	coreTypes "github.com/hoyci/ms-chat/core/types"
	coreUtils "github.com/hoyci/ms-chat/core/utils"
	"github.com/hoyci/ms-chat/message-service/types"
)

var validate = validator.New()

type RoomHandler struct {
	roomStore *RoomStore
}

func NewRoomHandler(roomStore *RoomStore) *RoomHandler {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})
	return &RoomHandler{roomStore: roomStore}
}

func (h *RoomHandler) HandleCreateRoom(w http.ResponseWriter, r *http.Request) {
	var requestPayload types.CreateRoomPayload
	if err := coreUtils.ParseJSON(r, &requestPayload); err != nil {
		coreUtils.WriteError(w, http.StatusBadRequest, err, "HandleCreateRoom", coreTypes.BadRequestResponse{Error: "Body is not a valid json"})
		return
	}

	if err := validate.Struct(requestPayload); err != nil {
		var errorMessages []string
		for _, e := range err.(validator.ValidationErrors) {
			errorMessages = append(errorMessages, fmt.Sprintf("Field '%s' is invalid: %s", e.Field(), e.Tag()))
		}

		coreUtils.WriteError(w, http.StatusBadRequest, err, "HandleCreateRoom", coreTypes.BadRequestStructResponse{Error: errorMessages})
		return
	}

	roomID, err := h.roomStore.Create(r.Context(), requestPayload)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(w, http.StatusServiceUnavailable, err, "HandleCreateRoom", coreTypes.ContextCanceledResponse{Error: "Request canceled"})
			return
		}

		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleCreateRoom", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
		return
	}

	coreUtils.WriteJSON(w, http.StatusCreated, types.CreateRoomResponse{Message: "Room successfully created", RoomID: roomID.Hex()})
}

func (h *RoomHandler) HandleGetRoomByID(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		coreUtils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing room_id on query params"), "HandleCreateRoom", coreTypes.BadRequestResponse{Error: "Missing room_id on query params"})
	}

	room, err := h.roomStore.GetByID(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			coreUtils.WriteError(w, http.StatusServiceUnavailable, err, "HandleCreateRoom", coreTypes.ContextCanceledResponse{Error: "Request canceled"})
			return
		}

		coreUtils.WriteError(w, http.StatusInternalServerError, err, "HandleCreateRoom", coreTypes.InternalServerErrorResponse{Error: "An unexpected error occurred"})
		return
	}

	coreUtils.WriteJSON(w, http.StatusCreated, room)
}
