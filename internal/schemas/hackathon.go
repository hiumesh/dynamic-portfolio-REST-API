package schemas

import "github.com/go-playground/validator/v10"

type SchemaHackathonLink struct {
	Platform string `json:"platform" validate:"required,oneof=Github Website Social"`
	Label    string `json:"label" validate:"required,min=2,max=100"`
	URL      string `json:"url" validate:"required,url"`
}

type SchemaHackathon struct {
	Avatar          string                `json:"avatar" validate:"omitempty,url"`
	Title           string                `json:"title" binding:"required" validate:"required,min=3,max=150"`
	Location        string                `json:"location" binding:"required" validate:"required,min=5,max=200"`
	StartDate       string                `json:"start_date" binding:"required" validate:"required,datetime=2006-01-02"`
	EndDate         string                `json:"end_date" validate:"required,datetime=2006-01-02"`
	Description     string                `json:"description" binding:"required" validate:"required,min=10,max=2000"`
	CertificateLink string                `json:"certificate_link" validate:"omitempty,url"`
	Links           []SchemaHackathonLink `json:"links" validate:"dive"`
}

func (s *SchemaHackathon) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SchemaHackathonMetadata struct {
	Heading     string `json:"heading" binding:"required" validate:"required,min=3,max=100"`
	Description string `json:"description" binding:"required" validate:"required,min=3,max=1000"`
}

func (s *SchemaHackathonMetadata) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
