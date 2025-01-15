package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type handlerPortfolio struct {
	service services.ServicePortfolio
}

func (h *handlerPortfolio) Get(ctx *gin.Context) {
	userId := utilities.GetClaims(ctx).Subject

	res, err := h.service.Get(userId)

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
