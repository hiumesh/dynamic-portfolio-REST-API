package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryUserTechProject interface {
	GetAll(userId *string, query *string, cursor int, limit int) (*[]schemas.SelectUserTechProject, error)
	GetUserTechProjects(userId string, query *string, cursor int, limit int) (*[]schemas.SelectUserTechProject, error)
	Get(userId string, id string) (any, error)
	Create(userId string, data *schemas.SchemaTechProject) (*models.TechProject, error)
	Update(userId string, id string, data *schemas.SchemaTechProject) (*models.TechProject, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
}

type repositoryUserTechProject struct {
	db *gorm.DB
}

func (r *repositoryUserTechProject) GetAll(userId *string, query *string, cursor int, limit int) (*[]schemas.SelectUserTechProject, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		select
			tech_projects.id,
			tech_projects.title,
			tech_projects.description,
			tech_projects.tech_used,
			tech_projects.attributes -> 'links' as links,
			user_profiles.user_id as publisher_id,
			user_profiles.avatar_url as publisher_avatar,
			user_profiles.full_name as publisher_name,
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
	`

	var args []interface{}
	args = append(args, limit, cursor)

	if query != nil && *query != "" {
		baseQuery += " where tech_projects.fts @@ to_tsquery(?)"
		args = append([]interface{}{*query}, args...)
	}

	baseQuery += `
		group by
			tech_projects.id,
			user_profiles.user_id
		order by
			tech_projects.created_at desc
		LIMIT ? OFFSET ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userTechProjects []schemas.SelectUserTechProject
	for rows.Next() {
		var project schemas.SelectUserTechProject
		var attachmentsJSON []byte

		// Scan into fields and JSON data
		err = rows.Scan(
			&project.ID,
			&project.Title,
			&project.Description,
			&project.TechUsed,
			&project.Links,
			&project.PublisherId,
			&project.PublisherAvatar,
			&project.PublisherName,
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

func (r *repositoryUserTechProject) GetUserTechProjects(userId string, query *string, cursor int, limit int) (*[]schemas.SelectUserTechProject, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
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
			left join attachments on attachments.parent_id = tech_projects.id
			and attachments.parent_table = 'tech_projects'
		where
			tech_projects.user_id = ?
	`

	var args []interface{}
	args = append(args, limit, cursor)

	if query != nil && *query != "" {
		baseQuery += " AND tech_projects.fts @@ to_tsquery(?)"
		args = append([]interface{}{*query}, args...)
	}

	args = append([]interface{}{userId}, args...)

	baseQuery += `
		group by
			tech_projects.id
		order by
			tech_projects.order_index desc
		LIMIT ? OFFSET ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *repositoryUserTechProject) Get(userId string, id string) (any, error) {
	var rows *sql.Rows
	var err error

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
			left join attachments on attachments.parent_id = tech_projects.id
			and attachments.parent_table = 'tech_projects'
		where
			tech_projects.user_id = ?
			and tech_projects.id = ?
		group by
			tech_projects.id
	`, userId, id).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var project schemas.SelectUserTechProject
	for rows.Next() {
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
	}

	if project.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &project, nil
}

func (r *repositoryUserTechProject) Create(userId string, data *schemas.SchemaTechProject) (*models.TechProject, error) {
	type MaxIndexResult struct {
		MaxIndex int16
	}
	maxIndexResult := MaxIndexResult{MaxIndex: 0}
	if err := r.db.Model(&models.TechProject{}).Select("max(order_index) as max_index").Where("user_id = ?", userId).Group("user_id").Take(&maxIndexResult).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	techProject := models.TechProject{
		UserId:      userUUID,
		OrderIndex:  maxIndexResult.MaxIndex + 1,
		Title:       data.Title,
		Description: data.Description,
		TechUsed:    data.TechUsed,
	}

	attributes := map[string]interface{}{}

	if data.Links != nil && len(data.Links) > 0 {
		attributes["links"] = data.Links
	}

	attributesJson, err := json.Marshal(attributes)
	if err != nil {
		return nil, errors.New("failed to marshal json")
	}

	techProject.Attributes = attributesJson

	if err = r.db.Create(&techProject).Error; err != nil {
		return nil, err
	}

	return &techProject, nil

}

func (r *repositoryUserTechProject) Update(userId string, id string, data *schemas.SchemaTechProject) (*models.TechProject, error) {

	techProject := map[string]interface{}{
		"title":       data.Title,
		"description": data.Description,
		"skills_used": pq.Array(data.TechUsed),
		"attributes": datatypes.JSONSet("attributes").
			Set("{links}", data.Links),
	}

	var updatedRows models.TechProjects

	if err := r.db.Model(&updatedRows).Clauses(clause.Returning{}).Where("id = ? and user_id = ?", id, userId).Updates(techProject).Error; err != nil {
		return nil, err
	}

	if len(updatedRows) == 0 {
		return nil, errors.New("record not found")
	}

	return &updatedRows[0], nil
}

func (r *repositoryUserTechProject) Delete(userId string, id string) error {
	var techProject models.TechProject

	if err := r.db.Where("user_id = ?", userId).Where("id = ?", id).First(&techProject).Error; err != nil {
		return err
	}

	if err := r.db.Model(&models.TechProject{}).Where("user_id = ? and order_index > ?", userId, techProject.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
		return err
	}

	if err := r.db.Unscoped().Delete(&techProject).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryUserTechProject) Reorder(userId string, id string, newIndex int) error {
	var techProject models.TechProject

	if err := r.db.First(&techProject, id).Error; err != nil {
		return err
	}

	if techProject.OrderIndex == int16(newIndex) {
		return nil
	}

	type CountResult struct {
		Count int16
	}
	countResult := CountResult{}
	if err := r.db.Model(&models.TechProject{}).Select("count(*) as count").Where("user_id = ?", userId).Group("user_id").Take(&countResult).Error; err != nil {
		return err
	}

	if countResult.Count < int16(newIndex) {
		return errors.New("invalid index for reordering")
	}

	if techProject.OrderIndex < int16(newIndex) {
		if err := r.db.Model(&models.TechProject{}).Where("user_id = ? and order_index > ? and order_index <= ?", userId, techProject.OrderIndex, newIndex).UpdateColumn("order_index", gorm.Expr("order_index - ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.TechProject{}).Where("id = ?", techProject.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	if techProject.OrderIndex > int16(newIndex) {
		if err := r.db.Model(&models.TechProject{}).Where("user_id = ? and order_index >= ? and order_index < ?", userId, newIndex, techProject.OrderIndex).UpdateColumn("order_index", gorm.Expr("order_index + ?", 1)).Error; err != nil {
			return err
		}
		if err := r.db.Model(&models.TechProject{}).Where("id = ?", techProject.ID).UpdateColumn("order_index", newIndex).Error; err != nil {
			return err
		}
	}

	return nil
}

func NewUserTechProjectRepository(db *gorm.DB) *repositoryUserTechProject {
	return &repositoryUserTechProject{
		db: db,
	}
}
