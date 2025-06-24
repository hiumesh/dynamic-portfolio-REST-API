package repositories

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"gorm.io/gorm"
)

type RepositoryComment interface {
	Get(userId *string, module string, slug string, cursor int, limit int, parentId *int) (*[]schemas.SelectComment, error)
	GetById(id uint) (*models.Comment, error)
	Create(userId string, data *schemas.SchemaComment) (*models.Comment, error)
	Reaction(commentId uint, userId uuid.UUID, data *schemas.SchemaReaction) (any, error)
	Reply(commentId uint, userId uuid.UUID, data *schemas.SchemaCommentReply) (any, error)
}

type repositoryComment struct {
	db *gorm.DB
}

func (r *repositoryComment) Get(userId *string, module string, slug string, cursor int, limit int, parentId *int) (*[]schemas.SelectComment, error) {

	var rows *sql.Rows
	var err error

	baseQuery := `
		select
			comments.id,
			comments.parent_id,
			comments.body,
			user_profiles.user_id as author_id,
			user_profiles.full_name as author_name,
			user_profiles.avatar_url as author_avatar,
			comments.attributes,
	`

	if userId != nil {
		baseQuery += `
			array_remove(array_agg(comment_reactions.type), NULL) as reactions,
		`
	} else {
		baseQuery += `
			null as reactions,
		`
	}

	baseQuery += `
			comments.created_at,
			comments.updated_at
		from
			comments
			inner join user_profiles on user_profiles.user_id = comments.user_id
	`

	var args []any

	if userId != nil {
		baseQuery += `
			left join comment_reactions on comment_reactions.comment_id = comments.id and comment_reactions.user_id = ?
		`
		args = append(args, *userId)
	}

	if module == "blog" {
		baseQuery += `
			inner join blog_comments on blog_comments.comment_id = comments.id
			inner join blogs on blogs.id = blog_comments.blog_id and blogs.slug = ?
		`
	} else {
		return nil, errors.New("invalid module")
	}

	baseQuery += `
		where
	`

	args = append(args, slug)

	if parentId != nil {
		baseQuery += "comments.parent_id = ?"
		args = append(args, parentId)
	} else {
		baseQuery += "comments.parent_id is null"
	}

	args = append(args, limit, cursor)

	baseQuery += `
		group by comments.id, user_profiles.user_id
		order by comments.created_at desc
		limit ?
		offset ?
	`

	rows, err = r.db.Raw(baseQuery, args...).Rows()

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var comments []schemas.SelectComment

	for rows.Next() {
		var comment schemas.SelectComment
		if err := rows.Scan(
			&comment.ID,
			&comment.ParentId,
			&comment.Body,
			&comment.AuthorId,
			&comment.AuthorName,
			&comment.AuthorAvatar,
			&comment.Attributes,
			&comment.Reactions,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return &comments, nil

}

func (r *repositoryComment) GetById(id uint) (*models.Comment, error) {
	var comment models.Comment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *repositoryComment) Create(userId string, data *schemas.SchemaCreateComment) (*models.Comment, error) {

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	comment := models.Comment{UserId: userUUID, Body: data.Body, Attributes: []byte("{}")}

	if err := r.db.Create(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (r *repositoryComment) Reaction(commentId uint, userId uuid.UUID, data *schemas.SchemaReaction) (any, error) {

	reaction := models.CommentReaction{
		UserId:    userId,
		CommentId: commentId,
		Type:      data.Reaction,
	}

	if data.Action == "remove" {
		if err := r.db.Delete(&models.CommentReaction{}, "comment_id = ? and user_id = ? and type = ?", commentId, userId, data.Reaction).Error; err != nil {
			return nil, err
		}
	}

	if data.Action == "add" {
		if err := r.db.Create(&reaction).Error; err != nil {
			return nil, err
		}
	}

	return reaction, nil
}

func (r *repositoryComment) Reply(commentId uint, userId uuid.UUID, data *schemas.SchemaCommentReply) (*models.Comment, error) {

	reply := models.Comment{
		UserId:     userId,
		ParentId:   &commentId,
		Body:       data.Body,
		Attributes: []byte("{}"),
	}

	if err := r.db.Create(&reply).Error; err != nil {
		return nil, err
	}

	return &reply, nil
}

func NewCommentRepository(db *gorm.DB) *repositoryComment {
	return &repositoryComment{
		db: db,
	}
}
