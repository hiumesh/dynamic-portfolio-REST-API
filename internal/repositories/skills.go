package repositories

import (
	"database/sql"
	"strings"

	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"gorm.io/gorm"
)

type RepositorySkill interface {
	GetAll(userId *string, query *string, cursor int, limit int) (*models.Skills, error)
	GetUserSkills(userId string) (*models.Skills, error)
}

type repositorySkill struct {
	db *gorm.DB
}

func (r *repositorySkill) GetAll(query *string, cursor int, limit int) (*models.Skills, error) {
	var skills models.Skills
	if query != nil && *query != "" {
		if err := r.db.Where("LOWER(name) ILIKE ?", "%"+strings.ToLower(*query)+"%").Offset(cursor).Limit(limit).Find(&skills).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Offset(cursor).Limit(limit).Find(&skills).Error; err != nil {
			return nil, err
		}
	}

	return &skills, nil
}

func (r *repositorySkill) GetUserSkills(userId string) (*models.Skills, error) {
	var rows *sql.Rows
	var err error

	rows, err = r.db.Raw(`
		with
			user_skills as (
				select
					jsonb_array_elements_text(attributes -> 'skills') as skill_name
				from
					user_profiles
				where
					user_id = ?
			)
		select
			us.skill_name as name,
			skills.image
		from
			skills
		right	join user_skills us on skills.name = us.skill_name;
	`, userId).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var skills models.Skills

	for rows.Next() {
		var skill models.Skill
		err = rows.Scan(&skill.Name, &skill.Image)
		if err != nil {
			return nil, err
		}

		skills = append(skills, skill)
	}

	return &skills, nil

}

func NewRepositorySkill(db *gorm.DB) *repositorySkill {
	return &repositorySkill{db}
}
