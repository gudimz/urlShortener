package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/gudimz/urlShortener/internal/app/repository/psql/models"
	"github.com/gudimz/urlShortener/internal/pkg/ds"
)

func (s *Service) CreateShorten(ctx context.Context, input ds.InputShorten) (*ds.Shorten, error) {
	var (
		id       = uuid.New().ID()
		shortUrl = input.ShortenUrl.OrElse(generateShortenURL(id))
	)

	shorten, err := s.repository.CreateShorten(ctx, &ds.Shorten{
		ShortUrl:  shortUrl,
		OriginUrl: input.OriginUrl,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, ds.ErrShortURLAlreadyExists
			}
		}

		return nil, err
	}

	return models.ModelFromDbShorten(shorten), nil
}

func (s *Service) GetShorten(ctx context.Context, shortUrl string) (*ds.Shorten, error) {
	shorten, err := s.repository.GetShorten(ctx, shortUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ds.ErrShortUrlNotFound
		}

		return nil, err
	}

	return models.ModelFromDbShorten(shorten), nil
}

func (s *Service) Redirect(ctx context.Context, shortUrl string) (string, error) {
	dbShorten, err := s.repository.GetShorten(ctx, shortUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ds.ErrShortUrlNotFound
		}

		return "", err
	}

	shorten := models.ModelFromDbShorten(dbShorten)

	err = s.repository.UpdateShorten(ctx, shortUrl)
	if err != nil {
		return "", err
	}

	return shorten.OriginUrl, nil
}

func (s *Service) DeleteShorten(ctx context.Context, shortUrl string) error {
	count, err := s.repository.DeleteShorten(ctx, shortUrl)
	if err == nil && count == 0 {
		return ds.ErrShortUrlNotFound
	}
	return err
}
