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

	blogService := services.NewBlogService(db)
	blogHandler := NewBlogHandler(blogService)

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
		portfolioRouter.GET("/", api.authenticateIfSessionPresent(), portfolioHandler.GetAll)
		portfolioRouter.GET("/user", api.requireAuthentication(), portfolioHandler.GetUserDetail)
		portfolioRouter.GET("/:slug/:module", api.authenticateIfSessionPresent(), portfolioHandler.GetSubModule)
		portfolioRouter.GET("/:slug", api.authenticateIfSessionPresent(), portfolioHandler.GetPortfolio)
		portfolioRouter.PUT("/skills", api.requireAuthentication(), portfolioHandler.UpsertSkills)
		portfolioRouter.GET("/status/:Status", api.requireAuthentication(), portfolioHandler.UpdateStatus)
		educationRouter := portfolioRouter.Group("/educations").Use(api.requireAuthentication())
		{
			educationRouter.GET("/", userEducationHandler.GetAll)
			educationRouter.POST("/", userEducationHandler.Create)
			educationRouter.PUT("/:Id", userEducationHandler.Update)
			educationRouter.PATCH("/:Id/reorder", userEducationHandler.Reorder)
			educationRouter.DELETE("/:Id", userEducationHandler.Delete)
			educationRouter.GET("/metadata", userEducationHandler.GetMetadata)
			educationRouter.PUT("/metadata", userEducationHandler.UpdateMetadata)
		}

		experienceRouter := portfolioRouter.Group("/experiences").Use(api.requireAuthentication())
		{
			experienceRouter.GET("/", userExperienceHandler.GetAll)
			experienceRouter.POST("/", userExperienceHandler.Create)
			experienceRouter.PUT("/:Id", userExperienceHandler.Update)
			experienceRouter.PATCH("/:Id/reorder", userExperienceHandler.Reorder)
			experienceRouter.DELETE("/:Id", userExperienceHandler.Delete)
			experienceRouter.GET("/metadata", userExperienceHandler.GetMetadata)
			experienceRouter.PUT("/metadata", userExperienceHandler.UpdateMetadata)
		}

		certificationRouter := portfolioRouter.Group("/certifications").Use(api.requireAuthentication())
		{
			certificationRouter.GET("/", userCertificationHandler.GetAll)
			certificationRouter.POST("/", userCertificationHandler.Create)
			certificationRouter.PUT("/:Id", userCertificationHandler.Update)
			certificationRouter.PATCH("/:Id/reorder", userCertificationHandler.Reorder)
			certificationRouter.DELETE("/:Id", userCertificationHandler.Delete)
			certificationRouter.GET("/metadata", userCertificationHandler.GetMetadata)
			certificationRouter.PUT("/metadata", userCertificationHandler.UpdateMetadata)
		}

		hackathonRouter := portfolioRouter.Group("/hackathons").Use(api.requireAuthentication())
		{
			hackathonRouter.GET("/", userHackathonHandler.GetAll)
			hackathonRouter.POST("/", userHackathonHandler.Create)
			hackathonRouter.PUT("/:Id", userHackathonHandler.Update)
			hackathonRouter.PATCH("/:Id/reorder", userHackathonHandler.Reorder)
			hackathonRouter.DELETE("/:Id", userHackathonHandler.Delete)
			hackathonRouter.GET("/metadata", userHackathonHandler.GetMetadata)
			hackathonRouter.PUT("/metadata", userHackathonHandler.UpdateMetadata)
		}
	}

	workGalleryRouter := router.Group("/work-gallery")
	{
		workGalleryRouter.GET("/", api.authenticateIfSessionPresent(), userWorkGalleryHandler.GetAll)
		workGalleryRouter.GET("/user", api.requireAuthentication(), userWorkGalleryHandler.GetUserWorkGallery)
		workGalleryRouter.GET("/user/:Id", api.requireAuthentication(), userWorkGalleryHandler.Get)
		workGalleryRouter.POST("/", api.requireAuthentication(), userWorkGalleryHandler.Create)
		workGalleryRouter.PUT("/:Id", api.requireAuthentication(), userWorkGalleryHandler.Update)
		workGalleryRouter.PATCH("/:Id/reorder", api.requireAuthentication(), userWorkGalleryHandler.Reorder)
		workGalleryRouter.DELETE("/:Id", api.requireAuthentication(), userWorkGalleryHandler.Delete)
		workGalleryRouter.GET("/metadata", api.requireAuthentication(), userWorkGalleryHandler.GetMetadata)
		workGalleryRouter.PUT("/metadata", api.requireAuthentication(), userWorkGalleryHandler.UpdateMetadata)

	}

	blogRouter := router.Group("/blogs")
	{
		blogRouter.GET("/", api.authenticateIfSessionPresent(), blogHandler.GetAll)
		blogRouter.GET("/user", api.requireAuthentication(), blogHandler.GetUserBlogs)
		blogRouter.GET("/user/:Id", api.requireAuthentication(), blogHandler.Get)
		blogRouter.GET("/:slug", api.authenticateIfSessionPresent(), blogHandler.GetBlogBySlug)
		blogRouter.PUT("/:Id/unpublish", blogHandler.Unpublish)
		blogRouter.POST("/", api.requireAuthentication(), blogHandler.Create)
		blogRouter.PUT("/:Id", api.requireAuthentication(), blogHandler.Update)
		blogRouter.DELETE("/:Id", api.requireAuthentication(), blogHandler.Delete)
		blogRouter.GET("/metadata", api.requireAuthentication(), blogHandler.GetMetadata)
		blogRouter.PUT("/metadata", api.requireAuthentication(), blogHandler.UpdateMetadata)
	}
}
