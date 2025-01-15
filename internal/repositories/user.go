package repositories

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RepositoryUser interface {
	GetProfile(userId string) (*models.UserProfile, error)
	GetPortfolio(userId string) (interface{}, error)
	UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error
	ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error
	UpsertSkills(userId string, data *schemas.SchemaSkills) error
	AddOrUpdateModuleMetadata(userId string, module string, data *interface{}) error
	UpdateStatus(userId string, status models.PortfolioStatus) error
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

type SelectGetPortfolio struct {
	ID                string         `json:"id"`
	Status            string         `json:"status"`
	BasicDetails      datatypes.JSON `json:"basic_details"`
	Skills            datatypes.JSON `json:"skills"`
	AdditionalDetails datatypes.JSON `json:"additional_details"`
}

func (r *repositoryUser) GetPortfolio(userId string) (interface{}, error) {
	var rows *sql.Rows
	var err error

	rows, err = r.db.Raw(`
		select
			user_profiles.user_id as id,
			user_profiles.portfolio_status as status,
			json_build_object(
				'email',
				user_profiles.email,
				'name',
				user_profiles.full_name,
				'avatar',
				user_profiles.avatar_url,
				'slug',
				user_profiles.slug,
				'about',
				user_profiles.attributes -> 'about',
				'tagline',
				user_profiles.attributes -> 'tagline',
				'college',
				user_profiles.attributes -> 'college',
				'graduation_year',
				user_profiles.attributes -> 'graduation_year',
				'work_domains',
				user_profiles.attributes -> 'work_domains',
				'social_profiles',
				user_profiles.attributes -> 'social_profiles'
			) as basic_details,
			user_profiles.attributes -> 'skills' as skills,
			json_build_object(
				'education_metadata',
				user_profiles.attributes -> 'education_metadata',
				'hackathon_metadata',
				user_profiles.attributes -> 'hackathon_metadata',
				'work_gallery_metadata',
				user_profiles.attributes -> 'work_gallery_metadata',
				'certification_metadata',
				user_profiles.attributes -> 'certification_metadata',
				'work_experience_metadata',
				user_profiles.attributes -> 'work_experience_metadata',
				'blog_metadata',
				user_profiles.attributes -> 'blog_metadata'
			) as additional_details
		from
			user_profiles
		where
			user_profiles.user_id = ?
			`, userId).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results SelectGetPortfolio

	for rows.Next() {
		if err := rows.Scan(&results.ID, &results.Status, &results.BasicDetails, &results.Skills, &results.AdditionalDetails); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if results.ID == "" {
		return nil, errors.New("failed to get portfolio")
	}

	return &results, nil
}

func (r *repositoryUser) UpdateStatus(userId string, status models.PortfolioStatus) error {
	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"portfolio_status": status,
	}).Error; err != nil {
		return err
	}

	return nil
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
