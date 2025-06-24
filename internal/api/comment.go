package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type handlerComment struct {
	service services.ServiceComment
}

func (h *handlerComment) GetAll(ctx *gin.Context) {

	claims := utilities.GetClaims(ctx)
	var userId *string
	if claims != nil {
		userId = &claims.Subject
	}

	slug := ctx.Query("slug")
	if slug == "" {
		HandleResponseError(ctx, ValidationError("Invalid object value. Object must be a non-empty string.", nil))
		return
	}

	moduleStr := ctx.Query("module")
	if moduleStr == "" {
		HandleResponseError(ctx, ValidationError("Invalid module value. Module must be a non-empty string.", nil))
		return
	}

	parentStr := ctx.Query("parent_id")
	var parentId *int

	if parentStr != "" {
		val, err := strconv.Atoi(parentStr)
		if err != nil {
			HandleResponseError(ctx, ValidationError("Invalid parent_id value. Parent ID must be a non-negative integer.", err))
			return
		}
		parentId = &val
	}

	limitStr := ctx.DefaultQuery("limit", "5")
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

	res, err := h.service.GetAll(userId, moduleStr, slug, cursor, limit, parentId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)

}

func (h *handlerComment) Create(ctx *gin.Context) {

	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaCreateComment
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

func (h *handlerComment) Reply(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	var data schemas.SchemaCommentReply
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	res, err := h.service.Reply(userId, id, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerComment) Reaction(ctx *gin.Context) {

	userId := utilities.GetClaims(ctx).Subject
	id := ctx.Param("Id")

	var data schemas.SchemaReaction
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	res, err := h.service.Reaction(id, userId, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func NewCommentHandler(service services.ServiceComment) *handlerComment {
	return &handlerComment{service}
}
