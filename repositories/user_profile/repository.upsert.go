package repositories

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"gorm.io/gorm"
)

type RepositoryUpsertProfile interface {
	UpsertProfileRepository(userId string, profile *schemas.SchemaProfileBasic) *helpers.DatabaseError
}

type repositoryUpsertProfile struct {
	db *gorm.DB
}

func (r *repositoryUpsertProfile) UpsertProfileRepository(userId string, profile *schemas.SchemaProfileBasic) *helpers.DatabaseError {
	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Update("attributes", *profile).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return nil
}

func NewUpsertProfileRepository(db *gorm.DB) *repositoryUpsertProfile {
	return &repositoryUpsertProfile{
		db: db,
	}
}
