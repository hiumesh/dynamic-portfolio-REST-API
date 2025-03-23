package seeds

import (
	"encoding/json"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/lib/pq"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TechProject struct {
	UserId      string         `json:"user_id"`
	OrderIndex  int16          `json:"order_index"`
	Title       string         `json:"title" fake:"{sentence:5}"`
	StartDate   time.Time      `json:"start_date" fake:"{date}"`
	EndDate     *time.Time     `json:"end_date" fake:"{date}"`
	Description string         `json:"description" fake:"{paragraph:2,3,10,\n\n}"`
	SkillsUsed  pq.StringArray `json:"skills_used" gorm:"type:text"`
	Attributes  datatypes.JSON `json:"attributes"`
	CreatedAt   time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt   time.Time      `json:"updated_at" fake:"{date}"`
}

func (TechProject) TableName() string {
	return "tech_projects"
}

func NewTechProject(user *User) (*[]TechProject, error) {
	k := rand.IntN(2) + 3

	projects := make([]TechProject, k)
	for i := 0; i < k; i++ {
		project := TechProject{}
		if err := gofakeit.Struct(&project); err != nil {
			return nil, err
		}

		attributes := map[string]interface{}{}

		h := rand.IntN(2) + 3
		links := make([]interface{}, h)
		for j := 0; j < h; j++ {
			links[j] = map[string]interface{}{
				"platform": gofakeit.RandomString([]string{"Github", "Website", "Social"}),
				"label":    gofakeit.Word(),
				"url":      gofakeit.URL(),
			}
		}

		attributes["links"] = links

		attributesJson, err := json.Marshal(attributes)
		if err != nil {
			return nil, err
		}

		project.UserId = user.Id
		project.OrderIndex = int16(i + 1)
		project.Title = strings.TrimSuffix(project.Title, ".")
		h = rand.IntN(2) + 3
		project.SkillsUsed = make([]string, h)
		for j := 0; j < h; j++ {
			project.SkillsUsed[j] = gofakeit.Word()
		}
		project.Attributes = attributesJson

		projects[i] = project
	}

	return &projects, nil
}

func SeedTechProjects(db *gorm.DB, users *[]User) error {
	var projects []TechProject

	for _, u := range *users {
		pro, err := NewTechProject(&u)
		if err != nil {
			return err
		}
		projects = append(projects, *pro...)
	}

	if err := db.Create(&projects).Error; err != nil {
		return err
	}

	return nil
}
