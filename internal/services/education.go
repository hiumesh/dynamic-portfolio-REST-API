package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceUserEducation interface {
	GetAll(userId string) (*models.Educations, error)
	Create(userId string, data *schemas.SchemaEducation) (*models.Education, error)
	Update(userId string, id string, data *schemas.SchemaEducation) (*models.Education, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
	GetMetadata(userId string) (interface{}, error)
	UpdateMetadata(userId string, data *schemas.SchemaEducationMetadata) error
}

type serviceUserEducation struct {
	db *gorm.DB
}

func (s *serviceUserEducation) GetAll(userId string) (*models.Educations, error) {
	userEducationRepository := repositories.NewUserEducationRepository(s.db)

	res, err := userEducationRepository.GetAll(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserEducation) Create(userId string, data *schemas.SchemaEducation) (*models.Education, error) {
	userEducationRepository := repositories.NewUserEducationRepository(s.db)

	edu, err := userEducationRepository.Create(userId, data)
	if err != nil {
		return nil, err
	}

	return edu, nil
}

func (s *serviceUserEducation) Update(userId string, id string, data *schemas.SchemaEducation) (*models.Education, error) {
	userEducationRepository := repositories.NewUserEducationRepository(s.db)

	edu, err := userEducationRepository.Update(userId, id, data)
	if err != nil {
		return nil, err
	}

	return edu, nil
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

func (s *serviceUserEducation) Delete(userId string, id string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		userEducationRepository := repositories.NewUserEducationRepository(tx)

		if err := userEducationRepository.Delete(userId, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serviceUserEducation) GetMetadata(userId string) (interface{}, error) {
	userRepository := repositories.NewUserRepository(s.db)

	res, err := userRepository.GetModuleMetadata(userId, "education")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserEducation) UpdateMetadata(userId string, data *schemas.SchemaEducationMetadata) error {
	userRepository := repositories.NewUserRepository(s.db)

	err := userRepository.AddOrUpdateModuleMetadata(userId, "education", data)
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
