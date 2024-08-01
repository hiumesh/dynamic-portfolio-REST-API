package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Status string

const (
	Draft    Status = "DRAFT"
	Active   Status = "ACTIVE"
	InActive Status = "IN_ACTIVE"
)

type ModelPortfolios struct {
	ID         string         `json:"id" gorm:"primary_key;"`
	Status     Status         `json:"status" gorm:"type:enum('DRAFT', 'ACTIVE', 'IN_ACTIVE')"`
	Attributes datatypes.JSON `json:"attributes"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func (model *ModelPortfolios) BeforeUpdate(db *gorm.DB) error {
	model.UpdatedAt = time.Now().Local()
	return nil
}
