package seeds

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {
	gofakeit.Seed(0)

	logrus.Info("Seeding users...")
	users, err := SeedUsers(db, 100)
	if err != nil {
		return err
	}

	logrus.Info("Seeding tags and blogs...")
	tags, err := SeedTags(db, 1000)
	if err != nil {
		return err
	}
	err = SeedBlogs(db, users, tags, 500)
	if err != nil {
		return err
	}

	logrus.Info("Seeding portfolio...")
	err = SeedPortfolio(db, users)
	if err != nil {
		return err
	}

	logrus.Info("Seeding tech projects...")
	err = SeedTechProjects(db, users)
	if err != nil {
		return err
	}
	return nil
}

func Truncate(db *gorm.DB) error {
	if db.Unscoped().Where("1 = 1").Delete(&User{}).Error != nil {
		return db.Error
	}
	if db.Unscoped().Where("1 = 1").Delete(&Tag{}).Error != nil {
		return db.Error
	}
	return nil
}
