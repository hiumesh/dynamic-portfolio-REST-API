package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceBlog interface {
	GetAll(userId *string, query *string, cursor int, limit int) (any, error)
	GetUserBlogs(userId string, query *string, cursor int, limit int) (any, error)
	Get(userId string, blogId string) (any, error)
	GetBlogBySlug(slug string) (any, error)
	Create(userId string, data *schemas.SchemaBlog, publish bool) (*models.Blog, error)
	Update(userId string, blogId string, data *schemas.SchemaBlog, publish bool) (*models.Blog, error)
	Unpublish(userId string, blogId string) error
	Delete(userId string, blogId string) error
	GetMetadata(userId string) (any, error)
	UpdateMetadata(userId string, data *schemas.SchemaBlogMetadata) error
}

type serviceBlog struct {
	db *gorm.DB
}

func (s *serviceBlog) GetAll(userId *string, query *string, cursor int, limit int) (any, error) {
	blogRepository := repositories.NewBlogRepository(s.db)

	res, err := blogRepository.GetAll(userId, query, cursor, limit)
	if err != nil {
		return nil, err
	}

	var nextCursor *int
	if len(*res) == limit {
		temp := cursor + limit
		nextCursor = &temp
	}

	return map[string]any{"list": res, "cursor": nextCursor}, nil
}

func (s *serviceBlog) GetUserBlogs(userId string, query *string, cursor int, limit int) (any, error) {
	blogRepository := repositories.NewBlogRepository(s.db)

	res, err := blogRepository.GetUserBlogs(userId, query, cursor, limit)
	if err != nil {
		return nil, err
	}

	var nextCursor *int
	if len(*res) == limit {
		temp := cursor + limit
		nextCursor = &temp
	}

	return map[string]any{"list": res, "cursor": nextCursor}, nil
}

func (s *serviceBlog) Get(userId string, id string) (any, error) {
	blogRepository := repositories.NewBlogRepository(s.db)

	res, err := blogRepository.Get(userId, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceBlog) GetBlogBySlug(slug string) (any, error) {
	blogRepository := repositories.NewBlogRepository(s.db)

	res, err := blogRepository.GetBlogBySlug(slug)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceBlog) Create(userId string, data *schemas.SchemaBlog, publish bool) (*models.Blog, error) {
	var blog *models.Blog

	err := s.db.Transaction(func(tx *gorm.DB) error {
		blogRepository := repositories.NewBlogRepository(tx)
		tagRepository := repositories.NewTagRepository(tx)

		tags, err := tagRepository.FindOrCreate(userId, data.Tags)
		if err != nil {
			return err
		}

		blog, err = blogRepository.Create(userId, tags, data, publish)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return blog, nil

}

func (s *serviceBlog) Update(userId string, id string, data *schemas.SchemaBlog, publish bool) (*models.Blog, error) {
	var blog *models.Blog

	err := s.db.Transaction(func(tx *gorm.DB) error {
		blogRepository := repositories.NewBlogRepository(tx)
		tagRepository := repositories.NewTagRepository(tx)

		tags, err := tagRepository.FindOrCreate(userId, data.Tags)
		if err != nil {
			return err
		}

		blog, err = blogRepository.Update(userId, id, tags, data, publish)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (s *serviceBlog) Unpublish(userId string, id string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		blogRepository := repositories.NewBlogRepository(tx)

		if err := blogRepository.Unpublish(userId, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serviceBlog) Delete(userId string, id string) error {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		blogRepository := repositories.NewBlogRepository(tx)

		if err := blogRepository.Delete(userId, id); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *serviceBlog) GetMetadata(userId string) (any, error) {
	userRepository := repositories.NewUserRepository(s.db)

	res, err := userRepository.GetModuleMetadata(userId, "blog")
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *serviceBlog) UpdateMetadata(userId string, data *schemas.SchemaBlogMetadata) error {
	userRepository := repositories.NewUserRepository(s.db)

	err := userRepository.AddOrUpdateModuleMetadata(userId, "blog", data)
	if err != nil {
		return err
	}

	return nil
}

func NewBlogService(db *gorm.DB) *serviceBlog {
	return &serviceBlog{
		db: db,
	}
}
