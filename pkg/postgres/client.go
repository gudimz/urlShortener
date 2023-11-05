package postgres

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gudimz/urlShortener/internal/pkg/ds"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewClient(ctx context.Context, pc ds.PostgresConfig) (*pgxpool.Pool, error) {
	hostPort := net.JoinHostPort(pc.Host, pc.Port)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s",
		pc.Username,
		pc.Password,
		hostPort,
		pc.DBName)
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return pool, nil
}
