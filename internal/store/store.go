package store

import (
	"context"
	"database/sql"
	"time"
)

var (
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	Users interface {
		Create(context.Context, *User) error
		GetUserWithEmail(context.Context, string) (*User, error)
	}
	Sessions interface {
		Create(context.Context, *Session) (*Session, error)
	}
}

func NewStore(db *sql.DB) Storage {
	return Storage{
		Users:    &UserStore{db},
		Sessions: &SessionStore{db},
	}
}
