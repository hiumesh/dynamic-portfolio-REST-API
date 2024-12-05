package api

import (
	"net/http"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/observability"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	defaultVersion = "unknown version"
)

func NewAPI(globalConfig *config.GlobalConfiguration, db *gorm.DB) *API {
	return NewAPIWithVersion(globalConfig, db, defaultVersion)
}

type API struct {
	handler *gin.Engine
	db      *gorm.DB
	config  *config.GlobalConfiguration
	version string
}

func NewAPIWithVersion(globalConfig *config.GlobalConfiguration, db *gorm.DB, version string) *API {
	api := &API{config: globalConfig, db: db, version: version}
	gin.SetMode(gin.ReleaseMode)
	app := gin.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))
	app.Use(helmet.Default())
	app.Use(gzip.Gzip(gzip.BestCompression))

	app.Use(observability.AddRequestID(globalConfig))
	app.Use(observability.NewStructuredLogger(logrus.StandardLogger(), globalConfig))
	app.Use(recoverer())

	logrus.Info(globalConfig.API.MaxRequestDuration)

	// TODO: Fix
	// if globalConfig.API.MaxRequestDuration > 0 {
	// 	app.Use(timeoutMiddleware(globalConfig.API.MaxRequestDuration))
	// }

	app.GET("/health", api.HealthCheck)

	router := app.Group("/rest/v1")

	setupRoutes(router, db, api, globalConfig)

	api.handler = app

	return api
}

type HealthCheckResponse struct {
	Version     string `json:"version"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (a *API) HealthCheck(ctx *gin.Context) {
	sendJSON(ctx, http.StatusOK, HealthCheckResponse{
		Version:     a.version,
		Name:        "GoTrue",
		Description: "GoTrue is a user registration and authentication API",
	})
}

// ServeHTTP implements the http.Handler interface by passing the request along
// to its underlying Handler.
func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}
