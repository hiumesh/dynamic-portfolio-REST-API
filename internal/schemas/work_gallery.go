package schemas

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type SchemaTechProjectLink struct {
	Platform string `json:"platform" validate:"required,oneof=Demo Website SourceCode"`
	Label    string `json:"label" validate:"required,min=2,max=100"`
	URL      string `json:"url" validate:"required,url"`
}

type SchemaTechProject struct {
	Title       string                  `json:"title" binding:"required" validate:"required,min=4,max=100"`
	Description string                  `json:"description" binding:"required" validate:"required,min=10,max=2000"`
	TechUsed    []string                `json:"tech_used" binding:"required" validate:"required,min=3,max=20,dive,min=1,max=100"`
	Links       []SchemaTechProjectLink `json:"links" validate:"dive"`
	Attachments Attachments             `json:"attachments" validate:"omitempty,max=20,dive"`
}

func (s SchemaTechProject) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SchemaTechProjectMetadata struct {
	Heading     string `json:"heading" binding:"required" validate:"required,min=3,max=100"`
	Description string `json:"description" binding:"required" validate:"required,min=3,max=1000"`
}

func (s *SchemaTechProjectMetadata) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SelectUserTechProject struct {
	ID              uint               `json:"id"`
	OrderIndex      *int16             `json:"order_index"`
	Title           string             `json:"title"`
	Description     string             `json:"description"`
	TechUsed        pq.StringArray     `json:"tech_used" gorm:"type:text"`
	Links           *datatypes.JSON    `json:"links"`
	PublisherId     *string            `json:"publisher_id"`
	PublisherAvatar *string            `json:"publisher_avatar"`
	PublisherName   *string            `json:"publisher_name"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
	Attachments     []SelectAttachment `json:"attachments"`
}
