package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (maker *JWTMaker) GenerateToken(email string, userID uuid.UUID, duration time.Duration) (string, *UserCLaims, error) {
	claims, err := NewUserClaims(email, userID, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, err
	}

	return tokenStr, claims, nil
}
