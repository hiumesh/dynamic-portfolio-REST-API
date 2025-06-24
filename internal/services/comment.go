package services

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type ServiceComment interface {
	GetAll(userId *string, module string, slug string, cursor int, limit int, parentId *int) (any, error)
	Create(userId string, data *schemas.SchemaCreateComment) (any, error)
	Reaction(commentId string, userId string, data *schemas.SchemaReaction) (any, error)
	Reply(userId string, commentId string, data *schemas.SchemaCommentReply) (any, error)
}

type serviceComment struct {
	db *gorm.DB
}

func (s *serviceComment) GetAll(userId *string, module string, slug string, cursor int, limit int, parentId *int) (any, error) {

	commentRepository := repositories.NewCommentRepository(s.db)

	res, err := commentRepository.Get(userId, module, slug, cursor, limit, parentId)
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

func (s *serviceComment) Create(userId string, data *schemas.SchemaCreateComment) (any, error) {

	commentRepository := repositories.NewCommentRepository(s.db)
	blogRepository := repositories.NewBlogRepository(s.db)

	err := s.db.Transaction(func(tx *gorm.DB) error {
		comment, err := commentRepository.Create(userId, data)
		if err != nil {
			return err
		}

		if data.Module == "blog" {
			blog, err := blogRepository.GetBlogBySlug(data.Slug)
			if err != nil {
				return err
			}

			_, err = blogRepository.CreateComment(blog.ID, comment.ID)
			if err != nil {
				return err
			}
		} else {
			return errors.New("module not supported")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *serviceComment) Reaction(commentId string, userId string, data *schemas.SchemaReaction) (any, error) {
	commentRepository := repositories.NewCommentRepository(s.db)

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	commentIdInt, err := strconv.Atoi(commentId)
	if err != nil {
		return nil, err
	}

	_, err = commentRepository.Reaction(uint(commentIdInt), userUUID, data)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *serviceComment) Reply(userId string, commentId string, data *schemas.SchemaCommentReply) (any, error) {
	commentRepository := repositories.NewCommentRepository(s.db)
	blogRepository := repositories.NewBlogRepository(s.db)

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	commentIdInt, err := strconv.Atoi(commentId)
	if err != nil {
		return nil, err
	}

	parentComment, err := commentRepository.GetById(uint(commentIdInt))
	if err != nil {
		return nil, err
	}

	newComment, err := commentRepository.Reply(parentComment.ID, userUUID, data)
	if err != nil {
		return nil, err
	}

	if data.Module == "blog" {
		commentBlog, err := blogRepository.GetCommentBlog(parentComment.ID)
		if err != nil {
			return nil, err
		}

		_, err = blogRepository.CreateComment(commentBlog.BlogId, newComment.ID)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("module not supported")
	}

	return nil, nil
}

func NewServiceComment(db *gorm.DB) *serviceComment {
	return &serviceComment{db: db}
}
