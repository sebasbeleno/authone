package main

import (
	"net/http"

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
		// handle errors messages
		app.writeJsonError(w, http.StatusInternalServerError, err)
		return
	}

	// return the user
	app.jsonResponse(w, http.StatusCreated, user)
}
