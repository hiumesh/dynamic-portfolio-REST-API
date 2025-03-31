package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceWorkGallery interface {
	GetAll(userId *string, query *string, cursor int, limit int) (any, error)
	GetUserWorkGallery(userId string, query *string, cursor int, limit int) (any, error)
	Get(userId string, id string) (any, error)
	Create(userId string, data *schemas.SchemaTechProject) (any, error)
	Update(userId string, id string, data *schemas.SchemaTechProject) (any, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
	GetMetadata(userId string) (any, error)
	UpdateMetadata(userId string, data *schemas.SchemaTechProjectMetadata) error
}

type serviceWorkGallery struct {
	db *gorm.DB
}

func (s *serviceWorkGallery) GetAll(userId *string, query *string, cursor int, limit int) (any, error) {
	userTechProjectRepository := repositories.NewUserTechProjectRepository(s.db)

	res, err := userTechProjectRepository.GetAll(userId, query, cursor, limit)
	if err != nil {
		return nil, err
	}

	var nextCursor *int
	if len(*res) == limit {
		temp := cursor + limit
		nextCursor = &temp
	}

	return map[string]any{"list": res, "cursor": nextCursor}, nil
}

func (s *serviceWorkGallery) GetUserWorkGallery(userId string, query *string, cursor int, limit int) (any, error) {
	userTechProjectRepository := repositories.NewUserTechProjectRepository(s.db)

	res, err := userTechProjectRepository.GetUserTechProjects(userId, query, cursor, limit)
	if err != nil {
		return nil, err
	}

	var nextCursor *int
	if len(*res) == limit {
		temp := cursor + limit
		nextCursor = &temp
	}

	return map[string]any{"list": res, "cursor": nextCursor}, nil
}

func (s *serviceWorkGallery) Get(userId string, id string) (any, error) {
	userTechProjectRepository := repositories.NewUserTechProjectRepository(s.db)

	res, err := userTechProjectRepository.Get(userId, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceWorkGallery) Create(userId string, data *schemas.SchemaTechProject) (any, error) {
	tx := s.db.Begin()
	userTechProjectRepository := repositories.NewUserTechProjectRepository(tx)
	userAttachmentRepository := repositories.NewAttachmentRepository(tx)

	tp, err := userTechProjectRepository.Create(userId, data)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var atts *models.Attachments = nil

	if len(data.Attachments) > 0 {
		atts, err = userAttachmentRepository.CreateMany(userId, models.TechProject{}.TableName(), tp.ID, &data.Attachments)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()
	response := struct {
		models.TechProject
		Attachments any `json:"attachments"`
	}{
		TechProject: *tp,
		Attachments: atts,
	}

	return response, nil

}

func (s *serviceWorkGallery) Update(userId string, id string, data *schemas.SchemaTechProject) (any, error) {
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
		Attachments any `json:"attachments"`
	}{
		TechProject: *tp,
		Attachments: atts,
	}

	return response, nil
}

func (s *serviceWorkGallery) Reorder(userId string, id string, newIndex int) error {
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

func (s *serviceWorkGallery) Delete(userId string, id string) error {
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

func (s *serviceWorkGallery) GetMetadata(userId string) (any, error) {
	userRepository := repositories.NewUserRepository(s.db)

	res, err := userRepository.GetModuleMetadata(userId, "work_gallery")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceWorkGallery) UpdateMetadata(userId string, data *schemas.SchemaTechProjectMetadata) error {
	userRepository := repositories.NewUserRepository(s.db)

	err := userRepository.AddOrUpdateModuleMetadata(userId, "work_gallery", data)
	if err != nil {
		return err
	}

	return nil
}

func NewWorkGalleryService(db *gorm.DB) *serviceWorkGallery {
	return &serviceWorkGallery{
		db: db,
	}
}
