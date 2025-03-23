package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Certification struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserId          uuid.UUID      `json:"user_id"`
	OrderIndex      int16          `json:"order_index"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	CompletionDate  time.Time      `json:"completion_date"`
	CertificateLink *string        `json:"certificate_link"`
	SkillsUsed      pq.StringArray `json:"skills_used" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Certification) TableName() string {
	return "certifications"
}

type Certifications []Certification
