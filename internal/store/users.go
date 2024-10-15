package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
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

func (p *password) Compare(plainText string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
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
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_address_key"`:
			return ErrDuplicateEmail
		}
	}

	return nil
}

func (s *UserStore) GetUserWithEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT user_id, email_address, password_hash, created_at, updated_at
		FROM auth.users WHERE email_address = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserId,
		&user.EmailAddress,
		&user.PasswordHash.hash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
