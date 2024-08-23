package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/handlers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/middlewares"
	"github.com/hiumesh/dynamic-portfolio-REST-API/services"
	"gorm.io/gorm"
)

func InitUserProfileRoutes(db *gorm.DB, router *gin.RouterGroup) {

	service := services.NewUserProfileService(db)
	handler := handlers.NewUserProfileHandler(service)

	groupRoute := router.Group("/users").Use(middlewares.Auth())
	{
		groupRoute.GET("/profile", handler.Get)
		groupRoute.PUT("/profile", handler.Upsert)
	}

}
