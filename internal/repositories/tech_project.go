package repositories

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryUserTechProject interface {
	GetAll(userId string) (*models.TechProjects, error)
	Create(userId string, data *schemas.SchemaTechProject) (*models.TechProject, error)
	Update(userId string, id string, data *schemas.SchemaTechProject) (*models.TechProject, error)
	Reorder(userId string, id string, newIndex int) error
	Delete(userId string, id string) error
}

type repositoryUserTechProject struct {
	db *gorm.DB
}

type SelectAttachment struct {
	ID       uint   `json:"id"`
	FileUrl  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}

type SelectUserTechProject struct {
	ID          uint               `json:"id"`
	OrderIndex  int16              `json:"order_index"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	StartDate   time.Time          `json:"start_date"`
	EndDate     *time.Time         `json:"end_date"`
	SkillsUsed  pq.StringArray     `json:"skills_used" gorm:"type:text"`
	Attributes  datatypes.JSON     `json:"attributes"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	Attachments []SelectAttachment `json:"attachments"`
}

func (r *repositoryUserTechProject) GetAll(userId string) (*[]SelectUserTechProject, error) {
	var rows *sql.Rows
	var err error

	// Execute raw query
	rows, err = r.db.Raw(`
		select
			tech_projects.id,
			tech_projects.order_index,
			tech_projects.title,
			tech_projects.description,
			tech_projects.start_date,
			tech_projects.end_date,
			tech_projects.skills_used,
			tech_projects.attributes,
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
		group by
			tech_projects.id
		order by
			order_index desc
	`, userId).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate and manually map rows
	var userTechProjects []SelectUserTechProject
	for rows.Next() {
		var project SelectUserTechProject
		var attachmentsJSON []byte

		// Scan into fields and JSON data
		err = rows.Scan(
			&project.ID,
			&project.OrderIndex,
			&project.Title,
			&project.Description,
			&project.StartDate,
			&project.EndDate,
			&project.SkillsUsed,
			&project.Attributes,
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

	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return nil, errors.New("failed to parse start date")
	}

	techProject := models.TechProject{
		UserId:      userUUID,
		OrderIndex:  maxIndexResult.MaxIndex + 1,
		Title:       data.Title,
		StartDate:   startDate,
		Description: data.Description,
		SkillsUsed:  data.SkillsUsed,
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

	if data.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", data.EndDate)

		if err != nil {
			return nil, errors.New("failed to parse end date")
		}
		techProject.EndDate = &endDate
	}

	if err = r.db.Create(&techProject).Error; err != nil {
		return nil, err
	}

	return &techProject, nil

}

func (r *repositoryUserTechProject) Update(userId string, id string, data *schemas.SchemaTechProject) (*models.TechProject, error) {
	startDate, err := time.Parse("2006-01-02", data.StartDate)
	if err != nil {
		return nil, errors.New("failed to parse start date")
	}

	techProject := map[string]interface{}{
		"title":       data.Title,
		"start_date":  startDate,
		"description": data.Description,
		"skills_used": pq.Array(data.SkillsUsed),
		"attributes": datatypes.JSONSet("attributes").
			Set("{links}", data.Links),
	}

	if data.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", data.EndDate)

		if err != nil {
			return nil, errors.New("failed to parse end date")
		}
		techProject["end_date"] = endDate
	} else {
		techProject["end_date"] = nil
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
