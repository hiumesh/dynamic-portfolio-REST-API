package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/services"
)

type handlerUserProfile struct {
	service services.ServiceUserProfile
}

func (h *handlerUserProfile) Get(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject

	res, err := h.service.Get(userId)

	if err != nil {
		if dbErr, ok := err.(*helpers.DatabaseError); ok {
			switch dbErr.Type {
			case "ErrRecordNotFound":
				helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusNotFound, Method: ctx.Request.Method, Message: "Profile Not Found"})
				return
			default:
				helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusInternalServerError, Method: ctx.Request.Method, Message: "Database Error"})
				return
			}
		}

		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusInternalServerError, Method: ctx.Request.Method, Message: err.Error()})
		return
	}

	helpers.APIResponse(ctx, "Profile Data Successfully Fetched", http.StatusOK, ctx.Request.Method, res)
}

func (h *handlerUserProfile) Upsert(ctx *gin.Context) {
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

	err := h.service.Upsert(userId, &data)

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

func NewUserProfileHandler(service services.ServiceUserProfile) *handlerUserProfile {
	return &handlerUserProfile{
		service: service,
	}
}
