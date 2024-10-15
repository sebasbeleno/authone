package main

import (
	"net/http"
	"time"

	"github.com/sebasbeleno/authone/internal/store"
)

type SignUpUserPayload struct {
	EmailAddress string `json:"email" validate:"required,email,max=255"`
	Password     string `json:"password" validate:required`
}

func (app *application) signUpUserWithEmailAddress(w http.ResponseWriter, r *http.Request) {
	var payload SignUpUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.writeJsonError(w, http.StatusBadRequest, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.writeJsonError(w, http.StatusBadRequest, err)
		return
	}

	user := &store.User{
		EmailAddress: payload.EmailAddress,
	}

	if err := user.PasswordHash.Set(payload.Password); err != nil {
		app.writeJsonError(w, http.StatusInternalServerError, err)
		return
	}

	err := app.store.Users.Create(r.Context(), user)

	if err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerErrorResponse(w, r, err)
		}

		return
	}

	// return the user
	app.jsonResponse(w, http.StatusCreated, user)
}

type SignInUseWithEmailAndPasswordPayload struct {
	EmailAddress string `json:"email" validate:"required,email,max=255"`
	Password     string `json:"password" validate:"required"`
}

type SignInResponse struct {
	AccessToken           string      `json:"access_token"`
	TokenType             string      `json:"token_type"`
	RereshToken           string      `json:"refresh_token"`
	AccessTokenExpiresAt  time.Time   `json:"access_token_expires_at"`
	RefreshTokenExpiresAt time.Time   `json:"refresh_token_expires_at"`
	SessionId             string      `json:"session_id"`
	User                  *store.User `json:"user"`
}

func (app *application) signInUserWithEmailAndPassword(w http.ResponseWriter, r *http.Request) {
	var payload SignInUseWithEmailAndPasswordPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.writeJsonError(w, http.StatusBadRequest, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.writeJsonError(w, http.StatusBadRequest, err)
		return
	}

	user, err := app.store.Users.GetUserWithEmail(r.Context(), payload.EmailAddress)

	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	if err := user.PasswordHash.Compare(payload.Password); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	accessToken, accessClaims, err := app.config.tokenMaker.GenerateToken(user.EmailAddress, user.UserId, time.Minute*15)

	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	refreshToken, refreshClaims, err := app.config.tokenMaker.GenerateToken(user.EmailAddress, user.UserId, time.Hour*24*7)

	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	session, err := app.store.Sessions.Create(r.Context(), &store.Session{
		UserEmailAddress: user.EmailAddress,
		ExpiryTime:       refreshClaims.ExpiresAt.Time,
		IsRevoked:        false,
		RefreshToken:     refreshToken,
	})
	if err != nil {
		app.internalServerErrorResponse(w, r, err)
		return
	}

	response := &SignInResponse{
		AccessToken:           accessToken,
		TokenType:             "Bearer",
		User:                  user,
		RereshToken:           refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
		SessionId:             session.SessionId.String(),
	}

	// return the user
	app.jsonResponse(w, http.StatusOK, response)
}
