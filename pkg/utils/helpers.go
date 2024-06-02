package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	ec "web-golang-101/pkg/errorcodes"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator"
)

func GetEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func InitializeSentry(dsn string) {
	if dsn == "" {
		Logger().Info().Msg("Sentry DSN is empty. Skipping Sentry initialization.")
		return
	}

	Logger().Info().Msg("Initializing Sentry...")
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: true,
		Debug:            true,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func InitializeAppKey(appKey string) {
	if appKey == "" {
		Logger().Info().Msg("Application Key is empty. Skipping initialization.")
		return
	}

	Logger().Info().Msg("Initializing Application Key...")
	os.Setenv("APP_KEY", appKey)
}

func GetHost(r *http.Request) string {
	if r.TLS != nil {
		return fmt.Sprintf("https://%s", r.Host)
	}
	return fmt.Sprintf("http://%s", r.Host)
}

// HandlerFunc is a type that defines a function that takes a byte slice and returns an interface and an error.
type HandlerFunc func([]byte) (interface{}, *ec.Error)

func CommonJsonHandler(r *http.Request, logic HandlerFunc) (interface{}, *ec.Error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ec.AsDefaultError(err)
	}
	return logic(body)
}

// commonHandler is a function that handles common logic for HTTP handlers.
func CommonHandler(w http.ResponseWriter, r *http.Request, msg string, logic HandlerFunc) {
	response := NewResponse(w)
	data, errc := CommonJsonHandler(r, logic)
	if errc != nil {
		response.WriteErrorResponse(*errc.Message, *errc.HTTPStatus)
		return
	}
	response.WriteSuccessResponse(msg, data)
}

func Validator() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})
	return validate
}

func AppEnv() string {
	return GetEnvWithDefault("APP_ENV", "production")
}

func IsDevelopment() bool {
	return strings.HasPrefix(AppEnv(), "dev")
}

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}
