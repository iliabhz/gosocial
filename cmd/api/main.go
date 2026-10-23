package main

import (
	"log"

	"github.com/iliabhz/gosocial/internal/db"
	"github.com/iliabhz/gosocial/internal/env"
	"github.com/iliabhz/gosocial/internal/store/postgres"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appCfg := newAppConfig(env.GetString("ADDR", ":8080"))
	db, err := db.New(env.GetString("DB_DSN", "postgres://admin:adminpass@localhost/gosocial?sslmode=disable&pool_max_conns=30&pool_max_conn_idle_time=15m"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Println("Database connection pool established")

	store := postgres.NewPostgresStorage(db)

	app := newApplication(store)

	mux := app.mount()
	log.Fatal(app.run(appCfg, mux))
}
