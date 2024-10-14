package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/sebasbeleno/authone/internal/store"
)

type application struct {
	config config
	store  store.Storage
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIddleTime string
}

type config struct {
	addr string
	db   dbConfig
}

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthHandler)

		router.Route("/auth", func(r chi.Router) {
			r.Post("/signup", app.signUpUserWithEmailAddress)
		})
	})

	return router

}

func (app *application) run(mux http.Handler) {
	// Start the server
	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: mux,
	}

	fmt.Printf("server has started at port%s \n", app.config.addr)

	err := srv.ListenAndServe()

	if err != nil {
		// handle error
		fmt.Print("Error starting server")
	}

}
