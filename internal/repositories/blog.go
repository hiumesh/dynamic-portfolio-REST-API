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
	GetAll(userId *string, query *string, cursor int, limit int) (any, error)
	GetUserBlogs(userId *string, query *string, cursor int, limit int) (any, error)
	Get(userId string, id string) (any, error)
	GetBlogBySlug(slug string) (*schemas.SchemaBlog, error)
	Create(userId string, tags *models.Tags, data *schemas.SchemaBlog, publish bool) (*models.Blog, error)
	Update(userId string, id string, tags *models.Tags, data *schemas.SchemaBlog, publish bool) (*models.Blog, error)
	Unpublish(userId string, id string) error
	Delete(userId string, id string) error
	GetCommentBlog(id uint) (*models.BlogComment, error)
	CreateComment(blogId uint, commentId uint) (any, error)
	Reaction(blogId uint, userId uuid.UUID, data *schemas.SchemaReaction) (any, error)
}

type repositoryBlog struct {
	db *gorm.DB
}

func (r *repositoryBlog) GetAll(userId *string, query *string, cursor int, limit int) (*[]schemas.SelectBlog, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		select
			blogs.id,
			blogs.cover_image,
			blogs.title,
			blogs.slug,
			user_profiles.user_id as publisher_id,
			user_profiles.avatar_url as publisher_avatar,
			user_profiles.full_name as publisher_name,
			blogs.attributes ->> 'comments_count' as comments_count,
			blogs.attributes -> 'reactions_metadata' as reactions_metadata,
			blogs.published_at,
			blogs.created_at,
			blogs.updated_at,
			array_remove(array_agg(tags.name), NULL) AS tags
		from
			blogs
			inner join user_profiles on user_profiles.user_id = blogs.user_id
			left join blog_tags on blog_tags.blog_id = blogs.id
			left join tags on tags.id = blog_tags.tag_id
		where
			blogs.published_at is not null
	`

	var args []any
	args = append(args, limit, cursor)

	if query != nil && *query != "" {
		baseQuery += " AND blogs.fts @@ to_tsquery(?)"
		args = append([]interface{}{*query}, args...)
	}

	baseQuery += `
		group by
			blogs.id,
			user_profiles.user_id
		ORDER BY blogs.updated_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []schemas.SelectBlog
	for rows.Next() {
		var blog schemas.SelectBlog

		err = rows.Scan(&blog.ID, &blog.CoverImage, &blog.Title, &blog.Slug, &blog.PublisherId, &blog.PublisherAvatar, &blog.PublisherName, &blog.CommentsCount, &blog.ReactionsMetadata, &blog.PublishedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Tags)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return &blogs, nil
}

func (r *repositoryBlog) GetUserBlogs(userId string, query *string, cursor int, limit int) (*[]schemas.SelectBlog, error) {
	var rows *sql.Rows
	var err error

	baseQuery := `
		select
			blogs.id,
			blogs.cover_image,
			blogs.title,
			blogs.slug,
			blogs.attributes ->> 'comments_count' as comments_count,
			blogs.attributes -> 'reactions_metadata' as reactions_metadata,
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
	`

	var args []interface{}
	args = append(args, limit, cursor)

	if query != nil && *query != "" {
		baseQuery += " AND blogs.fts @@ to_tsquery(?)"
		args = append([]interface{}{*query}, args...)
	}

	args = append([]interface{}{userId}, args...)

	baseQuery += `
		group by
			blogs.id
		ORDER BY blogs.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []schemas.SelectBlog
	for rows.Next() {
		var blog schemas.SelectBlog

		err = rows.Scan(&blog.ID, &blog.CoverImage, &blog.Title, &blog.Slug, &blog.CommentsCount, &blog.ReactionsMetadata, &blog.PublishedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Tags)
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
			blogs.attributes ->> 'comments_count' as comments_count,
			blogs.attributes -> 'reactions_metadata' as reactions_metadata,
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
		err = rows.Scan(&blog.ID, &blog.CoverImage, &blog.Title, &blog.Body, &blog.Slug, &blog.CommentsCount, &blog.ReactionsMetadata, &blog.PublishedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Tags)
		if err != nil {
			return nil, err
		}
	}

	if blog.ID == 0 {
		return nil, errors.New("record not found")
	}

	return &blog, nil
}

func (r *repositoryBlog) GetBlogBySlug(slug string) (*schemas.SelectBlog, error) {
	var rows *sql.Rows
	var err error

	rows, err = r.db.Raw(`
		select
			blogs.id,
			blogs.cover_image,
			blogs.title,
			blogs.body,
			blogs.slug,
			user_profiles.user_id as publisher_id,
			user_profiles.avatar_url as publisher_avatar,
			user_profiles.full_name as publisher_name,
			blogs.attributes ->> 'comments_count' as comments_count,
			blogs.attributes -> 'reactions_metadata' as reactions_metadata,
			blogs.published_at,
			blogs.created_at,
			blogs.updated_at,
			array_remove(array_agg(tags.name), NULL) AS tags
		from
			blogs
			inner join user_profiles on user_profiles.user_id = blogs.user_id
			left join blog_tags on blog_tags.blog_id = blogs.id
			left join tags on tags.id = blog_tags.tag_id
		where
			blogs.published_at is not null
			and blogs.slug = ?
		group by
			blogs.id,
			user_profiles.user_id
			`, slug).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blog schemas.SelectBlog
	for rows.Next() {
		err = rows.Scan(&blog.ID, &blog.CoverImage, &blog.Title, &blog.Body, &blog.Slug, &blog.PublisherId, &blog.PublisherAvatar, &blog.PublisherName, &blog.CommentsCount, &blog.ReactionsMetadata, &blog.PublishedAt, &blog.CreatedAt, &blog.UpdatedAt, &blog.Tags)
		if err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if blog.ID == 0 {
		return nil, errors.New("failed to get blog")
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

func (r *repositoryBlog) GetCommentBlog(id uint) (*models.BlogComment, error) {
	var commentBlog models.BlogComment

	if err := r.db.Where("comment_id = ?", id).First(&commentBlog).Error; err != nil {
		return nil, err
	}

	return &commentBlog, nil
}

func (r *repositoryBlog) CreateComment(blogId uint, commentId uint) (any, error) {

	comment := models.BlogComment{
		BlogId:    blogId,
		CommentId: commentId,
	}

	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (r *repositoryBlog) Reaction(blogId uint, userId uuid.UUID, data *schemas.SchemaReaction) (any, error) {

	reaction := models.BlogReaction{
		BlogId: blogId,
		UserId: userId,
		Type:   data.Reaction,
	}

	if err := r.db.Create(&reaction).Error; err != nil {
		return nil, err
	}

	return reaction, nil
}

func NewBlogRepository(db *gorm.DB) *repositoryBlog {
	return &repositoryBlog{
		db: db,
	}
}
