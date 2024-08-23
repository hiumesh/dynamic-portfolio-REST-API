package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"gorm.io/gorm"
)

type ServiceUserEducation interface {
	GetAll(userId string) (*models.UserEducations, error)
	Create(userId string, data *schemas.SchemaUserEducation) error
	Update(userId string, id string, data *schemas.SchemaUserEducation) error
	Reorder(userId string, id string, newIndex int) error
}

type serviceUserEducation struct {
	db *gorm.DB
}

func (s *serviceUserEducation) GetAll(userId string) (*models.UserEducations, error) {
	userEducationRepository := repositories.NewUserEducationRepository(s.db)

	res, err := userEducationRepository.GetAll(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserEducation) Create(userId string, data *schemas.SchemaUserEducation) error {
	userEducationRepository := repositories.NewUserEducationRepository(s.db)

	if err := userEducationRepository.Create(userId, data); err != nil {
		return err
	}

	return nil
}

func (s *serviceUserEducation) Update(userId string, id string, data *schemas.SchemaUserEducation) error {
	userEducationRepository := repositories.NewUserEducationRepository(s.db)

	if err := userEducationRepository.Update(userId, id, data); err != nil {
		return err
	}

	return nil
}

func (s *serviceUserEducation) Reorder(userId string, id string, newIndex int) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userEducationRepository := repositories.NewUserEducationRepository(tx)

		if err := userEducationRepository.Reorder(userId, id, newIndex); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func NewUserEducationService(db *gorm.DB) *serviceUserEducation {
	return &serviceUserEducation{
		db: db,
	}
}
