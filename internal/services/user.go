package services

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/models"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/pkg"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/repositories"
	"github.com/hiumesh/dynamic-portfolio-REST-API/internal/schemas"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ServiceUser interface {
	GetProfile(userId string) (*models.UserProfile, error)
	UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error
	ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error
	GetPostPresignedURLs(ctx context.Context, files []schemas.File) ([]any, error)
	GetFollowers(userId string, cursor int, limit int) (*[]schemas.SelectFollowers, error)
	GetFollowing(userId string, cursor int, limit int) (*[]schemas.SelectFollowing, error)
	FollowUser(userId string, followingUserId string) error
	UnfollowUser(userId string, followingUserId string) error
	FollowStatus(userId string, followingUserId string) (any, error)
}

type serviceUser struct {
	db        *gorm.DB
	presigner *pkg.Presigner
}

func (s *serviceUser) GetProfile(userId string) (*models.UserProfile, error) {
	repository := repositories.NewUserRepository(s.db)

	res, err := repository.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *serviceUser) ProfileSetup(userId string, profile *schemas.SchemaProfileBasic) error {
	repository := repositories.NewUserRepository(s.db)

	if err := repository.ProfileSetup(userId, profile); err != nil {
		return err
	}
	return nil
}

func (s *serviceUser) UpsertProfile(userId string, profile *schemas.SchemaProfileBasic) error {
	repository := repositories.NewUserRepository(s.db)

	if err := repository.UpsertProfile(userId, profile); err != nil {
		return err
	}
	return nil
}

func (s *serviceUser) GetPostPresignedURLs(ctx context.Context, files []schemas.File) ([]interface{}, error) {
	var wg sync.WaitGroup

	results := make(chan interface{}, len(files))
	errs := make(chan error, len(files))

	for _, file := range files {
		wg.Add(1)
		go func() {
			defer wg.Done()

			nameSplit := strings.Split(file.FileName, ".")
			if len(nameSplit) < 2 {
				errs <- errors.New("invalid file name")
				return
			}

			key := "public/" + strings.Join(strings.Fields(strings.ToLower(nameSplit[0])), "-") + "-" + strconv.FormatInt(time.Now().UnixMilli(), 10) + "." + nameSplit[len(nameSplit)-1]
			presignedPostRequest, err := s.presigner.PutObject(ctx, key, 120)

			if err != nil {
				errs <- err
			} else {
				results <- map[string]string{
					"file_name": file.FileName,
					"key":       key,
					"url":       presignedPostRequest.URL,
				}
			}
		}()
	}

	wg.Wait()

	close(results)
	close(errs)

	if len(errs) > 0 {
		for err := range errs {
			logrus.Error(err)
		}
		return nil, errors.New("failed to get presigned urls")
	}

	var urls []interface{}

	for url := range results {
		urls = append(urls, url)
	}

	return urls, nil
}

func (s *serviceUser) GetFollowers(userId string, cursor int, limit int) (*[]schemas.SelectFollowers, error) {
	repository := repositories.NewUserRepository(s.db)

	res, err := repository.GetFollowers(userId, cursor, limit)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *serviceUser) GetFollowing(userId string, cursor int, limit int) (*[]schemas.SelectFollowing, error) {
	repository := repositories.NewUserRepository(s.db)

	res, err := repository.GetFollowing(userId, cursor, limit)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *serviceUser) FollowUser(userId string, followingUserId string) error {
	repository := repositories.NewUserRepository(s.db)

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("failed to parse user id")
	}

	followingProfile, err := repository.GetProfileBySlug(followingUserId)
	if err != nil {
		return err
	}

	if err := repository.FollowUser(userUUID, followingProfile.UserId); err != nil {
		return err
	}
	return nil
}

func (s *serviceUser) UnfollowUser(userId string, followingUserId string) error {
	repository := repositories.NewUserRepository(s.db)

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return errors.New("failed to parse user id")
	}

	followingProfile, err := repository.GetProfileBySlug(followingUserId)
	if err != nil {
		return err
	}

	if err := repository.UnfollowUser(userUUID, followingProfile.UserId); err != nil {
		return err
	}
	return nil
}

func (s *serviceUser) FollowStatus(userId string, followingUserId string) (any, error) {
	repository := repositories.NewUserRepository(s.db)

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, errors.New("failed to parse user id")
	}

	followingProfile, err := repository.GetProfileBySlug(followingUserId)
	if err != nil {
		return nil, err
	}

	res, err := repository.FollowStatus(userUUID, followingProfile.UserId)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return map[string]any{"is_following": true}, nil
	}

	return map[string]any{"is_following": false}, nil

}
func NewUserService(db *gorm.DB, presigner *pkg.Presigner) *serviceUser {
	return &serviceUser{
		db:        db,
		presigner: presigner,
	}
}
