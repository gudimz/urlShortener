package shorten

import (
	"context"
	"github.com/gudimz/urlShortener/internal/model"
	"strings"
)

const (
	dictionary    = "abcdefghjkmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	lenDictionary = uint32(len(dictionary))
)

type Storage interface {
	CreateShorten(ctx context.Context, ms model.Shorten) error
	GetShorten(ctx context.Context, id string) (*model.Shorten, error)
	DeleteShorten(ctx context.Context, id string) error
	UpdateShorten(ctx context.Context, id string) error
}

func GenerateShortenUrl(id uint32) string {
	var (
		builder strings.Builder
		nums    []uint32
	)

	for id > 0 {
		nums = append(nums, id%lenDictionary)
		id /= lenDictionary
	}

	for _, num := range nums {
		builder.WriteByte(dictionary[num])
	}

	return builder.String()
}
