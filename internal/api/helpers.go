package api

import "github.com/gin-gonic/gin"

func sendJSON(ctx *gin.Context, status int, obj interface{}) {
	ctx.JSON(status, obj)
}
