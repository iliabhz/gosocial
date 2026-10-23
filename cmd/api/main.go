package main

import (
	"log"

	"github.com/iliabhz/gosocial/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appCfg := newAppConfig()
	initAppConfig(appCfg)
	dbCfg := newDBConfig()
	initDBConfig(dbCfg)

	store, err := db.NewDB(dbCfg.dsn)
	if err != nil {
		log.Fatal(err)
	}

	defer store.Close()

	log.Println("Database connection pool established")

	app := newApplication(store)

	mux := app.mount()
	log.Fatal(app.run(appCfg, mux))
}
