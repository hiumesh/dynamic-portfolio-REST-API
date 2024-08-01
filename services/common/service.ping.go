package services

import repositories "github.com/hiumesh/dynamic-portfolio-REST-API/repositories/common"

type ServicePing interface {
	PingService() string
}

type servicePing struct {
	repository repositories.RepositoryPing
}

func NewPingService(repository repositories.RepositoryPing) *servicePing {
	return &servicePing{repository: repository}
}

func (s *servicePing) PingService() string {
	res := s.repository.PingRepository()
	return res
}
