package seeds

import (
	"encoding/json"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var CollageDegrees = []string{
	"B. Voc",
	"B.A.",
	"B.Arch",
	"B.B.A.",
	"B.C.A.",
	"B.Com.",
	"B.E.",
	"B.F.Tech.",
	"B.Pharm.",
	"B.S.",
	"B.Sc.",
	"B.Tech",
	"B.Tech + M.Tech",
	"Bachelor of Fine Arts (BFA)",
	"Diploma",
	"M. Voc",
	"M.Arch",
	"M.B.A.",
	"M.C.A.",
	"M.Com.",
	"M.Des.",
	"M.E.",
	"M.Ed.",
	"M.F.Tech",
	"M.S.",
	"M.Sc.",
	"M.Tech",
	"PG Diploma",
}

type Education struct {
	UserId        string         `json:"user_id"`
	OrderIndex    int16          `json:"order_index"`
	Type          string         `json:"type" fake:"{randomstring:[SCHOOL,COLLEGE]}"`
	InstituteName string         `json:"institute_name" fake:"{school}"`
	Grade         string         `json:"grade"`
	Attributes    datatypes.JSON `json:"attributes"`
	CreatedAt     time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt     time.Time      `json:"updated_at" fake:"{date}"`
	// DeletedAt     *time.Time     `json:"deleted_at"`
}

func (Education) TableName() string {
	return "user_educations"
}

func NewEducation(user *User) (*[]Education, error) {
	k := rand.IntN(2) + 3

	educations := make([]Education, k)
	for i := 0; i < k; i++ {
		education := Education{}
		if err := gofakeit.Struct(&education); err != nil {
			return nil, err
		}

		attributes := map[string]interface{}{}

		if education.Type == "SCHOOL" {
			attributes["class"] = gofakeit.RandomString([]string{"X", "XII"})
			attributes["passing_year"] = gofakeit.Year()
		}

		if education.Type == "COLLEGE" {
			attributes["degree"] = gofakeit.RandomString(CollageDegrees)
			attributes["field_of_study"] = gofakeit.RandomString([]string{"Physics", "Chemistry", "Maths", "Biology", "Computer Science", "Chemical Engineering", "Electrical Engineering", "Mechanical Engineering"})
			attributes["start_year"] = gofakeit.Year()
			attributes["end_year"] = gofakeit.Year()
		}
		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return nil, err
		}

		education.UserId = user.Id
		education.OrderIndex = int16(i + 1)
		education.Grade = strconv.FormatFloat(gofakeit.Float64Range(5.0, 10.0), 'f', 2, 64)
		education.Attributes = attributesJson

		educations[i] = education
	}

	return &educations, nil
}

type Hackathon struct {
	UserId          string         `json:"user_id"`
	OrderIndex      int16          `json:"order_index"`
	Avatar          *string        `json:"avatar"`
	Title           string         `json:"title" fake:"{sentence:5}"`
	Location        string         `json:"location"`
	StartDate       time.Time      `json:"start_date" fake:"{date}"`
	EndDate         time.Time      `json:"end_date" fake:"{date}"`
	Description     string         `json:"description" fake:"{paragraph:2,3,10,\n\n}"`
	CertificateLink *string        `json:"certificate_link"`
	Attributes      datatypes.JSON `json:"attributes"`
	CreatedAt       time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt       time.Time      `json:"updated_at" fake:"{date}"`
	// DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}

func (Hackathon) TableName() string {
	return "user_hackathons"
}

func NewHackathon(user *User) (*[]Hackathon, error) {
	k := rand.IntN(2) + 3

	hackathons := make([]Hackathon, k)
	for i := 0; i < k; i++ {
		hackathon := Hackathon{}
		if err := gofakeit.Struct(&hackathon); err != nil {
			return nil, err
		}

		attributes := map[string]interface{}{}

		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return nil, err
		}

		hackathon.UserId = user.Id
		hackathon.OrderIndex = int16(i + 1)
		hackathon.Title = strings.TrimSuffix(hackathon.Title, ".")
		hackathon.Location = gofakeit.StreetName() + ", " + gofakeit.City() + ", " + gofakeit.State() + ", " + gofakeit.Zip() + ", " + gofakeit.Country()
		hackathon.Attributes = attributesJson
		hackathon.CertificateLink = nil
		hackathon.Avatar = nil

		hackathons[i] = hackathon
	}

	return &hackathons, nil
}

