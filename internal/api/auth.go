package api

import (
	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/pkg"
)

func (a *API) requireAuthentication() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		tokenStr, err := pkg.ExtractBearerToken(ctx)
		if err != nil {
			HandleResponseError(ctx, UnauthorizedError(ctx.Request.Method, err.Error()))
			return
		}

		token, err := pkg.ParseJWTClaims(tokenStr, ctx, config.AccessTokenClaims{}, a.config.JWT.Secret)
		if err != nil {
			HandleResponseError(ctx, UnauthorizedError(ctx.Request.Method, err.Error()))
			return
		}

		withToken(ctx, token)

		ctx.Next()
	})
}

func (a *API) authenticateIfSessionPresent() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		tokenStr, err := pkg.ExtractBearerToken(ctx)
		if err != nil {
			ctx.Next()
			return
		}

		token, err := pkg.ParseJWTClaims(tokenStr, ctx, config.AccessTokenClaims{}, a.config.JWT.Secret)
		if err != nil {
			ctx.Next()
			return
		}

		withToken(ctx, token)
		ctx.Next()
	})
}
