package contacts

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hoyci/ms-chat/contacts-service/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateContact_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)

	contactID := "contact1"
	ownerID := "owner1"

	mock.ExpectQuery("INSERT INTO contacts").
		WithArgs(ownerID, contactID, types.StatusPending, sqlmock.AnyArg()).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id", "owner_id", "contact_id", "status", "created_at", "updated_at", "deleted_at",
				},
			).
				AddRow(1, ownerID, contactID, types.StatusPending, time.Now(), time.Now(), nil),
		)

	contact, err := store.CreateContact(context.Background(), contactID, ownerID)
	assert.NoError(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, ownerID, contact.OwnerID)
	assert.Equal(t, contactID, contact.ContactID)
}

func TestCreateContact_Failure(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)

	contactID := "contact1"
	ownerID := "owner1"

	mock.ExpectQuery("INSERT INTO contacts").
		WithArgs(ownerID, contactID, types.StatusPending, sqlmock.AnyArg()).
		WillReturnError(errors.New("insert error"))

	contact, err := store.CreateContact(context.Background(), contactID, ownerID)
	assert.Error(t, err)
	assert.Nil(t, contact)
}

func TestGetContactByOwnerID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)
	ownerID := "owner1"
	contactID := "contact1"

	mock.ExpectQuery("SELECT \\* FROM contacts WHERE owner_id = \\$1 AND contact_id = \\$2 AND deleted_at IS null").
		WithArgs(ownerID, contactID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id", "owner_id", "contact_id", "status", "created_at", "updated_at", "deleted_at",
				},
			).
				AddRow(1, ownerID, contactID, types.StatusPending, time.Now(), time.Now(), nil),
		)

	contact, err := store.GetContactByOwnerID(context.Background(), ownerID, contactID)
	assert.NoError(t, err)
	assert.NotNil(t, contact)
	assert.Equal(t, ownerID, contact.OwnerID)
	assert.Equal(t, contactID, contact.ContactID)
}

func TestGetContactByOwnerID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)
	ownerID := "owner1"
	contactID := "contact1"

	mock.ExpectQuery("SELECT \\* FROM contacts WHERE owner_id = \\$1 AND contact_id = \\$2 AND deleted_at IS null").
		WithArgs(ownerID, contactID).
		WillReturnError(sql.ErrNoRows)

	contact, err := store.GetContactByOwnerID(context.Background(), ownerID, contactID)
	assert.Error(t, err)
	assert.Nil(t, contact)
}

func TestGetAllContactsByOwnerID_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)
	ownerID := "owner1"

	mock.ExpectQuery("SELECT id, owner_id, contact_id, status, created_at, updated_at, deleted_at FROM contacts WHERE owner_id = \\$1 AND deleted_at IS null").
		WithArgs(ownerID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id", "owner_id", "contact_id", "status", "created_at", "updated_at", "deleted_at",
				},
			).
				AddRow(1, ownerID, "contact1", types.StatusPending, time.Now(), time.Now(), nil).
				AddRow(2, ownerID, "contact2", types.StatusPending, time.Now(), time.Now(), nil),
		)

	contacts, err := store.GetAllContactsByOwnerID(context.Background(), ownerID)
	assert.NoError(t, err)
	assert.Len(t, contacts, 2)
}

func TestGetAllContactsByOwnerID_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)
	ownerID := "owner1"

	mock.ExpectQuery("SELECT id, owner_id, contact_id, status, created_at, updated_at, deleted_at FROM contacts WHERE owner_id = \\$1 AND deleted_at IS null").
		WithArgs(ownerID).
		WillReturnRows(
			sqlmock.NewRows(
				[]string{
					"id", "owner_id", "contact_id", "status", "created_at", "updated_at", "deleted_at",
				},
			),
		)

	contacts, err := store.GetAllContactsByOwnerID(context.Background(), ownerID)
	assert.NoError(t, err)
	assert.Empty(t, contacts)
}

func TestGetAllContactsByOwnerID_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)
	ownerID := "owner1"

	mock.ExpectQuery("SELECT id, owner_id, contact_id, status, created_at, updated_at, deleted_at FROM contacts WHERE owner_id = \\$1 AND deleted_at IS null").
		WithArgs(ownerID).
		WillReturnError(errors.New("query error"))

	contacts, err := store.GetAllContactsByOwnerID(context.Background(), ownerID)
	assert.Error(t, err)
	assert.Nil(t, contacts)
}

func TestGetAllContactsByOwnerID_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)
	ownerID := "owner1"

	mock.ExpectQuery("SELECT id, owner_id, contact_id, status, created_at, updated_at, deleted_at FROM contacts WHERE owner_id = \\$1 AND deleted_at IS null").
		WithArgs(ownerID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "owner_id"}).
				AddRow("texto_invalido", "owner1"),
		)

	contacts, err := store.GetAllContactsByOwnerID(context.Background(), ownerID)
	assert.Error(t, err)
	assert.Nil(t, contacts)
}

func TestGetAllContactsByOwnerID_RowsError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	store := NewContactStore(db)
	ownerID := "owner1"

	rows := sqlmock.NewRows(
		[]string{
			"id", "owner_id", "contact_id", "status", "created_at", "updated_at", "deleted_at",
		},
	).
		AddRow(1, ownerID, "contact1", types.StatusPending, time.Now(), time.Now(), nil).
		AddRow(2, ownerID, "contact2", types.StatusPending, time.Now(), time.Now(), nil).
		RowError(1, errors.New("rows error"))

	mock.ExpectQuery("SELECT id, owner_id, contact_id, status, created_at, updated_at, deleted_at FROM contacts WHERE owner_id = \\$1 AND deleted_at IS null").
		WithArgs(ownerID).
		WillReturnRows(rows)

	contacts, err := store.GetAllContactsByOwnerID(context.Background(), ownerID)
	assert.Error(t, err)
	assert.Nil(t, contacts)
}
