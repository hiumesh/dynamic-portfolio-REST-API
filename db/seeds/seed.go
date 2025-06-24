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

	logrus.Info("Seeding tags...")
	tags, err := SeedTags(db, 1000)
	if err != nil {
		return err
	}
	logrus.Info("Seeding blogs...")
	err = SeedBlogs(db, users, tags, 30)
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
	if err := db.Unscoped().Where("1 = 1").Delete(&BlogComment{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&BlogReaction{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&BlogTag{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&Blog{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&TechProject{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&Education{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&Hackathon{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&Certification{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&Experience{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&Comment{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&UserProfile{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&Tag{}).Error; err != nil {
		return err
	}
	if err := db.Unscoped().Where("1 = 1").Delete(&User{}).Error; err != nil {
		return err
	}

	return nil
}
