package routes

import (
	"github.com/gin-gonic/gin"
	handlers "github.com/hiumesh/dynamic-portfolio-REST-API/handlers/user_profile"
	"github.com/hiumesh/dynamic-portfolio-REST-API/middlewares"
	repositories "github.com/hiumesh/dynamic-portfolio-REST-API/repositories/user_profile"
	services "github.com/hiumesh/dynamic-portfolio-REST-API/services/user_profile"
	"gorm.io/gorm"
)

func InitUserProfileRoutes(db *gorm.DB, router *gin.RouterGroup) {

	getProfileRepository := repositories.NewGetProfileRepository(db)
	getProfileService := services.NewGetProfileService(getProfileRepository)
	getProfileHandler := handlers.NewGetProfileHandler(getProfileService).GetProfileHandler

	upsertProfileRepository := repositories.NewUpsertProfileRepository(db)
	upsertProfileService := services.NewUpsertProfileService(upsertProfileRepository)
	upsertProfileHandler := handlers.NewUpsertProfileHandler(upsertProfileService).UpsertProfileHandler

	groupRoute := router.Group("/users").Use(middlewares.Auth())
	{
		groupRoute.GET("/profile", getProfileHandler)
		groupRoute.PUT("/profile", upsertProfileHandler)
	}

}
