package api

import (
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type handlerPortfolio struct {
	service services.ServicePortfolio
}

func (h *handlerPortfolio) GetUserDetail(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetUserPortfolio(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerPortfolio) GetAll(ctx *gin.Context) {
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

func (h *handlerPortfolio) GetPortfolio(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		HandleResponseError(ctx, ValidationError("Invalid slug value. Slug must be a non-empty string.", nil))
		return
	}

	res, err := h.service.GetPortfolio(slug)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerPortfolio) GetSubModule(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		HandleResponseError(ctx, ValidationError("Invalid slug value. Slug must be a non-empty string.", nil))
		return
	}

	module := ctx.Param("module")
	if module == "" {
		HandleResponseError(ctx, ValidationError("Invalid module value. Module must be a non-empty string.", nil))
		return
	}
	modules := []string{"educations", "work_experiences", "certifications", "hackathons", "works", "skills"}

	if !slices.Contains(modules, module) {
		HandleResponseError(ctx, ValidationError("Invalid module value. Module must be one of educations, experiences, certifications, hackathons.", nil))
		return
	}

	res, err := h.service.GetSubModule(slug, module)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerPortfolio) GetUserSkills(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.GetSkills(userId)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, res)
}

func (h *handlerPortfolio) UpsertSkills(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaSkills
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	err := h.service.UpsertSkills(userId, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)

}

func (h *handlerPortfolio) UpsertResume(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaResume
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	err := h.service.UpsertResume(userId, &data.ResumeUrl)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerPortfolio) UpdateProfileAttachment(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	var data schemas.SchemaProfileAttachment
	if err := ctx.ShouldBindJSON(&data); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	if err := data.Validate(); err != nil {
		HandleResponseError(ctx, err)
		return
	}

	err := h.service.UpdateProfileAttachment(userId, &data)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func (h *handlerPortfolio) UpdateStatus(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject
	status := ctx.Param("Status")

	err := h.service.UpdateStatus(userId, status)

	if err != nil {
		HandleResponseError(ctx, err)
		return
	}

	sendJSON(ctx, http.StatusOK, nil)
}

func NewPortfolioHandler(service services.ServicePortfolio) *handlerPortfolio {
	return &handlerPortfolio{service: service}
}
