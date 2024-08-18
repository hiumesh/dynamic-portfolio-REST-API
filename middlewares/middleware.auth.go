package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/config"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/pkg"
	"github.com/sirupsen/logrus"
)

func Auth() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {
		tokenStr, err := pkg.ExtractBearerToken(ctx)
		if err != nil {
			logrus.Errorf("authentication error: %v", err)
			helpers.APIErrorResponse(ctx, helpers.UnauthorizedError(ctx.Request.Method, err.Error()))
			return
		}

		token, err := pkg.ParseJWTClaims(tokenStr, ctx, config.AccessTokenClaims{})
		if err != nil {
			logrus.Errorf("authentication error: %v", err)
			helpers.APIErrorResponse(ctx, helpers.UnauthorizedError(ctx.Request.Method, err.Error()))
			return
		}

		helpers.WithToken(ctx, token)

		ctx.Next()
	})
}
