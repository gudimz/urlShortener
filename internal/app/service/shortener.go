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
		shortURL = input.ShortenURL.OrElse(generateShortenURL(id))
	)

	shorten, err := s.repository.CreateShorten(ctx, &ds.Shorten{
		ShortURL:  shortURL,
		OriginURL: input.OriginURL,
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

	return models.ModelFromDBShorten(shorten), nil
}

func (s *Service) GetShorten(ctx context.Context, shortURL string) (*ds.Shorten, error) {
	shorten, err := s.repository.GetShorten(ctx, shortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ds.ErrShortURLNotFound
		}

		return nil, err
	}

	return models.ModelFromDBShorten(shorten), nil
}

func (s *Service) Redirect(ctx context.Context, shortURL string) (string, error) {
	dbShorten, err := s.repository.GetShorten(ctx, shortURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ds.ErrShortURLNotFound
		}

		return "", err
	}

	shorten := models.ModelFromDBShorten(dbShorten)

	err = s.repository.UpdateShorten(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return shorten.OriginURL, nil
}

func (s *Service) DeleteShorten(ctx context.Context, shortURL string) error {
	count, err := s.repository.DeleteShorten(ctx, shortURL)
	if err == nil && count == 0 {
		return ds.ErrShortURLNotFound
	}
	return err
}
