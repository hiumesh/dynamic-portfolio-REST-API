package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Education struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	UserId        uuid.UUID      `json:"user_id"`
	OrderIndex    int16          `json:"order_index"`
	Type          string         `json:"type" gorm:"type:user_education_type_enum"`
	InstituteName string         `json:"institute_name"`
	Grade         string         `json:"grade"`
	Attributes    datatypes.JSON `json:"attributes"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Education) TableName() string {
	return "educations"
}

type Educations []Education
