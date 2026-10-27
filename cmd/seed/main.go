package main

import (
	"log"

	"github.com/iliabhz/gosocial/db"
	"github.com/iliabhz/gosocial/internal/env"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error reading .env file: %s", err.Error())
	}

	dsn := env.GetString("DB_DSN", "postgres://admin:adminpass@localhost/gosocial?sslmode=disable&pool_max_conns=30&pool_max_conn_idle_time=15m")
	store, err := db.NewDB(dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer store.Close()

	db.Seed(store)
}
