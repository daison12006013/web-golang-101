package env

import (
	"os"
	"testing"
)

func TestWithDefault(t *testing.T) {
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	if WithDefault("TEST_KEY", "default_value") != "test_value" {
		t.Errorf("WithDefault() did not return the correct value")
	}

	if WithDefault("NON_EXISTENT_KEY", "default_value") != "default_value" {
		t.Errorf("WithDefault() did not return the default value for non-existent key")
	}
}

func TestAppEnv(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	if AppEnv() != "development" {
		t.Errorf("AppEnv() did not return the correct value")
	}
}

func TestAppKey(t *testing.T) {
	os.Setenv("APP_KEY", "test_key")
	defer os.Unsetenv("APP_KEY")

	if AppKey() != "test_key" {
		t.Errorf("AppKey() did not return the correct value")
	}
}

func TestIsDevelopment(t *testing.T) {
	os.Setenv("APP_ENV", "development")
	defer os.Unsetenv("APP_ENV")

	if !IsDevelopment() {
		t.Errorf("IsDevelopment() did not return true for development environment")
	}

	os.Setenv("APP_ENV", "production")
	if IsDevelopment() {
		t.Errorf("IsDevelopment() did not return false for non-development environment")
	}
}

func TestIpClientHeaderKey(t *testing.T) {
	os.Setenv("IP_CLIENT_HEADER_KEY", "test_key")
	defer os.Unsetenv("IP_CLIENT_HEADER_KEY")

	if IpClientHeaderKey() != "test_key" {
		t.Errorf("IpClientHeaderKey() did not return the correct value")
	}
}
