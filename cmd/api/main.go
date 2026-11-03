package main

import (
	"github.com/iliabhz/gosocial/db"
	"github.com/iliabhz/gosocial/internal/mailer"
	"github.com/iliabhz/gosocial/internal/store/cache"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const version = "0.0.1"

var JWTAuth *JWTAuthenticator

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
	jwtCfg := newJWTConfig()
	cacheCfg := newCacheConfig()

	if jwtCfg.secret == "" {
		logger.Fatal("No JWT_SECRET was set in the environment variables")
	}

	store, err := db.NewDB(dbCfg.dsn)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Database connection pool established")

	cacheStore, err := cache.NewRedisCacheStorage(cacheCfg.dsn, cacheCfg.userCacheExpiry, cacheCfg.roleCacheExpiry, store, logger)
	if err != nil {
		cacheStore.Close()
		logger.Fatal(err)
	}
	defer cacheStore.Close()

	mailer := mailer.NewMailGunMailer(logger, mailCfg.domain, mailCfg.fromEmail, mailCfg.mailGun.apiKey)
	JWTAuth = NewJWTAuthenticator(jwtCfg.secret, jwtCfg.audience, jwtCfg.issuer, jwtCfg.expiry)

	app := newApplication(appCfg, cacheStore, mailer, logger)

	app.mount(appCfg, mailCfg)
	err = app.run(appCfg)
	if err != nil {
		logger.Fatal(err)
	}
}
