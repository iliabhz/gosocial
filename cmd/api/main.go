package main

import (
	"log"

	"github.com/iliabhz/tiago-backend/internal/env"
	"github.com/iliabhz/tiago-backend/internal/store/postgres"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := newConfig(env.GetString("ADDR", ":8080"))
	store := postgres.NewPostgresStorage(nil)

	app := newApplication(cfg, store)

	mux := app.mount()
	log.Fatal(app.run(mux))
}
