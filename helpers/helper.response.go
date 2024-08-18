package helpers

import (
	"github.com/gin-gonic/gin"
)

type Responses struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func APIResponse(ctx *gin.Context, Message string, StatusCode int, Method string, Data interface{}) {
	jsonResponse := Responses{
		StatusCode: StatusCode,
		Method:     Method,
		Message:    Message,
		Data:       Data,
	}

	if StatusCode >= 400 {
		ctx.AbortWithStatusJSON(StatusCode, jsonResponse)
	} else {
		ctx.JSON(StatusCode, jsonResponse)
	}
}

type ErrorResponse struct {
	StatusCode int         `json:"statusCode"`
	Method     string      `json:"method"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error,omitempty"`
}

func APIErrorResponse(ctx *gin.Context, Error *HTTPError) {
	errResponse := ErrorResponse{
		StatusCode: Error.StatusCode,
		Method:     Error.Method,
		Message:    Error.Message,
		Error:      Error.ErrorData,
	}
	ctx.AbortWithStatusJSON(Error.StatusCode, errResponse)
}
