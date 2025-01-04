package repositories

import (
	"strings"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RepositoryUser interface {
	GetProfile(userId string) (*models.UserProfile, error)
	UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error
	ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error
	UpsertSkills(userId string, data *schemas.SchemaSkills) error
	AddOrUpdateModuleMetadata(userId string, module string, data *interface{}) error
}

type repositoryUser struct {
	db *gorm.DB
}

func (r *repositoryUser) GetProfile(userId string) (*models.UserProfile, error) {
	var profile models.UserProfile

	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *repositoryUser) ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error {
	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"full_name": profile.FullName,
		"slug":      strings.Join(strings.Fields(strings.ToLower(profile.FullName)), "-") + "-" + strings.Split(userId, "-")[0],
		// "attributes": gorm.Expr("JSONB_SET(attributes, '{college}', to_jsonb(?))", profile.College),
		"attributes": datatypes.JSONSet("attributes").Set("{college}", profile.College).Set("{graduation_year}", profile.GraduationYear).Set("{work_domains}", profile.WorkDomains),
		// "attributes": datatypes.JSONSet("attributes").Set("college", profile.College).Set("graduation_year", profile.GraduationYear).Set("work_domains", profile.WorkDomains),
	}).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error {
	data := map[string]interface{}{
		"full_name":  profile.FullName,
		"avatar_url": profile.ProfilePicture,
		"attributes": datatypes.JSONSet("attributes").
			Set("{about}", profile.About).
			Set("{tagline}", profile.Tagline).
			Set("{college}", profile.College).
			Set("{graduation_year}", profile.GraduationYear).
			Set("{work_domains}", profile.WorkDomains).
			Set("{social_profiles}", profile.SocialProfileLinks),
		// "attributes": map[string]interface{}{
		// 	"college":         profile.College,
		// 	"graduation_year": profile.GraduationYear,
		// 	"work_domains":    profile.WorkDomains,
		// 	"social_profiles": profile.SocialProfileLinks,
		// },
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) UpsertSkills(userId string, data *schemas.SchemaSkills) error {
	t := map[string]interface{}{
		"attributes": datatypes.JSONSet("attributes").
			Set("{skills}", data.Skills),
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) AddOrUpdateModuleMetadata(userId string, module string, data interface{}) error {
	t := map[string]interface{}{
		"attributes": datatypes.JSONSet("attributes").
			Set("{"+module+"_metadata}", data),
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

func NewUserRepository(db *gorm.DB) *repositoryUser {
	return &repositoryUser{
		db: db,
	}
}
