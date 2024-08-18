package helpers

import (
	"fmt"
	"net/http"
)

type DatabaseError struct {
	Type      string      `json:"type"`
	ErrorData interface{} `json:"error,omitempty"`
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.ErrorData)
}

type HTTPError struct {
	StatusCode      int         `json:"statusCode"`
	Method          string      `json:"method"`
	Message         string      `json:"message"`
	ErrorData       interface{} `json:"error,omitempty"`
	InternalError   error       `json:"-"`
	InternalMessage string      `json:"-"`
	ErrorID         string      `json:"error_id,omitempty"`
}

func (e HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
}

func (e *HTTPError) Is(target error) bool {
	return e.Error() == target.Error()
}

// Cause returns the root cause error
func (e *HTTPError) Cause() error {
	if e.InternalError != nil {
		return e.InternalError
	}
	return e
}

// WithInternalError adds internal error information to the error
func (e *HTTPError) WithInternalError(err error) *HTTPError {
	e.InternalError = err
	return e
}

// WithInternalMessage adds internal message information to the error
func (e *HTTPError) WithInternalMessage(fmtString string, args ...interface{}) *HTTPError {
	e.InternalMessage = fmt.Sprintf(fmtString, args...)
	return e
}

func httpError(statusCode int, method string, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Method:     method,
		Message:    fmt.Sprintf(fmtString, args...),
	}
}

func httpValidationError(statusCode int, method string, fmtString string, errorData interface{}, args ...interface{}) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Method:     method,
		ErrorData:  errorData,
		Message:    fmt.Sprintf(fmtString, args...),
	}
}

func ValidationError(method string, fmtString string, errorData interface{}, args ...interface{}) *HTTPError {
	return httpValidationError(http.StatusForbidden, method, fmtString, errorData, args...)
}

func UnauthorizedError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnauthorized, method, fmtString, args...)
}

func BadRequestError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, method, fmtString, args...)
}

func InternalServerError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusInternalServerError, method, fmtString, args...)
}

func NotFoundError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusNotFound, method, fmtString, args...)
}

func ExpiredTokenError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnauthorized, method, fmtString, args...)
}

func ForbiddenError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusForbidden, method, fmtString, args...)
}

func UnprocessableEntityError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnprocessableEntity, method, fmtString, args...)
}

func TooManyRequestsError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusTooManyRequests, method, fmtString, args...)
}

func ConflictError(method string, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusConflict, method, fmtString, args...)
}
