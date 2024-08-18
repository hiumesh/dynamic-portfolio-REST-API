package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hiumesh/dynamic-portfolio-REST-API/config"
)

type contextKey string

func (c contextKey) String() string {
	return "gotrue api context key " + string(c)
}

const (
	tokenKey     = contextKey("jwt")
	requestIDKey = contextKey("request_id")
)

func WithToken(ctx *gin.Context, token *jwt.Token) {
	ctx.Set(string(tokenKey), token)
}

func GetToken(ctx *gin.Context) *jwt.Token {
	obj, exists := ctx.Get(string(tokenKey))
	if !exists || obj == nil {
		return nil
	}

	return obj.(*jwt.Token)
}

func GetClaims(ctx *gin.Context) *config.AccessTokenClaims {
	token := GetToken(ctx)
	if token == nil {
		return nil
	}
	return token.Claims.(*config.AccessTokenClaims)
}

func WithRequestID(ctx *gin.Context, id string) {
	ctx.Set(string(requestIDKey), id)
}

func GetRequestID(ctx *gin.Context) string {
	obj, exists := ctx.Get(string(requestIDKey))
	if !exists || obj == nil {
		return ""
	}
	return obj.(string)
}
