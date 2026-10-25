package main

import (
	"log"

	"github.com/iliabhz/gosocial/db"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	appCfg := newAppConfig()
	dbCfg := newDBConfig()

	store, err := db.NewDB(dbCfg.dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer store.Close()

	log.Println("Database connection pool established")

	app := newApplication(appCfg, store)

	app.mount()
	log.Fatal(app.run(appCfg))
}
