package env

import (
	"os"
	"strings"
)

func WithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func AppEnv() string {
	return WithDefault("APP_ENV", "production")
}

func AppKey() string {
	key := os.Getenv("APP_KEY")
	if key == "" {
		panic("APP_KEY is not set")
	}

	return key
}

func IsDevelopment() bool {
	return strings.HasPrefix(AppEnv(), "dev")
}

func IpClientHeaderKey() string {
	return WithDefault("IP_CLIENT_HEADER_KEY", "CF-Connecting-IP")
}
