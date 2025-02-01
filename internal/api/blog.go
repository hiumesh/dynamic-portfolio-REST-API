package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type handlerBlog struct {
	service services.ServiceBlog
}

func (h *handlerBlog) GetAll(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetAll(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerBlog) Get(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	res, err := h.service.Get(userId, id)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerBlog) Create(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	status := ctx.Query("status")

	var data schemas.SchemaBlog
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	res, err := h.service.Create(userId, &data, status == "publish")

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)

}

func (h *handlerBlog) Update(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")
	status := ctx.Query("status")

	var data schemas.SchemaBlog
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	res, err := h.service.Update(userId, id, &data, status == "publish")

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)

}

func (h *handlerBlog) Unpublish(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	err := h.service.Unpublish(userId, id)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerBlog) Delete(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	err := h.service.Delete(userId, id)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerBlog) GetMetadata(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetMetadata(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerBlog) UpdateMetadata(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaBlogMetadata
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

func NewBlogHandler(service services.ServiceBlog) *handlerBlog {
	return &handlerBlog{service: service}
}
