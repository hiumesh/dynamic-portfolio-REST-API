package services

import (
	"errors"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServicePortfolio interface {
	GetAll(userId *string, query *string, cursor int, limit int) (interface{}, error)
	GetPortfolio(slug string) (interface{}, error)
	GetSubModule(slug string, module string) (interface{}, error)
	GetUserPortfolio(userId string) (interface{}, error)
	UpsertSkills(userId string, data *schemas.SchemaSkills) error
	UpdateStatus(userId string, status string) error
}

type servicePortfolio struct {
	db *gorm.DB
}

func (s *servicePortfolio) GetAll(userId *string, query *string, cursor int, limit int) (interface{}, error) {
	portfolioRepository := repositories.NewPortfolioRepository(s.db)

	res, err := portfolioRepository.GetAll(userId, query, cursor, limit)
	if err != nil {
		return nil, err
	}

	var nextCursor *int
	if len(*res) == limit {
		temp := cursor + limit
		nextCursor = &temp
	}

	return map[string]interface{}{"list": res, "cursor": nextCursor}, nil
}

func (s *servicePortfolio) GetPortfolio(slug string) (interface{}, error) {
	portfolioRepository := repositories.NewPortfolioRepository(s.db)

	res, err := portfolioRepository.GetPortfolio(slug)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *servicePortfolio) GetSubModule(slug string, module string) (interface{}, error) {
	portfolioRepository := repositories.NewPortfolioRepository(s.db)

	var res any
	var err error
	switch module {
	case "educations":
		res, err = portfolioRepository.GetEducations(slug)
	case "work_experiences":
		res, err = portfolioRepository.GetWorkExperiences(slug)
	case "certifications":
		res, err = portfolioRepository.GetCertifications(slug)
	case "hackathons":
		res, err = portfolioRepository.GetHackathons(slug)
	case "works":
		res, err = portfolioRepository.GetTechProjects(slug)
	default:
		return nil, errors.New("invalid module")
	}

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *servicePortfolio) GetUserPortfolio(userId string) (interface{}, error) {
	portfolioRepository := repositories.NewPortfolioRepository(s.db)

	res, err := portfolioRepository.GetUserPortfolio(userId)
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
