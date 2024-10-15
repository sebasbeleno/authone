package main

import "net/http"

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: LOGGER

	app.writeJsonError(w, http.StatusBadRequest, err)
}

func (app *application) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.writeJsonError(w, http.StatusInternalServerError, err)
}
