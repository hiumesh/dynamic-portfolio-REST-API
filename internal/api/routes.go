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

	portfolioService := services.NewPortfolioService(db)
	portfolioHandler := NewPortfolioHandler(portfolioService)

	userEducationService := services.NewUserEducationService(db)
	userEducationHandler := NewUserEducationHandler(userEducationService)

	userExperienceService := services.NewUserExperienceService(db)
	userExperienceHandler := NewUserExperienceHandler(userExperienceService)

	userCertificationService := services.NewUserCertificationService(db)
	userCertificationHandler := NewUserCertificationHandler(userCertificationService)

	userHackathonService := services.NewUserHackathonService(db)
	userHackathonHandler := NewUserHackathonHandler(userHackathonService)

	userWorkGalleryService := services.NewWorkGalleryService(db)
	userWorkGalleryHandler := NewWorkGalleryHandler(userWorkGalleryService)

	userRouter := router.Group("/users")
	{
		userRouter.POST("/presigned-urls", userHandler.GetPresignedURLs)
		profileRouter := userRouter.Group("/profile").Use(api.requireAuthentication())
		{
			profileRouter.GET("/", userHandler.GetProfile)
			profileRouter.PUT("/setup", userHandler.ProfileSetup)
			profileRouter.PUT("/", userHandler.UpsertProfile)
		}
	}

	portfolioRouter := router.Group("/portfolio")
	{
		portfolioRouter.Use(api.requireAuthentication()).GET("/user", portfolioHandler.Get)
		portfolioRouter.Use(api.requireAuthentication()).PUT("/skills", portfolioHandler.UpsertSkills)
		portfolioRouter.Use(api.requireAuthentication()).GET("/status/:Status", portfolioHandler.UpdateStatus)
		educationRouter := portfolioRouter.Group("/educations").Use(api.requireAuthentication())
		{
			educationRouter.GET("/", userEducationHandler.GetAll)
			educationRouter.POST("/", userEducationHandler.Create)
			educationRouter.PUT("/:Id", userEducationHandler.Update)
			educationRouter.PATCH("/:Id/reorder", userEducationHandler.Reorder)
			educationRouter.DELETE("/:Id", userEducationHandler.Delete)
			educationRouter.PUT("/metadata", userEducationHandler.UpdateMetadata)
		}

		experienceRouter := portfolioRouter.Group("/experiences").Use(api.requireAuthentication())
		{
			experienceRouter.GET("/", userExperienceHandler.GetAll)
			experienceRouter.POST("/", userExperienceHandler.Create)
			experienceRouter.PUT("/:Id", userExperienceHandler.Update)
			experienceRouter.PATCH("/:Id/reorder", userExperienceHandler.Reorder)
			experienceRouter.DELETE("/:Id", userExperienceHandler.Delete)
			experienceRouter.PUT("/metadata", userExperienceHandler.UpdateMetadata)
		}

		certificationRouter := portfolioRouter.Group("/certifications").Use(api.requireAuthentication())
		{
			certificationRouter.GET("/", userCertificationHandler.GetAll)
			certificationRouter.POST("/", userCertificationHandler.Create)
			certificationRouter.PUT("/:Id", userCertificationHandler.Update)
			certificationRouter.PATCH("/:Id/reorder", userCertificationHandler.Reorder)
			certificationRouter.DELETE("/:Id", userCertificationHandler.Delete)
			certificationRouter.PUT("/metadata", userCertificationHandler.UpdateMetadata)
		}

		hackathonRouter := portfolioRouter.Group("/hackathons").Use(api.requireAuthentication())
		{
			hackathonRouter.GET("/", userHackathonHandler.GetAll)
			hackathonRouter.POST("/", userHackathonHandler.Create)
			hackathonRouter.PUT("/:Id", userHackathonHandler.Update)
			hackathonRouter.PATCH("/:Id/reorder", userHackathonHandler.Reorder)
			hackathonRouter.DELETE("/:Id", userHackathonHandler.Delete)
			hackathonRouter.PUT("/metadata", userHackathonHandler.UpdateMetadata)
		}
	}

	workGalleryRouter := router.Group("/work-gallery").Use(api.requireAuthentication())
	{

		workGalleryRouter.GET("/", userWorkGalleryHandler.GetAll)
		workGalleryRouter.POST("/", userWorkGalleryHandler.Create)
		workGalleryRouter.PUT("/:Id", userWorkGalleryHandler.Update)
		workGalleryRouter.PATCH("/:Id/reorder", userWorkGalleryHandler.Reorder)
		workGalleryRouter.DELETE("/:Id", userWorkGalleryHandler.Delete)
		workGalleryRouter.PUT("/metadata", userWorkGalleryHandler.UpdateMetadata)

	}
}
