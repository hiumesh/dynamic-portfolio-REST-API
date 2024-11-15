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
	ProfilePicture     string                    `json:"profile_picture" validate:"omitempty,url"`
	College            string                    `json:"college" validate:"omitempty,min=10,max=100"`
	GraduationYear     string                    `json:"graduation_year" validate:"omitempty,number,min=4,max=4,year_in_range=200 5"`
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
