package seeds

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

type Tag struct {
	ID        int64      `gorm:"primaryKey;autoIncrement"`
	Name      string     `json:"name" fake:"{word}"`
	CreatedAt time.Time  `json:"created_at" fake:"{date}"`
	UpdatedAt time.Time  `json:"updated_at" fake:"{date}"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (Tag) TableName() string {
	return "tags"
}

func NewTag() (*Tag, error) {
	tag := Tag{}
	if err := gofakeit.Struct(&tag); err != nil {
		return nil, err
	}
	tag.ID = 0
	return &tag, nil
}

func SeedTags(db *gorm.DB, count int) (*[]Tag, error) {
	if count == 0 {
		count = 10
	}
	uniqueTags := make(map[string]bool)
	tags := make([]Tag, count)
	for len(uniqueTags) < count {
		tag, err := NewTag()
		if err != nil {
			return nil, err
		}
		if _, ok := uniqueTags[tag.Name]; !ok {
			tags[len(uniqueTags)] = *tag
			uniqueTags[tag.Name] = true
		}
	}
	if err := db.Create(&tags).Error; err != nil {
		return nil, err
	}

	return &tags, nil
}
