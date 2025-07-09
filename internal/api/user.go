package api

import (
	"net/http"
	"strconv"

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

func (h *handlerUser) GetFollowers(ctx *gin.Context) {
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

	res, err := h.service.GetFollowers(userId, cursor, limit)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerUser) GetFollowing(ctx *gin.Context) {
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

	res, err := h.service.GetFollowing(userId, cursor, limit)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerUser) Follow(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	slug := ctx.Param("slug")

	err := h.service.FollowUser(userId, slug)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerUser) Unfollow(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	slug := ctx.Param("slug")

	err := h.service.UnfollowUser(userId, slug)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerUser) FollowStatus(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	slug := ctx.Param("slug")

	res, err := h.service.FollowStatus(userId, slug)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func NewUserHandler(service services.ServiceUser) *handlerUser {
	return &handlerUser{
		service: service,
	}
}
