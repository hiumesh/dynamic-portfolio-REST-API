package cmd

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/db/seeds"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/observability"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var seedCmd = cobra.Command{
	Use:  "seed",
	Long: "Seed database with test data",
	Run:  seed,
}

func seed(cmd *cobra.Command, args []string) {
	globalConfig := loadGlobalConfig(cmd.Context())

	db, err := gorm.Open(postgres.Open(globalConfig.DB.URL), &gorm.Config{Logger: observability.NewGormLogrusLogger(globalConfig.LOGGING.Level, globalConfig.LOGGING.SQL)})

	if err != nil {
		logrus.Fatalf("error opening database: %+v", err)
	}

	logrus.Info("Truncating...")
	err = seeds.Truncate(db)
	if err != nil {
		logrus.Fatalf("error truncating tables: %+v", err)
	}

	logrus.Info("Seeding...")
	db.Transaction(func(tx *gorm.DB) error {
		err = seeds.Seed(db)

		if err != nil {
			logrus.Fatalf("error seeding the tables: %+v", err)
			return err
		}

		return nil
	})

}
