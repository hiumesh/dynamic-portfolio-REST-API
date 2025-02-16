package seeds

import (
	"encoding/json"
	"errors"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var SocialPlatforms = []string{
	"Github",
	"LinkedIn",
	"Hacker Rank",
	"GeeksforGeeks",
	"CodeChef",
	"LeetCode",
	"Codeforces",
	"Topcoder",
	"Hackerearth",
	"Behance",
	"Blog",
	"Portfolio Website",
	"Dribble",
	"Other Profile Link",
}

type User struct {
	Id                string         `fake:"{uuid}" json:"id"`
	InstanceId        string         `json:"instance_id" default:"00000000-0000-0000-0000-000000000000"`
	Aud               string         `json:"aud" default:"authenticated"`
	Role              string         `json:"role" default:"authenticated"`
	Email             string         `fake:"{email}" json:"email"`
	EncryptedPassword string         `json:"encrypted_password"`
	LastSignInAt      time.Time      `json:"last_sign_in_at" fake:"{date}"`
	EmailConfirmedAt  time.Time      `json:"email_confirmed_at" fake:"{date}"`
	RawAppMetaData    datatypes.JSON `json:"app_metadata"`
	RawUserMetaData   datatypes.JSON `json:"user_metadata"`
	CreatedAt         time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt         time.Time      `json:"updated_at" fake:"{date}"`
	DeletedAt         *time.Time     `json:"deleted_at"`
}

func (User) TableName() string {
	return "auth.users"
}

type Identities struct {
	UserId       string         `json:"user_id"`
	IdentityData datatypes.JSON `json:"identity_data"`
	Provider     string         `json:"provider"`
	ProviderId   string         `json:"provider_id"`
	LastSignInAt time.Time      `json:"last_sign_in_at" fake:"{date}"`
	CreatedAt    time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt    time.Time      `json:"updated_at" fake:"{date}"`
}

func (Identities) TableName() string {
	return "auth.identities"
}

func NewIdentity(user *User) (*Identities, error) {
	identity := Identities{}
	if err := gofakeit.Struct(&identity); err != nil {
		logrus.Error(err)
		return nil, err
	}
	identity.UserId = user.Id
	identity.IdentityData = user.RawUserMetaData
	identity.Provider = "email"
	identity.ProviderId = user.Id

	return &identity, nil
}

type UserProfile struct {
	UserId          string         `json:"user_id"`
	Email           string         `json:"email"`
	FullName        string         `json:"full_name"`
	Slug            string         `json:"slug"`
	PortfolioStatus string         `json:"status"`
	Attributes      datatypes.JSON `json:"attributes"`
	CreatedAt       time.Time      `json:"created_at" fake:"{date}"`
	UpdatedAt       time.Time      `json:"updated_at" fake:"{date}"`
	DeletedAt       *time.Time     `json:"deleted_at"`
}

func hashPassword(password string) (string, error) {
	const cost = 10
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func NewUser(password string) (*User, error) {
	user := User{}
	if err := gofakeit.Struct(&user); err != nil {
		logrus.Error(err)
		return nil, err
	}

	user.EncryptedPassword = password
	user.InstanceId = "00000000-0000-0000-0000-000000000000"
	user.Aud = "authenticated"
	user.Role = "authenticated"
	rawAppMetaData := map[string]interface{}{
		"provider":  "email",
		"providers": []string{"email"},
	}
	rawAppMetaDataJson, err := json.Marshal(rawAppMetaData)
	if err != nil {
		return nil, errors.New("failed to marshal json")
	}
	user.RawAppMetaData = rawAppMetaDataJson

	rawUserMetaData := map[string]interface{}{
		"sub":            user.Id,
		"email":          user.Email,
		"email_verified": true,
		"phone_verified": false,
	}

	rawUserMetaDataJson, err := json.Marshal(rawUserMetaData)
	if err != nil {
		return nil, errors.New("failed to marshal json")
	}
	user.RawUserMetaData = rawUserMetaDataJson
	return &user, nil
}

func NewUserProfile(user *User) (*UserProfile, error) {
	profile := UserProfile{}
	if err := gofakeit.Struct(&profile); err != nil {
		logrus.Error(err)
		return nil, err
	}
	profile.UserId = user.Id
	profile.Email = user.Email
	profile.FullName = gofakeit.Name()
	profile.Slug = strings.Join(strings.Fields(strings.ToLower(profile.FullName)), "-") + "-" + strings.Split(user.Id, "-")[0]
	profile.PortfolioStatus = "ACTIVE"
	k := rand.IntN(2) + 3
	workDomains := make([]string, k)
	for j := 0; j < k; j++ {
		workDomains[j] = gofakeit.Word()
	}
	k = rand.IntN(2) + 3
	socialLinks := make([]interface{}, k)
	for j := 0; j < k; j++ {
		socialLinks[j] = map[string]interface{}{
			"platform": gofakeit.RandomString(SocialPlatforms),
			"url":      gofakeit.URL(),
		}
	}

	attributes := map[string]interface{}{
		"about":           gofakeit.Paragraph(2, 4, 500, "\n\n"),
		"tagline":         gofakeit.Sentence(50),
		"college":         gofakeit.School(),
		"graduation_year": gofakeit.Year(),
		"work_domains":    workDomains,
		"social_links":    socialLinks,
	}

	attributesJson, err := json.Marshal(attributes)
	if err != nil {
		return nil, errors.New("failed to marshal json")
	}
	profile.Attributes = attributesJson
	return &profile, nil
}

func SeedUsers(db *gorm.DB, count int) (*[]User, error) {
	if count == 0 {
		count = 100
	}

	password, err := hashPassword("uswag007")
	if err != nil {
		return nil, err
	}

	users := make([]User, count)
	identities := make([]Identities, count)
	profiles := make([]UserProfile, count)
	for i := 0; i < count; i++ {
		user, err := NewUser(password)
		if err != nil {
			return nil, err
		}
		users[i] = *user

		identity, err := NewIdentity(user)
		if err != nil {
			return nil, err
		}
		identities[i] = *identity

		profile, err := NewUserProfile(user)
		if err != nil {
			return nil, err
		}
		profiles[i] = *profile
	}

	if err := db.Create(&users).Error; err != nil {
		return nil, err
	}
	if err := db.Create(&identities).Error; err != nil {
		return nil, err
	}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"full_name", "slug", "attributes", "portfolio_status"}),
	}).Create(&profiles).Error; err != nil {
		return nil, err
	}

	return &users, nil
}
