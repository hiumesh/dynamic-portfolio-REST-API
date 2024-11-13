package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/services"
)

type handlerUserEducation struct {
	service services.ServiceUserEducation
}

func (h *handlerUserEducation) GetAll(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject

	res, err := h.service.GetAll(userId)

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

	helpers.APIResponse(ctx, "Data Successfully Fetched", http.StatusOK, ctx.Request.Method, res)
}

func (h *handlerUserEducation) Create(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject

	var data schemas.SchemaUserEducation
	if err := ctx.ShouldBindJSON(&data); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Invalid request payload", ErrorData: err})
		return
	}

	if err := data.Validate(); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Validation Error", ErrorData: helpers.ValidationErrorsToJSON(err)})
		return
	}

	res, err := h.service.Create(userId, &data)

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

	helpers.APIResponse(ctx, "Successful", http.StatusOK, ctx.Request.Method, res)

}

func (h *handlerUserEducation) Update(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	var data schemas.SchemaUserEducation
	if err := ctx.ShouldBindJSON(&data); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Invalid request payload", ErrorData: err})
		return
	}

	if err := data.Validate(); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Validation Error", ErrorData: helpers.ValidationErrorsToJSON(err)})
		return
	}

	res, err := h.service.Update(userId, id, &data)

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

	helpers.APIResponse(ctx, "Update Successful", http.StatusOK, ctx.Request.Method, res)

}

func (h *handlerUserEducation) Reorder(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	var data schemas.SchemaReorderUserEducation
	if err := ctx.ShouldBindJSON(&data); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Invalid request payload", ErrorData: err})
		return
	}

	if err := data.Validate(); err != nil {
		helpers.APIErrorResponse(ctx, &helpers.HTTPError{StatusCode: http.StatusBadRequest, Method: ctx.Request.Method, Message: "Validation Error", ErrorData: helpers.ValidationErrorsToJSON(err)})
		return
	}

	err := h.service.Reorder(userId, id, int(data.NewIndex))

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

	helpers.APIResponse(ctx, "Reordering Successfully", http.StatusOK, ctx.Request.Method, nil)

}

func (h *handlerUserEducation) Delete(ctx *gin.Context) {
	userId := helpers.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	err := h.service.Delete(userId, id)

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

	helpers.APIResponse(ctx, "Deleted Successfully", http.StatusOK, ctx.Request.Method, nil)

}

func NewUserEducationHandler(service services.ServiceUserEducation) *handlerUserEducation {
	return &handlerUserEducation{
		service: service,
	}
}
