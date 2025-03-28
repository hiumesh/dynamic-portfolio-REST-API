package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceUserExperience interface {
	GetAll(userId string) (*models.WorkExperiences, error)
	Create(userId string, data *schemas.SchemaWorkExperience) (*models.WorkExperience, error)
	Update(userId string, id string, data *schemas.SchemaWorkExperience) (*models.WorkExperience, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
	GetMetadata(userId string) (interface{}, error)
	UpdateMetadata(userId string, data *schemas.SchemaWorkExperienceMetadata) error
}

type serviceUserExperience struct {
	db *gorm.DB
}

func (s *serviceUserExperience) GetAll(userId string) (*models.WorkExperiences, error) {
	userExperienceRepository := repositories.NewUserExperienceRepository(s.db)

	res, err := userExperienceRepository.GetAll(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserExperience) Create(userId string, data *schemas.SchemaWorkExperience) (*models.WorkExperience, error) {
	userExperienceRepository := repositories.NewUserExperienceRepository(s.db)

	exp, err := userExperienceRepository.Create(userId, data)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (s *serviceUserExperience) Update(userId string, id string, data *schemas.SchemaWorkExperience) (*models.WorkExperience, error) {
	userExperienceRepository := repositories.NewUserExperienceRepository(s.db)

	exp, err := userExperienceRepository.Update(userId, id, data)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (s *serviceUserExperience) Reorder(userId string, id string, newIndex int) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userExperienceRepository := repositories.NewUserExperienceRepository(tx)

		if err := userExperienceRepository.Reorder(userId, id, newIndex); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (s *serviceUserExperience) Delete(userId string, id string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		userExperienceRepository := repositories.NewUserExperienceRepository(tx)

		if err := userExperienceRepository.Delete(userId, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serviceUserExperience) GetMetadata(userId string) (interface{}, error) {
	userRepository := repositories.NewUserRepository(s.db)

	res, err := userRepository.GetModuleMetadata(userId, "work_experience")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserExperience) UpdateMetadata(userId string, data *schemas.SchemaWorkExperienceMetadata) error {
	userRepository := repositories.NewUserRepository(s.db)

	err := userRepository.AddOrUpdateModuleMetadata(userId, "work_experience", data)
	if err != nil {
		return err
	}

	return nil
}

func NewUserExperienceService(db *gorm.DB) *serviceUserExperience {
	return &serviceUserExperience{
		db: db,
	}
}
