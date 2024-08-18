package repositories

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	"gorm.io/gorm"
)

type RepositoryGetProfile interface {
	GetProfileRepository(userId string) (*models.UserProfile, *helpers.DatabaseError)
}

type repositoryGetProfile struct {
	db *gorm.DB
}

func NewGetProfileRepository(db *gorm.DB) *repositoryGetProfile {
	return &repositoryGetProfile{
		db: db,
	}
}

func (r *repositoryGetProfile) GetProfileRepository(userId string) (*models.UserProfile, *helpers.DatabaseError) {
	var profile models.UserProfile

	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return &profile, nil
}
