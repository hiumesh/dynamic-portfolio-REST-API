package repositories

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryUserEducation interface {
	GetAll(userId string) (*models.Educations, error)
	Create(userId string, data *schemas.SchemaEducation) (*models.Education, error)
	Update(userId string, id string, data *schemas.SchemaEducation) (*models.Education, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
}

type repositoryUserEducation struct {
	db *gorm.DB
}

func (r *repositoryUserEducation) GetAll(userId string) (*models.Educations, error) {
	var userEducations models.Educations

	if err := r.db.Where("user_id = ?", userId).Order("order_index desc").Find(&userEducations).Error; err != nil {
		return nil, err
	}

	return &userEducations, nil
}

func (r *repositoryUserEducation) Create(userId string, data *schemas.SchemaEducation) (*models.Education, error) {
	type MaxIndexResult struct {
		MaxIndex int16
	}
	maxIndexResult := MaxIndexResult{MaxIndex: 0}
	if err := r.db.Model(&models.Education{}).Select("max(order_index) as max_index").Where("user_id = ?", userId).Group("user_id").Take(&maxIndexResult).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}
	edu := models.Education{
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
			return nil, errors.New("failed to marshal json")
		}
		edu.Type = data.Type
		edu.Attributes = attributesJson
	}

	if data.Type == "COLLEGE" {
		attributes := map[string]string{
			"degree":         data.Degree,
			"field_of_study": data.FieldOfStudy,
			"start_year":     data.StartYear,
			"end_year":       data.EndYear,
		}
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return nil, errors.New("failed to marshal json")
		}
		edu.Type = data.Type
		edu.Attributes = attributesJson
	}

	if err = r.db.Create(&edu).Error; err != nil {
		return nil, err
	}

	return &edu, nil
}

func (r *repositoryUserEducation) Update(userId string, id string, data *schemas.SchemaEducation) (*models.Education, error) {
	education := models.Education{
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
			return nil, errors.New("failed to marshal json")
		}
		education.Type = data.Type
		education.Attributes = attributesJson
	}

	if data.Type == "COLLEGE" {
		attributes := map[string]string{
			"degree":         data.Degree,
			"field_of_study": data.FieldOfStudy,
			"start_year":     data.StartYear,
			"end_year":       data.EndYear,
		}
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return nil, errors.New("failed to marshal json")
		}
		education.Type = data.Type
		education.Attributes = attributesJson
	}

	var updatedRows models.Educations

	if err := r.db.Model(&updatedRows).Clauses(clause.Returning{}).Where("id = ? and user_id = ?", id, userId).Updates(education).Error; err != nil {
		return nil, err
	}

	if len(updatedRows) == 0 {
		return nil, errors.New("record not found")
	}

	return &updatedRows[0], nil
}

func (r *repositoryUserEducation) Reorder(userId string, id string, newIndex int) error {
	var education models.Education

	if err := r.db.First(&education, id).Error; err != nil {
		return err
	}

	if education.OrderIndex == int16(newIndex) {
		return nil
	}

	type CountResult struct {
		Count int16
	}
	countResult := CountResult{}
	if err := r.db.Model(&models.Education{}).Select("count(*) as count").Where("user_id = ?", userId).Group("user_id").Take(&countResult).Error; err != nil {
		return err
	}

	if countResult.Count < int16(newIndex) {
		return errors.New("invalid index for reordering")
	}

	if education.OrderIndex < int16(newIndex) {
		if err := r.db.Model(&models.Education{}).Where("user_id = ? and order_index > ? and order_index <= ?", userId, education.OrderIndex, newIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.Education{}).Where("id = ?", education.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	if education.OrderIndex > int16(newIndex) {
		if err := r.db.Model(&models.Education{}).Where("user_id = ? and order_index >= ? and order_index < ?", userId, newIndex, education.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index + ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.Education{}).Where("id = ?", education.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *repositoryUserEducation) Delete(userId string, id string) error {

	var education models.Education

	if err := r.db.Where("user_id = ?", userId).Where("id = ?", id).First(&education).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.Education{}).Where("user_id = ? and order_index > ?", userId, education.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Delete(&education).Error; err != nil {
		return err
	}

	return nil
}

func NewUserEducationRepository(db *gorm.DB) *repositoryUserEducation {
	return &repositoryUserEducation{
		db: db,
	}
}
