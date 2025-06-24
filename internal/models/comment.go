package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Comment struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	UserId     uuid.UUID      `json:"user_id"`
	ParentId   *uint          `json:"parent_id"`
	Body       string         `json:"body"`
	Attributes datatypes.JSON `json:"attributes"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (Comment) TableName() string {
	return "comments"
}

type Comments []Comment

type CommentReaction struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserId    uuid.UUID `json:"user_id"`
	CommentId uint      `json:"comment_id"`
	Type      string    `json:"type" gorm:"type:user_reaction_type_enum"`
}

func (CommentReaction) TableName() string {
	return "comment_reactions"
}
