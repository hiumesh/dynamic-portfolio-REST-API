package seeds

import (
	"math/rand/v2"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"gorm.io/gorm"
)

type Comment struct {
	ID        int64     `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"user_id"`
	ParentId  *int64    `json:"parent_id"`
	Body      string    `json:"body" fake:"{paragraph:2,3,10,\n\n}"`
	CreatedAt time.Time `json:"created_at" fake:"{date}"`
	UpdatedAt time.Time `json:"updated_at" fake:"{date}"`
}

func (Comment) TableName() string {
	return "comments"
}

type CommentReaction struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	UserId    string `json:"user_id"`
	CommentId int64  `json:"comment_id"`
	Type      string `json:"type" fake:"{randomstring:[like,clap,heart]}"`
}

func (CommentReaction) TableName() string {
	return "comment_reactions"
}

func NewBlogComments(blogs *[]Blog, users *[]User) (*[]Comment, map[int64][]int64, error) {

	n := len(*users)

	blogIdToCommentsIdxMap := make(map[int64][]int64)
	comments := []Comment{}

	for _, blog := range *blogs {

		k := rand.IntN(3) + 4

		blogIdToCommentsIdxMap[blog.ID] = make([]int64, k)
		uniqueIdx := make(map[int]bool)
		if k < n {
			for len(uniqueIdx) < k {
				idx := rand.IntN(n)
				if _, ok := uniqueIdx[idx]; !ok {
					comment := Comment{}
					if err := gofakeit.Struct(&comment); err != nil {
						return nil, nil, err
					}

					comment.ID = 0
					comment.UserId = (*users)[idx].Id
					comment.ParentId = nil
					comments = append(comments, comment)

					blogIdToCommentsIdxMap[blog.ID][len(uniqueIdx)] = int64(len(comments) - 1)
					uniqueIdx[idx] = true
				}
			}
		}

	}

	return &comments, blogIdToCommentsIdxMap, nil

}

func NewBlogChildComments(comments *[]BlogComment, users *[]User) (*[]Comment, map[int64][]int64, error) {

	blogIdToChildCommentsIdxMap := make(map[int64][]int64)
	childComments := []Comment{}
	n := len(*users)

	for _, comment := range *comments {

		blogIdToChildCommentsIdxMap[comment.BlogId] = make([]int64, 0)
	}

	for _, comment := range *comments {

		k := rand.IntN(3) + 4

		uniqueIdx := make(map[int]bool)
		if k < n {
			for len(uniqueIdx) < k {
				idx := rand.IntN(n)
				if _, ok := uniqueIdx[idx]; !ok {
					childComment := Comment{}

					if err := gofakeit.Struct(&childComment); err != nil {
						return nil, nil, err
					}

					childComment.ID = 0
					childComment.UserId = (*users)[idx].Id
					childComment.ParentId = &comment.CommentId
					childComments = append(childComments, childComment)

					blogIdToChildCommentsIdxMap[comment.BlogId] = append(blogIdToChildCommentsIdxMap[comment.BlogId], int64(len(childComments)-1))
					uniqueIdx[idx] = true
				}
			}
		}
	}

	return &childComments, blogIdToChildCommentsIdxMap, nil

}

func NewCommentsReactions(db *gorm.DB, comments *[]Comment, users *[]User) error {

	n := len(*users)
	commentsReactions := make([]CommentReaction, 0)

	for _, comment := range *comments {

		k := rand.IntN(3) + 4

		uniqueIdx := make(map[int]bool)
		if k < n {
			for len(uniqueIdx) < k {
				idx := rand.IntN(n)
				if _, ok := uniqueIdx[idx]; !ok {
					commentReaction := CommentReaction{}
					if err := gofakeit.Struct(&commentReaction); err != nil {
						return err
					}

					commentReaction.ID = 0
					commentReaction.UserId = (*users)[idx].Id
					commentReaction.CommentId = comment.ID
					commentsReactions = append(commentsReactions, commentReaction)
					uniqueIdx[idx] = true
				}
			}
		}
	}

	if err := db.Create(&commentsReactions).Error; err != nil {
		return err
	}

	return nil
}
