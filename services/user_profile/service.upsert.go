package services

import (
	repositories "github.com/hiumesh/dynamic-portfolio-REST-API/repositories/user_profile"
	"github.com/hiumesh/dynamic-portfolio-REST-API/schemas"
)

type ServiceUpsertProfile interface {
	UpsertProfileService(string, *schemas.SchemaProfileBasic) error
}

type serviceUpsertProfile struct {
	repository repositories.RepositoryUpsertProfile
}

func (s *serviceUpsertProfile) UpsertProfileService(userId string, profile *schemas.SchemaProfileBasic) error {
	if err := s.repository.UpsertProfileRepository(userId, profile); err != nil {
		return err
	}
	return nil
}

func NewUpsertProfileService(repository repositories.RepositoryUpsertProfile) *serviceUpsertProfile {
	return &serviceUpsertProfile{
		repository: repository,
	}
}
