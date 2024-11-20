package pkg

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
)

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

func ExtractBearerToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.Request.Header.Get("Authorization")
	matches := bearerRegexp.FindStringSubmatch(authHeader)
	if len(matches) != 2 {
		return "", errors.New("This endpoint requires a Bearer token")
	}

	return matches[1], nil
}

func ParseJWTClaims(bearer string, ctx *gin.Context, claims interface{}, secret string) (*jwt.Token, error) {

	p := jwt.Parser{ValidMethods: []string{jwt.SigningMethodHS256.Name}}
	token, err := p.ParseWithClaims(bearer, &config.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, errors.New("Invalid JWT: token may be expired or system is unable to parse and verify signature")
	}

	return token, nil
}
