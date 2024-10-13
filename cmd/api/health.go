package main

import "net/http"

func (app *application) healthHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok"}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.writeJsonError(w, http.StatusInternalServerError, err)
	}
}
