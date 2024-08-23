package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/services"
)

type handlerCommon struct {
	service services.ServiceCommon
}

func (h *handlerCommon) Ping(ctx *gin.Context) {
	res := h.service.Ping()
	helpers.APIResponse(ctx, res, http.StatusOK, ctx.Request.Method, res)
}

func NewCommonHandler(service services.ServiceCommon) *handlerCommon {
	return &handlerCommon{
		service: service,
	}
}
