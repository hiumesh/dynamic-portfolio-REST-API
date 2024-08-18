package services

import (
	"github.com/hiumesh/dynamic-portfolio-REST-API/models"
	repositories "github.com/hiumesh/dynamic-portfolio-REST-API/repositories/user_profile"
)

type ServiceGetProfile interface {
	GetProfileService(string) (*models.UserProfile, error)
}

type serviceGetProfile struct {
	repository repositories.RepositoryGetProfile
}

func NewGetProfileService(repository repositories.RepositoryGetProfile) *serviceGetProfile {
	return &serviceGetProfile{
		repository: repository,
	}
}

func (r *serviceGetProfile) GetProfileService(userId string) (*models.UserProfile, error) {
	res, err := r.repository.GetProfileRepository(userId)
	if err != nil {
		return nil, err
	}
	return res, nil
}
