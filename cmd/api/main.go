package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sebasbeleno/authone/internal/env"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}
	app := &application{
		config: cfg,
	}

	mux := app.mount()

	app.run(mux)

}
