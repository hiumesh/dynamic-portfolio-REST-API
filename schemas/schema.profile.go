package schemas

import (
	"github.com/go-playground/validator/v10"
	"github.com/hiumesh/dynamic-portfolio-REST-API/helpers"
)

type SchemaSocialProfileLink struct {
	Platform string `json:"platform" validate:"required,social_platform"`
	URL      string `json:"url" validate:"required,url"`
}

type SchemaProfileBasic struct {
	FullName           string                    `json:"full_name" binding:"required" validate:"required,min=4,max=30"`
	ProfilePicture     string                    `json:"profile_picture"`
	Collage            string                    `json:"collage" validate:"min=10,max=100"`
	GraduationYear     string                    `json:"graduation_year" validate:"number,min=4,max=4,year_in_range=200 5"`
	WorkDomains        []string                  `json:"work_domains" binding:"required" validate:"required,min=1,max=5,work_domain,unique,dive,required,min=3,max=100"`
	SocialProfileLinks []SchemaSocialProfileLink `json:"social_profiles" validate:"unique=Platform,dive"`
}

func (s *SchemaProfileBasic) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("year_in_range", helpers.YearWithinValidRangeValidator)
	validate.RegisterValidation("work_domain", helpers.WorkDomainsValidator)
	validate.RegisterValidation("social_platform", helpers.SocialPlatformValidator)

	return validate.Struct(s)

	// err := validate.Struct(s)

	// structKeysToLabel := map[string]string{
	// 	"FullName":           "Name",
	// 	"ProfilePicture":     "Profile Picture",
	// 	"Collage":            "Collage",
	// 	"GraduationYear":     "Graduation Year",
	// 	"WorkDomains":        "Work Domain",
	// 	"SocialProfileLinks": "Social Profile Link",
	// }

	// errorsMap := make(map[string]string)

	// errs := err.(validator.ValidationErrors)
	// for _, fieldErr := range errs {
	// 	errorsMap[fieldErr.Field()] = err
	// }
}

type SchemaProfileEducation struct {
	Type          string `json:"type" validate:"required,oneof=SCHOOL COLLAGE"`
	InstituteName string `json:"institute_name" validate:"required,min=4,max=100"`
	Degree        string `json:"degree" validate:"required_if=Type COLLAGE,min=4,max=100,collage_degree"`
	FieldOfStudy  string `json:"field_of_study" validate:"required_if=Type COLLAGE,min=4,max=100,field_of_study"`
	Class         string `json:"class" validate:"required_if=Type SCHOOL,oneof=X XII"`
	Grade         string `json:"grade" validate:"required,numeric,min=1,max=10"`
	StartYear     string `json:"start_year" validate:"required_if=Type COLLAGE,number,min=4,max=4,year_in_range=20 10"`
	EndYear       string `json:"end_year" validate:"required_if=Type COLLAGE,number,min=4,max=4,year_in_range=20 10"`
	PassingYear   string `json:"passing_year" validate:"required_if=Type SCHOOL,number,min=4,max=4,year_in_range=20 10"`
}

func (s *SchemaProfileEducation) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("year_in_range", helpers.YearWithinValidRangeValidator)
	validate.RegisterValidation("collage_degree", helpers.CollageDegreeValidator)
	validate.RegisterValidation("field_of_study", helpers.FieldOfStudyValidator)

	return validate.Struct(s)
}
