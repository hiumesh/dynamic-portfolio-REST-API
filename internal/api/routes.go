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

	userExperienceService := services.NewUserExperienceService(db)
	userExperienceHandler := NewUserExperienceHandler(userExperienceService)

	userCertificationService := services.NewUserCertificationService(db)
	userCertificationHandler := NewUserCertificationHandler(userCertificationService)

	userRouter := router.Group("/users")
	{
		userRouter.POST("/presigned-urls", userHandler.GetPresignedURLs)
		profileRouter := userRouter.Group("/profile").Use(api.requireAuthentication())
		{
			profileRouter.GET("/", userHandler.GetProfile)
			profileRouter.PUT("/setup", userHandler.ProfileSetup)
			profileRouter.PUT("/", userHandler.UpsertProfile)
			profileRouter.PUT("/skills", userHandler.UpsertSkills)
		}

		educationRouter := userRouter.Group("/educations").Use(api.requireAuthentication())
		{
			educationRouter.GET("/", userEducationHandler.GetAll)
			educationRouter.POST("/", userEducationHandler.Create)
			educationRouter.PUT("/:Id", userEducationHandler.Update)
			educationRouter.PATCH("/:Id/reorder", userEducationHandler.Reorder)
			educationRouter.DELETE("/:Id", userEducationHandler.Delete)
		}

		experienceRouter := userRouter.Group("/experiences").Use(api.requireAuthentication())
		{
			experienceRouter.GET("/", userExperienceHandler.GetAll)
			experienceRouter.POST("/", userExperienceHandler.Create)
			experienceRouter.PUT("/:Id", userExperienceHandler.Update)
			experienceRouter.PATCH("/:Id/reorder", userExperienceHandler.Reorder)
			experienceRouter.DELETE("/:Id", userExperienceHandler.Delete)
		}

		certificationRouter := userRouter.Group("/certifications").Use(api.requireAuthentication())
		{
			certificationRouter.GET("/", userCertificationHandler.GetAll)
			certificationRouter.POST("/", userCertificationHandler.Create)
			certificationRouter.PUT("/:Id", userCertificationHandler.Update)
			certificationRouter.PATCH("/:Id/reorder", userCertificationHandler.Reorder)
			certificationRouter.DELETE("/:Id", userCertificationHandler.Delete)
		}
	}
}
