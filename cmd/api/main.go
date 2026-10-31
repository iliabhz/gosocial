package main

import (
	"github.com/iliabhz/gosocial/db"
	"github.com/iliabhz/gosocial/internal/mailer"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
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
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatalf("Error loading .env file: %s", err.Error())
	}

	appCfg := newAppConfig()
	dbCfg := newDBConfig()
	mailCfg := newMailConfig()

	store, err := db.NewDB(dbCfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer store.Close()

	logger.Info("Database connection pool established")

	mailer := mailer.NewMailGunMailer(logger, mailCfg.domain, mailCfg.fromEmail, mailCfg.mailGun.apiKey)

	app := newApplication(appCfg, store, mailer, logger)

	app.mount(appCfg, mailCfg)
	logger.Fatal(app.run(appCfg))
}
