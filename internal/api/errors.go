package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/observability"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type DatabaseError struct {
	Type      string      `json:"type"`
	ErrorData interface{} `json:"error,omitempty"`
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.ErrorData)
}

type HTTPError struct {
	HTTPStatus      int         `json:"statusCode"`
	ErrorCode       string      `json:"error_code,omitempty"`
	Message         string      `json:"message"`
	ErrorData       interface{} `json:"error,omitempty"`
	InternalError   error       `json:"-"`
	InternalMessage string      `json:"-"`
	ErrorID         string      `json:"error_id,omitempty"`
}

func (e *HTTPError) Error() string {
	if e.InternalMessage != "" {
		return e.InternalMessage
	}
	return fmt.Sprintf("%d: %s", e.HTTPStatus, e.Message)
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

func httpError(httpStatus int, errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return &HTTPError{
		HTTPStatus: httpStatus,
		ErrorCode:  errorCode,
		Message:    fmt.Sprintf(fmtString, args...),
	}
}

func httpValidationError(httpStatus int, errorCode ErrorCode, fmtString string, errorData interface{}, args ...interface{}) *HTTPError {
	return &HTTPError{
		HTTPStatus: httpStatus,
		ErrorCode:  errorCode,
		ErrorData:  errorData,
		Message:    fmt.Sprintf(fmtString, args...),
	}
}

func ValidationError(fmtString string, errorData interface{}, args ...interface{}) *HTTPError {
	return httpValidationError(http.StatusForbidden, ErrorCodeValidationFailed, fmtString, errorData, args...)
}

func UnauthorizedError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnauthorized, ErrorCodeNoAuthorization, fmtString, args...)
}

func BadRequestError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusBadRequest, errorCode, fmtString, args...)
}

func InternalServerError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusInternalServerError, ErrorCodeUnexpectedFailure, fmtString, args...)
}

func CotFoundError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusNotFound, errorCode, fmtString, args...)
}

func ForbiddenError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusForbidden, errorCode, fmtString, args...)
}

func UnprocessableEntityError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusUnprocessableEntity, errorCode, fmtString, args...)
}

func TooManyRequestsError(errorCode ErrorCode, fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusTooManyRequests, errorCode, fmtString, args...)
}

func ConflictError(fmtString string, args ...interface{}) *HTTPError {
	return httpError(http.StatusConflict, ErrorCodeConflict, fmtString, args...)
}

type ErrorCause interface {
	Cause() error
}

type HTTPErrorResponse20240101 struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func HandleResponseError(ctx *gin.Context, err error) {
	log := observability.GetLogEntry(ctx).Entry
	errorID := utilities.GetRequestID(ctx)

	apiVersion, averr := DetermineClosestAPIVersion(ctx.Request.Header.Get(APIVersionHeaderName))

	if averr != nil {
		log.WithError(averr).Warn("Invalid version passed to " + APIVersionHeaderName + " header, defaulting to initial version")
	} else if apiVersion != APIVersionInitial {
		// Echo back the determined API version from the request
		ctx.Writer.Header().Set(APIVersionHeaderName, FormatAPIVersion(apiVersion))
	}

	switch e := err.(type) {
	case *HTTPError:
		switch {
		case e.HTTPStatus >= http.StatusInternalServerError:
			e.ErrorID = errorID
			// this will get us the stack trace too
			log.WithError(e.Cause()).Error(e.Error())
		case e.HTTPStatus == http.StatusTooManyRequests:
			log.WithError(e.Cause()).Warn(e.Error())
		default:
			log.WithError(e.Cause()).Info(e.Error())
		}

		if e.ErrorCode != "" {
			ctx.Writer.Header().Set("x-sb-error-code", e.ErrorCode)
		}

		if apiVersion.Compare(APIVersion20240101) >= 0 {
			resp := HTTPErrorResponse20240101{
				Code:    e.ErrorCode,
				Message: e.Message,
			}

			if resp.Code == "" {
				if e.HTTPStatus == http.StatusInternalServerError {
					resp.Code = ErrorCodeUnexpectedFailure
				} else {
					resp.Code = ErrorCodeUnknown
				}
			}

			ctx.JSON(e.HTTPStatus, resp)

		} else {
			if e.ErrorCode == "" {
				if e.HTTPStatus == http.StatusInternalServerError {
					e.ErrorCode = ErrorCodeUnexpectedFailure
				} else {
					e.ErrorCode = ErrorCodeUnknown
				}
			}

			// Provide better error messages for certain user-triggered Postgres errors.
			if pgErr := utilities.NewPostgresError(e.InternalError); pgErr != nil {
				ctx.JSON(pgErr.HttpStatusCode, pgErr)
				return
			}
			ctx.JSON(e.HTTPStatus, e)
		}

	case *json.InvalidUnmarshalError:
		unmarshalErr, _ := err.(*json.UnmarshalTypeError)
		httpError := HTTPError{
			HTTPStatus: http.StatusBadRequest,
			ErrorCode:  ErrorCodeValidationFailed,
			Message:    fmt.Sprintf("Field %s is of incorrect type", unmarshalErr.Field),
		}
		ctx.JSON(http.StatusInternalServerError, httpError)

	case validator.ValidationErrors:
		errJson := utilities.ValidationErrorsToJSON(err)
		validationError := httpValidationError(http.StatusBadRequest, ErrorCodeValidationFailed, "Validation Error", errJson)

		ctx.JSON(http.StatusBadRequest, validationError)

	case ErrorCause:
		HandleResponseError(ctx, e.Cause())

	default:
		log.WithError(e).Errorf("Unhandled server error: %s", e.Error())

		if apiVersion.Compare(APIVersion20240101) >= 0 {
			resp := HTTPErrorResponse20240101{
				Code:    ErrorCodeUnexpectedFailure,
				Message: "Unexpected failure, please check server logs for more information",
			}

			ctx.JSON(http.StatusInternalServerError, resp)

		} else {
			httpError := HTTPError{
				HTTPStatus: http.StatusInternalServerError,
				ErrorCode:  ErrorCodeUnexpectedFailure,
				Message:    "Unexpected failure, please check server logs for more information",
			}

			ctx.JSON(http.StatusInternalServerError, httpError)

		}
	}
}

func recoverer() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := observability.GetLogEntry(ctx)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}

				se := &HTTPError{
					HTTPStatus: http.StatusInternalServerError,
					Message:    http.StatusText(http.StatusInternalServerError),
				}

				HandleResponseError(ctx, se)
			}

		}()

		ctx.Next()

	})
}
