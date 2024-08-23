package repositories

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"gorm.io/gorm"
)

type RepositoryUserEducation interface {
	GetAll(userId string) (*models.UserEducations, error)
	Create(userId string, data *schemas.SchemaUserEducation) error
	Update(userId string, id string, data *schemas.SchemaUserEducation) error
	Reorder(userId string, id string, newIndex int) error
}

type repositoryUserEducation struct {
	db *gorm.DB
}

func (r *repositoryUserEducation) GetAll(userId string) (*models.UserEducations, error) {
	var userEducations models.UserEducations

	if err := r.db.Where("user_id = ?", userId).Order("order_index asc").Find(&userEducations).Error; err != nil {
		return nil, &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return &userEducations, nil
}

func (r *repositoryUserEducation) Create(userId string, data *schemas.SchemaUserEducation) error {
	type MaxIndexResult struct {
		MaxIndex int16
	}
	maxIndexResult := MaxIndexResult{MaxIndex: 0}
	if err := r.db.Model(&models.UserEducation{}).Select("max(order_index) as max_index").Where("user_id = ?", userId).Group("user_id").Take(&maxIndexResult).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("failed to parse user id")
	}
	edu := models.UserEducation{
		UserId:        userUUID,
		InstituteName: data.InstituteName,
		Grade:         data.Grade,
		OrderIndex:    maxIndexResult.MaxIndex + 1,
	}

	if data.Type == "SCHOOL" {
		attributes := map[string]string{
			"class":        data.Class,
			"passing_year": data.PassingYear,
		}
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return errors.New("failed to marshal json")
		}
		edu.Type = data.Type
		edu.Attributes = attributesJson
	}

	if data.Type == "COLLAGE" {
		attributes := map[string]string{
			"field_of_study": data.FieldOfStudy,
			"start_year":     data.StartYear,
			"end_year":       data.EndYear,
		}
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return errors.New("failed to marshal json")
		}
		edu.Type = data.Type
		edu.Attributes = attributesJson
	}

	if err = r.db.Create(&edu).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return nil
}

func (r *repositoryUserEducation) Update(userId string, id string, data *schemas.SchemaUserEducation) error {
	education := models.UserEducation{
		InstituteName: data.InstituteName,
		Grade:         data.Grade,
	}

	if data.Type == "SCHOOL" {
		attributes := map[string]string{
			"class":        data.Class,
			"passing_year": data.PassingYear,
		}
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return errors.New("failed to marshal json")
		}
		education.Type = data.Type
		education.Attributes = attributesJson
	}

	if data.Type == "COLLAGE" {
		attributes := map[string]string{
			"field_of_study": data.FieldOfStudy,
			"start_year":     data.StartYear,
			"end_year":       data.EndYear,
		}
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return errors.New("failed to marshal json")
		}
		education.Type = data.Type
		education.Attributes = attributesJson
	}

	if err := r.db.Model(&models.UserEducation{}).Where("id = ? and user_id = ?", id, userId).Updates(education).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return nil
}

func (r *repositoryUserEducation) Reorder(userId string, id string, newIndex int) error {
	var education models.UserEducation

	if err := r.db.First(&education, id).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	if education.OrderIndex == int16(newIndex) {
		return nil
	}

	type CountResult struct {
		Count int16
	}
	countResult := CountResult{}
	if err := r.db.Model(&models.UserEducation{}).Select("count(*) as count").Where("user_id = ?", userId).Group("user_id").Take(&countResult).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	if countResult.Count < int16(newIndex) {
		return errors.New("invalid index for reordering")
	}

	r.db.Model(&models.UserEducation{}).Where("id = ?", education.ID).UpdateColumn("order_index", 9999)

	if education.OrderIndex < int16(newIndex) {
		r.db.Model(&models.UserEducation{}).Where("user_id = ? and order_index > ? and order_index <= ?", userId, education.OrderIndex, newIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1))
		r.db.Model(&models.UserEducation{}).Where("id = ?", education.ID).UpdateColumn("order_index", newIndex)
	}

	if education.OrderIndex > int16(newIndex) {
		r.db.Model(&models.UserEducation{}).Where("user_id = ? and order_index >= ? and order_index < ?", userId, newIndex, education.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index + ?", 1))
		r.db.Model(&models.UserEducation{}).Where("id = ?", education.ID).UpdateColumn("order_index", newIndex)
	}

	return nil
}

func NewUserEducationRepository(db *gorm.DB) *repositoryUserEducation {
	return &repositoryUserEducation{
		db: db,
	}
}
