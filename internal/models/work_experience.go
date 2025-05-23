package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type WorkExperience struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserId          uuid.UUID      `json:"user_id"`
	OrderIndex      int16          `json:"order_index"`
	CompanyName     string         `json:"company_name"`
	CompanyUrl      string         `json:"company_url"`
	JobType         string         `json:"job_type" gorm:"type:work_experiences_job_type_enum"`
	JobTitle        string         `json:"job_title"`
	Location        string         `json:"location"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         *time.Time     `json:"end_date"`
	Description     string         `json:"description"`
	SkillsUsed      pq.StringArray `json:"skills_used" gorm:"type:text"`
	CertificateLink *string        `json:"certificate_link"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (WorkExperience) TableName() string {
	return "work_experiences"
}

type WorkExperiences []WorkExperience