type Certification struct {
	UserId          string         `json:"user_id"`
	OrderIndex      int16          `json:"order_index"`
	Title           string         `json:"title" fake:"{sentence:5}"`
	Description     pq.StringArray `json:"description" gorm:"type:text"`
	CompletionDate  time.Time      `json:"completion_date" fake:"{date}"`
	CertificateLink *string        `json:"certificate_link"`
	SkillsUsed      pq.StringArray `json:"skills_used" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt       time.Time      `json:"updated_at" fake:"{date}"`
	// DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

func (Certification) TableName() string {
	return "user_certifications"
}

func NewCertification(user *User) (*[]Certification, error) {
	k := rand.IntN(2) + 3

	certifications := make([]Certification, k)
	for i := 0; i < k; i++ {
		certification := Certification{}
		if err := gofakeit.Struct(&certification); err != nil {
			return nil, err
		}

		certification.UserId = user.Id
		certification.OrderIndex = int16(i + 1)
		certification.Title = strings.TrimSuffix(certification.Title, ".")
		certification.CertificateLink = nil
		k := rand.IntN(2) + 3
		certification.Description = make([]string, k)
		for j := 0; j < k; j++ {
			certification.Description[j] = gofakeit.Sentence(10)
		}
		k = rand.IntN(2) + 3
		certification.SkillsUsed = make([]string, k)
		for j := 0; j < k; j++ {
			certification.SkillsUsed[j] = gofakeit.Word()
		}

		certifications[i] = certification
	}

	return &certifications, nil
}

type Experience struct {
	UserId          string         `json:"user_id"`
	OrderIndex      int16          `json:"order_index"`
	CompanyName     string         `json:"company_name"`
	CompanyUrl      string         `json:"company_url" fake:"{url}"`
	JobType         string         `json:"job_type" fake:"{randomstring:[PART_TIME,SEMI_FULL_TIME,FULL_TIME]}"`
	JobTitle        string         `json:"job_title" fake:"{jobtitle}"`
	Location        string         `json:"location"`
	StartDate       time.Time      `json:"start_date" fake:"{date}"`
	EndDate         *time.Time     `json:"end_date" fake:"{date}"`
	Description     pq.StringArray `json:"description" gorm:"type:text"`
	CertificateLink *string        `json:"certificate_link"`
	SkillsUsed      pq.StringArray `json:"skills_used" gorm:"type:text"`
	CreatedAt       time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt       time.Time      `json:"updated_at" fake:"{date}"`
	// DeletedAt       gorm.DeletedAt `json:"deleted_at"`
}

func (Experience) TableName() string {
	return "user_experiences"
}

func NewExperience(user *User) (*[]Experience, error) {
	k := rand.IntN(2) + 3

	experiences := make([]Experience, k)
	for i := 0; i < k; i++ {
		experience := Experience{}
		if err := gofakeit.Struct(&experience); err != nil {
			return nil, err
		}

		experience.UserId = user.Id
		experience.OrderIndex = int16(i + 1)
		experience.Location = gofakeit.StreetName() + ", " + gofakeit.City() + ", " + gofakeit.State() + ", " + gofakeit.Zip() + ", " + gofakeit.Country()
		experience.CertificateLink = nil
		experience.CompanyName = gofakeit.Company() + " " + gofakeit.CompanySuffix()
		k := rand.IntN(2) + 3
		experience.Description = make([]string, k)
		for j := 0; j < k; j++ {
			experience.Description[j] = gofakeit.Sentence(10)
		}
		k = rand.IntN(2) + 3
		experience.SkillsUsed = make([]string, k)
		for j := 0; j < k; j++ {
			experience.SkillsUsed[j] = gofakeit.Word()
		}

		experiences[i] = experience
	}

	return &experiences, nil
}

func SeedPortfolio(db *gorm.DB, users *[]User) error {
	var educations []Education
	var hackathons []Hackathon
	var certifications []Certification
	var experiences []Experience

	for _, u := range *users {
		edus, err := NewEducation(&u)
		if err != nil {
			return err
		}

		hack, err := NewHackathon(&u)
		if err != nil {
			return err
		}

		cert, err := NewCertification(&u)
		if err != nil {
			return err
		}

		exps, err := NewExperience(&u)
		if err != nil {
			return err
		}

		educations = append(educations, *edus...)
		hackathons = append(hackathons, *hack...)
		certifications = append(certifications, *cert...)
		experiences = append(experiences, *exps...)
	}

	if err := db.Create(&educations).Error; err != nil {
		return err
	}

	if err := db.Create(&hackathons).Error; err != nil {
		return err
	}

	if err := db.Create(&certifications).Error; err != nil {
		return err
	}

	if err := db.Create(&experiences).Error; err != nil {
		return err
	}

	return nil
}
