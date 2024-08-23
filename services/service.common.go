package services

import "gorm.io/gorm"

type ServiceCommon interface {
	Ping() string
}

type serviceCommon struct {
	db *gorm.DB
}

func (s *serviceCommon) Ping() string {
	return "ping service"
}

func NewCommonService(db *gorm.DB) *serviceCommon {
	return &serviceCommon{
		db: db,
	}
}
