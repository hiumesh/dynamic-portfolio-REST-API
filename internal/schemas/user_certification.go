package schemas

import "github.com/go-playground/validator/v10"

type SchemaUserCertification struct {
	Title           string   `json:"title" binding:"required" validate:"required,min=4,max=100"`
	Description     []string `json:"description" binding:"required" validate:"required,min=1,max=20,dive,min=5,max=500"`
	SkillsUsed      []string `json:"skills_used" binding:"required" validate:"required,min=3,max=20,dive,min=1,max=100"`
	CompletionDate  string   `json:"completion_date" binding:"required" validate:"required,datetime=2006-01-02"`
	CertificateLink string   `json:"certificate_link" validate:"omitempty,url"`
}

func (s *SchemaUserCertification) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
