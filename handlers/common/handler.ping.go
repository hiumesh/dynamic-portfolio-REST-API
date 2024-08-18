package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	helpers "github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	services "github.com/hiumesh/dynamic-portfolio-REST-API/services/common"
)

type handlerPing struct {
	service services.ServicePing
}

func NewHandlerPing(service services.ServicePing) *handlerPing {
	return &handlerPing{service: service}
}

func (h *handlerPing) PingHandler(ctx *gin.Context) {
	res := h.service.PingService()
	helpers.APIResponse(ctx, res, http.StatusOK, ctx.Request.Method, res)
}
