package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/config"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/pkg"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/services"
	"gorm.io/gorm"
)

func setupRoutes(router *gin.RouterGroup, db *gorm.DB, api *API, globalConfig *config.GlobalConfiguration) {
	presigner := pkg.NewPresigner(context.TODO(), &globalConfig.AWS)

	userService := services.NewUserService(db, presigner)
	userHandler := NewUserHandler(userService)

	userEducationService := services.NewUserEducationService(db)
	userEducationHandler := NewUserEducationHandler(userEducationService)

	userRouter := router.Group("/users")
	{
		userRouter.POST("/presigned-urls", userHandler.GetPresignedURLs)
		profileRouter := userRouter.Group("/profile").Use(api.requireAuthentication())
		{
			profileRouter.GET("/", userHandler.GetProfile)
			profileRouter.PUT("/setup", userHandler.ProfileSetup)
			profileRouter.PUT("/", userHandler.UpsertProfile)
		}

		educationRouter := userRouter.Group("/educations").Use(api.requireAuthentication())
		{
			educationRouter.GET("/", userEducationHandler.GetAll)
			educationRouter.POST("/", userEducationHandler.Create)
			educationRouter.PUT("/:Id", userEducationHandler.Update)
			educationRouter.PATCH("/:Id/reorder", userEducationHandler.Reorder)
			educationRouter.DELETE("/:Id", userEducationHandler.Delete)
		}
	}
}
