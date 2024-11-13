package schemas

import (
	"github.com/go-playground/validator/v10"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
)

type SchemaUserEducation struct {
	Type          string `json:"type" binding:"required" validate:"required,oneof=SCHOOL COLLEGE"`
	InstituteName string `json:"institute_name" binding:"required" validate:"required,min=4,max=100"`
	Degree        string `json:"degree" validate:"required_if=Type COLLEGE,omitempty,min=4,max=100,college_degree"`
	FieldOfStudy  string `json:"field_of_study" validate:"required_if=Type COLLEGE,omitempty,min=4,max=100"`
	Class         string `json:"class" validate:"required_if=Type SCHOOL,omitempty,oneof=X XII"`
	Grade         string `json:"grade" binding:"required" validate:"required,numeric,min=1,max=10"`
	StartYear     string `json:"start_year" validate:"required_if=Type COLLEGE,omitempty,number,min=4,max=4,year_in_range=20 10"`
	EndYear       string `json:"end_year" validate:"required_if=Type COLLEGE,omitempty,number,min=4,max=4,year_in_range=20 10"`
	PassingYear   string `json:"passing_year" validate:"required_if=Type SCHOOL,omitempty,number,min=4,max=4,year_in_range=20 10"`
}

func (s *SchemaUserEducation) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("year_in_range", helpers.YearWithinValidRangeValidator)
	validate.RegisterValidation("college_degree", helpers.CollegeDegreeValidator)

	return validate.Struct(s)
}

type SchemaReorderUserEducation struct {
	NewIndex int16 `json:"new_index" binding:"required" validate:"required,number,min=1,max=20"`
}

func (s *SchemaReorderUserEducation) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
