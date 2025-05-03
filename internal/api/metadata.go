package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
)

type handlerMetadata struct {
	service services.ServiceMetadata
}

func (h *handlerMetadata) GetAllSkills(ctx *gin.Context) {
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

	res, err := h.service.GetAllSkills(query, cursor, limit)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func NewMetadataHandler(service services.ServiceMetadata) *handlerMetadata {
	return &handlerMetadata{service: service}
}
