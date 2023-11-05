package service

import (
	"context"

	"github.com/gudimz/urlShortener/internal/app/repository/psql/models"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
)

//go:generate mockgen -source=service.go -destination=service_mock.go -package=service
type Repository interface {
	CreateShorten(context.Context, *ds.Shorten) (*models.DBShorten, error)
	GetShorten(context.Context, string) (*models.DBShorten, error)
	DeleteShorten(context.Context, string) (int64, error)
	UpdateShorten(context.Context, string) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
