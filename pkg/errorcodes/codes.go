package errorcodes

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed "codes.json"
var errorCodesJSON []byte

type Error struct {
	Err        error
	HTTPStatus *int
	Message    *string
}

var errorCodes map[string]int

func init() {
	// Load error codes from embedded JSON
	json.Unmarshal(errorCodesJSON, &errorCodes)
}

func Use(err error, msg string) *Error {
	e := &Error{Err: err, Message: &msg}
	if httpStatus, ok := errorCodes[msg]; ok {
		e.HTTPStatus = &httpStatus
	}
	return e
}

func (e *Error) OriginalError() error {
	return e.Err
}

func (e *Error) Error() error {
	if e.HTTPStatus != nil && e.Message != nil {
		return fmt.Errorf("%d | %s | %w", *e.HTTPStatus, *e.Message, e.Err)
	}
	return e.Err
}

func AsDefaultError(err error) *Error {
	return Use(err, "Default Error")
}

func AsDatabaseConnection(err error) *Error {
	return Use(err, "Database Connection")
}

func AsQueryError(err error) *Error {
	return Use(err, "Query Error")
}

func AsRecordNotFound(err error) *Error {
	return Use(err, "Record Not Found")
}

func AsPermissionError(err error) *Error {
	return Use(err, "Permission Error")
}

func AsBadRequest(err error) *Error {
	return Use(err, "Bad Request")
}

func AsConflict(err error) *Error {
	return Use(err, "Conflict")
}
