package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Blog struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	CoverImage  *string        `json:"cover_image"`
	UserId      uuid.UUID      `json:"user_id"`
	Title       string         `json:"title"`
	Body        *string        `json:"body"`
	Slug        string         `json:"slug"`
	Attributes  datatypes.JSON `json:"attributes"`
	PublishedAt *time.Time     `json:"published_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Tags        []Tag          `json:"tags" gorm:"many2many:blog_tags"`
}

func (Blog) TableName() string {
	return "blogs"
}

type Blogs []Blog

type BlogTag struct {
	BlogId uint `json:"blog_id"`
	TagId  uint `json:"tag_id"`
}

func (BlogTag) TableName() string {
	return "blog_tags"
}

type BlogTags []BlogTag

type BlogComment struct {
	BlogId    uint `json:"blog_id"`
	CommentId uint `json:"comment_id"`
}

func (BlogComment) TableName() string {
	return "blog_comments"
}

type BlogComments []BlogComment

type BlogReaction struct {
	ID     uint      `json:"id" gorm:"primaryKey"`
	BlogId uint      `json:"blog_id"`
	UserId uuid.UUID `json:"user_id"`
	Type   string    `json:"type" gorm:"type:user_reaction_type_enum"`
}

func (BlogReaction) TableName() string {
	return "blog_reactions"
}

type BlogReactions []BlogReaction
