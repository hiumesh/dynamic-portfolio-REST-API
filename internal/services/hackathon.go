package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceUserHackathon interface {
	GetAll(userId string) (*models.Hackathons, error)
	Create(userId string, data *schemas.SchemaHackathon) (*models.Hackathon, error)
	Update(userId string, id string, data *schemas.SchemaHackathon) (*models.Hackathon, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
	GetMetadata(userId string) (interface{}, error)
	UpdateMetadata(userId string, data *schemas.SchemaHackathonMetadata) error
}

type serviceUserHackathon struct {
	db *gorm.DB
}

func (s *serviceUserHackathon) GetAll(userId string) (*models.Hackathons, error) {
	userHackathonRepository := repositories.NewUserHackathonRepository(s.db)

	res, err := userHackathonRepository.GetAll(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserHackathon) Create(userId string, data *schemas.SchemaHackathon) (*models.Hackathon, error) {
	userHackathonRepository := repositories.NewUserHackathonRepository(s.db)

	exp, err := userHackathonRepository.Create(userId, data)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (s *serviceUserHackathon) Update(userId string, id string, data *schemas.SchemaHackathon) (*models.Hackathon, error) {
	userHackathonRepository := repositories.NewUserHackathonRepository(s.db)

	exp, err := userHackathonRepository.Update(userId, id, data)
	if err != nil {
		return nil, err
	}

	return exp, nil
}

func (s *serviceUserHackathon) Reorder(userId string, id string, newIndex int) error {

	err := s.db.Transaction(func(tx *gorm.DB) error {
		userHackathonRepository := repositories.NewUserHackathonRepository(tx)

		if err := userHackathonRepository.Reorder(userId, id, newIndex); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (s *serviceUserHackathon) Delete(userId string, id string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		userHackathonRepository := repositories.NewUserHackathonRepository(tx)

		if err := userHackathonRepository.Delete(userId, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serviceUserHackathon) GetMetadata(userId string) (interface{}, error) {
	userRepository := repositories.NewUserRepository(s.db)

	res, err := userRepository.GetModuleMetadata(userId, "hackathon")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceUserHackathon) UpdateMetadata(userId string, data *schemas.SchemaHackathonMetadata) error {
	userRepository := repositories.NewUserRepository(s.db)

	err := userRepository.AddOrUpdateModuleMetadata(userId, "hackathon", data)
	if err != nil {
		return err
	}

	return nil
}

func NewUserHackathonService(db *gorm.DB) *serviceUserHackathon {
	return &serviceUserHackathon{
		db: db,
	}
}
