package services

import (
	"errors"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServicePortfolio interface {
	Get(userId string) (interface{}, error)
	UpsertSkills(userId string, data *schemas.SchemaSkills) error
	UpdateStatus(userId string, status string) error
}

type servicePortfolio struct {
	db *gorm.DB
}

func (s *servicePortfolio) Get(userId string) (interface{}, error) {
	userRepository := repositories.NewUserRepository(s.db)

	res, err := userRepository.GetPortfolio(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *servicePortfolio) UpsertSkills(userId string, data *schemas.SchemaSkills) error {
	repository := repositories.NewUserRepository(s.db)

	if err := repository.UpsertSkills(userId, data); err != nil {
		return err
	}
	return nil
}

func (s *servicePortfolio) UpdateStatus(userId string, status string) error {
	repository := repositories.NewUserRepository(s.db)

	switch status {
	case "publish":
		return repository.UpdateStatus(userId, "ACTIVE")
	case "takedown":
		return repository.UpdateStatus(userId, "IN_ACTIVE")
	default:
		return errors.New("invalid status")
	}
}

func NewPortfolioService(db *gorm.DB) *servicePortfolio {
	return &servicePortfolio{db: db}
}
