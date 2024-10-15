package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	SessionId        uuid.UUID `json:"session_id"`
	UserEmailAddress string    `json:"user_id"`
	ExpiryTime       time.Time `json:"expiry_time"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	IsRevoked        bool      `json:"is_revoked"`
	RefreshToken     string    `json:"refresh_token"`
}

type SessionStore struct {
	db *sql.DB
}

func (s *SessionStore) Create(ctx context.Context, session *Session) (*Session, error) {
	query := `
		INSERT INTO auth.sessions (user_email, expires_at, is_revoked, refresh_token)
		VALUES ($1, $2, $3, $4)
		RETURNING session_id, created_at, updated_at, is_revoked, refresh_token, expires_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	sessionCreated := &Session{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		session.UserEmailAddress,
		session.ExpiryTime,
		session.IsRevoked,
		session.RefreshToken,
	).Scan(
		&sessionCreated.SessionId,
		&sessionCreated.CreatedAt,
		&sessionCreated.UpdatedAt,
		&sessionCreated.IsRevoked,
		&sessionCreated.RefreshToken,
		&sessionCreated.ExpiryTime,
	)

	if err != nil {
		return nil, err
		//TODO: handle the error
	}

	return sessionCreated, nil
}
