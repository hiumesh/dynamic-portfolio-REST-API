package repositories

import (
	"database/sql"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RepositoryBlog interface {
	GetAll(userId string) (interface{}, error)
	Get(userId string, id string) (interface{}, error)
	Create(userId string, tags *models.Tags, data *schemas.SchemaBlog, publish bool) (*models.Blog, error)
	Update(userId string, id string, tags *models.Tags, data *schemas.SchemaBlog, publish bool) (*models.Blog, error)
	Unpublish(userId string, id string) error
	Delete(userId string, id string) error
}

type repositoryBlog struct {
	db *gorm.DB
}

func (r *repositoryBlog) GetAll(userId string) (interface{}, error) {
	var rows *sql.Rows
	var err error

	rows, err = r.db.Raw(`
		select
			blogs.id,
			blogs.cover_image,
			blogs.title,
			blogs.slug,
			blogs.published_at,
			blogs.created_at,
			blogs.updated_at,
			array_remove(array_agg(tags.name), NULL) AS tags
		from
			blogs
			left join blog_tags on blog_tags.blog_id = blogs.id
			left join tags on tags.id = blog_tags.tag_id
		where
			blogs.user_id = ?
		group by
			blogs.id
		order by
			created_at desc

	`, userId).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []schemas.SelectBlog
	for rows.Next() {
		var blog schemas.SelectBlog

		err = rows.Scan(&blog.ID, &blog.CoverImage, &blog.Title, &blog.Slug, &blog.PublishedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Tags)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return &blogs, nil
}

func (r *repositoryBlog) Get(userId string, id string) (interface{}, error) {
	var rows *sql.Rows
	var err error

	rows, err = r.db.Raw(`
		select
			blogs.id,
			blogs.cover_image,
			blogs.title,
			blogs.body,
			blogs.slug,
			blogs.published_at,
			blogs.created_at,
			blogs.updated_at,
			array_remove(array_agg(tags.name), NULL) AS tags
		from
			blogs
			left join blog_tags on blog_tags.blog_id = blogs.id
			left join tags on tags.id = blog_tags.tag_id
		where
			blogs.user_id = ? and blogs.id = ?
		group by
			blogs.id
	`, userId, id).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blog schemas.SelectBlog
	for rows.Next() {
		err = rows.Scan(&blog.ID, &blog.CoverImage, &blog.Title, &blog.Body, &blog.Slug, &blog.PublishedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Tags)
		if err != nil {
			return nil, err
		}
	}

	if blog.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &blog, nil
}

func (r *repositoryBlog) Create(userId string, tags *models.Tags, data *schemas.SchemaBlog, publish bool) (*models.Blog, error) {
	var uid = uuid.NewString()
	var slug = url.QueryEscape(strings.Join(strings.Fields(strings.ToLower(data.Title)), "-") + "-" + strings.Split(uid, "-")[0])

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	var blog = models.Blog{
		UserId:     userUUID,
		Title:      data.Title,
		Body:       &data.Body,
		Slug:       slug,
		Attributes: []byte("{}"),
	}

	if data.CoverImage != "" {
		blog.CoverImage = &data.CoverImage
	}

	if tags != nil {
		blog.Tags = *tags
	}

	if publish {
		now := time.Now()
		blog.PublishedAt = &now
	}

	if err := r.db.Create(&blog).Error; err != nil {
		return nil, err
	}

	blog.Tags = *tags

	return &blog, nil
}

func (r *repositoryBlog) Update(userId string, id string, tags *models.Tags, data *schemas.SchemaBlog, publish bool) (*models.Blog, error) {
	var blog models.Blog
	if err := r.db.Where("id = ? and user_id = ?", id, userId).First(&blog).Error; err != nil {
		return nil, err
	}

	if blog.ID == 0 {
		return nil, errors.New("record not found")
	}

	var blogData = map[string]interface{}{
		"title":       data.Title,
		"body":        data.Body,
		"cover_image": data.CoverImage,
	}

	if publish && blog.PublishedAt == nil {
		now := time.Now()
		blogData["published_at"] = now
	}

	var updatedRows models.Blogs
	if err := r.db.Model(&updatedRows).Clauses(clause.Returning{}).Where("id = ? and user_id = ?", id, userId).Updates(blogData).Error; err != nil {
		return nil, err
	}

	if len(updatedRows) == 0 {
		return nil, errors.New("record not found")
	}

	blog = updatedRows[0]

	if err := r.db.Where("blog_id = ?", id).Delete(&models.BlogTags{}).Error; err != nil {
		return nil, err
	}

	var blogTags models.BlogTags
	if tags != nil {
		for _, tag := range *tags {
			var blogTag = models.BlogTag{
				BlogId: blog.ID,
				TagId:  tag.ID,
			}

			blogTags = append(blogTags, blogTag)
		}
	}

	if len(blogTags) > 0 {
		if err := r.db.Create(&blogTags).Error; err != nil {
			return nil, err
		}
	}

	blog.Tags = *tags

	return &blog, nil

}

func (r *repositoryBlog) Unpublish(userId string, id string) error {
	var blog models.Blog
	if err := r.db.Where("id = ? and user_id = ?", id, userId).First(&blog).Error; err != nil {
		return err
	}

	if blog.ID == 0 {
		return errors.New("record not found")
	}

	if blog.PublishedAt == nil {
		return errors.New("blog is already unpublished")
	}

	if err := r.db.Model(&models.Blog{}).Where("id = ? and user_id = ?", id, userId).Update("published_at", nil).Error; err != nil {
		return err
	}

	return nil
}

func (r *repositoryBlog) Delete(userId string, id string) error {
	var blog models.Blog

	if err := r.db.Where("user_id = ? and id = ?", userId, id).Delete(&blog).Error; err != nil {
		return err
	}

	return nil
}

func NewBlogRepository(db *gorm.DB) *repositoryBlog {
	return &repositoryBlog{
		db: db,
	}
}
