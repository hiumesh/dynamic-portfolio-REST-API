package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"gorm.io/gorm"
)

type ServiceUserProfile interface {
	Get(userId string) (*models.UserProfile, error)
	Upsert(userId string, profile *schemas.SchemaProfileBasic) error
}

type serviceUserProfile struct {
	db *gorm.DB
}

func (s *serviceUserProfile) Get(userId string) (*models.UserProfile, error) {
	repository := repositories.NewUserProfileRepository(s.db)

	res, err := repository.Get(userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *serviceUserProfile) Upsert(userId string, profile *schemas.SchemaProfileBasic) error {
	repository := repositories.NewUserProfileRepository(s.db)

	if err := repository.Upsert(userId, profile); err != nil {
		return err
	}
	return nil
}

func NewUserProfileService(db *gorm.DB) *serviceUserProfile {
	return &serviceUserProfile{
		db: db,
	}
}
