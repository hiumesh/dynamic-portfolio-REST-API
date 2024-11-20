package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceUser interface {
	GetProfile(userId string) (*models.UserProfile, error)
	UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error
	ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error
}

type serviceUser struct {
	db *gorm.DB
}

func (s *serviceUser) GetProfile(userId string) (*models.UserProfile, error) {
	repository := repositories.NewUserRepository(s.db)

	res, err := repository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *serviceUser) ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error {
	repository := repositories.NewUserRepository(s.db)

	if err := repository.ProfileSetup(userId, profile); err != nil {
		return err
	}
	return nil
}

func (s *serviceUser) UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error {
	repository := repositories.NewUserRepository(s.db)

	if err := repository.UpsertProfile(userId, profile); err != nil {
		return err
	}
	return nil
}

func NewUserService(db *gorm.DB) *serviceUser {
	return &serviceUser{
		db: db,
	}
}
