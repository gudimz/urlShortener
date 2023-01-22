package shorten

import (
	"context"
	"github.com/google/uuid"
	"github.com/gudimz/urlShortener/internal/model"
	"time"
)

type Storage interface {
	CreateShorten(ctx context.Context, ms model.Shorten) error
	GetShorten(ctx context.Context, shortUrl string) (*model.Shorten, error)
	DeleteShorten(ctx context.Context, shortUrl string) (int64, error)
	UpdateShorten(ctx context.Context, shortUrl string) error
}

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) CreateShorten(ctx context.Context, input model.InputShorten) (*model.Shorten, error) {
	var (
		id       = uuid.New().ID()
		shortUrl = input.ShortenUrl.OrElse(GenerateShortenUrl(id))
	)

	shorten := model.Shorten{
		ShortUrl:    shortUrl,
		OriginUrl:   input.OriginUrl,
		Visits:      0,
		DateCreated: time.Now().UTC(),
		DateUpdated: time.Now().UTC(),
	}

	err := s.storage.CreateShorten(ctx, shorten)
	if err != nil {
		return nil, err
	}

	return &shorten, nil
}

func (s *Service) GetShorten(ctx context.Context, shortUrl string) (*model.Shorten, error) {
	shorten, err := s.storage.GetShorten(ctx, shortUrl)
	if err != nil {
		return nil, err
	}
	return shorten, nil
}

func (s *Service) Redirect(ctx context.Context, shortUrl string) (string, error) {
	shorten, err := s.storage.GetShorten(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	err = s.storage.UpdateShorten(ctx, shortUrl)
	if err != nil {
		return shorten.OriginUrl, err
	}

	return shorten.OriginUrl, nil
}

func (s *Service) DeleteShorten(ctx context.Context, shortUrl string) (int64, error) {
	return s.storage.DeleteShorten(ctx, shortUrl)
}
