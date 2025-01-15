package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type handlerUser struct {
	service services.ServiceUser
}

func (h *handlerUser) GetProfile(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetProfile(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerUser) ProfileSetup(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaProfileBasic
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	err := h.service.ProfileSetup(userId, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)

}

func (h *handlerUser) UpsertProfile(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaProfileBasic
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	err := h.service.UpsertProfile(userId, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)

}

func (h *handlerUser) GetPresignedURLs(ctx *gin.Context) {
	// userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaPresignedURL
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	urls, err := h.service.GetPostPresignedURLs(ctx, data.Files)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, urls)
}

func NewUserHandler(service services.ServiceUser) *handlerUser {
	return &handlerUser{
		service: service,
	}
}
