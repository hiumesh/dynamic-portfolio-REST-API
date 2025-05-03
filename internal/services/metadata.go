package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"gorm.io/gorm"
)

type ServiceMetadata interface {
	GetAllSkills(query *string, cursor int, limit int) (any, error)
}

type serviceMetadata struct {
	db *gorm.DB
}

func (s *serviceMetadata) GetAllSkills(query *string, cursor int, limit int) (any, error) {
	skillRepository := repositories.NewRepositorySkill(s.db)

	res, err := skillRepository.GetAll(query, cursor, limit)
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

func NewMetadataService(db *gorm.DB) *serviceMetadata {
	return &serviceMetadata{db: db}
}
