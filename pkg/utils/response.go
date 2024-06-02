package utils

import (
	"encoding/json"
	"net/http"
	"os"

	ec "web-golang-101/pkg/errorcodes"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	w       http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *Response {
	return &Response{w: w}
}

func (r *Response) WriteSuccessResponse(message string, data interface{}) {
	r.Success = true
	r.Message = message
	r.Data = data
	r.writeResponse(http.StatusOK)
}

func (r *Response) WriteValidationResponse(message string, data interface{}) {
	r.Success = false
	r.Message = message
	r.Data = data
	r.writeResponse(http.StatusBadRequest)
}

func (r *Response) HandleErrorCode(errc *ec.Error) {
	if ok := r.WriteValidationError(errc); ok {
		return
	}

	go func() {
		if errc.HTTPStatus != nil && *errc.HTTPStatus >= 500 {
			Logger().Err(errc.Error()).Msg("Error response")
		}
	}()

	r.WriteErrorResponse(errc.Error().Error(), *errc.HTTPStatus)
}

func (r *Response) WriteErrorResponse(message string, statusCode int) {
	if os.Getenv("APP_DEBUG") == "true" {
		r.Message = message
	} else {
		r.Message = "Internal Server Error"
	}

	r.Success = false
	r.Data = nil
	r.writeResponse(statusCode)
}

func (r *Response) WriteValidationError(errc *ec.Error) bool {
	err := errc.OriginalError()

	if errValidator, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]any)
		english := en.New()
		uni := ut.New(english, english)
		trans, _ := uni.GetTranslator("en")

		for _, err := range errValidator {
			errors[err.Field()] = map[string]any{
				"tag":         err.Tag(),
				"actualtag":   err.ActualTag(),
				"param":       err.Param(),
				"translation": err.Translate(trans),
				// "kind":            err.Kind(),
				// "type":            err.Type(),
				// "value":           err.Value(),
				// "namespace":       err.Namespace(),
				// "structnamespace": err.StructNamespace(),
				// "structfield":     err.StructField(),
			}
		}

		r.WriteValidationResponse("Validation error", errors)
		return true
	}

	return false
}

func (r *Response) writeResponse(statusCode int) {
	js, err := json.Marshal(r)
	if err != nil {
		http.Error(r.w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(statusCode)
	r.w.Write(js)
}
