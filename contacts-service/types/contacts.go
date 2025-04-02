package types

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ContactStore interface {
	CreateContact(ctx context.Context, contactID string, ownerID string) (*Contact, error)
	GetAllContactsByOwnerID(ctx context.Context, ownerID string) ([]*Contact, error)
	GetContactByOwnerID(ctx context.Context, contactID string, ownerID string) (*Contact, error)
}

type ContactStatus string

const (
	StatusPending  ContactStatus = "pending"
	StatusAccepted ContactStatus = "accepted"
	StatusRejected ContactStatus = "rejected"
)

func (cs *ContactStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case string(StatusPending), string(StatusAccepted), string(StatusRejected):
		*cs = ContactStatus(s)
		return nil
	default:
		return fmt.Errorf("invalid status: %s", s)
	}
}

type Contact struct {
	ID        string        `json:"id"`
	OwnerID   string        `json:"owner_id"`
	ContactID string        `json:"contact_id"`
	Status    ContactStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	DeletedAt *time.Time    `json:"deletedAt"`
	UpdatedAt *time.Time    `json:"updatedAt"`
}

type CreateContactPayload struct {
	ContactID string `json:"contact_id" validate:"required,uuid4"`
}

type CreateContactResponse struct {
	Contact *Contact `json:"contact"`
}

type GetContactResponse struct {
	Contact []*Contact `json:"contact"`
}
