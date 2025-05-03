package models

type Skill struct {
	ID    uint    `json:"id" gorm:"primaryKey"`
	Name  string  `json:"name"`
	Image *string `json:"image"`
}

func (Skill) TableName() string {
	return "skills"
}

type Skills []Skill
