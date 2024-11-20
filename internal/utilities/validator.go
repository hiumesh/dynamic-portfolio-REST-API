package utilities

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

func ValidationErrorsToJSON(err error) map[string]string {
	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = fmt.Sprintf("field %s: %s\n", err.Error(), err.Tag())
	}
	return errors
}

func YearWithinValidRangeValidator(f1 validator.FieldLevel) bool {
	params := strings.Split(f1.Param(), " ")

	if len(params) != 2 {
		return false
	}

	minOffset, err := strconv.Atoi(params[0])
	if err != nil {
		return false
	}

	maxOffset, err := strconv.Atoi(params[1])
	if err != nil {
		return false
	}

	yearStr := f1.Field().String()
	year, err := strconv.Atoi(yearStr)

	if err != nil {
		return false
	}

	currentYear := time.Now().Year()
	minYear := currentYear - minOffset
	maxYear := currentYear + maxOffset

	return year >= minYear && year <= maxYear

}

func WorkDomainsValidator(f1 validator.FieldLevel) bool {
	predefinedWorkDomains := map[string]bool{
		"Backend Developer":            true,
		"Big Data Engineer":            true,
		"Blockchain Developer":         true,
		"Brand Specialist":             true,
		"Business Analyst":             true,
		"Business Development":         true,
		"Cloud Engineer":               true,
		"Community Management":         true,
		"Content Analyst":              true,
		"Content Writer":               true,
		"Customer success":             true,
		"Data Analyst":                 true,
		"Data Engineer":                true,
		"Data Entry":                   true,
		"Data Science":                 true,
		"DevOps Engineer":              true,
		"Developer Advocate (DevRel)":  true,
		"Digital Marketing":            true,
		"Embedded Software Engineer":   true,
		"Fashion Design":               true,
		"Finance":                      true,
		"Founder's office":             true,
		"Frontend Developer":           true,
		"Fullstack Developer":          true,
		"Game Developer":               true,
		"Graphic Design":               true,
		"Growth":                       true,
		"Human Resources (HR)":         true,
		"Integration Support Engineer": true,
		"Law":                          true,
		"Legal Consultant":             true,
		"MLOps Engineer":               true,
		"Machine Learning Engineer":    true,
		"Management Consultant":        true,
		"Market/Business Research":     true,
		"Marketing":                    true,
		"Marketing/Sales":              true,
		"Mobile App Developer":         true,
		"NGO":                          true,
		"Network Engineer":             true,
		"Operations":                   true,
		"Product Management":           true,
		"Product Operations":           true,
		"Product marketing":            true,
		"Prompt Engineer":              true,
		"Public Relations (PR)":        true,
		"QA Engineer":                  true,
		"Social Media Marketing":       true,
		"Software Developer":           true,
		"Software Development Engineer in Test (SDET)": true,
		"Solution Analyst":              true,
		"Strategy":                      true,
		"Subject Matter Expert (SME)":   true,
		"Supply Chain Management (SCM)": true,
		"Teaching Assistant (TA)":       true,
		"Technical Content Engineer":    true,
		"Technical Content Writer":      true,
		"Technical Operations":          true,
		"Technical Product Manager":     true,
		"Testing Engineer":              true,
		"UI/UX Designer":                true,
		"Video/Graphics Editing":        true,
		"Volunteer":                     true,
		"Wordpress Developer":           true,
	}

	array := f1.Field()

	if array.Kind() != reflect.Slice {
		return false
	}

	for i := 0; i < array.Len(); i++ {
		element := array.Index(i).String()
		if _, exists := predefinedWorkDomains[element]; !exists {
			return false
		}
	}

	return true
}

func SocialPlatformValidator(f1 validator.FieldLevel) bool {
	predefinedSocialPlatforms := map[string]bool{
		"Hacker Rank":        true,
		"GeeksforGeeks":      true,
		"CodeChef":           true,
		"LeetCode":           true,
		"Codeforces":         true,
		"Topcoder":           true,
		"Hackerearth":        true,
		"Behance":            true,
		"Blog":               true,
		"Portfolio Website":  true,
		"Dribble":            true,
		"Other Profile Link": true,
	}

	socialPlatform := f1.Field().String()

	if _, exists := predefinedSocialPlatforms[socialPlatform]; exists {
		return true
	}

	return false
}

func CollegeDegreeValidator(f1 validator.FieldLevel) bool {
	predefinedCollegeDegrees := map[string]bool{
		"B. Voc":                      true,
		"B.A.":                        true,
		"B.Arch":                      true,
		"B.B.A.":                      true,
		"B.C.A.":                      true,
		"B.Com.":                      true,
		"B.E.":                        true,
		"B.F.Tech.":                   true,
		"B.Pharm.":                    true,
		"B.S.":                        true,
		"B.Sc.":                       true,
		"B.Tech":                      true,
		"B.Tech + M.Tech":             true,
		"Bachelor of Fine Arts (BFA)": true,
		"Diploma":                     true,
		"M. Voc":                      true,
		"M.Arch":                      true,
		"M.B.A.":                      true,
		"M.C.A.":                      true,
		"M.Com.":                      true,
		"M.Des.":                      true,
		"M.E.":                        true,
		"M.Ed.":                       true,
		"M.F.Tech":                    true,
		"M.S.":                        true,
		"M.Sc.":                       true,
		"M.Tech":                      true,
		"PG Diploma":                  true,
	}

	degree := f1.Field().String()

	if _, exists := predefinedCollegeDegrees[degree]; exists {
		return true
	}

	return false
}
