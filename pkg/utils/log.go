package utils

import (
	"os"
	"testing"

	"github.com/daison12006013/web-golang-101/pkg/env"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func DebugTest(t *testing.T, data interface{}, name string) {
	t.Logf("\n\n [%v]:\n\t\t > %+v\n\n", name, data)
}

var isLoggerInit = false

func Logger() *zerolog.Logger {
	if isLoggerInit {
		return &log.Logger
	}

	envLogLevel := env.WithDefault("LOG_APP_LEVEL", "warn")
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(getLogLevel(envLogLevel))

	if os.Getenv("SENTRY_DSN") != "" {
		envSentryLevel := env.WithDefault("LOG_SENTRY_LEVEL", "warn")
		sentryLevel := getLogLevel(envSentryLevel)
		log.Logger = log.Hook(sentryHook{minLevel: sentryLevel})
	}

	isLoggerInit = true
	log.Debug().Msg("Logger initialized")

	return &log.Logger
}

func getLogLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.DebugLevel
	}
}

type sentryHook struct {
	minLevel zerolog.Level
}

func (sh sentryHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level >= sh.minLevel {
		sentry.CaptureMessage(msg)
	}
}
