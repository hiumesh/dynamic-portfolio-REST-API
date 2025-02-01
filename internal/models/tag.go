package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Tag struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserId      *uuid.UUID     `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Attributes  datatypes.JSON `json:"attributes"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Tag) TableName() string {
	return "tags"
}

type Tags []Tag
