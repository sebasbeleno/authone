package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// TODO: DEFINE THE User STRUCT
type User struct {
	UserId       uuid.UUID `json:"user_id"`
	EmailAddress string    `json:"email_address"`
	PasswordHash password  `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.text = &plainText
	p.hash = hash

	return nil
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO auth.users (email_address, password_hash)
		VALUES ($1, $2)
		RETURNING user_id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.EmailAddress,
		user.PasswordHash.hash,
	).Scan(
		&user.UserId,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}
