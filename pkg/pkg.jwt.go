package pkg

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AMREntry struct {
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
	Provider  string `json:"provider,omitempty"`
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	Email                         string                 `json:"email"`
	Phone                         string                 `json:"phone"`
	AppMetaData                   map[string]interface{} `json:"app_metadata"`
	UserMetaData                  map[string]interface{} `json:"user_metadata"`
	Role                          string                 `json:"role"`
	AuthenticatorAssuranceLevel   string                 `json:"aal,omitempty"`
	AuthenticationMethodReference []AMREntry             `json:"amr,omitempty"`
	SessionId                     string                 `json:"session_id,omitempty"`
}

func extractBearerToken(ctx *gin.Context) (string, *utils.HTTPError) {
	authHeader := ctx.Request.Header.Get("Authorization")
	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", utils.UnauthorizedError("This endpoint requires a Bearer token")
	}

	return matches[1], nil
}

func parseJWTClaims(bearer string, ctx *gin.Context) (context.Context, *utils.HTTPError) {
	config := a.config

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	token, err := p.ParseWithClaims(bearer, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.Secret), nil
	})
	if err != nil {
		return ctx, utils.UnauthorizedError("invalid JWT: unable to parse or verify signature, %v", err)
	}
	utils.WithToken(ctx, token)
	return ctx, nil
}
