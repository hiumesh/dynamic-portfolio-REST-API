package schemas

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type SchemaBlog struct {
	Title       string      `json:"title" binding:"required" validate:"required,min=4,max=100"`
	Body        string      `json:"body" validate:"omitempty,min=10,max=2000"`
	CoverImage  string      `json:"cover_image" validate:"omitempty,url"`
	Tags        []string    `json:"tags" validate:"omitempty,min=0,max=10,dive,min=1,max=100"`
	Attachments Attachments `json:"attachments" validate:"omitempty,max=20,dive"`
}

func (s SchemaBlog) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SchemaBlogMetadata struct {
	Heading     string `json:"heading" binding:"required" validate:"required,min=3,max=100"`
	Description string `json:"description" binding:"required" validate:"required,min=3,max=1000"`
}

func (s *SchemaBlogMetadata) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SelectBlog struct {
	ID                uint            `json:"id"`
	CoverImage        *string         `json:"cover_image"`
	Title             string          `json:"title"`
	Body              *string         `json:"body"`
	Slug              string          `json:"slug"`
	PublisherId       *string         `json:"publisher_id"`
	PublisherAvatar   *string         `json:"publisher_avatar"`
	PublisherName     *string         `json:"publisher_name"`
	PublishedAt       *time.Time      `json:"published_at"`
	CommentsCount     *int            `json:"comments_count"`
	ReactionsMetadata *datatypes.JSON `json:"reactions_metadata"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	Tags              *pq.StringArray `json:"tags" gorm:"type:text"`
}
