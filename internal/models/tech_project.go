package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TechProject struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserId      uuid.UUID      `json:"user_id"`
	OrderIndex  int16          `json:"order_index"`
	Title       string         `json:"title"`
	StartDate   time.Time      `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	Description string         `json:"description"`
	SkillsUsed  pq.StringArray `json:"skills_used" gorm:"type:text"`
	Attributes  datatypes.JSON `json:"attributes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	// Attachments []Attachment   `json:"attachments" gorm:"foreignkey:ParentId"`
}

func (TechProject) TableName() string {
	return "tech_projects"
}

type TechProjects []TechProject
