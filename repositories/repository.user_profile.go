package repositories

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"gorm.io/gorm"
)

type RepositoryUserProfile interface {
	Get(userId string) (*models.UserProfile, error)
	Upsert(userId string, profile *schemas.SchemaProfileBasic) error
}

type repositoryUserProfile struct {
	db *gorm.DB
}

func (r *repositoryUserProfile) Get(userId string) (*models.UserProfile, error) {
	var profile models.UserProfile

	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return &profile, nil
}

func (r *repositoryUserProfile) Upsert(userId string, profile *schemas.SchemaProfileBasic) error {
	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Update("attributes", *profile).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return nil
}

func NewUserProfileRepository(db *gorm.DB) *repositoryUserProfile {
	return &repositoryUserProfile{
		db: db,
	}
}
