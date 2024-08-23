package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/handlers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/services"

	"gorm.io/gorm"
)

func InitCommonRoutes(db *gorm.DB, router *gin.RouterGroup) {

	service := services.NewCommonService(db)
	handler := handlers.NewCommonHandler(service)

	groupRoute := router.Group("/common")
	{
		groupRoute.GET("/ping", handler.Ping)
	}

}
