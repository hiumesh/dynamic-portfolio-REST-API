package schemas

import (
	"github.com/go-playground/validator/v10"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
)

type SchemaUserEducation struct {
	Type          string `json:"type" binding:"required" validate:"required,oneof=SCHOOL COLLAGE"`
	InstituteName string `json:"institute_name" binding:"required" validate:"required,min=4,max=100"`
	Degree        string `json:"degree" validate:"required_if=Type COLLAGE,omitempty,min=4,max=100,collage_degree"`
	FieldOfStudy  string `json:"field_of_study" validate:"required_if=Type COLLAGE,omitempty,min=4,max=100,field_of_study"`
	Class         string `json:"class" validate:"required_if=Type SCHOOL,oneof=X XII"`
	Grade         string `json:"grade" binding:"required" validate:"required,numeric,min=1,max=10"`
	StartYear     string `json:"start_year" validate:"required_if=Type COLLAGE,omitempty,number,min=4,max=4,year_in_range=20 10"`
	EndYear       string `json:"end_year" validate:"required_if=Type COLLAGE,omitempty,number,min=4,max=4,year_in_range=20 10"`
	PassingYear   string `json:"passing_year" validate:"required_if=Type SCHOOL,omitempty,number,min=4,max=4,year_in_range=20 10"`
}

func (s *SchemaUserEducation) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("year_in_range", helpers.YearWithinValidRangeValidator)
	validate.RegisterValidation("collage_degree", helpers.CollageDegreeValidator)
	validate.RegisterValidation("field_of_study", helpers.FieldOfStudyValidator)

	return validate.Struct(s)
}

type SchemaReorderUserEducation struct {
	NewIndex int16 `json:"new_index" binding:"required" validate:"required,number,min=1,max=20"`
}

func (s *SchemaReorderUserEducation) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
