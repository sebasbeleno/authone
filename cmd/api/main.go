package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/sebasbeleno/authone/internal/db"
	"github.com/sebasbeleno/authone/internal/env"
	"github.com/sebasbeleno/authone/internal/store"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgresql://postgres:3161697mi@localhost:5432/authone?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIddleTime: env.GetString("DB_MAX_IDLE_TIME", "15min"),
		},
	}

	db, err := db.NewDB(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIddleTime)

	if err != nil {
		log.Panic(err)
	}

	store := store.NewStore(db)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()

	app.run(mux)

}
