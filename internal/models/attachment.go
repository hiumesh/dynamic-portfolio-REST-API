package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Attachment struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	ParentTable string         `json:"parent_table"`
	ParentId    uint           `json:"parent_id"`
	UserId      uuid.UUID      `json:"user_id"`
	FileUrl     string         `json:"file_url"`
	FileName    string         `json:"file_name"`
	FileType    string         `json:"file_type"`
	FileSize    int64          `json:"file_size"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Attachment) TableName() string {
	return "attachments"
}

type Attachments []Attachment
