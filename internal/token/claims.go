package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserCLaims struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewUserClaims(email string, userID uuid.UUID, duration time.Duration) (*UserCLaims, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &UserCLaims{
		Email:  email,
		UserID: userID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId.String(),
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
