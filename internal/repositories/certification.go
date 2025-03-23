package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryUserCertification interface {
	GetAll(userId string) (*models.Certifications, error)
	Create(userId string, data *schemas.SchemaCertification) (*models.Certification, error)
	Update(userId string, id string, data *schemas.SchemaCertification) (*models.Certification, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
}

type repositoryUserCertification struct {
	db *gorm.DB
}

func (r *repositoryUserCertification) GetAll(userId string) (*models.Certifications, error) {
	var userCertifications models.Certifications

	if err := r.db.Where("user_id = ?", userId).Order("order_index desc").Find(&userCertifications).Error; err != nil {
		return nil, err
	}

	return &userCertifications, nil
}

func (r *repositoryUserCertification) Create(userId string, data *schemas.SchemaCertification) (*models.Certification, error) {
	type MaxIndexResult struct {
		MaxIndex int16
	}
	maxIndexResult := MaxIndexResult{MaxIndex: 0}
	if err := r.db.Model(&models.Certification{}).Select("max(order_index) as max_index").Where("user_id = ?", userId).Group("user_id").Take(&maxIndexResult).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	completionDate, err := time.Parse("2006-01-02", data.CompletionDate)
	if err != nil {
		return nil, errors.New("failed to parse completion date")
	}

	ctl := models.Certification{
		UserId:          userUUID,
		OrderIndex:      maxIndexResult.MaxIndex + 1,
		Title:           data.Title,
		Description:     data.Description,
		CompletionDate:  completionDate,
		CertificateLink: &data.CertificateLink,
		SkillsUsed:      data.SkillsUsed,
	}

	if err = r.db.Create(&ctl).Error; err != nil {
		return nil, err
	}

	return &ctl, nil
}

func (r *repositoryUserCertification) Update(userId string, id string, data *schemas.SchemaCertification) (*models.Certification, error) {
	completionDate, err := time.Parse("2006-01-02", data.CompletionDate)
	if err != nil {
		return nil, errors.New("failed to parse completion date")
	}

	certificate := models.Certification{
		Title:           data.Title,
		Description:     data.Description,
		CompletionDate:  completionDate,
		CertificateLink: &data.CertificateLink,
		SkillsUsed:      data.SkillsUsed,
	}

	var updatedRows models.Certifications

	if err := r.db.Model(&updatedRows).Clauses(clause.Returning{}).Where("id = ? and user_id = ?", id, userId).Updates(certificate).Error; err != nil {
		return nil, err
	}

	if len(updatedRows) == 0 {
		return nil, errors.New("record not found")
	}

	return &updatedRows[0], nil
}

func (r *repositoryUserCertification) Reorder(userId string, id string, newIndex int) error {
	var certificate models.Certification

	if err := r.db.First(&certificate, id).Error; err != nil {
		return err
	}

	if certificate.OrderIndex == int16(newIndex) {
		return nil
	}

	type CountResult struct {
		Count int16
	}
	countResult := CountResult{}
	if err := r.db.Model(&models.Certification{}).Select("count(*) as count").Where("user_id = ?", userId).Group("user_id").Take(&countResult).Error; err != nil {
		return err
	}

	if countResult.Count < int16(newIndex) {
		return errors.New("invalid index for reordering")
	}

	if certificate.OrderIndex < int16(newIndex) {
		if err := r.db.Model(&models.Certification{}).Where("user_id = ? and order_index > ? and order_index <= ?", userId, certificate.OrderIndex, newIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.Certification{}).Where("id = ?", certificate.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	if certificate.OrderIndex > int16(newIndex) {
		if err := r.db.Model(&models.Certification{}).Where("user_id = ? and order_index >= ? and order_index < ?", userId, newIndex, certificate.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index + ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.Certification{}).Where("id = ?", certificate.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *repositoryUserCertification) Delete(userId string, id string) error {

	var certificate models.Certification

	if err := r.db.Where("user_id = ?", userId).Where("id = ?", id).First(&certificate).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.Certification{}).Where("user_id = ? and order_index > ?", userId, certificate.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Delete(&certificate).Error; err != nil {
		return err
	}

	return nil
}

func NewUserCertificationRepository(db *gorm.DB) *repositoryUserCertification {
	return &repositoryUserCertification{
		db: db,
	}
}
