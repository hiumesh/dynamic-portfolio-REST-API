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

type RepositoryUserExperience interface {
	GetAll(userId string) (*models.UserExperience, error)
	Create(userId string, data *schemas.SchemaUserExperience) (*models.UserExperience, error)
	Update(userId string, id string, data *schemas.SchemaUserExperience) (*models.UserExperience, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
}

type repositoryUserExperience struct {
	db *gorm.DB
}

func (r *repositoryUserExperience) GetAll(userId string) (*models.UserExperiences, error) {
	var userExperiences models.UserExperiences

	if err := r.db.Where("user_id = ?", userId).Order("order_index desc").Find(&userExperiences).Error; err != nil {
		return nil, err
	}

	return &userExperiences, nil
}

func (r *repositoryUserExperience) Create(userId string, data *schemas.SchemaUserExperience) (*models.UserExperience, error) {
	type MaxIndexResult struct {
		MaxIndex int16
	}
	maxIndexResult := MaxIndexResult{MaxIndex: 0}
	if err := r.db.Model(&models.UserExperience{}).Select("max(order_index) as max_index").Where("user_id = ?", userId).Group("user_id").Take(&maxIndexResult).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	experience := models.UserExperience{
		UserId:          userUUID,
		OrderIndex:      maxIndexResult.MaxIndex + 1,
		CompanyName:     data.CompanyName,
		CompanyUrl:      data.CompanyUrl,
		JobType:         data.JobType,
		JobTitle:        data.JobTitle,
		Location:        data.Location,
		StartDate:       startDate,
		Description:     data.Description,
		SkillsUsed:      data.SkillsUsed,
		CertificateLink: &data.CertificateLink,
	}

	if data.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", data.EndDate)

		if err != nil {
			return nil, errors.New("failed to parse end date")
		}
		experience.EndDate = &endDate
	}

	if err = r.db.Create(&experience).Error; err != nil {
		return nil, err
	}

	return &experience, nil

}

func (r *repositoryUserExperience) Update(userId string, id string, data *schemas.SchemaUserExperience) (*models.UserExperience, error) {
	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return nil, errors.New("failed to parse start date")
	}

	experience := models.UserExperience{
		CompanyName:     data.CompanyName,
		CompanyUrl:      data.CompanyUrl,
		JobType:         data.JobType,
		JobTitle:        data.JobTitle,
		Location:        data.Location,
		StartDate:       startDate,
		Description:     data.Description,
		SkillsUsed:      data.SkillsUsed,
		CertificateLink: &data.CertificateLink,
	}

	if data.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", data.EndDate)

		if err != nil {
			return nil, errors.New("failed to parse end date")
		}
		experience.EndDate = &endDate
	}

	var updatedRows models.UserExperiences

	if err := r.db.Model(&updatedRows).Clauses(clause.Returning{}).Where("id = ? and user_id = ?", id, userId).Updates(experience).Error; err != nil {
		return nil, err
	}

	if len(updatedRows) == 0 {
		return nil, errors.New("record not found")
	}

	return &updatedRows[0], nil
}

func (r *repositoryUserExperience) Delete(userId string, id string) error {
	var experience models.UserExperience

	if err := r.db.Where("user_id = ?", userId).Where("id = ?", id).First(&experience).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.UserExperience{}).Where("user_id = ? and order_index > ?", userId, experience.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Delete(&experience).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUserExperience) Reorder(userId string, id string, newIndex int) error {
	var experience models.UserExperience

	if err := r.db.First(&experience, id).Error; err != nil {
		return err
	}

	if experience.OrderIndex == int16(newIndex) {
		return nil
	}

	type CountResult struct {
		Count int16
	}
	countResult := CountResult{}
	if err := r.db.Model(&models.UserExperience{}).Select("count(*) as count").Where("user_id = ?", userId).Group("user_id").Take(&countResult).Error; err != nil {
		return err
	}

	if countResult.Count < int16(newIndex) {
		return errors.New("invalid index for reordering")
	}

	if err := r.db.Model(&models.UserExperience{}).Where("id = ?", experience.ID).UpdateColumn("order_index", 9999).Error; err != nil {
		return err
	}

	if experience.OrderIndex < int16(newIndex) {
		if err := r.db.Model(&models.UserExperience{}).Where("user_id = ? and order_index > ? and order_index <= ?", userId, experience.OrderIndex, newIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.UserExperience{}).Where("id = ?", experience.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	if experience.OrderIndex > int16(newIndex) {
		if err := r.db.Model(&models.UserExperience{}).Where("user_id = ? and order_index >= ? and order_index < ?", userId, newIndex, experience.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index + ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.UserExperience{}).Where("id = ?", experience.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	return nil
}

func NewUserExperienceRepository(db *gorm.DB) *repositoryUserExperience {
	return &repositoryUserExperience{
		db: db,
	}
}
