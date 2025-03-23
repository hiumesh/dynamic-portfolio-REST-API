package schemas

import "github.com/go-playground/validator/v10"

type SchemaTechProjectLink struct {
	Platform string `json:"platform" validate:"required,oneof=Github Website Social"`
	Label    string `json:"label" validate:"required,min=2,max=100"`
	URL      string `json:"url" validate:"required,url"`
}

type SchemaTechProject struct {
	Title       string                  `json:"title" binding:"required" validate:"required,min=4,max=100"`
	Description string                  `json:"description" binding:"required" validate:"required,min=10,max=2000"`
	StartDate   string                  `json:"start_date" binding:"required" validate:"required,datetime=2006-01-02"`
	EndDate     string                  `json:"end_date" validate:"omitempty,datetime=2006-01-02"`
	SkillsUsed  []string                `json:"skills_used" binding:"required" validate:"required,min=3,max=20,dive,min=1,max=100"`
	Links       []SchemaTechProjectLink `json:"links" validate:"dive"`
	Attachments Attachments             `json:"attachments" validate:"required,min=1,max=20,dive"`
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
