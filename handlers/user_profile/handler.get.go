package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	services "github.com/hiumesh/dynamic-portfolio-REST-API/services/user_profile"
)

type handlerGetProfile struct {
	service services.ServiceGetProfile
}

func (h *handlerGetProfile) GetProfileHandler(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject

	res, err := h.service.GetProfileService(userId)

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

func NewGetProfileHandler(service services.ServiceGetProfile) *handlerGetProfile {
	return &handlerGetProfile{
		service: service,
	}
}
