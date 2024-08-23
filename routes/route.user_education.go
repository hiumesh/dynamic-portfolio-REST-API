package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/handlers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/middlewares"
	"github.com/hiumesh/dynamic-portfolio-REST-API/services"
	"gorm.io/gorm"
)

func InitUserEducationRoutes(db *gorm.DB, router *gin.RouterGroup) {

	service := services.NewUserEducationService(db)
	handler := handlers.NewUserEducationHandler(service)

	groupRoute := router.Group("/users").Use(middlewares.Auth())
	{
		groupRoute.GET("/educations", handler.GetAll)
		groupRoute.POST("/educations", handler.Create)
		groupRoute.PUT("/educations/:Id", handler.Update)
		groupRoute.PATCH("/educations/:Id/reorder", handler.Reorder)
	}

}
