package helpers

import (
	"github.com/gin-gonic/gin"
	schemas "github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
)

func APIResponse(ctx *gin.Context, Message string, StatusCode int, Method string, Data interface{}) {
	jsonResponse := schemas.SchemaResponses{
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

func ValidationErrorResponse(ctx *gin.Context, StatusCode int, Method string, Error interface{}) {
	errResponse := schemas.SchemaErrorResponse{
		StatusCode: StatusCode,
		Method:     Method,
		Error:      Error,
	}

	ctx.AbortWithStatusJSON(StatusCode, errResponse)
}
