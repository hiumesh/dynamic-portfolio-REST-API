package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceUserCertification interface {
	GetAll(userId string) (*models.UserCertifications, error)
	Create(userId string, data *schemas.SchemaUserCertification) (*models.UserCertification, error)
	Update(userId string, id string, data *schemas.SchemaUserCertification) (*models.UserCertification, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
	UpdateMetadata(userId string, data *schemas.SchemaUserCertificationMetadata) error
}

type serviceUserCertification struct {
	db *gorm.DB
}

func (s *serviceUserCertification) GetAll(userId string) (*models.UserCertifications, error) {
	userExperienceRepository := repositories.NewUserCertificationRepository(s.db)

	res, err := userExperienceRepository.GetAll(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserCertification) Create(userId string, data *schemas.SchemaUserCertification) (*models.UserCertification, error) {
	userExperienceRepository := repositories.NewUserCertificationRepository(s.db)

	res, err := userExperienceRepository.Create(userId, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserCertification) Update(userId string, id string, data *schemas.SchemaUserCertification) (*models.UserCertification, error) {
	userExperienceRepository := repositories.NewUserCertificationRepository(s.db)

	res, err := userExperienceRepository.Update(userId, id, data)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserCertification) Reorder(userId string, id string, newIndex int) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userExperienceRepository := repositories.NewUserCertificationRepository(tx)

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

func (s *serviceUserCertification) Delete(userId string, id string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		userExperienceRepository := repositories.NewUserCertificationRepository(tx)

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

func (s *serviceUserCertification) UpdateMetadata(userId string, data *schemas.SchemaUserCertificationMetadata) error {
	userRepository := repositories.NewUserRepository(s.db)

	err := userRepository.AddOrUpdateModuleMetadata(userId, "certification", data)
	if err != nil {
		return err
	}

	return nil
}

func NewUserCertificationService(db *gorm.DB) *serviceUserCertification {
	return &serviceUserCertification{
		db: db,
	}
}
