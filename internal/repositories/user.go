package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RepositoryUser interface {
	GetProfile(userId string) (*models.UserProfile, error)
	GetProfileBySlug(slug string) (*models.UserProfile, error)
	UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error
	ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error
	UpsertSkills(userId string, data *schemas.SchemaSkills) error
	UpsertResume(userId string, url *string) error
	AddOrUpdateModuleMetadata(userId string, module string, data *any) error
	GetModuleMetadata(userId string, module string) (any, error)
	UpdateStatus(userId string, status models.PortfolioStatus) error
	UpdateProfileAttachment(userId string, module string, url *string) error
	GetFollowers(userId string, cursor int, limit int) (*[]schemas.SelectFollowers, error)
	GetFollowing(userId string, cursor int, limit int) (*[]schemas.SelectFollowing, error)
	FollowUser(userId uuid.UUID, followingUserId uuid.UUID) error
	UnfollowUser(userId uuid.UUID, followingUserId uuid.UUID) error
	FollowStatus(userId uuid.UUID, followingUserId uuid.UUID) (*models.UserFollow, error)
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

func (r *repositoryUser) GetProfileBySlug(slug string) (*models.UserProfile, error) {
	var profile models.UserProfile

	if err := r.db.Where("slug = ?", slug).First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
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
	data := map[string]any{
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
	t := map[string]any{
		"attributes": datatypes.JSONSet("attributes").
			Set("{skills}", data.Skills),
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) UpsertResume(userId string, url *string) error {
	t := map[string]any{
		"attributes": datatypes.JSONSet("attributes").
			Set("{resume}", url),
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) UpdateProfileAttachment(userId string, module string, url *string) error {
	t := map[string]any{
		"attributes": datatypes.JSONSet("attributes").
			Set("{"+module+"}", url),
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) GetModuleMetadata(userId string, module string) (interface{}, error) {
	var rows *sql.Rows
	var err error

	metadataField := "'" + module + "_metadata" + "'"

	query := `
    SELECT COALESCE(attributes->%s, '{}') AS metadata 
    FROM user_profiles 
    WHERE user_id = ?;
	`

	rows, err = r.db.Raw(fmt.Sprintf(query, metadataField), userId).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var data datatypes.JSON

	for rows.Next() {
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if data == nil {
		return nil, errors.New("failed to get metadata")
	}

	return data, nil
}

func (r *repositoryUser) AddOrUpdateModuleMetadata(userId string, module string, data interface{}) error {
	t := map[string]any{
		"attributes": datatypes.JSONSet("attributes").
			Set("{"+module+"_metadata}", data),
	}

	if err := r.db.Model(&models.UserProfile{}).Where("user_id = ?", userId).Updates(t).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) GetFollowers(userId string, cursor int, limit int) (*[]schemas.SelectFollowers, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT
			user_id AS id,
			full_name AS name,
			email,
			avatar_url AS avatar,
			slug,
			attributes ->> 'tagline' AS tagline
		FROM
			user_profiles
			inner join user_follows on user_profiles.user_id = user_follows.follower_id
		WHERE
			user_follows.following_id = ?
	`

	var args []any
	args = append(args, userId, limit, cursor)

	baseQuery += `
		ORDER BY user_profiles.updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []schemas.SelectFollowers

	for rows.Next() {
		var portfolio schemas.SelectFollowers

		err = rows.Scan(
			&portfolio.ID,
			&portfolio.Name,
			&portfolio.Email,
			&portfolio.Avatar,
			&portfolio.Slug,
			&portfolio.Tagline,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, portfolio)
	}

	return &results, nil
}

func (r *repositoryUser) GetFollowing(userId string, cursor int, limit int) (*[]schemas.SelectFollowing, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT
			user_id AS id,
			full_name AS name,
			email,
			avatar_url AS avatar,
			slug,
			attributes ->> 'tagline' AS tagline
		FROM
			user_profiles
			inner join user_follows on user_profiles.user_id = user_follows.following_id
		WHERE
			user_follows.follower_id = ?
	`

	var args []any
	args = append(args, userId, limit, cursor)

	baseQuery += `
		ORDER BY user_profiles.updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []schemas.SelectFollowing

	for rows.Next() {
		var portfolio schemas.SelectFollowing

		err = rows.Scan(
			&portfolio.ID,
			&portfolio.Name,
			&portfolio.Email,
			&portfolio.Avatar,
			&portfolio.Slug,
			&portfolio.Tagline,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, portfolio)
	}

	return &results, nil
}

func (r *repositoryUser) FollowUser(userId uuid.UUID, followingUserId uuid.UUID) error {
	userFollower := models.UserFollow{
		FollowerId:  userId,
		FollowingId: followingUserId,
	}

	if err := r.db.Create(&userFollower).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUser) UnfollowUser(userId uuid.UUID, followingUserId uuid.UUID) error {
	return r.db.Where("follower_id = ? and following_id = ?", userId, followingUserId).Delete(&models.UserFollow{}).Error
}

func (r *repositoryUser) FollowStatus(userId uuid.UUID, followingUserId uuid.UUID) (*models.UserFollow, error) {
	var follow models.UserFollow

	if err := r.db.Where("follower_id = ? and following_id = ?", userId, followingUserId).First(&follow).Error; err != nil {
		return nil, err
	}

	return &follow, nil
}

func NewUserRepository(db *gorm.DB) *repositoryUser {
	return &repositoryUser{
		db: db,
	}
}
