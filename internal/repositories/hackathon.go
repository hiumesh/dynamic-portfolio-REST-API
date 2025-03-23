package repositories

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryUserHackathon interface {
	GetAll(userId string) (*models.Hackathons, error)
	Create(userId string, data *schemas.SchemaHackathon) (*models.Hackathon, error)
	Update(userId string, id string, data *schemas.SchemaHackathon) (*models.Hackathon, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
}

type repositoryUserHackathon struct {
	db *gorm.DB
}

func (r *repositoryUserHackathon) GetAll(userId string) (*models.Hackathons, error) {
	var userHackathons models.Hackathons

	if err := r.db.Where("user_id = ?", userId).Order("order_index desc").Find(&userHackathons).Error; err != nil {
		return nil, err
	}

	return &userHackathons, nil
}

func (r *repositoryUserHackathon) Create(userId string, data *schemas.SchemaHackathon) (*models.Hackathon, error) {
	type MaxIndexResult struct {
		MaxIndex int16
	}
	maxIndexResult := MaxIndexResult{MaxIndex: 0}
	if err := r.db.Model(&models.Hackathon{}).Select("max(order_index) as max_index").Where("user_id = ?", userId).Group("user_id").Take(&maxIndexResult).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return nil, errors.New("failed to parse start date")
	}

	endDate, err := time.Parse("2006-01-02", data.EndDate)
	if err != nil {
		return nil, errors.New("failed to parse end date")
	}

	hackathon := models.Hackathon{
		UserId:          userUUID,
		OrderIndex:      maxIndexResult.MaxIndex + 1,
		Avatar:          &data.Avatar,
		Title:           data.Title,
		Location:        data.Location,
		StartDate:       startDate,
		EndDate:         endDate,
		Description:     data.Description,
		CertificateLink: &data.CertificateLink,
	}

	if data.Links != nil && len(data.Links) > 0 {
		attributes := map[string]interface{}{
			"links": data.Links,
		}

		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return nil, errors.New("failed to marshal json")
		}

		hackathon.Attributes = attributesJson
	} else {
		hackathon.Attributes = []byte("{}")
	}

	if err = r.db.Create(&hackathon).Error; err != nil {
		return nil, err
	}

	return &hackathon, nil

}

func (r *repositoryUserHackathon) Update(userId string, id string, data *schemas.SchemaHackathon) (*models.Hackathon, error) {
	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return nil, errors.New("failed to parse start date")
	}

	endDate, err := time.Parse("2006-01-02", data.EndDate)
	if err != nil {
		return nil, errors.New("failed to parse end date")
	}

	hackathon := map[string]interface{}{
		"avatar":           data.Avatar,
		"title":            data.Title,
		"location":         data.Location,
		"start_date":       startDate,
		"end_date":         endDate,
		"description":      data.Description,
		"certificate_link": data.CertificateLink,
		"attributes": datatypes.JSONSet("attributes").
			Set("{links}", data.Links),
	}

	var updatedRows models.Hackathons

	if err := r.db.Model(&updatedRows).Clauses(clause.Returning{}).Where("id = ? and user_id = ?", id, userId).Updates(hackathon).Error; err != nil {
		return nil, err
	}

	if len(updatedRows) == 0 {
		return nil, errors.New("record not found")
	}

	return &updatedRows[0], nil
}

func (r *repositoryUserHackathon) Delete(userId string, id string) error {
	var hackathon models.Hackathon

	if err := r.db.Where("user_id = ?", userId).Where("id = ?", id).First(&hackathon).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.Hackathon{}).Where("user_id = ? and order_index > ?", userId, hackathon.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Delete(&hackathon).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUserHackathon) Reorder(userId string, id string, newIndex int) error {
	var hackathon models.Hackathon

	if err := r.db.First(&hackathon, id).Error; err != nil {
		return err
	}

	if hackathon.OrderIndex == int16(newIndex) {
		return nil
	}

	type CountResult struct {
		Count int16
	}
	countResult := CountResult{}
	if err := r.db.Model(&models.Hackathon{}).Select("count(*) as count").Where("user_id = ?", userId).Group("user_id").Take(&countResult).Error; err != nil {
		return err
	}

	if countResult.Count < int16(newIndex) {
		return errors.New("invalid index for reordering")
	}

	if hackathon.OrderIndex < int16(newIndex) {
		if err := r.db.Model(&models.Hackathon{}).Where("user_id = ? and order_index > ? and order_index <= ?", userId, hackathon.OrderIndex, newIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.Hackathon{}).Where("id = ?", hackathon.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	if hackathon.OrderIndex > int16(newIndex) {
		if err := r.db.Model(&models.Hackathon{}).Where("user_id = ? and order_index >= ? and order_index < ?", userId, newIndex, hackathon.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index + ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.Hackathon{}).Where("id = ?", hackathon.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	return nil
}

func NewUserHackathonRepository(db *gorm.DB) *repositoryUserHackathon {
	return &repositoryUserHackathon{
		db: db,
	}
}
