package errorcodes

import (
	"errors"
	"testing"
)

func TestUse(t *testing.T) {
	// Initialize error codes
	errorCodes = map[string]int{
		"Test error": 400,
	}

	// Create an error
	err := errors.New("This is a test error")

	// Use the Use function to wrap the error
	wrappedErr := Use(err, "Test error")

	// Check if the HTTP status is correct
	if *wrappedErr.HTTPStatus != 400 {
		t.Errorf("HTTP status is incorrect, got: %d, want: %d.", *wrappedErr.HTTPStatus, 400)
	}

	// Check if the message is correct
	if *wrappedErr.Message != "Test error" {
		t.Errorf("Message is incorrect, got: %s, want: %s.", *wrappedErr.Message, "Test error")
	}
}
