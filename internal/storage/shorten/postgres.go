package shorten

import (
	"github.com/gudimz/urlShortener/internal/db/postgres"
	"github.com/gudimz/urlShortener/pkg/logging"
	"strings"
)

type psql struct {
	db     postgres.Client
	logger *logging.Logger
}

func queryForLogger(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}
