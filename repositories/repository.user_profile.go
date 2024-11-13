package repositories

import (
	"strings"

	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RepositoryUserProfile interface {
	Get(userId string) (*models.UserProfile, error)
	Upsert(userId string, profile *schemas.SchemaProfileBasic) error
	ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error
}

type repositoryUserProfile struct {
	db *gorm.DB
}

func (r *repositoryUserProfile) Get(userId string) (*models.UserProfile, error) {
	var profile models.UserProfile

	if err := r.db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
		return nil, &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return &profile, nil
}

func (r *repositoryUserProfile) ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error {
	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"full_name": profile.FullName,
		"slug":      strings.Join(strings.Fields(strings.ToLower(profile.FullName)), "-") + "-" + strings.Split(userId, "-")[0],
		// "attributes": gorm.Expr("JSONB_SET(attributes, '{college}', to_jsonb(?))", profile.College),
		"attributes": datatypes.JSONSet("attributes").Set("{college}", profile.College).Set("{graduation_year}", profile.GraduationYear).Set("{work_domains}", profile.WorkDomains),
		// "attributes": datatypes.JSONSet("attributes").Set("college", profile.College).Set("graduation_year", profile.GraduationYear).Set("work_domains", profile.WorkDomains),
	}).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return nil
}

func (r *repositoryUserProfile) Upsert(userId string, profile *schemas.SchemaProfileBasic) error {
	data := map[string]interface{}{
		"full_name":  profile.FullName,
		"avatar_url": profile.ProfilePicture,
		"attributes": datatypes.JSONSet("attributes").Set("college", profile.College).Set("graduation_year", profile.GraduationYear).Set("work_domains", profile.WorkDomains).Set("social_profiles", profile.SocialProfileLinks),
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(data).Error; err != nil {
		return &helpers.DatabaseError{Type: err.Error(), ErrorData: err}
	}

	return nil
}

func NewUserProfileRepository(db *gorm.DB) *repositoryUserProfile {
	return &repositoryUserProfile{
		db: db,
	}
}
