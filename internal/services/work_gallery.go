package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceWorkGallery interface {
	GetAll(userId string) (interface{}, error)
	Create(userId string, data *schemas.SchemaTechProject) (interface{}, error)
	Update(userId string, id string, data *schemas.SchemaTechProject) (interface{}, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
	GetMetadata(userId string) (interface{}, error)
	UpdateMetadata(userId string, data *schemas.SchemaTechProjectMetadata) error
}

type service struct {
	db *gorm.DB
}

func (s *service) GetAll(userId string) (interface{}, error) {
	userTechProjectRepository := repositories.NewUserTechProjectRepository(s.db)

	res, err := userTechProjectRepository.GetAll(userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) Create(userId string, data *schemas.SchemaTechProject) (interface{}, error) {
	tx := s.db.Begin()
	userTechProjectRepository := repositories.NewUserTechProjectRepository(tx)
	userAttachmentRepository := repositories.NewAttachmentRepository(tx)

	tp, err := userTechProjectRepository.Create(userId, data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	atts, err := userAttachmentRepository.CreateMany(userId, models.TechProject{}.TableName(), tp.ID, &data.Attachments)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	response := struct {
		models.TechProject
		Attachments interface{} `json:"attachments"`
	}{
		TechProject: *tp,
		Attachments: atts,
	}

	return response, nil

}

func (s *service) Update(userId string, id string, data *schemas.SchemaTechProject) (interface{}, error) {
	tx := s.db.Begin()
	userTechProjectRepository := repositories.NewUserTechProjectRepository(tx)
	userAttachmentRepository := repositories.NewAttachmentRepository(tx)

	tp, err := userTechProjectRepository.Update(userId, id, data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	atts, err := userAttachmentRepository.UpdateOrCreate(userId, models.TechProject{}.TableName(), tp.ID, &data.Attachments)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	response := struct {
		models.TechProject
		Attachments interface{} `json:"attachments"`
	}{
		TechProject: *tp,
		Attachments: atts,
	}

	return response, nil
}

func (s *service) Reorder(userId string, id string, newIndex int) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		userTechProjectRepository := repositories.NewUserTechProjectRepository(tx)

		if err := userTechProjectRepository.Reorder(userId, id, newIndex); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil

}

func (s *service) Delete(userId string, id string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		userTechProjectRepository := repositories.NewUserTechProjectRepository(tx)
		userAttachmentRepository := repositories.NewAttachmentRepository(tx)

		if err := userTechProjectRepository.Delete(userId, id); err != nil {
			return err
		}

		if err := userAttachmentRepository.DeleteMany(userId, models.TechProject{}.TableName(), id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetMetadata(userId string) (interface{}, error) {
	userRepository := repositories.NewUserRepository(s.db)

	res, err := userRepository.GetModuleMetadata(userId, "work_gallery")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *service) UpdateMetadata(userId string, data *schemas.SchemaTechProjectMetadata) error {
	userRepository := repositories.NewUserRepository(s.db)

	err := userRepository.AddOrUpdateModuleMetadata(userId, "work_gallery", data)
	if err != nil {
		return err
	}

	return nil
}

func NewWorkGalleryService(db *gorm.DB) *service {
	return &service{
		db: db,
	}
}
