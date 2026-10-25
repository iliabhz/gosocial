package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iliabhz/gosocial/db"
	"github.com/joho/godotenv"
)

const version = "0.0.1"

var Validate *validator.Validate

type HTTPError struct {
	msg    string
	status int
}

func NewHTTPError(msg string, status int) *HTTPError {
	return &HTTPError{
		msg:    msg,
		status: status,
	}
}

func (e *HTTPError) Error() string {
	return e.msg
}

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
	Validate.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		if name == "" {
			return f.Name
		}

		return name
	})
}

func validationMessages(err error) map[string]string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	messages := make(map[string]string, len(validationErrors))

	for _, fieldErr := range validationErrors {
		field := fieldErr.Field()

		switch fieldErr.Tag() {
		case "required":
			messages[field] = fmt.Sprintf("%s is required", field)

		case "email":
			messages[field] = "Must be a valid email address"

		case "min":
			messages[field] = fmt.Sprintf(
				"Must be at least %s characters long",
				fieldErr.Param(),
			)

		case "max":
			messages[field] = fmt.Sprintf(
				"Must be at max %s characters long",
				fieldErr.Param(),
			)

		default:
			messages[field] = "Invalid value"
		}
	}

	return messages
}

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
