package schemas

import "github.com/go-playground/validator/v10"

type SchemaUserExperience struct {
	CompanyName     string   `json:"company_name" binding:"required" validate:"required,min=3,max=150"`
	CompanyUrl      string   `json:"company_url" validate:"omitempty,url"`
	JobType         string   `json:"job_type" binding:"required" validate:"required,oneof= PART_TIME SEMI_FULL_TIME FULL_TIME"`
	JobTitle        string   `json:"job_title" binding:"required" validate:"required,min=5,max=200"`
	Location        string   `json:"location" binding:"required" validate:"required,min=5,max=200"`
	StartDate       string   `json:"start_date" binding:"required" validate:"required,datetime=DateOnly"`
	EndDate         string   `json:"end_date" validate:"omitempty,datetime=DateOnly"`
	Description     []string `json:"description" binding:"required" validate:"required,min=3,max=20,dive,min=10,max=500"`
	SkillsUsed      []string `json:"skills_used" binding:"required" validate:"required,min=3,max=20,dive,min=10,max=500"`
	CertificateLink string   `json:"certificate_link" validate:"omitempty,url"`
}

func (s *SchemaUserExperience) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
