package main

import (
	"log"

	"github.com/iliabhz/gosocial/db"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

//	@title			GoSocial API
//	@description	API for GoSocial which is a social network built in go
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
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

	app.mount(appCfg)
	log.Fatal(app.run(appCfg))
}
