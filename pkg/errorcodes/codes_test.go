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

func TestOriginalError(t *testing.T) {
	err := errors.New("original error")
	e := &Error{Err: err}
	if e.OriginalError() != err {
		t.Errorf("OriginalError() = %v, want %v", e.OriginalError(), err)
	}
}

func TestError(t *testing.T) {
	status := 400
	message := "error message"
	err := errors.New("original error")
	e := &Error{HTTPStatus: &status, Message: &message, Err: err}
	want := "400 | error message | original error"
	if e.Error().Error() != want {
		t.Errorf("Error() = %v, want %v", e.Error().Error(), want)
	}
}

func TestErrorMethod(t *testing.T) {
	err := errors.New("original error")
	e := &Error{Err: err}

	result := e.Error()
	if !errors.Is(result, err) {
		t.Errorf("Error() = %v, want %v", result, err)
	}
}

func TestAsDefaultError(t *testing.T) {
	err := errors.New("original error")
	e := AsDefaultError(err)
	if *e.Message != "Default Error" {
		t.Errorf("AsDefaultError() = %v, want %v", *e.Message, "Default Error")
	}
}

// Repeat the pattern for other error types
func TestAsDatabaseConnection(t *testing.T) {
	err := errors.New("original error")
	e := AsDatabaseConnection(err)
	if *e.Message != "Database Connection" {
		t.Errorf("AsDatabaseConnection() = %v, want %v", *e.Message, "Database Connection")
	}
}

func TestAsQueryError(t *testing.T) {
	err := errors.New("original error")
	e := AsQueryError(err)
	if *e.Message != "Query Error" {
		t.Errorf("AsQueryError() = %v, want %v", *e.Message, "Query Error")
	}
}

func TestAsRecordNotFound(t *testing.T) {
	err := errors.New("original error")
	e := AsRecordNotFound(err)
	if *e.Message != "Record Not Found" {
		t.Errorf("AsRecordNotFound() = %v, want %v", *e.Message, "Record Not Found")
	}
}

func TestAsPermissionError(t *testing.T) {
	err := errors.New("original error")
	e := AsPermissionError(err)
	if *e.Message != "Permission Error" {
		t.Errorf("AsPermissionError() = %v, want %v", *e.Message, "Permission Error")
	}
}

func TestBadRequestError(t *testing.T) {
	err := errors.New("original error")
	e := AsBadRequest(err)
	if *e.Message != "Bad Request" {
		t.Errorf("BadRequest() = %v, want %v", *e.Message, "Bad Request")
	}
}

func TestConflictError(t *testing.T) {
	err := errors.New("original error")
	e := AsConflict(err)
	if *e.Message != "Conflict" {
		t.Errorf("Conflict() = %v, want %v", *e.Message, "Conflict")
	}
}
