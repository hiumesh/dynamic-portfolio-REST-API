package schemas

import "gorm.io/datatypes"

type SelectGetPortfolio struct {
	ID                string         `json:"id"`
	Status            string         `json:"status"`
	BasicDetails      datatypes.JSON `json:"basic_details"`
	Skills            datatypes.JSON `json:"skills"`
	AdditionalDetails datatypes.JSON `json:"additional_details"`
}

type SelectPortfoliosItem struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	Avatar         *string         `json:"avatar"`
	Slug           string          `json:"slug"`
	College        *string         `json:"college"`
	Skills         *datatypes.JSON `json:"skills"`
	Tagline        *string         `json:"tagline"`
	WorkDomains    *datatypes.JSON `json:"work_domains"`
	SocialProfiles *datatypes.JSON `json:"social_profiles"`
}
