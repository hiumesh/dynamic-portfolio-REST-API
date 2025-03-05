package seeds

import (
	"math/rand/v2"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var resc = regexp.MustCompile(`[^\w\s]`)

type Blog struct {
	ID          int64      `gorm:"primaryKey;autoIncrement"`
	UserId      string     `json:"user_id"`
	Title       string     `json:"title" fake:"{sentence:5}"`
	Body        string     `json:"body" fake:"{markdown}"`
	Slug        string     `json:"slug"`
	PublishedAt *time.Time `json:"published_at" fake:"{date}"`
	CreatedAt   time.Time  `json:"created_at" fake:"{date}"`
	UpdatedAt   time.Time  `json:"updated_at" fake:"{date}"`
	DeletedAt   *time.Time `json:"deleted_at"`
}

func (Blog) TableName() string {
	return "blogs"
}

type BlogTag struct {
	BlogId int64 `json:"blog_id"`
	TagId  int64 `json:"tag_id"`
}

func (BlogTag) TableName() string {
	return "blog_tags"
}

func NewBlog(users *[]User) (*Blog, error) {
	var uid = uuid.NewString()

	userIdx := rand.IntN(len(*users))
	user := (*users)[userIdx]

	blog := Blog{}
	if err := gofakeit.Struct(&blog); err != nil {
		return nil, err
	}
	blog.ID = 0
	blog.UserId = user.Id
	blog.Title = strings.TrimSuffix(blog.Title, ".")
	blog.Slug = url.QueryEscape(strings.Join(strings.Fields(strings.ToLower(resc.ReplaceAllString(blog.Title, ""))), "-") + "-" + strings.Split(uid, "-")[0])

	return &blog, nil
}

func NewBlogTags(blog *Blog, tags *[]Tag) (*[]BlogTag, error) {

	n := len(*tags)
	k := rand.IntN(3) + 4
	blogTags := make([]BlogTag, k)
	uniqueIdx := make(map[int]bool)
	if k < n {
		for len(uniqueIdx) < k {
			idx := rand.IntN(n)
			if _, ok := uniqueIdx[idx]; !ok {
				blogTag := BlogTag{
					BlogId: blog.ID,
					TagId:  (*tags)[idx].ID,
				}
				blogTags[len(uniqueIdx)] = blogTag
				uniqueIdx[idx] = true
			}

		}
	}

	return &blogTags, nil
}

func SeedBlogs(db *gorm.DB, users *[]User, tags *[]Tag, count int) error {

	if count == 0 {
		count = 100
	}

	blogs := make([]Blog, count)
	for i := 0; i < count; i++ {
		blog, err := NewBlog(users)
		if err != nil {
			return err
		}
		blogs[i] = *blog
	}

	if err := db.Create(&blogs).Error; err != nil {
		return err
	}

	var blogsTags []BlogTag
	for i := 0; i < count; i++ {
		blogTags, err := NewBlogTags(&blogs[i], tags)
		if err != nil {
			return err
		}
		blogsTags = append(blogsTags, *blogTags...)
	}

	if err := db.Create(&blogsTags).Error; err != nil {
		return err
	}

	return nil
}
