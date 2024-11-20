package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Status string

const (
	Draft    Status = "DRAFT"
	Active   Status = "ACTIVE"
	InActive Status = "IN_ACTIVE"
)

type UserProfile struct {
	UserId          uuid.UUID       `json:"user_id" gorm:"primaryKey"`
	Email           string          `json:"email" gorm:"uniqueIndex"`
	FullName        *string         `json:"full_name"`
	AvatarUrl       *string         `json:"avatar_url"`
	Slug            string          `json:"slug" gorm:"uniqueIndex"`
	PortfolioStatus Status          `json:"status" gorm:"type:user_profiles_portfolio_status_enum"`
	Attributes      *datatypes.JSON `json:"attributes"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
}

func (UserProfile) TableName() string {
	return "user_profiles"
}

type UserProfiles []UserProfile
