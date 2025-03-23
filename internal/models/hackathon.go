package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Hackathon struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	UserId          uuid.UUID      `json:"user_id"`
	OrderIndex      int16          `json:"order_index"`
	Avatar          *string        `json:"avatar"`
	Title           string         `json:"title"`
	Location        string         `json:"location"`
	StartDate       time.Time      `json:"start_date"`
	EndDate         time.Time      `json:"end_date"`
	Description     string         `json:"description"`
	CertificateLink *string        `json:"certificate_link"`
	Attributes      datatypes.JSON `json:"attributes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Hackathon) TableName() string {
	return "hackathons"
}

type Hackathons []Hackathon
