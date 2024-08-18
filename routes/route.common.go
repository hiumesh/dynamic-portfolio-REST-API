package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/hiumesh/dynamic-portfolio-REST-API/handlers/common"
	repositories "github.com/hiumesh/dynamic-portfolio-REST-API/repositories/common"
	services "github.com/hiumesh/dynamic-portfolio-REST-API/services/common"
	"gorm.io/gorm"
)

func InitCommonRoutes(db *gorm.DB, router *gin.RouterGroup) {

	pingRepository := repositories.NewRepositoryPing(db)
	pingService := services.NewPingService(pingRepository)
	pingHandler := handlers.NewHandlerPing(pingService)

	groupRoute := router.Group("/common")
	{
		groupRoute.GET("/ping", pingHandler.PingHandler)
	}

}
