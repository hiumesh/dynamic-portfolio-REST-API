package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type handlerUserExperience struct {
	service services.ServiceUserExperience
}

func (h *handlerUserExperience) GetAll(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetAll(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerUserExperience) Create(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaWorkExperience
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	res, err := h.service.Create(userId, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)

}

func (h *handlerUserExperience) Update(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	var data schemas.SchemaWorkExperience
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	res, err := h.service.Update(userId, id, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)

}

func (h *handlerUserExperience) Reorder(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	var data schemas.SchemaReorderItem
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	err := h.service.Reorder(userId, id, int(data.NewIndex))

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)

}

func (h *handlerUserExperience) Delete(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	err := h.service.Delete(userId, id)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerUserExperience) GetMetadata(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetMetadata(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerUserExperience) UpdateMetadata(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaWorkExperienceMetadata
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	err := h.service.UpdateMetadata(userId, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func NewUserExperienceHandler(service services.ServiceUserExperience) *handlerUserExperience {
	return &handlerUserExperience{
		service: service,
	}
}
