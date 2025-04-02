package contacts

import (
	"context"
	"database/sql"
	"time"

	"github.com/hoyci/ms-chat/contacts-service/types"
)

type ContactStore struct {
	db *sql.DB
}

func NewContactStore(db *sql.DB) *ContactStore {
	return &ContactStore{db: db}
}

func (s *ContactStore) CreateContact(ctx context.Context, contactID string, ownerID string) (
	*types.Contact, error,
) {
	contact := &types.Contact{}
	err := s.db.QueryRowContext(
		ctx,
		"INSERT INTO contacts (owner_id, contact_id, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id, owner_id, contact_id, status, created_at, updated_at, deleted_at",
		ownerID,
		contactID,
		types.StatusPending,
		time.Now(),
	).Scan(
		&contact.ID,
		&contact.OwnerID,
		&contact.ContactID,
		&contact.Status,
		&contact.CreatedAt,
		&contact.UpdatedAt,
		&contact.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactStore) GetContactByOwnerID(ctx context.Context, ownerID string, contactID string) (
	*types.Contact, error,
) {
	contact := &types.Contact{}
	err := s.db.QueryRowContext(
		ctx, "SELECT *  FROM contacts WHERE owner_id = $1 AND contact_id = $2 AND deleted_at IS null", ownerID,
		contactID,
	).
		Scan(
			&contact.ID,
			&contact.OwnerID,
			&contact.ContactID,
			&contact.Status,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.DeletedAt,
		)
	if err != nil {
		return nil, err
	}

	return contact, nil
}

func (s *ContactStore) GetAllContactsByOwnerID(ctx context.Context, ownerID string) ([]*types.Contact, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, owner_id, contact_id, status, created_at, updated_at, deleted_at FROM contacts WHERE owner_id = $1 AND deleted_at IS null",
		ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*types.Contact
	for rows.Next() {
		contact := &types.Contact{}
		if err := rows.Scan(
			&contact.ID,
			&contact.OwnerID,
			&contact.ContactID,
			&contact.Status,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.DeletedAt,
		); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}
