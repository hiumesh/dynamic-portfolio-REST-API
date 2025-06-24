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

type BlogReaction struct {
	ID     int64  `json:"id" gorm:"primaryKey;autoIncrement"`
	BlogId int64  `json:"blog_id"`
	UserId string `json:"user_id"`
	Type   string `json:"type" fake:"{randomstring:[like,clap,heart]}"`
}

func (BlogReaction) TableName() string {
	return "blog_reactions"
}

type BlogComment struct {
	BlogId    int64 `json:"blog_id"`
	CommentId int64 `json:"comment_id"`
}

func (BlogComment) TableName() string {
	return "blog_comments"
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

func NewBlogReactions(blog *Blog, users *[]User) (*[]BlogReaction, error) {

	n := len(*users)
	k := rand.IntN(3) + 7
	blogReactions := make([]BlogReaction, k)
	uniqueIdx := make(map[int]bool)
	if k < n {
		for len(uniqueIdx) < k {
			idx := rand.IntN(n)
			if _, ok := uniqueIdx[idx]; !ok {
				blogReaction := BlogReaction{}
				if err := gofakeit.Struct(&blogReaction); err != nil {
					return nil, err
				}

				blogReaction.ID = 0
				blogReaction.BlogId = blog.ID
				blogReaction.UserId = (*users)[idx].Id

				blogReactions[len(uniqueIdx)] = blogReaction
				uniqueIdx[idx] = true
			}
		}
	}

	return &blogReactions, nil
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

	var blogsReactions []BlogReaction
	for i := 0; i < count; i++ {
		blogReactions, err := NewBlogReactions(&blogs[i], users)
		if err != nil {
			return err
		}
		blogsReactions = append(blogsReactions, *blogReactions...)
	}

	if err := db.Create(&blogsReactions).Error; err != nil {
		return err
	}

	comments, blogIdToCommentsIdxMap, err := NewBlogComments(&blogs, users)
	if err != nil {
		return err
	}

	if err := db.Create(&comments).Error; err != nil {
		return err
	}

	blogComments := []BlogComment{}
	for _, blog := range blogs {

		if _, ok := blogIdToCommentsIdxMap[blog.ID]; ok {
			for _, idx := range blogIdToCommentsIdxMap[blog.ID] {
				blogComments = append(blogComments, BlogComment{
					BlogId:    blog.ID,
					CommentId: (*comments)[idx].ID,
				})
			}
		}
	}

	if err := db.Create(&blogComments).Error; err != nil {
		return err
	}

	blogChildComments, blogIdToChildCommentsIdxMap, err := NewBlogChildComments(&blogComments, users)
	if err != nil {
		return err
	}

	if err := db.Create(&blogChildComments).Error; err != nil {
		return err
	}

	moreBlogComments := []BlogComment{}
	for _, blog := range blogs {
		if _, ok := blogIdToChildCommentsIdxMap[blog.ID]; ok {
			for _, idx := range blogIdToChildCommentsIdxMap[blog.ID] {
				moreBlogComments = append(moreBlogComments, BlogComment{
					BlogId:    blog.ID,
					CommentId: (*blogChildComments)[idx].ID,
				})
			}
		}
	}
	if err := db.Create(&moreBlogComments).Error; err != nil {
		return err
	}

	totalComments := append(*comments, *blogChildComments...)

	if err := NewCommentsReactions(db, &totalComments, users); err != nil {
		return err
	}

	return nil
}
