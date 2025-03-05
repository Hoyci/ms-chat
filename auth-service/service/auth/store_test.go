package auth

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hoyci/ms-chat/auth-service/types"
	"github.com/stretchr/testify/assert"
)

func TestGetRefreshTokenByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewAuthStore(db)

	t.Run("database did not find any row", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM refresh_tokens WHERE user_id = \$1`).
			WithArgs(1).
			WillReturnError(sql.ErrNoRows)

		refreshToken, err := store.GetRefreshTokenByUserID(context.Background(), 1)

		assert.Nil(t, refreshToken)
		assert.Error(t, err)
		assert.Equal(t, err, sql.ErrNoRows)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("database unexpected error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM refresh_tokens WHERE user_id = \$1`).
			WithArgs(1).
			WillReturnError(fmt.Errorf("database connection error"))

		refreshToken, err := store.GetRefreshTokenByUserID(context.Background(), 1)

		assert.Error(t, err)
		assert.Zero(t, refreshToken)
		assert.NotEqual(t, err, sql.ErrNoRows)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("successfully get user by ID", func(t *testing.T) {
		mock.ExpectQuery(`SELECT \* FROM refresh_tokens WHERE user_id = \$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "jti", "expires_at", "created_at"}).
				AddRow(1, 1, "31a0641b-e109-4467-b78c-13b72d0242a5", time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC)))

		refreshToken, err := store.GetRefreshTokenByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, refreshToken)
		assert.Equal(t, 1, refreshToken.UserID)
		assert.Equal(t, "31a0641b-e109-4467-b78c-13b72d0242a5", refreshToken.Jti)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}

func TestUpdateRefreshTokenByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	store := NewAuthStore(db)
	payload := types.UpdateRefreshTokenPayload{
		UserID:    1,
		Jti:       "31a0641b-e109-4467-b78c-13b72d0242a5",
		ExpiresAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	t.Run("database did not find any row to update", func(t *testing.T) {
		mock.ExpectQuery(`
				INSERT INTO refresh_tokens \(user_id, jti, expires_at\) 
				VALUES \(\$1, \$2, \$3\) 
				ON CONFLICT \(user_id\) 
				DO UPDATE 
				SET jti = EXCLUDED\.jti, expires_at = EXCLUDED\.expires_at 
				RETURNING id, user_id, jti, created_at, expires_at
			`).
			WithArgs(payload.UserID, payload.Jti, payload.ExpiresAt).
			WillReturnError(sql.ErrNoRows)

		err := store.UpsertRefreshToken(context.Background(), payload)

		assert.Error(t, err)
		assert.Equal(t, err, sql.ErrNoRows)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("database unexpected error", func(t *testing.T) {
		mock.ExpectQuery(`
				INSERT INTO refresh_tokens \(user_id, jti, expires_at\) 
				VALUES \(\$1, \$2, \$3\) 
				ON CONFLICT \(user_id\) 
				DO UPDATE 
				SET jti = EXCLUDED\.jti, expires_at = EXCLUDED\.expires_at 
				RETURNING id, user_id, jti, created_at, expires_at
			`).
			WithArgs(payload.UserID, payload.Jti, payload.ExpiresAt).
			WillReturnError(fmt.Errorf("database connection error"))

		err := store.UpsertRefreshToken(context.Background(), payload)

		assert.Error(t, err)
		assert.NotEqual(t, err, sql.ErrNoRows)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("conflict error on upsert", func(t *testing.T) {
		mock.ExpectQuery(`
				INSERT INTO refresh_tokens \(user_id, jti, expires_at\) 
				VALUES \(\$1, \$2, \$3\) 
				ON CONFLICT \(user_id\) 
				DO UPDATE 
				SET jti = EXCLUDED\.jti, expires_at = EXCLUDED\.expires_at 
				RETURNING id, user_id, jti, created_at, expires_at
			`).
			WithArgs(payload.UserID, payload.Jti, payload.ExpiresAt).
			WillReturnError(fmt.Errorf("unique constraint violation"))

		err := store.UpsertRefreshToken(context.Background(), payload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unique constraint violation")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})

	t.Run("successfully update refresh token", func(t *testing.T) {
		mock.ExpectQuery(`
				INSERT INTO refresh_tokens \(user_id, jti, expires_at\) 
				VALUES \(\$1, \$2, \$3\) 
				ON CONFLICT \(user_id\) 
				DO UPDATE 
				SET jti = EXCLUDED\.jti, expires_at = EXCLUDED\.expires_at 
				RETURNING id, user_id, jti, created_at, expires_at
			`).
			WithArgs(payload.UserID, payload.Jti, payload.ExpiresAt).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "jti", "created_at", "expires_at"}).
				AddRow(1, payload.UserID, payload.Jti, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), payload.ExpiresAt))

		err := store.UpsertRefreshToken(context.Background(), payload)

		assert.NoError(t, err)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet expectations: %v", err)
		}
	})
}
