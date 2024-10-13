package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// TODO: DEFINE THE User STRUCT
type User struct {
	UserId       uuid.UUID
	EmailAddress string
	PasswordHash string
	PasswordSalt string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (user_id, email_address, password_hash, password_salt, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.UserId,
		user.EmailAddress,
		user.PasswordHash,
		user.PasswordSalt,
	).Scan(&user.UserId, &user.CreatedAt)

	if err != nil {
		// TODO: Handle error

		return err
	}

	return nil
}
