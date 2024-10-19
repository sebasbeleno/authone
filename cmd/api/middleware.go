package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/sebasbeleno/authone/internal/token"
)

type authKey struct{}

func GetAuthMiddlewareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)

			if err != nil {
				http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), authKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func verifyClaimsFromAuthHeader(r *http.Request, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("no authorization header")
	}

	fields := strings.Fields(authHeader)
	if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
		return nil, fmt.Errorf("invalid authorization header")
	}

	tokenStr := fields[1]
	claims, err := tokenMaker.VerifyToken(tokenStr)

	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
