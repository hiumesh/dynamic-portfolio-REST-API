package schemas

import (
	"github.com/go-playground/validator/v10"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/utilities"
)

type SchemaSocialProfileLink struct {
	Platform string `json:"platform" validate:"required,social_platform"`
	URL      string `json:"url" validate:"required,url"`
}

type SchemaProfileBasic struct {
	FullName           string                    `json:"full_name" binding:"required" validate:"required,min=4,max=30"`
	Tagline            string                    `json:"tagline" binding:"required" validate:"omitempty,min=10,max=200"`
	About              string                    `json:"about" validate:"omitempty,min=10,max=1000"`
	ProfilePicture     string                    `json:"profile_picture" validate:"omitempty,url"`
	College            string                    `json:"college" validate:"omitempty,min=10,max=100"`
	GraduationYear     string                    `json:"graduation_year" validate:"omitempty,number,min=4,max=4,year_in_range=200 5"`
	WorkDomains        []string                  `json:"work_domains" binding:"required" validate:"required,min=1,max=5,work_domain,unique,dive,required,min=3,max=100"`
	SocialProfileLinks []SchemaSocialProfileLink `json:"social_profiles" validate:"unique=Platform,dive"`
}

func (s *SchemaProfileBasic) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("year_in_range", utilities.YearWithinValidRangeValidator)
	validate.RegisterValidation("work_domain", utilities.WorkDomainsValidator)
	validate.RegisterValidation("social_platform", utilities.SocialPlatformValidator)

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

type File struct {
	FileName string `json:"file_name" validate:"required,min=3,max=100"`
	FileSize int64  `json:"file_size" validate:"required"`
	FileType string `json:"file_type" validate:"required,min=3,max=100"`
}

type SchemaPresignedURL struct {
	Files []File `json:"files" validate:"required,min=1,max=5"`
}

func (s *SchemaPresignedURL) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SchemaReorderItem struct {
	NewIndex int16 `json:"new_index" binding:"required" validate:"required,number,min=1,max=20"`
}

func (s *SchemaReorderItem) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}

type SchemaSkills struct {
	Skills []string `json:"skills" binding:"required" validate:"required,min=1,max=70,unique,dive,required,min=3,max=50"`
}

func (s *SchemaSkills) Validate() error {
	validate := validator.New()
	return validate.Struct(s)
}
