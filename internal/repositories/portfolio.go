package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type RepositoryPortfolio interface {
	GetAll(userId string, query *string, cursor int, limit int) (*[]schemas.SelectPortfoliosItem, error)
	GetPortfolio(slug string) (interface{}, error)
	GetEducations(slug string) (interface{}, error)
	GetWorkExperiences(slug string) (interface{}, error)
	GetCertifications(slug string) (interface{}, error)
	GetHackathons(slug string) (interface{}, error)
	GetTechProjects(slug string) (interface{}, error)
	GetUserPortfolio(userId string) (interface{}, error)
}

type repositoryPortfolio struct {
	db *gorm.DB
}

func (r *repositoryPortfolio) GetAll(userId *string, query *string, cursor int, limit int) (*[]schemas.SelectPortfoliosItem, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		SELECT
			user_id AS id,
			full_name AS name,
			email,
			avatar_url AS avatar,
			slug,
			attributes ->> 'college' AS college,
			attributes -> 'skills' AS skills,
			attributes ->> 'tagline' AS tagline,
			attributes -> 'work_domains' AS work_domains,
			attributes -> 'social_profiles' AS social_profiles
		FROM
			user_profiles
		WHERE
			user_profiles.portfolio_status = 'ACTIVE'
	`

	var args []interface{}
	args = append(args, limit, cursor)

	if query != nil && *query != "" {
		baseQuery += " AND fts @@ to_tsquery(?)"
		args = append([]interface{}{*query}, args...)
	}

	baseQuery += `
		ORDER BY user_profiles.updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []schemas.SelectPortfoliosItem

	for rows.Next() {
		var portfolio schemas.SelectPortfoliosItem

		err = rows.Scan(
			&portfolio.ID,
			&portfolio.Name,
			&portfolio.Email,
			&portfolio.Avatar,
			&portfolio.Slug,
			&portfolio.College,
			&portfolio.Skills,
			&portfolio.Tagline,
			&portfolio.WorkDomains,
			&portfolio.SocialProfiles,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, portfolio)
	}

	return &results, nil
}

func (r *repositoryPortfolio) GetUserPortfolio(userId string) (interface{}, error) {
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

	var results schemas.SelectGetPortfolio

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

func (r *repositoryPortfolio) GetPortfolio(slug string) (interface{}, error) {
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
			user_profiles.slug = ?
			`, slug).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results schemas.SelectGetPortfolio

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

func (r *repositoryPortfolio) GetEducations(slug string) (interface{}, error) {
	var userEducations models.Educations

	if err := r.db.Joins("inner join user_profiles on user_profiles.user_id = educations.user_id").Where("user_profiles.slug = ?", slug).Order("order_index desc").Find(&userEducations).Error; err != nil {
		return nil, err
	}

	return &userEducations, nil
}

func (r *repositoryPortfolio) GetHackathons(slug string) (*models.Hackathons, error) {
	var userHackathons models.Hackathons

	if err := r.db.Joins("inner join user_profiles on user_profiles.user_id = hackathons.user_id").Where("user_profiles.slug = ?", slug).Order("order_index desc").Find(&userHackathons).Error; err != nil {
		return nil, err
	}

	return &userHackathons, nil
}

func (r *repositoryPortfolio) GetWorkExperiences(slug string) (*models.WorkExperiences, error) {
	var userExperiences models.WorkExperiences

	if err := r.db.Joins("inner join user_profiles on user_profiles.user_id = work_experiences.user_id").Where("user_profiles.slug = ?", slug).Order("order_index desc").Find(&userExperiences).Error; err != nil {
		return nil, err
	}

	return &userExperiences, nil
}

func (r *repositoryPortfolio) GetCertifications(slug string) (*models.Certifications, error) {
	var userCertifications models.Certifications

	if err := r.db.Joins("inner join user_profiles on user_profiles.user_id = certifications.user_id").Where("user_profiles.slug = ?", slug).Order("order_index desc").Find(&userCertifications).Error; err != nil {
		return nil, err
	}

	return &userCertifications, nil
}

func (r *repositoryPortfolio) GetTechProjects(slug string) (*[]schemas.SelectUserTechProject, error) {
	var rows *sql.Rows
	var err error

	// Execute raw query
	rows, err = r.db.Raw(`
		select
			tech_projects.id,
			tech_projects.order_index,
			tech_projects.title,
			tech_projects.description,
			tech_projects.tech_used,
			tech_projects.attributes -> 'links' as links,
			tech_projects.created_at,
			tech_projects.updated_at,
			coalesce(json_agg(
				json_build_object(
					'id', attachments.id,
					'file_url', attachments.file_url,
					'file_name', attachments.file_name,
					'file_type', attachments.file_type,
					'file_size', attachments.file_size
				)
			) filter (where attachments.id is not null), '[]') as attachments
		from
			tech_projects
			inner join user_profiles on user_profiles.user_id = tech_projects.user_id
			left join attachments on attachments.parent_id = tech_projects.id
			and attachments.parent_table = 'tech_projects'
		where
			user_profiles.slug = ?
		group by
			tech_projects.id
		order by
			order_index desc
		limit 6
	`, slug).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate and manually map rows
	var userTechProjects []schemas.SelectUserTechProject
	for rows.Next() {
		var project schemas.SelectUserTechProject
		var attachmentsJSON []byte

		// Scan into fields and JSON data
		err = rows.Scan(
			&project.ID,
			&project.OrderIndex,
			&project.Title,
			&project.Description,
			&project.TechUsed,
			&project.Links,
			&project.CreatedAt,
			&project.UpdatedAt,
			&attachmentsJSON,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal JSON into Attachments slice
		err = json.Unmarshal(attachmentsJSON, &project.Attachments)
		if err != nil {
			return nil, err
		}

		userTechProjects = append(userTechProjects, project)
	}

	return &userTechProjects, nil
}

func NewPortfolioRepository(db *gorm.DB) *repositoryPortfolio {
	return &repositoryPortfolio{
		db: db,
	}
}
