package main

import (
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/hiumesh/dynamic-portfolio-REST-API/pkg"
	"github.com/hiumesh/dynamic-portfolio-REST-API/routes"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	app := SetupRouter()
	logrus.Fatal(app.Run(":" + pkg.GodotEnv("GO_PORT")))
}

func SetupRouter() *gin.Engine {
	db := SetupDatabase()
	app := gin.Default()

	if pkg.GodotEnv("GO_ENV") != "production" && pkg.GodotEnv("GO_ENV") != "test" {
		gin.SetMode(gin.DebugMode)
	} else if pkg.GodotEnv("GO_ENV") == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	app.Use(helmet.Default())
	app.Use(gzip.Gzip(gzip.BestCompression))

	router := app.Group("/api/v1")

	routes.InitCommonRoutes(db, router)
	routes.InitUserProfileRoutes(db, router)

	return app
}

func SetupDatabase() *gorm.DB {
	db, err := gorm.Open(postgres.Open(pkg.GodotEnv("DATABASE_URI")), &gorm.Config{})

	if err != nil {
		defer logrus.Info("Connect into Database Failed")
		logrus.Fatal(err.Error())
	}

	if pkg.GodotEnv("GO_ENV") != "production" {
		logrus.Info("Connect into Database Successfully")
	}

	// err = db.AutoMigrate(
	// 	&models.Portfolio{},
	// )

	// if err != nil {
	// 	defer logrus.Info("Auto Migration Failed")
	// 	logrus.Fatal(err.Error())
	// }

	return db
}
