package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type handlerWorkGallery struct {
	service services.ServiceWorkGallery
}

func (h *handlerWorkGallery) GetAll(ctx *gin.Context) {
	claims := utilities.GetClaims(ctx)
	var userId *string
	if claims != nil {
		userId = &claims.Subject
	}
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		HandleResponseError(ctx, ValidationError("Invalid limit value. Limit must be a positive integer.", err))
		return
	}

	cursorStr := ctx.DefaultQuery("cursor", "0")
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil || cursor < 0 {
		HandleResponseError(ctx, ValidationError("Invalid cursor value. Cursor must be a non-negative integer.", err))
		return
	}

	var query *string
	queryStr := ctx.Query("query")
	if len(queryStr) > 0 {
		query = &queryStr
	}

	res, err := h.service.GetAll(userId, query, cursor, limit)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerWorkGallery) GetUserWorkGallery(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		HandleResponseError(ctx, ValidationError("Invalid limit value. Limit must be a positive integer.", err))
		return
	}

	cursorStr := ctx.DefaultQuery("cursor", "0")
	cursor, err := strconv.Atoi(cursorStr)
	if err != nil || cursor < 0 {
		HandleResponseError(ctx, ValidationError("Invalid cursor value. Cursor must be a non-negative integer.", err))
		return
	}

	var query *string
	queryStr := ctx.Query("query")
	if len(queryStr) > 0 {
		query = &queryStr
	}

	res, err := h.service.GetUserWorkGallery(userId, query, cursor, limit)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerWorkGallery) Get(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	res, err := h.service.Get(userId, id)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerWorkGallery) Create(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaTechProject
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

func (h *handlerWorkGallery) Update(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	var data schemas.SchemaTechProject
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

func (h *handlerWorkGallery) Reorder(ctx *gin.Context) {
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

func (h *handlerWorkGallery) Delete(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	err := h.service.Delete(userId, id)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerWorkGallery) GetMetadata(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetMetadata(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerWorkGallery) UpdateMetadata(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaTechProjectMetadata
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

func NewWorkGalleryHandler(service services.ServiceWorkGallery) *handlerWorkGallery {
	return &handlerWorkGallery{
		service: service,
	}
}
