package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	services "github.com/hiumesh/dynamic-portfolio-REST-API/services/user_profile"
)

type handlerUpsertProfile struct {
	service services.ServiceUpsertProfile
}

func (h *handlerUpsertProfile) UpsertProfileHandler(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject

	var data schemas.SchemaProfileBasic
	if err := ctx.ShouldBindJSON(&data); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Invalid request payload", ErrorData: err})
		return
	}

	if err := data.Validate(); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Validation Error", ErrorData: helpers.ValidationErrorsToJSON(err)})
		return
	}

	err := h.service.UpsertProfileService(userId, &data)

	if err != nil {
		if dbErr, ok := err.(*helpers.DatabaseError); ok {
			switch dbErr.Type {
			default:
				helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusInternalServerError, Method: ctx.Request.Method, Message: "Database Error"})
				return
			}
		}

		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusInternalServerError, Method: ctx.Request.Method, Message: err.Error()})
		return
	}

	helpers.APIResponse(ctx, "Profile Updated", http.StatusOK, ctx.Request.Method, nil)

}

func NewUpsertProfileHandler(service services.ServiceUpsertProfile) *handlerUpsertProfile {
	return &handlerUpsertProfile{
		service: service,
	}
}
