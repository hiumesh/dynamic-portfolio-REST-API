package schemas

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type SchemaComment struct {
	Body string `json:"body" binding:"required" validate:"required,min=4,max=1000"`
}

func (s SchemaComment) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SchemaCreateComment struct {
	Module string `json:"module" binding:"required" validate:"required"`
	Slug   string `json:"slug" binding:"required" validate:"required"`
	Body   string `json:"body" binding:"required" validate:"required,min=4,max=1000"`
}

func (s SchemaCreateComment) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SchemaCommentReply struct {
	Module string `json:"module" binding:"required" validate:"required"`
	Body   string `json:"body" binding:"required" validate:"required,min=4,max=1000"`
}

func (s SchemaCommentReply) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SelectComment struct {
	ID           uint            `json:"id"`
	Body         string          `json:"body"`
	AuthorId     uuid.UUID       `json:"author_id"`
	AuthorName   *string         `json:"author_name"`
	AuthorAvatar *string         `json:"author_avatar"`
	ParentId     *uint           `json:"parent_id"`
	Attributes   datatypes.JSON  `json:"attributes"`
	Reactions    *pq.StringArray `json:"reactions"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}
