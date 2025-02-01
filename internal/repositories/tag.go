package repositories

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"gorm.io/gorm"
)

type RepositoryTag interface {
	FindOrCreate(userId string, tags []string) (*models.Tags, error)
}

type repositoryTag struct {
	db *gorm.DB
}

func (r *repositoryTag) FindOrCreate(userId string, tags []string) (*models.Tags, error) {

	var existingTags models.Tags
	r.db.Where("name in ?", tags).Find(&existingTags)

	existingTagNames := map[string]bool{}
	for _, tag := range existingTags {
		existingTagNames[tag.Name] = true
	}

	var newTags models.Tags
	for _, tagName := range tags {
		if !existingTagNames[tagName] {
			newTags = append(newTags, models.Tag{Name: tagName, Attributes: []byte("{}")})
		}
	}

	if len(newTags) > 0 {
		result := r.db.Create(&newTags)
		if result.Error != nil {
			panic(result.Error)
		}
	}

	combinedTags := append(existingTags, newTags...)

	return &combinedTags, nil
}

func NewTagRepository(db *gorm.DB) *repositoryTag {
	return &repositoryTag{
		db: db,
	}
}
